package httpserver

import (
	"context"
	"log"
	"net/http"
	"shop_api/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewServer(cfg *config.ServerConfig, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}
}

func RegisterServerLifecycle(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Server starting on %s", server.Addr)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("Server failed to start:", err)
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			log.Println("Closing server")
			if err := server.Shutdown(ctx); err != nil {
				return err
			}

			log.Println("Server closed")
			return nil
		},
	})
}
