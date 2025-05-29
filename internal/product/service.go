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

	// if len(req.Variants) == 0 {
	// 	return fmt.Errorf("variants are required")
	// }

	// if len(req.Variants) != len(variantFiles) {
	// 	return fmt.Errorf("variant count mismatch with files")
	// }

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

	return s.repository.Create(ctx, product)
}
