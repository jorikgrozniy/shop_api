package v1

import (
	"shop_api/internal/controller/restapi/v1/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	clientHandler   *handler.ClientHandler
	imageHandler    *handler.ImageHandler
	productHandler  *handler.ProductHandler
	supplierHandler *handler.SupplierHandler
}

func NewRouter(
	clientHandler *handler.ClientHandler,
	imageHandler *handler.ImageHandler,
	productHandler *handler.ProductHandler,
	supplierHandler *handler.SupplierHandler,
) *Router {
	return &Router{
		clientHandler:   clientHandler,
		imageHandler:    imageHandler,
		productHandler:  productHandler,
		supplierHandler: supplierHandler,
	}
}

func (r *Router) Register(api *gin.RouterGroup) {
	clients := api.Group("/clients")
	{
		clients.POST("", r.clientHandler.AddClient)
		clients.DELETE("/:id", r.clientHandler.DeleteClient)
		clients.GET("", r.clientHandler.GetClientsWithParams)
		clients.PATCH("/:id/address", r.clientHandler.ChangeClientAddress)
	}

	suppliers := api.Group("/suppliers")
	{
		suppliers.POST("", r.supplierHandler.AddSupplier)
		suppliers.DELETE("/:id", r.supplierHandler.DeleteSupplier)
		suppliers.GET("/:id", r.supplierHandler.GetSupplier)
		suppliers.GET("", r.supplierHandler.GetSuppliersWithParams)
		suppliers.PATCH("/:id/address", r.supplierHandler.ChangeSupplierAddress)
	}

	images := api.Group("/images")
	{
		images.PUT("/:id", r.imageHandler.ChangeImage)
		images.DELETE("/:id", r.imageHandler.DeleteImage)
		images.GET("/:id", r.imageHandler.GetImage)
	}

	products := api.Group("/products")
	{
		products.POST("", r.productHandler.AddProduct)
		products.POST("/:id/decrease-stock", r.productHandler.DecreaseProductStock)
		products.GET("/:id", r.productHandler.GetProduct)
		products.GET("", r.productHandler.GetProductsWithParams)
		products.DELETE("/:id", r.productHandler.DeleteProduct)
		products.POST("/:id/image", r.productHandler.AddProductImage)
		products.GET("/:id/image", r.productHandler.GetProductImage)
	}
}
