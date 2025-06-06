package product

import (
	"context"
	"fmt"
	"modular_monolith/internal/cloudinaryutil"
	fileutil "modular_monolith/internal/common"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest, productFiles ProductFiles) error
	GetAllProducts(ctx context.Context) ([]*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	UpdateProduct(ctx context.Context, id string, req *UpdateProductRequest, productFiles ProductFiles) error
	DeleteProduct(ctx context.Context, id string) error
}

type productService struct {
	repository    ProductRepository
	cloudUploader *cloudinaryutil.CloudinaryUploader
}

func NewProductService(repository ProductRepository, uploader *cloudinaryutil.CloudinaryUploader) ProductService {
	return &productService{
		repository:    repository,
		cloudUploader: uploader,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *CreateProductRequest, productFiles ProductFiles) error {

	if req.ProductName == "" {
		return fmt.Errorf("product name is required")
	}

	if req.ProductDescription == "" {
		return fmt.Errorf("product description is required")
	}

	if req.Color == "" {
		return fmt.Errorf("color is required")
	}

	if req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if req.Discount < 0 {
		return fmt.Errorf("discount must be greater than or equal to 0")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category id: %v", err)
	}

	if len(req.Sizes) == 0 {
		return fmt.Errorf("sizes are required")
	}

	var sizes []SizeOptions

	for i, s := range req.Sizes {

		if s.Size == "" {
			return fmt.Errorf("invalid size option at index %d", i)
		}

		if s.Stock < 0 {
			return fmt.Errorf("invalid stock for size option at index %d", i)
		}


		sizes = append(sizes, SizeOptions{
			Size:     s.Size,
			Stock:    s.Stock,
		})
	}

	// Create product object
	product := &Product{
		ProductName:        req.ProductName,
		ProductDescription: req.ProductDescription,
		CategoryID:         categoryID,
		Color:              req.Color,
		Price:              req.Price,
		Discount:           req.Discount,
		Currency:           req.Currency,
		Sizes:              sizes,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Handle main image upload
	if productFiles.MainImage != nil {
		tempPath := "/tmp/" + productFiles.MainImage.Filename
		if err := fileutil.SaveUploadedFile(productFiles.MainImage, tempPath); err != nil {
			return fmt.Errorf("failed to save main image: %w", err)
		}
		defer os.Remove(tempPath)

		mainImageURL, mainImagePublicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "products")
		if err != nil {
			return fmt.Errorf("failed to upload main image: %w", err)
		}

		product.MainImage = mainImageURL
		product.MainImagePublicID = mainImagePublicID
	}

	// Handle sub images upload
	var subImages []SubImage
	for i, subImage := range productFiles.SubImages {
		tempPath := "/tmp/" + subImage.Filename
		if err := fileutil.SaveUploadedFile(subImage, tempPath); err != nil {
			return fmt.Errorf("failed to save sub image %d: %w", i, err)
		}
		defer os.Remove(tempPath)

		subImageURL, subImagePublicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "products")
		if err != nil {
			return fmt.Errorf("failed to upload sub image %d: %w", i, err)
		}

		subImages = append(subImages, SubImage{
			SubImagePublicID: subImagePublicID,
			Url:              subImageURL,
		})
	}
	product.SubImages = subImages

	fmt.Printf("Product: %+v\n", product)

	return s.repository.Create(ctx, product)
}

func (s *productService) GetAllProducts(ctx context.Context) ([]*Product, error) {
	return s.repository.FindAll(ctx)
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*Product, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product id: %v", err)
	}

	return s.repository.FindByID(ctx, objectID)

}

