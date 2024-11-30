package main

import (
	"context"
	"emogpt/middleware"
	"emogpt/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, trell-auth-token, trell-app-version-int, creator-space-auth-token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// GracefulShutdown handles graceful shutdown of the server and ticker
func GracefulShutdown(server *http.Server) {
	stopper := make(chan os.Signal, 1)
	// Listen for interrupt and SIGTERM signals
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopper
		zap.L().Info("Shutting down gracefully...")

		// Create a context with a timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shut down the server
		if err := server.Shutdown(ctx); err != nil {
			zap.L().Error("Server shutdown failed", zap.Error(err))
			return
		}
		zap.L().Info("Server exited gracefully")
	}()
}

func main() {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)

	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())

	router.Use(CORSMiddleware())

	routes.Routes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// Create a server instance using gin engine as handler
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Call GracefulShutdown with the server and ticker
	GracefulShutdown(server)

	// Start the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

}
