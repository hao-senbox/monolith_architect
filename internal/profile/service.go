package profile

import (
	"context"
	"fmt"
	"mime/multipart"
	"modular_monolith/helper"
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
	cloudUploader *helper.CloudinaryUploader
}

func NewProfileService(repository ProfileRepository,
	uploader *helper.CloudinaryUploader) ProfileService {
	return &profileService{
		repository:    repository,
		cloudUploader: uploader,
	}
}

func (s *profileService) CreateProfile(ctx context.Context, req *CreateProfileRequest, file *multipart.FileHeader) error {
	fmt.Printf("check: %v", req.UserID)
	fmt.Println(req.UserID)
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user_id: %v", err)
	}

	existingUser, _ := s.repository.FindByUserID(ctx, userID)
	if existingUser != nil {
		return fmt.Errorf("profile of this user already exists")
	}

	if req.Gender == "" {
		return fmt.Errorf("gender is required")
	}

	birthDay, err := time.Parse("2006-01-02", req.BirthDay)
	if err != nil {
		return fmt.Errorf("invalid birth day format: %v", err)
	}

	if req.Address == "" {
		return fmt.Errorf("address is required")
	}

	var avatarURL, publicID string
	if file == nil {
		avatarURL = ""
		publicID = ""
	} else {
		tempPath := "/tmp/" + file.Filename
		if err := helper.SaveUploadedFile(file, tempPath); err != nil {
			return fmt.Errorf("failed to save avatar: %w", err)
		}
		defer os.Remove(tempPath)

		avatarURL, publicID, err = s.cloudUploader.UploadImage(ctx, tempPath, "profiles")
		if err != nil {
			return fmt.Errorf("failed to upload to Cloudinary: %w", err)
		}
	}

	profile := &Profile{
		UserID:   userID,
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

	profile := &Profile{
		UserID:   existingUser.UserID,
		Gender:   existingUser.Gender,
		BirthDay: existingUser.BirthDay,
		Avatar:   existingUser.Avatar,
		Address:  existingUser.Address,
		PublicID: existingUser.PublicID,
		Bio:      existingUser.Bio,
	}

	if req.Gender != "" {
		profile.Gender = req.Gender
	}

	if req.BirthDay != "" {
		birthDay, err := time.Parse("2006-01-02", req.BirthDay)
		if err != nil {
			return fmt.Errorf("invalid birth day format: %v", err)
		}
		profile.BirthDay = birthDay
	}

	if req.Address != "" {
		profile.Address = req.Address
	}

	if req.Bio != "" {
		profile.Bio = &req.Bio
	}

	if file != nil {
		tempPath := "/tmp/" + file.Filename
		if err := helper.SaveUploadedFile(file, tempPath); err != nil {
			return fmt.Errorf("failed to save avatar: %w", err)
		}
		defer os.Remove(tempPath)

		if existingUser.PublicID != "" {
			if err = s.cloudUploader.DeleteImage(ctx, existingUser.PublicID); err != nil {
				return fmt.Errorf("failed to delete old image: %w", err)
			}
		}

		avatarURL, publicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "profiles")
		if err != nil {
			return fmt.Errorf("failed to upload to Cloudinary: %w", err)
		}
		profile.Avatar = avatarURL
		profile.PublicID = publicID
	}

	return s.repository.UpdateByID(ctx, profile)
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
