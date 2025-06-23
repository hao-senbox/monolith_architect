package product

import (
	"fmt"
	"mime/multipart"
	"modular_monolith/helper"
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

	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	var req CreateProductRequest
	req.ProductName = c.PostForm("product_name")
	req.ProductDescription = c.PostForm("product_description")
	req.CategoryID = c.PostForm("category_id")
	req.Color = c.PostForm("color")
	if priceStr := c.PostForm("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			req.Price = price
		}
	}
	if discountStr := c.PostForm("discount"); discountStr != "" {
		if discount, err := strconv.ParseFloat(discountStr, 64); err == nil {
			req.Discount = discount
		}
	}
	req.Currency = c.PostForm("currency")

	if req.ProductName == "" || req.ProductDescription == "" || req.CategoryID == "" || req.Color == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("invalid request: product_name, product_description, category_id or color is missing"), helper.ErrInvalidRequest)
		return
	}

	sizeCountStr := c.PostForm("size_count")
	fmt.Printf("size_count: %s\n", sizeCountStr)
	sizeCount, err := strconv.Atoi(sizeCountStr)
	if err != nil || sizeCount == 0 {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("invalid size_count"), helper.ErrInvalidRequest)
		return
	}

	var sizes []CreateSizeOptionsRequest
	for i := 0; i < sizeCount; i++ {

		var size CreateSizeOptionsRequest

		size.Size = c.PostForm(fmt.Sprintf("sizes[%d][size]", i))

		if stockStr := c.PostForm(fmt.Sprintf("sizes[%d][stock]", i)); stockStr != "" {
			if stock, err := strconv.Atoi(stockStr); err == nil {
				size.Stock = stock
			}
		}

		sizes = append(sizes, size)
	}

	req.Sizes = sizes

	var productFiles ProductFiles

	if mainImage, err := c.FormFile("main_image"); err == nil {
		productFiles.MainImage = mainImage
	}

	var subImages []*multipart.FileHeader
	if subImagesArray := c.Request.MultipartForm.File["sub_images"]; len(subImagesArray) > 0 {
		subImages = subImagesArray
	}
	productFiles.SubImages = subImages
    
	err = h.ProductService.CreateProduct(c.Request.Context(), &req, productFiles)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", nil)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {

	products, err := h.ProductService.GetAllProducts(c)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", products)

}

func (h *ProductHandler) GetProductByID(c *gin.Context) {

	id := c.Param("id")

	product, err := h.ProductService.GetProductByID(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", product)

}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {

	id := c.Param("id")

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	var req UpdateProductRequest
	req.ProductName = c.PostForm("product_name")
	req.ProductDescription = c.PostForm("product_description")
	req.CategoryID = c.PostForm("category_id")
	req.Color = c.PostForm("color")
	if priceStr := c.PostForm("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			req.Price = price
		}
	}
	if discountStr := c.PostForm("discount"); discountStr != "" {
		if discount, err := strconv.ParseFloat(discountStr, 64); err == nil {
			req.Discount = discount
		}
	}
	req.Currency = c.PostForm("currency")

	if req.ProductName == "" || req.ProductDescription == "" || req.CategoryID == "" || req.Color == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("invalid request: product_name, product_description, category_id or color is missing"), helper.ErrInvalidRequest)
		return
	}

	// Get existing product to know current size count
	product, err := h.ProductService.GetProductByID(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	// Parse size options
	var sizes []CreateSizeOptionsRequest
	for i := 0; i < len(product.Sizes); i++ {
		var size CreateSizeOptionsRequest
		size.Size = c.PostForm(fmt.Sprintf("sizes[%d][size]", i))
		if stockStr := c.PostForm(fmt.Sprintf("sizes[%d][stock]", i)); stockStr != "" {
			if stock, err := strconv.Atoi(stockStr); err == nil {
				size.Stock = stock
			}
		}

		sizes = append(sizes, size)
	}

	req.Sizes = sizes

	// Handle file uploads
	var productFiles ProductFiles

	// Main image
	if mainImage, err := c.FormFile("main_image"); err == nil {
		productFiles.MainImage = mainImage
	}

	// Sub images
	var subImages []*multipart.FileHeader
	if subImagesArray := c.Request.MultipartForm.File["sub_images"]; len(subImagesArray) > 0 {
		subImages = subImagesArray
	}
	productFiles.SubImages = subImages

	err = h.ProductService.UpdateProduct(c.Request.Context(), id, &req, productFiles)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {

	id := c.Param("id")

	err := h.ProductService.DeleteProduct(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
}
