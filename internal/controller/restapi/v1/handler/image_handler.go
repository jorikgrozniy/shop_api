package handler

import (
	"encoding/base64"
	"net/http"
	"shop_api/internal/controller/restapi/v1/dto"
	"shop_api/internal/controller/restapi/v1/mapper"
	"shop_api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageHandler struct {
	imageService *service.ImageService
}

func NewImageHandler(imageService *service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

// ChangeImage godoc
// @Summary Change image
// @Description Change product image data
// @Tags images
// @Accept json
// @Produce json
// @Param id path string true "Image ID"
// @Param request body dto.ImageRequest true "Image data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /images/{id} [put]
func (h *ImageHandler) ChangeImage(c *gin.Context) {
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
	if err := h.imageService.ChangeImage(id, image.Image); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "updated",
	})
}

// DeleteImage godoc
// @Summary Delete image
// @Description Delete product image by id
// @Tags images
// @Produce json
// @Param id path string true "Image ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /images/{id} [delete]
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.imageService.RemoveImage(id); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "deleted",
	})
}

// GetImage godoc
// @Summary Get image
// @Description Get product image by image id
// @Tags images
// @Produce application/octet-stream
// @Param id path string true "Image ID"
// @Success 200 {file} file
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /images/{id} [get]
func (h *ImageHandler) GetImage(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	image, err := h.imageService.GetImage(id)
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
