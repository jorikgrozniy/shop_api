package handler

import (
	"net/http"
	"shop_api/internal/service"
)

func getServiceErrorCode(err error) int {
	switch err {
	case service.ErrServerInternal:
		return http.StatusInternalServerError
	case service.ErrDependentEntity:
		return http.StatusBadRequest
	case service.ErrClientNotFound:
		return http.StatusNotFound
	case service.ErrSupplierNotFound:
		return http.StatusNotFound
	case service.ErrProductNotFound:
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
