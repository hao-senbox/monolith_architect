package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"modular_monolith/internal/profile"
	"os"
	"strings"
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
}

type userService struct {
	repository UserRepository
	profileService profile.ProfileService
}

func NewUserService(repository UserRepository, profileService profile.ProfileService) UserService {
	return &userService{
		repository: repository,
		profileService: profileService,
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
	fmt.Printf("RegisterUser: %+v\n", req)
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
	token, refreshToken := s.GenerateToken(newUserID.String())

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


	token, refreshToken := s.GenerateToken(user.ID.String())
	
	updateFields := bson.M{
		"token":        token,
		"refreshToken": refreshToken,
		"updatedAt":    time.Now().Format(time.RFC3339),
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

	if strings.HasPrefix(userID, "ObjectID(") {
		userID = strings.TrimPrefix(userID, "ObjectID(\"")
		userID = strings.TrimSuffix(userID, "\")")
	}
	
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	
	updateFields := bson.M{
		"token":        "",
		"refreshToken": "",
		"updatedAt":    time.Now().Format(time.RFC3339),
	}
	
	err = s.repository.UpdateByID(ctx, objectID, updateFields)
	if err != nil {
		return fmt.Errorf("failed to update user tokens: %w", err)
	}
	
	return nil
}