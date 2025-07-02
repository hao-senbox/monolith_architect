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
		TotalLike:     0,
		TotalDislike:  0,
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