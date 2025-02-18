package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/open-tfe/tfe-service/internal/api/router"
	"github.com/open-tfe/tfe-service/internal/initialize"
	"github.com/open-tfe/tfe-service/internal/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Initialize configuration and database
	initialize.Config(logger)
	db := initialize.Database(logger)

	// Initialize services
	service := service.NewService(db, logger)

	jwtSecret := viper.GetString("jwt_secret")
	// Initialize router
	r := router.NewRouter(jwtSecret, service, logger)

	// Configure server
	addr := fmt.Sprintf("%s:%d",
		viper.GetString("server.host"),
		viper.GetInt("server.port"),
	)

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	logger.Info("Server starting on", zap.String("addr", addr))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
