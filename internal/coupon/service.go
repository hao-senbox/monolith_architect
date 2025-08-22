package coupon

import (
	"fmt"
	"math/rand"
	"modular_monolith/internal/user"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CouponService interface {
	CreateCoupon(c *gin.Context, req *CreateCouponRequest) error
	GetAllCoupons(c *gin.Context) ([]*CouponResponse, error)
	GetCouponByCode(c *gin.Context, code string) (*CouponResponse, error)
	GetCouponByUserID(c *gin.Context, userID string) ([]*CouponResponse, error)
	CanUseCoupon(c *gin.Context, req *CanUseCouponRequest) (*Coupon, error)
	DeleteCoupon(c *gin.Context, id string) error
}

type couponService struct {
	couponRepository CouponRepository
	userService      user.UserService
}

func NewCouponService(couponRepository CouponRepository, userService user.UserService) CouponService {
	return &couponService{
		couponRepository: couponRepository,
		userService:      userService,
	}
}

func (s *couponService) CreateCoupon(c *gin.Context, req *CreateCouponRequest) error {

	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	if req.Discount <= 0 {
		return fmt.Errorf("discount must be greater than 0")
	}

	var allowedUsers []primitive.ObjectID

	switch req.Type {
	case "public":
		if req.MaximumUse <= 0 {
			return fmt.Errorf("maximum use must be greater than 0")
		}
	case "private":
		for _, userIDStr := range req.AllowedUsers {
			objectUserID, err := primitive.ObjectIDFromHex(userIDStr)
			if err != nil {
				return fmt.Errorf("invalid user id: %v", err)
			}
			allowedUsers = append(allowedUsers, objectUserID)
		}

		if len(allowedUsers) == 0 {
			return fmt.Errorf("at least one allowed user is required for private coupon")
		}
	default:
		return fmt.Errorf("invalid type: %s", req.Type)
	}

	if req.ExpiredAt == "" {
		return fmt.Errorf("expired at is required")
	}

	parseTime, err := time.Parse("2006-01-02T15:04:05-07:00", req.ExpiredAt)
	if err != nil {
		return fmt.Errorf("invalid expired at format: %w", err)
	}

	const maxAttempts = 5
	var codeCoupon string
	for i := 0; i < maxAttempts; i++ {

		codeCoupon = s.generateCodeCoupon(9)

		check, err := s.couponRepository.CheckCodeCoupon(c, codeCoupon)
		if err != nil {
			return fmt.Errorf("failed to check code: %w", err)
		}

		if !check {
			break
		}

		if i == maxAttempts-1 {
			return fmt.Errorf("could not generate unique coupon code after %d attempts", maxAttempts)
		}

	}

	coupon := &Coupon{
		ID:           primitive.NewObjectID(),
		CodeCoupon:   codeCoupon,
		Name:         req.Name,
		Discount:     req.Discount,
		MaximumUse:   &req.MaximumUse,
		UserIsUsed:   []primitive.ObjectID{},
		AllowedUsers: allowedUsers,
		Type:         req.Type,
		ExpiredAt:    parseTime,
	}

	return s.couponRepository.Create(c, coupon)

}

func (s *couponService) GetAllCoupons(c *gin.Context) ([]*CouponResponse, error) {

	coupons, err := s.couponRepository.FindAll(c)
	if err != nil {
		return nil, err
	}

	var results []*CouponResponse

	for _, coupon := range coupons {

		result, err := s.getUserInfor(c, coupon)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}


func (s *couponService) generateCodeCoupon(length int) string {

	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func (s *couponService) GetCouponByCode(c *gin.Context, code string) (*CouponResponse, error) {

	if code == "" {
		return nil, fmt.Errorf("code is required")
	}

	coupon, err := s.couponRepository.FindByCode(c, code)
	if err != nil {
		return nil, err
	}

	result, err := s.getUserInfor(c, coupon)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (s *couponService) GetCouponByUserID(c *gin.Context, userID string) ([]*CouponResponse, error) {

	if userID == "" {
		return nil, fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	var results []*CouponResponse

	coupon, err := s.couponRepository.FindAllCouponsByUserID(c, objectID)
	if err != nil {
		return nil, err
	}

	for _, coupon := range coupon {

		result, err := s.getUserInfor(c, coupon)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *couponService) CanUseCoupon(c *gin.Context, req *CanUseCouponRequest) (*Coupon, error) {

	if req.CouponCode == "" {
		return nil, fmt.Errorf("code coupon is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid code coupon: %v", err)
	}

	coupon, err := s.couponRepository.FindByCode(c, req.CouponCode)
	if err != nil {
		return nil, fmt.Errorf("invalid code coupon: %v", err)
	}

	check, err := s.canUseCoupon(objectID, *coupon)
	if err != nil {
		return nil, err
	}

	if !check {
		return nil, fmt.Errorf("invalid code coupon")
	}

	coupon, err = s.couponRepository.FindByCode(c, req.CouponCode)
	if err != nil {
		return nil, fmt.Errorf("invalid code coupon: %v", err)
	}

	return coupon, nil

}

func (s *couponService) canUseCoupon(userID primitive.ObjectID, coupon Coupon) (bool, error) {

	if len(coupon.AllowedUsers) > 0 {
		found := false
		for _, allowedUser := range coupon.AllowedUsers {
			if allowedUser == userID {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Errorf("you are not allowed to use this coupon")
		}
	}

	for _, userIsUsed := range coupon.UserIsUsed {
		if userIsUsed == userID {
			return false, fmt.Errorf("you have already used this coupon")
		}
	}

	if coupon.MaximumUse != nil {
		if len(coupon.UserIsUsed) >= *coupon.MaximumUse {
			return false, fmt.Errorf("this coupon has been used %d times", *coupon.MaximumUse)
		}
	}

	if time.Now().After(coupon.ExpiredAt) {
		return false, fmt.Errorf("this coupon has expired")
	}

	return true, nil

}

func (s *couponService) DeleteCoupon(c *gin.Context, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	return s.couponRepository.Delete(c, objectID)

}

func (s *couponService) getUserInfor(c *gin.Context, coupon *Coupon) (*CouponResponse, error) {

	var userUsed []*UserInfor
	for _, userID := range coupon.UserIsUsed {
		user, err := s.userService.GetUserByID(c, userID.Hex())
		if err != nil {
			return nil, fmt.Errorf("cannot get user from UserIsUsed: %w", err)
		}
		userUsed = append(userUsed, &UserInfor{
			ID:       userID,
			FullName: user.FirstName + user.LastName,
		})
	}

	var allowedUsers []*UserInfor
	for _, userID := range coupon.AllowedUsers {
		user, err := s.userService.GetUserByID(c, userID.Hex())
		if err != nil {
			return nil, fmt.Errorf("cannot get user from AllowedUsers: %w", err)
		}
		allowedUsers = append(allowedUsers, &UserInfor{
			ID:       userID,
			FullName: user.FirstName + user.LastName,
		})
	}

	return &CouponResponse{
		ID:           coupon.ID,
		Name:         coupon.Name,
		CodeCoupon:   coupon.CodeCoupon,
		Discount:     coupon.Discount,
		MaximumUse:   coupon.MaximumUse,
		UserIsUsed:   userUsed,
		AllowedUsers: allowedUsers,
		Type:         coupon.Type,
		ExpiredAt:    coupon.ExpiredAt,
		CreatedAt:    coupon.CreatedAt,
		UpdatedAt:    coupon.UpdatedAt,
	}, nil
}