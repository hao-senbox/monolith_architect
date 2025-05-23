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
}

type profileService struct {
	repository ProfileRepository
	cloudUploader *cloudinaryutil.CloudinaryUploader
}

func NewProfileService(repository ProfileRepository, 
	uploader *cloudinaryutil.CloudinaryUploader) ProfileService {
	return &profileService{
		repository: repository,
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

	avatarURL, err := s.cloudUploader.UploadImage(ctx, tempPath, "profiles")
	if err != nil {
		return fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	profile := &Profile {
		UserID:    userID,
		FullName:  req.FullName,
		Gender:    req.Gender,
		BirthDay:  birthDay,
		Avatar:    avatarURL,
		Address:   req.Address,
		Bio:       &req.Bio,
	}
	
	return s.repository.Create(ctx, profile)
}

