package main

import (
	"auth_service/internal/config"
	"auth_service/internal/repository"
	httpHandler "auth_service/internal/transport/http"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	fmt.Println("Starting auth_service on port", cfg.AppPort)

	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	userRepo := repository.NewUserRepository(db)
	httpHandler.SetUserRepo(userRepo)
	httpHandler.SetConfig(cfg)

	r := gin.Default()
	httpHandler.RegisterRoutes(r)

	r.Run(":" + cfg.AppPort)
}
