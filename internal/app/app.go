package app

import (
	"shop_api/config"
	"shop_api/internal/controller/restapi"
	v1 "shop_api/internal/controller/restapi/v1"
	"shop_api/internal/controller/restapi/v1/handler"
	"shop_api/internal/repository/postgres"
	"shop_api/internal/service"
	"shop_api/pkg/httpserver"
	pg "shop_api/pkg/postgres"

	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.NewConfig,
			config.NewServerConfig,
			config.NewDatabaseConfig,

			pg.NewPostgres,

			postgres.NewAddressRepoPostgres,
			postgres.NewClientRepoPostgres,
			postgres.NewImageRepoPostgres,
			postgres.NewSupplierRepoPostgres,
			postgres.NewProductRepoPostgres,

			service.NewAddressService,
			service.NewClientService,
			service.NewImageService,
			service.NewSupplierService,
			service.NewProductService,

			handler.NewClientHandler,
			handler.NewImageHandler,
			handler.NewSupplierHandler,
			handler.NewProductHandler,

			v1.NewRouter,

			restapi.NewRouter,

			httpserver.NewServer,
		),

		fx.Invoke(
			httpserver.RegisterServerLifecycle,
		),
	)
}
