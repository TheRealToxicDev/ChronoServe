package api

import (
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/toxic-development/sysmanix/config"
)

// registerSwaggerRoutes registers Swagger documentation routes
func registerSwaggerRoutes(mux *http.ServeMux) {
	cfg := config.GetConfig()

	// Normalize swagger path - make sure it starts with a slash and doesn't end with one
	swaggerPath := "/" + strings.Trim(cfg.API.SwaggerPath, "/")

	// Create the Swagger handler
	handler := httpSwagger.Handler(
		httpSwagger.URL(swaggerPath+"/doc.json"), // The URL where the Swagger JSON will be served
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
	)

	// Register routes - without authentication
	mux.HandleFunc(swaggerPath, func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	mux.HandleFunc(swaggerPath+"/", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}
