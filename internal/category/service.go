package category

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) error
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, categoryID string) (*Category, error)
	UpdateCategory(ctx context.Context, req *UpdateCategoryRequest, categoryID string) error
	DeleteCategory(ctx context.Context, categoryID string) error
}

type categoryService struct {
	categoryRepository CategoryRepository
}

func NewCategoryService(categoryRepository CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}

func (s *categoryService) GetCategory(ctx context.Context, categoryID string) (*Category, error) {

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, err
	}
	return s.categoryRepository.FindByID(ctx, objectID)

}

func (s *categoryService) GetCategories(ctx context.Context) ([]*Category, error) {
	return s.categoryRepository.FindAll(ctx)
}

func (s *categoryService) CreateCategory(ctx context.Context, req *CreateCategoryRequest) error {

	var parentID *primitive.ObjectID

	if req.CategoryName == "" {
		return fmt.Errorf("category name is required")
	}

	if req.ParentID != nil {
		objectID, err := primitive.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return fmt.Errorf("invalid parent id: %v", err)
		}
		parentID = &objectID
	} else {
		parentID = nil
	} 

	category := Category{
		CategoryName: req.CategoryName,
		ParentID:     parentID,
	}

	return s.categoryRepository.Create(ctx, &category)

}

func (s *categoryService) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest, categoryID string) error {
	
	var parentID *primitive.ObjectID

	if categoryID == "" {
		return fmt.Errorf("category id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return fmt.Errorf("invalid category id: %v", err)
	}

	if req.CategoryName == "" {
		return fmt.Errorf("category name is required")
	}

	if req.ParentID != nil {
		objectID, err := primitive.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return fmt.Errorf("invalid parent id: %v", err)
		}
		parentID = &objectID
	} else {
		parentID = nil
	}
	
	category := Category{
		ID:           objectID,
		CategoryName: req.CategoryName,
		ParentID:     parentID,
	}

	return s.categoryRepository.UpdateByID(ctx, &category, objectID)

}

func (s *categoryService) DeleteCategory(ctx context.Context, categoryID string) error {
	
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return fmt.Errorf("invalid category id: %v", err)
	}

	return s.categoryRepository.DeleteByID(ctx, objectID)
}