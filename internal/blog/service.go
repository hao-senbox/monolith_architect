package blog

import (
	"fmt"
	"mime/multipart"
	"modular_monolith/helper"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService interface {
	CreateBlog(c *gin.Context, req *CreateBlogRequest, mainImage *multipart.FileHeader) error
	GetAllBlog(c *gin.Context) ([]*Blog, error)
	GetBlogByID(c *gin.Context, id string) (*Blog, error)
	UpdateBlog(c *gin.Context, id string, req *UpdateBlogRequest, mainImage *multipart.FileHeader) error
	DeleteBlog(c *gin.Context, id string) error
	LikeBlog(c *gin.Context, id string, userID primitive.ObjectID) error
	ViewBlog(c *gin.Context, id string) error
}

type blogService struct {
	blogRepo BlogRepository
	cloudUploader *helper.CloudinaryUploader
}

func NewBlogService(repo BlogRepository, uploader *helper.CloudinaryUploader) BlogService {
	return &blogService{
		blogRepo: repo,
		cloudUploader: uploader,
	}
}

func (s *blogService) CreateBlog(c *gin.Context, req *CreateBlogRequest, mainImage *multipart.FileHeader) error {
	
	var mainImageURL, mainImagePublicID string
	var err error
	
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}

	if req.Content == "" {
		return fmt.Errorf("content is required")
	}

	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if mainImage != nil {
		tempPath := "/tmp/" + mainImage.Filename
		if err := helper.SaveUploadedFile(mainImage, tempPath); err != nil {
			return fmt.Errorf("failed to save main image: %w", err)
		}
		defer os.Remove(tempPath)

		mainImageURL, mainImagePublicID, err = s.cloudUploader.UploadImage(c, tempPath, "blogs")
		if err != nil {
			return fmt.Errorf("failed to upload main image: %w", err)
		}

	}

	objectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user_id: %v", err)
	}

	blog := &Blog{
		ID:            primitive.NewObjectID(),
		Title:         req.Title,
		Content:       req.Content,
		ImageURL:      mainImageURL,
		ImagePublicID: mainImagePublicID,
		AuthorID:      objectID,
		TotalView:     0,
		LikedUsers:    []primitive.ObjectID{},
		DislikedUsers: []primitive.ObjectID{},
		Created:       time.Now(),
		Updated:       time.Now(),
	}

	return s.blogRepo.Create(c, blog)

}

func (s *blogService) GetAllBlog(c *gin.Context) ([]*Blog, error) {
	return s.blogRepo.FindAll(c)
}

func (s *blogService) GetBlogByID(c *gin.Context, id string) (*Blog, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	return s.blogRepo.FindID(c, objectID)

}

func (s *blogService) UpdateBlog(c *gin.Context, id string, req *UpdateBlogRequest, mainImage *multipart.FileHeader) error{

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	blog, err := s.blogRepo.FindID(c, objectID)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	if req.Content == "" {
		return fmt.Errorf("content is required")
	}

	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if req.Title == "" {
		return fmt.Errorf("title is required")
	}

	var mainImageURL, mainImagePublicID string
	if mainImage != nil {

		tempPath := "/tmp/" + mainImage.Filename
		if err := helper.SaveUploadedFile(mainImage, tempPath); err != nil {
			return fmt.Errorf("failed to save main image: %w", err)
		}
		defer os.Remove(tempPath)

		err = s.cloudUploader.DeleteImage(c, blog.ImagePublicID)
		if err != nil {
			return fmt.Errorf("failed to upload to Cloudinary: %w", err)
		}

		mainImageURL, mainImagePublicID, err = s.cloudUploader.UploadImage(c, tempPath, "blogs")
		if err != nil {
			return fmt.Errorf("failed to upload main image: %w", err)
		}
	} else {

		mainImageURL = blog.ImageURL
		mainImagePublicID = blog.ImagePublicID
	}

	blog.Title = req.Title
	blog.Content = req.Content
	blog.ImageURL = mainImageURL
	blog.ImagePublicID = mainImagePublicID
	blog.Updated = time.Now()

	return s.blogRepo.UpdateByID(c, objectID, blog)

}

func (s *blogService) DeleteBlog(c *gin.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	blog, err := s.blogRepo.FindID(c, objectID)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	err = s.cloudUploader.DeleteImage(c, blog.ImagePublicID)
	if err != nil {
		return fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	return s.blogRepo.DeleteByID(c, objectID)

}

func (s *blogService) LikeBlog(c *gin.Context, id string, userID primitive.ObjectID) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	blog, err := s.blogRepo.FindID(c, objectID)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	liked := s.containsObjectID(blog.LikedUsers, userID)
	disliked := s.containsObjectID(blog.DislikedUsers, userID)

	if liked {
		blog.LikedUsers = s.removeObjectID(blog.LikedUsers, userID)
		blog.DislikedUsers = append(blog.DislikedUsers, userID)
	} else {
		blog.LikedUsers = append(blog.LikedUsers, userID)
		if disliked {
			blog.DislikedUsers = s.removeObjectID(blog.DislikedUsers, userID)
		}
	}

	return s.blogRepo.UpdateByID(c, objectID, blog)
}

func (s *blogService) containsObjectID(list []primitive.ObjectID, id primitive.ObjectID) bool {

	for _, item := range list {
		if item == id {
			return true
		}
	}

	return false

}

func (s *blogService) removeObjectID(list []primitive.ObjectID, id primitive.ObjectID) []primitive.ObjectID {

	newList := []primitive.ObjectID{}
	
	for _, v := range list {
		if v != id {
			newList = append(newList, v)
		}
	}

	return newList


}

func(s *blogService) ViewBlog(c *gin.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	return s.blogRepo.IncrementViews(c, objectID)
}