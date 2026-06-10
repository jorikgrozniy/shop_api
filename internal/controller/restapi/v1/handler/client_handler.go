package handler

import (
	"net/http"
	"shop_api/internal/controller/restapi/v1/dto"
	"shop_api/internal/controller/restapi/v1/mapper"
	"shop_api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientHandler struct {
	clientService  *service.ClientService
	addressService *service.AddressService
}

func NewClientHandler(clientService *service.ClientService,
	addressService *service.AddressService) *ClientHandler {
	return &ClientHandler{
		clientService:  clientService,
		addressService: addressService,
	}
}

// AddClient godoc
// @Summary Create client
// @Description Create new client
// @Tags clients
// @Accept json
// @Produce json
// @Param request body dto.ClientRequest true "Client data"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clients [post]
func (h *ClientHandler) AddClient(c *gin.Context) {
	var req dto.ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid JSON format",
		})
		return
	}

	client, err := mapper.ClientDTOtoDAO(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	address := mapper.AddressDTOtoDAO(req.Address)

	if err := h.clientService.AddClient(client, address); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Status: "created",
	})
}

// DeleteClient godoc
// @Summary Delete client
// @Description Delete client by id
// @Tags clients
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clients/{id} [delete]
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.clientService.RemoveClient(id); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "deleted",
	})
}

// GetClientsWithParams godoc
// @Summary Get clients with params
// @Description Get clients list with optional params
// @Tags clients
// @Produce json
// @Param name query string false "Client name"
// @Param surname query string false "Client surname"
// @Param limit query int false "Limit" minimum(1) maximum(100) default(100)
// @Param offset query int false "Offset" minimum(0) default(0)
// @Success 200 {object} dto.ClientListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clients [get]
func (h *ClientHandler) GetClientsWithParams(c *gin.Context) {
	var req dto.ClientListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid JSON format",
		})
		return
	}

	clients, limit, offset, err := h.clientService.GetClientsWithParams(&req.Name, &req.Surname, req.Limit, req.Offset)
	if err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	clientsResponse := mapper.ClientListDAOtoDTO(clients)
	clientsResponse.Limit = limit
	clientsResponse.Offset = offset

	for i, client := range clients {
		address, err := h.addressService.GetAddress(client.AddressID)
		if err != nil {
			c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		addressResponse := mapper.AddressDAOtoDTO(address)
		clientsResponse.Data[i].Address = addressResponse
	}

	c.JSON(http.StatusOK, clientsResponse)
}

// ChangeClientAddress godoc
// @Summary Change client address
// @Description Change address of client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param request body dto.AddressRequest true "Address data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clients/{id}/address [patch]
func (h *ClientHandler) ChangeClientAddress(c *gin.Context) {
	clientIDStr := c.Param("id")

	clientID, err := uuid.Parse(clientIDStr)
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

	if err := h.clientService.ChangeClientAddress(clientID, address); err != nil {
		c.JSON(getServiceErrorCode(err), dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Status: "updated",
	})
}
