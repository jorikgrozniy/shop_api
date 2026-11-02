//go:build integration

package integration

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"shop_api/internal/controller/restapi"
	v1 "shop_api/internal/controller/restapi/v1"
	"shop_api/internal/controller/restapi/v1/handler"
	repo "shop_api/internal/repository/postgres"
	"shop_api/internal/service"
	pg "shop_api/pkg/postgres"
)

var (
	pool    *pgxpool.Pool
	baseURL string
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	container, err := tcpostgres.Run(ctx, "postgres:17-alpine",
		tcpostgres.WithDatabase("shop_test"),
		tcpostgres.WithUsername("test_user"),
		tcpostgres.WithPassword("test_password"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres container: %v\n", err)
		os.Exit(1)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection string: %v\n", err)
		terminate(ctx, container)
		os.Exit(1)
	}

	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect to postgres: %v\n", err)
		terminate(ctx, container)
		os.Exit(1)
	}
	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "ping postgres: %v\n", err)
		pool.Close()
		terminate(ctx, container)
		os.Exit(1)
	}

	if err := applyMigrations(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "apply migrations: %v\n", err)
		pool.Close()
		terminate(ctx, container)
		os.Exit(1)
	}

	server := httptest.NewServer(buildEngine())
	baseURL = server.URL + "/api/v1"

	code := m.Run()

	server.Close()
	pool.Close()
	terminate(ctx, container)
	os.Exit(code)
}

func terminate(ctx context.Context, container *tcpostgres.PostgresContainer) {
	if err := container.Terminate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "terminate postgres container: %v\n", err)
	}
}

func applyMigrations(ctx context.Context) error {
	files, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.up.sql"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .up.sql migrations found")
	}
	sort.Strings(files)

	for _, file := range files {
		sql, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("%s: %w", filepath.Base(file), err)
		}
	}
	return nil
}

func buildEngine() *gin.Engine {
	db := &pg.Postgres{Pool: pool}

	addressRepo := repo.NewAddressRepoPostgres(db)
	clientRepo := repo.NewClientRepoPostgres(db)
	imageRepo := repo.NewImageRepoPostgres(db)
	productRepo := repo.NewProductRepoPostgres(db)
	supplierRepo := repo.NewSupplierRepoPostgres(db)

	addressService := service.NewAddressService(addressRepo)
	imageService := service.NewImageService(imageRepo)
	txManager := service.NewTxManager(db)
	clientService := service.NewClientService(clientRepo, addressService, txManager)
	supplierService := service.NewSupplierService(supplierRepo, addressService, txManager)
	productService := service.NewProductService(productRepo, supplierService, imageService, txManager)

	clientHandler := handler.NewClientHandler(clientService, addressService)
	imageHandler := handler.NewImageHandler(imageService)
	productHandler := handler.NewProductHandler(productService, imageService)
	supplierHandler := handler.NewSupplierHandler(supplierService, addressService)

	return restapi.NewRouter(v1.NewRouter(clientHandler, imageHandler, productHandler, supplierHandler))
}
