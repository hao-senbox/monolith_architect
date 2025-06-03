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

	err := c.Request.ParseMultipartForm(32 << 20); 
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	var req CreateProductRequest
	req.ProductName = c.PostForm("product_name")
	req.ProductDescription = c.PostForm("product_description")
	req.CategoryID = c.PostForm("category_id")

	if req.ProductName == "" || req.ProductDescription == "" || req.CategoryID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("invalid request product_name, product_description or category_id"), helper.ErrInvalidRequest)
		return
	}

	variantCountStr := c.PostForm("variant_count")
	variantCount, err := strconv.Atoi(variantCountStr)
	if err != nil || variantCount == 0 {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	sizeOptionsStr := c.PostForm("size_options")
	sizeCount, err := strconv.Atoi(sizeOptionsStr)
	if err != nil || sizeCount == 0 {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	var variants []CreateProductVariantRequest
	var variantFiles []VariantFiles

	for i := 0; i < variantCount; i++ {

		var sizes []CreateSizeOptionsRequest

		for j := 0; j < sizeCount; j++ {

			var size CreateSizeOptionsRequest

			size.SKU = c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][sku]", i, j))
			size.Size = c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][size]", i, j))
			size.Currency = c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][currency]", i, j))

			if priceStr := c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][price]", i, j)); priceStr != "" {
				if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
					size.Price = price
				}
			}

			if discountStr := c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][discount]", i, j)); discountStr != "" {
				if discount, err := strconv.ParseFloat(discountStr, 64); err == nil {
					size.Discount = discount
				}
			}

			if stockStr := c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][stock]", i, j)); stockStr != "" {
				if stock, err := strconv.Atoi(stockStr); err == nil {
					size.Stock = stock
				}
			}

			sizes = append(sizes, size)

		}

		variant := CreateProductVariantRequest{
			Color:    c.PostForm(fmt.Sprintf("variants[%d][color]", i)),
			Sizes:    sizes,
		}

		variants = append(variants, variant)

		var files VariantFiles

		if mainImage, err := c.FormFile(fmt.Sprintf("variants[%d][main_image]", i)); err == nil {
			files.MainImage = mainImage
		}

		var subImages []*multipart.FileHeader

		if subImagesArray  := c.Request.MultipartForm.File[fmt.Sprintf("variants[%d][sub_image]", i)]; len(subImagesArray) > 0 {
			subImages = subImagesArray
		}

		files.SubImages = subImages
		variantFiles = append(variantFiles, files)
	}

	req.Variants = variants
	
	err = h.ProductService.CreateProduct(c.Request.Context(), &req, variantFiles) 
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

	if req.ProductName == "" || req.ProductDescription == "" || req.CategoryID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("invalid request product_name, product_description or category_id"), helper.ErrInvalidRequest)
		return
	}

	product, err := h.ProductService.GetProductByID(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	var variants []CreateProductVariantRequest
	var variantFiles []VariantFiles

	for i := 0; i < len(product.Variants); i++ {

		var sizes []CreateSizeOptionsRequest

		for j := 0; j < len(product.Variants[i].Sizes); j++ {

			var size CreateSizeOptionsRequest
			size.Size = c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][size]", i, j))
			size.SKU = c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][sku]", i, j))
			size.Currency = c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][currency]", i, j))

			if priceStr := c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][price]", i, j)); priceStr != "" {
				if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
					size.Price = price
				}
			}

			if discountStr := c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][discount]", i, j)); discountStr != "" {
				if discount, err := strconv.ParseFloat(discountStr, 64); err == nil {
					size.Discount = discount
				}
			}

			if stockStr := c.PostForm(fmt.Sprintf("variants[%d][sizes][%d][stock]", i, j)); stockStr != "" {
				if stock, err := strconv.Atoi(stockStr); err == nil {
					size.Stock = stock
				}
			}

			sizes = append(sizes, size)
		}

		variant := CreateProductVariantRequest{
			Color:    c.PostForm(fmt.Sprintf("variants[%d][color]", i)),
			Sizes:    sizes,
		}


		variants = append(variants, variant)

		var files VariantFiles

		if mainImage, err := c.FormFile((fmt.Sprintf("variants[%d][main_image]", i))); err == nil {
			files.MainImage = mainImage
		}

		var subImages []*multipart.FileHeader

		if subImagesArray  := c.Request.MultipartForm.File[fmt.Sprintf("variants[%d][sub_image]", i)]; len(subImagesArray) > 0 {
			subImages = subImagesArray
		}

		files.SubImages = subImages
		variantFiles = append(variantFiles, files)
	}

	req.Variants = variants
	
	err = h.ProductService.UpdateProduct(c.Request.Context(), id, &req, variantFiles) 
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