func (s *productService) UpdateProduct(ctx context.Context, id string, req *UpdateProductRequest, productFiles ProductFiles) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	if req.ProductName == "" {
		return fmt.Errorf("product name is required")
	}

	if req.ProductDescription == "" {
		return fmt.Errorf("product description is required")
	}

	if req.Color == "" {
		return fmt.Errorf("color is required")
	}

	if req.Price < 0 {
		return fmt.Errorf("price must be greater than or equal to 0")
	}

	if req.Discount < 0 {
		return fmt.Errorf("discount must be greater than or equal to 0")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category id: %v", err)
	}

	if len(req.Sizes) == 0 {
		return fmt.Errorf("sizes are required")
	}

	// Get existing product
	existingProduct, err := s.repository.FindByID(ctx, objectID)
	if err != nil {
		return fmt.Errorf("failed to get existing product: %w", err)
	}

	// Validate and convert sizes
	var sizes []SizeOptions
	for i, s := range req.Sizes {

		if s.Size == "" {
			return fmt.Errorf("invalid size option at index %d", i)
		}

		if s.Stock < 0 {
			return fmt.Errorf("invalid stock for size option at index %d", i)
		}

		sizes = append(sizes, SizeOptions{
			Size:     s.Size,
			Stock:    s.Stock,
		})
	}

	// Update product object
	existingProduct.ProductName = req.ProductName
	existingProduct.ProductDescription = req.ProductDescription
	existingProduct.CategoryID = categoryID
	existingProduct.Color = req.Color
	existingProduct.Price = req.Price
	existingProduct.Discount = req.Discount
	existingProduct.Currency = req.Currency
	existingProduct.Sizes = sizes
	existingProduct.UpdatedAt = time.Now()

	// Handle main image upload (if new image provided)
	if productFiles.MainImage != nil {
		// Delete old main image if exists
		if existingProduct.MainImagePublicID != "" {
			if err := s.cloudUploader.DeleteImage(ctx, existingProduct.MainImagePublicID); err != nil {
				// Log error but don't fail the update
				fmt.Printf("Warning: failed to delete old main image: %v\n", err)
			}
		}

		tempPath := "/tmp/" + productFiles.MainImage.Filename
		if err := fileutil.SaveUploadedFile(productFiles.MainImage, tempPath); err != nil {
			return fmt.Errorf("failed to save main image: %w", err)
		}
		defer os.Remove(tempPath)

		mainImageURL, mainImagePublicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "products")
		if err != nil {
			return fmt.Errorf("failed to upload main image: %w", err)
		}

		existingProduct.MainImage = mainImageURL
		existingProduct.MainImagePublicID = mainImagePublicID
	}

	// Handle sub images upload (if new images provided)
	if len(productFiles.SubImages) > 0 {
		// Delete old sub images
		for _, subImg := range existingProduct.SubImages {
			if err := s.cloudUploader.DeleteImage(ctx, subImg.SubImagePublicID); err != nil {
				// Log error but don't fail the update
				fmt.Printf("Warning: failed to delete old sub image: %v\n", err)
			}
		}

		var subImages []SubImage
		for i, subImage := range productFiles.SubImages {
			tempPath := "/tmp/" + subImage.Filename
			if err := fileutil.SaveUploadedFile(subImage, tempPath); err != nil {
				return fmt.Errorf("failed to save sub image %d: %w", i, err)
			}
			defer os.Remove(tempPath)

			subImageURL, subImagePublicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "products")
			if err != nil {
				return fmt.Errorf("failed to upload sub image %d: %w", i, err)
			}

			subImages = append(subImages, SubImage{
				SubImagePublicID: subImagePublicID,
				Url:              subImageURL,
			})
		}
		existingProduct.SubImages = subImages
	}

	return s.repository.UpdateByID(ctx, objectID, existingProduct)
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	product, err := s.repository.FindByID(ctx, objectID)
	if err != nil || product == nil {
		return fmt.Errorf("product not found")
	}

	err = s.cloudUploader.DeleteImage(ctx, product.MainImagePublicID)
	if err != nil {
		return fmt.Errorf("failed to delete main image: %w", err)
	}

	for _, sub := range product.SubImages {
		err := s.cloudUploader.DeleteImage(ctx, sub.SubImagePublicID)
		if err != nil {
			return fmt.Errorf("failed to delete sub image: %w", err)
		}
	}

	return s.repository.DeleteByID(ctx, objectID)
}
