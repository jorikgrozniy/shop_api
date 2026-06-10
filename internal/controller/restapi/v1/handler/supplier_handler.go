package handler

import (
	"net/http"
	"shop_api/internal/controller/restapi/v1/dto"
	"shop_api/internal/controller/restapi/v1/mapper"
	"shop_api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SupplierHandler struct {
	supplierService *service.SupplierService
	addressService  *service.AddressService
}

func NewSupplierHandler(supplierService *service.SupplierService,
	addressService *service.AddressService) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
		addressService:  addressService,
	}
}

// AddSupplier godoc
// @Summary Create supplier
// @Description Create new supplier
// @Tags suppliers
// @Accept json
// @Produce json
// @Param request body dto.SupplierRequest true "Supplier data"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /suppliers [post]
func (h *SupplierHandler) AddSupplier(c *gin.Context) {
	var req dto.SupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid JSON format",
		})
		return
	}

	supplier := mapper.SupplierDTOtoDAO(req)
	address := mapper.AddressDTOtoDAO(req.Address)

	if err := h.supplierService.AddSupplier(supplier, address); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Status: "created",
	})
}

// DeleteSupplier godoc
// @Summary Delete supplier
// @Description Delete supplier by id
// @Tags suppliers
// @Produce json
// @Param id path string true "Supplier ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.supplierService.RemoveSupplier(id); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "deleted",
	})
}

// GetSupplier godoc
// @Summary Get supplier
// @Description Get supplier by id
// @Tags suppliers
// @Produce json
// @Param id path string true "Supplier ID"
// @Success 200 {object} dto.SupplierResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /suppliers/{id} [get]
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	supplier, err := h.supplierService.GetSupplier(id)
	if err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	address, err := h.addressService.GetAddress(supplier.AddressID)
	if err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	supplierResponse := mapper.SupplierDAOtoDTO(supplier)
	addressResponse := mapper.AddressDAOtoDTO(address)
	supplierResponse.Address = addressResponse

	c.JSON(http.StatusOK, supplierResponse)
}

// GetSuppliersWithParams godoc
// @Summary Get suppliers with params
// @Description Get suppliers list with optional params
// @Tags suppliers
// @Produce json
// @Param limit query int false "Limit" minimum(1) maximum(100) default(100)
// @Param offset query int false "Offset" minimum(0) default(0)
// @Success 200 {object} dto.SupplierListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /suppliers [get]
func (h *SupplierHandler) GetSuppliersWithParams(c *gin.Context) {
	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid query format",
		})
		return
	}

	suppliers, limit, offset, err := h.supplierService.GetSuppliersWithParams(req.Limit, req.Offset)
	if err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	suppliersResponse := mapper.SupplierListDAOtoDTO(suppliers)
	suppliersResponse.Limit = limit
	suppliersResponse.Offset = offset

	for i, supplier := range suppliers {
		address, err := h.addressService.GetAddress(supplier.AddressID)
		if err != nil {
			c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		addressResponse := mapper.AddressDAOtoDTO(address)
		suppliersResponse.Data[i].Address = addressResponse
	}

	c.JSON(http.StatusOK, suppliersResponse)
}

// ChangeSupplierAddress godoc
// @Summary Change supplier address
// @Description Change address of supplier
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Param request body dto.AddressRequest true "Address data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /suppliers/{id}/address [patch]
func (h *SupplierHandler) ChangeSupplierAddress(c *gin.Context) {
	supplierIDStr := c.Param("id")

	supplierID, err := uuid.Parse(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req dto.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid JSON format",
		})
		return
	}

	address := mapper.AddressDTOtoDAO(req)

	if err := h.supplierService.ChangeSupplierAddress(supplierID, address); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "updated",
	})
}
