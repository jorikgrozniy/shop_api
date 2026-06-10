package handler

import (
	"encoding/base64"
	"net/http"
	"shop_api/internal/controller/restapi/v1/dto"
	"shop_api/internal/controller/restapi/v1/mapper"
	"shop_api/internal/entity"
	"shop_api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService *service.ProductService
	imageService   *service.ImageService
}

func NewProductHandler(productService *service.ProductService,
	imageService *service.ImageService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		imageService:   imageService,
	}
}

// AddProduct godoc
// @Summary Create product
// @Description Create new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.ProductRequest true "Product data"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products [post]
func (h *ProductHandler) AddProduct(c *gin.Context) {
	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid JSON format",
		})
		return
	}

	product, err := mapper.ProductDTOtoDAO(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var image *entity.Image = nil
	if req.Image.JPGbase64 != "" {
		if _, err := base64.StdEncoding.DecodeString(req.Image.JPGbase64); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "invalid image format: decode jpg to base64",
			})
			return
		}

		image = mapper.ImageDTOtoDAO(req.Image)
	}

	if err := h.productService.AddProduct(product, image); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Status: "created",
	})
}

// DecreaseProductStock godoc
// @Summary Decrease product stock
// @Description Decrease product stock by N
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.ProductStockRequest true "Product stock data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id}/decrease-stock [post]
func (h *ProductHandler) DecreaseProductStock(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req dto.ProductStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid JSON format",
		})
		return
	}

	if err := h.productService.DecreaseAvailableStock(id, req.Amount); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "updated",
	})
}

// GetProduct godoc
// @Summary Get product
// @Description Get product by id
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	product, err := h.productService.GetProduct(id)
	if err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	productResponse := mapper.ProductDAOtoDTO(product)

	c.JSON(http.StatusOK, productResponse)
}

// GetProductsWithParams godoc
// @Summary Get products with params
// @Description Get products list with optional params
// @Tags products
// @Produce json
// @Param limit query int false "Limit" minimum(1) maximum(100) default(100)
// @Param offset query int false "Offset" minimum(0) default(0)
// @Success 200 {object} dto.ProductListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetProductsWithParams(c *gin.Context) {
	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid query format",
		})
		return
	}

	products, limit, offset, err := h.productService.GetProductsWithParams(req.Limit, req.Offset)
	if err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	productsResponse := mapper.ProductListDAOtoDTO(products)
	productsResponse.Limit = limit
	productsResponse.Offset = offset

	c.JSON(http.StatusOK, productsResponse)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by id
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.productService.RemoveProduct(id); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "deleted",
	})
}

// AddProductImage godoc
// @Summary Add product image
// @Description Add image for product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.ImageRequest true "Image data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id}/image [post]
func (h *ProductHandler) AddProductImage(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req dto.ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid JSON format",
		})
		return
	}

	if _, err := base64.StdEncoding.DecodeString(req.JPGbase64); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid image format: decode jpg to base64",
		})
		return
	}

	image := mapper.ImageDTOtoDAO(req)
	if err := h.productService.AddProductImage(id, image); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "updated",
	})
}

// GetProductImage godoc
// @Summary Get product image
// @Description Get product image by product id
// @Tags products
// @Produce application/octet-stream
// @Param id path string true "Product ID"
// @Success 200 {file} file
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id}/image [get]
func (h *ProductHandler) GetProductImage(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	image, err := h.productService.GetProductImage(id)
	if err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	imageDTO := mapper.ImageDAOtoDTO(image)
	imageBytes, err := base64.StdEncoding.DecodeString(imageDTO.JPGbase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", `attachment; filename="image.jpg"`)
	c.Header("Content-Length", strconv.Itoa(len(imageBytes)))

	c.Data(http.StatusOK, "application/octet-stream", imageBytes)
}
