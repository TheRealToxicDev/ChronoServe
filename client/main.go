package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/toxic-development/sysmanix/api"
	"github.com/toxic-development/sysmanix/config"
	_ "github.com/toxic-development/sysmanix/docs"
	"github.com/toxic-development/sysmanix/middleware"
	"github.com/toxic-development/sysmanix/utils"
)

var (
	configFile string
)

func init() {
	showVersion := flag.Bool("version", false, "Show version information")
	flag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")
	flag.Parse()

	if *showVersion {
		utils.PrintVersionInfo()
		os.Exit(0)
	}
}

func main() {
	// Initialize configuration
	if err := config.InitConfig(configFile); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	config := config.GetConfig()

	// Initialize logger
	logger, err := utils.NewLogger(utils.LoggerOptions{
		Level:      utils.GetLogLevel(config.Logging.Level),
		Directory:  config.Logging.Directory,
		MaxSize:    config.Logging.MaxSize,
		MaxBackups: config.Logging.MaxBackups,
		Filename:   "app.log",
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	utils.CheckVersionInBackground(logger)
	logger.Info("Version checker started (current: v%s)", utils.Version)

	// Initialize auth middleware
	middleware.InitAuth(middleware.AuthConfig{
		SecretKey:     config.Auth.SecretKey,
		TokenDuration: config.Auth.TokenDuration,
		IssuedBy:      config.Auth.IssuedBy,
	})

	// Setup routes
	router := api.SetupRoutes()

	// Create server with configuration
	readTimeout, _ := time.ParseDuration(config.Server.ReadTimeout)
	writeTimeout, _ := time.ParseDuration(config.Server.WriteTimeout)

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: config.Server.MaxHeaderBytes,
	}

	// Graceful shutdown setup
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Reduced from 30s
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Could not gracefully shutdown the server: %v", err)
		}
		close(done)

		// Force exit if done channel isn't processed
		time.Sleep(2 * time.Second)
		logger.Warn("Forcing exit after timeout")
		os.Exit(0)
	}()

	// Start server
	logger.Info("SysManix is online and awaiting requests")
	logger.Info("Listening on %s:%d", config.Server.Host, config.Server.Port)
	logger.Info("Swagger documentation available at http://%s:%d%s", config.Server.Host, config.Server.Port, config.API.SwaggerPath)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed to start: %v", err)
		os.Exit(1)
	}

	<-done
	logger.Info("Server stopped")
}
