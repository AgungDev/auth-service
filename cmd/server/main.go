// @title UMMU Auth Service API
// @version 1.0.0
// @description Authentication and Authorization Service for UMMU Microservices
// @contact.name UMMU Development Team
// @host localhost:9001
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer {token}
package main

import (
	"auth_service/internal/config"
	httpServer "auth_service/internal/transport"
	"fmt"

	_ "auth_service/docs"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %v", err))
	}

	server, err := httpServer.NewServer(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize server: %v", err))
	}

	server.RegisterSwaggerRoutes()
	server.Run()
}
