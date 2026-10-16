package handler

import (
	"errors"
	"log"
	"net/http"

	"shop_api/internal/controller/restapi/v1/dto"
	"shop_api/internal/service"

	"github.com/gin-gonic/gin"
)

func getServiceErrorCode(err error) int {
	switch {
	case errors.Is(err, service.ErrServerInternal),
		errors.Is(err, service.ErrNilEntity):
		return http.StatusInternalServerError

	case errors.Is(err, service.ErrDependentEntity),
		errors.Is(err, service.ErrInvalidNameLength),
		errors.Is(err, service.ErrInvalidSurnameLength),
		errors.Is(err, service.ErrNotAdult),
		errors.Is(err, service.ErrInvalidGender),
		errors.Is(err, service.ErrInvalidCountryLength),
		errors.Is(err, service.ErrInvalidCityLength),
		errors.Is(err, service.ErrInvalidStreetLength),
		errors.Is(err, service.ErrInvalidPhoneLength),
		errors.Is(err, service.ErrInvalidCategoryLength),
		errors.Is(err, service.ErrInvalidPrice),
		errors.Is(err, service.ErrInvalidAvailableStock),
		errors.Is(err, service.ErrInvalidAmount),
		errors.Is(err, service.ErrInvalidSupplierID):
		return http.StatusBadRequest

	case errors.Is(err, service.ErrClientNotFound),
		errors.Is(err, service.ErrSupplierNotFound),
		errors.Is(err, service.ErrProductNotFound),
		errors.Is(err, service.ErrImageNotFound),
		errors.Is(err, service.ErrAddressNotFound):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}

func respondServiceError(c *gin.Context, err error) {
	code := getServiceErrorCode(err)

	if code >= http.StatusInternalServerError {
		log.Printf("handler error: %v", err)
		c.JSON(code, dto.ErrorResponse{
			Error: service.ErrServerInternal.Error(),
		})
		return
	}

	c.JSON(code, dto.ErrorResponse{
		Error: err.Error(),
	})
}
