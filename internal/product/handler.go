package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductService ProductService
}

func NewProductHandler(productService ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	var req CreateProductRequest	
	req.ProductName = c.PostForm("product_name")
	req.ProductDescription = c.PostForm("product_description")
	req.CategoryID = c.PostForm("category_id")

	if req.ProductName == "" || req.ProductDescription == "" || req.CategoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required product information"})
		return
	}

	variantCountStr := c.PostForm("variant_count")
	variantCount, err := strconv.Atoi(variantCountStr)
	if err != nil || variantCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant count"})
		return
	}

	var variants []CreateProductVariant
	var variantFiles []VariantFiles

	for i := 0; i < variantCount; i++ {

		variant := CreateProductVariant{
			SKU:      c.PostForm(fmt.Sprintf("variants[%d][sku]", i)),
			Size:     c.PostForm(fmt.Sprintf("variants[%d][size]", i)),
			Color:    c.PostForm(fmt.Sprintf("variants[%d][color]", i)),
			Currency: c.PostForm(fmt.Sprintf("variants[%d][currency]", i)),
		}


		if priceStr := c.PostForm(fmt.Sprintf("variants[%d][price]", i)); priceStr != "" {
			if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
				variant.Price = price
			}
		}

		if discountStr := c.PostForm(fmt.Sprintf("variants[%d][discount]", i)); discountStr != "" {
			if discount, err := strconv.ParseFloat(discountStr, 64); err == nil {
				variant.Discount = discount
			}
		}

		if stockStr := c.PostForm(fmt.Sprintf("variants[%d][stock]", i)); stockStr != "" {
			if stock, err := strconv.Atoi(stockStr); err == nil {
				variant.Stock = stock
			}
		}

		variants = append(variants, variant)

		var files VariantFiles
		
		if mainImage, err := c.FormFile(fmt.Sprintf("variants[%d][main_image]", i)); err == nil {
			files.MainImage = mainImage
		}

		if subImages := c.Request.MultipartForm.File[fmt.Sprintf("variants[%d][sub_images]", i)]; len(subImages) > 0 {
			files.SubImages = subImages
		}

		variantFiles = append(variantFiles, files)
	}

	if err := h.ProductService.CreateProduct(c.Request.Context(), &req, variantFiles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}	