package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"modular_monolith/internal/profile"
	"modular_monolith/pkg/email"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *RegisterRequest) (*User, error)
	LoginUser(ctx context.Context, email, password string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*UserWithProfile, error)
	GetAllUsers(ctx context.Context) ([]*UserWithProfile, error)
	DeleteUser(ctx context.Context, userID string) error
	ValidateToken(tokenString string) (*jwt.Token, error)
	RefreshToken(refreshToken string) (string, string, error)
	LogoutUser(ctx context.Context, userID string) error
	ChangePassword(ctx context.Context, req ChangePasswordRequest, userID string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
}

type userService struct {
	repository     UserRepository
	profileService profile.ProfileService
	EmailService   *email.EmailService
}

func NewUserService(repository UserRepository, profileService profile.ProfileService) UserService {
	emailService := email.NewEmailService()
	return &userService{
		repository:     repository,
		profileService: profileService,
		EmailService:   emailService,
	}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*UserWithProfile, error) {

	var userAll []*UserWithProfile

	user, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, u := range user {
		userWithProfile, err := s.repository.FindByID(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		userAll = append(userAll, userWithProfile)
	}

	return userAll, nil

}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	err = s.repository.DeleteByID(ctx, objectID)
	if err != nil {
		return err
	}

	err = s.profileService.DeleteProfile(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*UserWithProfile, error) {

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return s.repository.FindByID(ctx, objectID)

}

func (s *userService) RegisterUser(ctx context.Context, req *RegisterRequest) (*User, error) {

	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if req.Phone == "" {
		return nil, fmt.Errorf("phone is required")
	}

	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	user, err := s.repository.FindByEmail(ctx, req.Email)

	if user != nil {
		return nil, fmt.Errorf("user already exists")
	}

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	hashedPassword := s.HashPassword(req.Password)
	newUserID := primitive.NewObjectID()
	token, refreshToken := s.GenerateToken(newUserID.Hex())

	now := time.Now().Format(time.RFC3339)
	user = &User{
		ID:           newUserID,
		FristName:    req.FristName,
		LastName:     req.LastName,
		Email:        req.Email,
		Password:     hashedPassword,
		Phone:        req.Phone,
		Token:        token,
		RefreshToken: refreshToken,
		UserType:     "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	createdUser.Password = ""
	return createdUser, nil
}

func (s *userService) LoginUser(ctx context.Context, email, password string) (*User, error) {

	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	isValid, _ := s.VerifyPassword(user.Password, password)
	if !isValid {
		return nil, fmt.Errorf("invalid email or password")
	}

	token, refreshToken := s.GenerateToken(user.ID.Hex())

	updateFields := bson.M{
		"token":         token,
		"refresh_token": refreshToken,
		"updatedAt":     time.Now().Format(time.RFC3339),
	}

	err = s.repository.UpdateByID(ctx, user.ID, updateFields)
	if err != nil {
		return nil, fmt.Errorf("failed to update user tokens: %w", err)
	}

	user.Token = token
	user.RefreshToken = refreshToken
	user.Password = "" // Don't return the password

	return user, nil
}

func (s *userService) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func (s *userService) VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Login or password is incorrect"
		check = false
	}
	return check, msg
}

func (s *userService) GenerateToken(userID string) (string, string) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Panic("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Panic(err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		log.Panic(err)
	}

	return tokenString, refreshTokenString
}

func (s *userService) GenerateResetToken(userID string) (string, error) {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (s *userService) ValidateToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, errors.New("token is required")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (s *userService) RefreshToken(refreshToken string) (string, string, error) {
	token, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	user_id, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("invalid email in token")
	}

	newToken, newRefreshToken := s.GenerateToken(user_id)
	return newToken, newRefreshToken, nil
}

func (s *userService) LogoutUser(ctx context.Context, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	updateFields := bson.M{
		"token":         "",
		"refresh_token": "",
		"updated_at":    time.Now().Format(time.RFC3339),
	}

	err = s.repository.UpdateByID(ctx, objectID, updateFields)
	if err != nil {
		return fmt.Errorf("failed to update user tokens: %w", err)
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, req ChangePasswordRequest, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	existingUser, _ := s.repository.FindByUserID(ctx, objectID)
	if existingUser == nil {
		return fmt.Errorf("profile of this user does not exist")
	}

	check, _ := s.VerifyPassword(existingUser.Password, req.OldPassword)
	if !check {
		return fmt.Errorf("old password is incorrect")
	}

	updateFields := bson.M{
		"password":   s.HashPassword(req.NewPassword),
		"updated_at": time.Now().Format(time.RFC3339),
	}

	err = s.repository.UpdateByID(ctx, objectID, updateFields)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil
}

func (s *userService) ForgotPassword(ctx context.Context, email string) error {

	if email == "" {
		return fmt.Errorf("email is required")
	}

	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return fmt.Errorf("invalid email")
	}

	token, err := s.GenerateResetToken(user.ID.Hex())
	if err != nil {
		return err
	}

	err = s.EmailService.SendEmailChangePassword(email, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {

	if req.Token == "" {
		return fmt.Errorf("token is required")
	}

	if req.NewPassword == "" {
		return fmt.Errorf("new password is required")
	}

	token, err := s.ValidateToken(req.Token)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	user_id, ok := claims["user_id"].(string)
	if !ok {
		return errors.New("invalid email in token")
	}

	objectID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return err
	}

	updateFields := bson.M{
		"password":   s.HashPassword(req.NewPassword),
		"updated_at": time.Now().Format(time.RFC3339),
	}

	err = s.repository.UpdateByID(ctx, objectID, updateFields)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil

}
