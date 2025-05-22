package profile

import (
	"context"
	"fmt"
	"time"
)

type ProfileService interface {
	CreateProfile(ctx context.Context, req *CreateProfileRequest) error
}

type profileService struct {
	repository ProfileRepository
}

func NewProfileService(repository ProfileRepository) ProfileService {
	return &profileService{repository: repository}
}

func (s *profileService) CreateProfile(ctx context.Context, req *CreateProfileRequest) error {
	
	if req.UserID.IsZero() {
		return fmt.Errorf("user ID is required")
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

	profile := &Profile {
		UserID:    req.UserID,
		FullName:  req.FullName,
		Gender:    req.Gender,
		BirthDay:  birthDay,
		Avatar:    req.Avatar,
		Address:   req.Address,
		Bio:       &req.Bio,
	}
	
	return s.repository.Create(ctx, profile)
}

