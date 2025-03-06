package api

import (
	"net/http"
	"path"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/toxic-development/sysmanix/config"
)

// registerSwaggerRoutes registers Swagger documentation routes
func registerSwaggerRoutes(mux *http.ServeMux) {
	cfg := config.GetConfig()

	// Strip any leading/trailing slashes and ensure it starts with a slash
	swaggerPath := cfg.API.SwaggerPath
	swaggerPath = strings.Trim(swaggerPath, "/")

	// Register the Swagger handler at the configured path
	registerRouteWithMiddleware(mux, path.Join(swaggerPath, "*"), serveSwagger, false, nil)
}

// serveSwagger serves Swagger UI using configuration settings
func serveSwagger(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()

	// Determine the correct path for doc.json based on the SwaggerPath
	docJSONPath := path.Join(cfg.API.SwaggerPath, "doc.json")
	if !strings.HasPrefix(docJSONPath, "/") {
		docJSONPath = "/" + docJSONPath
	}

	// Configure the Swagger UI
	httpSwagger.Handler(
		httpSwagger.URL(docJSONPath),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	).ServeHTTP(w, r)
}
