package profile

import (
	"context"
	"fmt"
	"mime/multipart"
	"modular_monolith/internal/cloudinaryutil"
	fileutil "modular_monolith/internal/common"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileService interface {
	CreateProfile(ctx context.Context, req *CreateProfileRequest, file *multipart.FileHeader) error
	UpdateProfile(ctx context.Context, req *UpdateProfileRequest, file *multipart.FileHeader) error
	DeleteProfile(ctx context.Context, userID string) error
}

type profileService struct {
	repository    ProfileRepository
	cloudUploader *cloudinaryutil.CloudinaryUploader
}

func NewProfileService(repository ProfileRepository,
	uploader *cloudinaryutil.CloudinaryUploader) ProfileService {
	return &profileService{
		repository:    repository,
		cloudUploader: uploader,
	}
}

func (s *profileService) CreateProfile(ctx context.Context, req *CreateProfileRequest, file *multipart.FileHeader) error {

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user_id: %v", err)
	}

	existingUser, _ := s.repository.FindByUserID(ctx, userID)
	if existingUser != nil {
		return fmt.Errorf("profile of this user already exists")
	}

	if req.FullName == "" {
		return fmt.Errorf("full name is required")
	}

	if req.Gender == "" {
		return fmt.Errorf("gender is required")
	}

	birthDay, err := time.Parse(time.RFC3339, req.BirthDay)
	if err != nil {
		return fmt.Errorf("invalid birth day format: %v", err)
	}

	if req.Address == "" {
		return fmt.Errorf("address is required")
	}

	tempPath := "/tmp/" + file.Filename
	if err := fileutil.SaveUploadedFile(file, tempPath); err != nil {
		return fmt.Errorf("failed to save avatar: %w", err)
	}
	defer os.Remove(tempPath)

	avatarURL, publicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "profiles")
	if err != nil {
		return fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	profile := &Profile{
		UserID:   userID,
		FullName: req.FullName,
		Gender:   req.Gender,
		BirthDay: birthDay,
		Avatar:   avatarURL,
		Address:  req.Address,
		PublicID: publicID,
		Bio:      &req.Bio,
	}

	return s.repository.Create(ctx, profile)
}

func (s *profileService) UpdateProfile(ctx context.Context, req *UpdateProfileRequest, file *multipart.FileHeader) error {
	
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user_id: %v", err)
	}

	existingUser, _ := s.repository.FindByUserID(ctx, userID)
	if existingUser == nil {
		return fmt.Errorf("profile of this user does not exist")
	}

	if req.FullName == "" {
		return fmt.Errorf("full name is required")
	}

	if req.Gender == "" {
		return fmt.Errorf("gender is required")
	}

	birthDay, err := time.Parse(time.RFC3339, req.BirthDay)
	if err != nil {
		return fmt.Errorf("invalid birth day format: %v", err)
	}

	if req.Address == "" {
		return fmt.Errorf("address is required")
	}

	tempPath := "/tmp/" + file.Filename
	if err := fileutil.SaveUploadedFile(file, tempPath); err != nil {
		return fmt.Errorf("failed to save avatar: %w", err)
	}
	defer os.Remove(tempPath)

	avatarURL, err := s.cloudUploader.EditImage(ctx, tempPath, existingUser.PublicID)
	if err != nil {
		return fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	profile := &Profile{
		UserID:   userID,
		FullName: req.FullName,
		Gender:   req.Gender,
		BirthDay: birthDay,
		Avatar:   avatarURL,
		Address:  req.Address,
		PublicID: existingUser.PublicID,
		Bio:      &req.Bio,
	}

	return s.repository.Update(ctx, profile)

}

func (s *profileService) DeleteProfile(ctx context.Context, userID string) error {
	
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	existingUser, _ := s.repository.FindByUserID(ctx, objectID)
	if existingUser == nil {
		return fmt.Errorf("profile of this user does not exist")
	}

	err = s.cloudUploader.DeleteImage(ctx, existingUser.PublicID)
	if err != nil {
		return err
	}
	
	return s.repository.DeleteByID(ctx, objectID)
}
