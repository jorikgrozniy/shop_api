package restapi

import (
	_ "shop_api/docs"
	v1 "shop_api/internal/controller/restapi/v1"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(v1Router *v1.Router) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		v1Router.Register(v1)

		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}
