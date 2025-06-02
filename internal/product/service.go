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
	CreateProduct(ctx context.Context, req *CreateProductRequest, variantFiles []VariantFiles) error
	GetAllProducts(ctx context.Context) ([]*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	UpdateProduct(ctx context.Context, id string, req *UpdateProductRequest, variantFiles []VariantFiles) error
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

func (s *productService) CreateProduct(ctx context.Context, req *CreateProductRequest, variantFiles []VariantFiles) error {

	if req.ProductName == "" {
		return fmt.Errorf("product name is required")
	}

	if req.ProductDescription == "" {
		return fmt.Errorf("product description is required")
	}

	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category id: %v", err)
	}

	if len(req.Variants) == 0 {
		return fmt.Errorf("variants are required")
	}

	if len(req.Variants) != len(variantFiles) {
		return fmt.Errorf("variant count mismatch with files")
	}

	var variants []ProductVariant


	for i, v := range req.Variants {

		if v.Size == "" {
			return fmt.Errorf("size is required for variant %d", i)
		}

		if v.Color == "" {
			return fmt.Errorf("color is required for variant %d", i)
		}

		if v.Price <= 0 {
			return fmt.Errorf("price is required for variant %d", i)
		}

		if v.Stock < 0 {
			return fmt.Errorf("stock cannot be negative for variant %d", i)
		}

		if v.Currency == "" {
			return fmt.Errorf("currency is required for variant %d", i)
		}

		variant := ProductVariant{
			SKU:      v.SKU,
			Color:    v.Color,
			Size:     v.Size,
			Stock:    v.Stock,
			Price:    v.Price,
			Discount: v.Discount,
			Currency: v.Currency,
		}

		if variantFiles[i].MainImage != nil {

			tempPath := "/tmp/" + variantFiles[i].MainImage.Filename
			if err := fileutil.SaveUploadedFile(variantFiles[i].MainImage, tempPath); err != nil {
				return fmt.Errorf("failed to save avatar: %w", err)
			}

			defer os.Remove(tempPath)

			mainImageURL, mainImagePublicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "products")
			if err != nil {
				return fmt.Errorf("failed to upload main image for variant %d: %w", i, err)
			}

			variant.MainImage = mainImageURL
			variant.MainImagePublicID = mainImagePublicID
		}

		var subImages []SubImage
		for j, subImage := range variantFiles[i].SubImages {

			tempPath := "/tmp/" + subImage.Filename
			if err := fileutil.SaveUploadedFile(subImage, tempPath); err != nil {
				return fmt.Errorf("failed to save avatar: %w", err)
			}
			defer os.Remove(tempPath)
			
			subImageURL, subImagePublicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "products")
			if err != nil {
				return fmt.Errorf("failed to upload sub image %d for variant %d: %w", j, i, err)
			}

			subImages = append(subImages, SubImage{
				SubImagePublicID: subImagePublicID,
				Url:              subImageURL,
			})
		}
		variant.SubImages = subImages

		variants = append(variants, variant)
	}

	product := &Product{
		ProductName:        req.ProductName,
		ProductDescription: req.ProductDescription,
		CategoryID:         categoryID,
		Variants:           variants,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

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

func (s *productService) UpdateProduct(ctx context.Context, id string, req *UpdateProductRequest, variantFiles []VariantFiles) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	product, err := s.repository.FindByID(ctx, objectID)
	if err != nil || product == nil {
		return fmt.Errorf("product not found")
	}

	if req.ProductName == "" || req.ProductDescription == "" || req.CategoryID == "" {
		return fmt.Errorf("product name, description and category are required")
	}

	objectCategoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category id: %v", err)
	}

	var variants []ProductVariant

	for i, v := range req.Variants {
		if v.Size == "" || v.Color == "" || v.Currency == "" || v.Price <= 0 || v.Stock < 0 {
			return fmt.Errorf("invalid data for variant %d", i)
		}

		variant := ProductVariant{
			SKU:      v.SKU,
			Color:    v.Color,
			Size:     v.Size,
			Stock:    v.Stock,
			Price:    v.Price,
			Discount: v.Discount,
			Currency: v.Currency,
		}

		if variantFiles[i].MainImage != nil {
			tempPath := "/tmp/" + variantFiles[i].MainImage.Filename
			if err := fileutil.SaveUploadedFile(variantFiles[i].MainImage, tempPath); err != nil {
				return fmt.Errorf("failed to save main image: %w", err)
			}
			defer os.Remove(tempPath)

			_ = s.cloudUploader.DeleteImage(ctx, product.Variants[i].MainImagePublicID)

			mainImageURL, mainImagePublicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "products")
			if err != nil {
				return fmt.Errorf("failed to upload main image: %w", err)
			}

			variant.MainImage = mainImageURL
			variant.MainImagePublicID = mainImagePublicID
		} else {
			variant.MainImage = product.Variants[i].MainImage
			variant.MainImagePublicID = product.Variants[i].MainImagePublicID
		}

		if len(variantFiles[i].SubImages) > 0 {

			for _, sub := range product.Variants[i].SubImages {
				_ = s.cloudUploader.DeleteImage(ctx, sub.SubImagePublicID)
			}

			for _, subImage := range variantFiles[i].SubImages {
				tempPath := "/tmp/" + subImage.Filename
				if err := fileutil.SaveUploadedFile(subImage, tempPath); err != nil {
					return fmt.Errorf("failed to save sub image: %w", err)
				}
				defer os.Remove(tempPath)

				url, publicID, err := s.cloudUploader.UploadImage(ctx, tempPath, "products")
				if err != nil {
					return fmt.Errorf("failed to upload sub image: %w", err)
				}

				variant.SubImages = append(variant.SubImages, SubImage{
					Url:              url,
					SubImagePublicID: publicID,
				})
			}
		} else {
			variant.SubImages = product.Variants[i].SubImages
		}

		variants = append(variants, variant)
	}

	productUpdate := &Product{
		ProductName:        req.ProductName,
		ProductDescription: req.ProductDescription,
		CategoryID:         objectCategoryID,
		Variants:           variants,
		UpdatedAt:          time.Now(),
	}

	return s.repository.UpdateByID(ctx, objectID, productUpdate)

}
