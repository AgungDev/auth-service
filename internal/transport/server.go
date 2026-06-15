package transport

import (
	"auth_service/internal/config"
	"auth_service/internal/repositories"
	"auth_service/internal/security"
	controllersAPI "auth_service/internal/transport/http"
	middleware_http "auth_service/internal/transport/http/middleware_http"
	"auth_service/internal/usecases"
	loggerpkg "auth_service/pkg/logger"
	seederpkg "auth_service/pkg/seeder"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine          *gin.Engine
	host            string
	cfg             *config.Config
	logger          loggerpkg.Logger
	securityService security.SecurityService

	// Repositories
	userRepository       repositories.UserRepositoryInterface
	clientRepository     repositories.ClientRepositoryInterface
	roleRepository       repositories.RoleRepositoryInterface
	permissionRepository repositories.PermissionRepositoryInterface
	tokenRepository      repositories.TokenRepositoryInterface

	// Usecases
	userUsecase       usecases.UserUsecaseInterface
	clientUsecase     usecases.ClientUsecaseInterface
	roleUsecase       usecases.RoleUsecaseInterface
	permissionUsecase usecases.PermissionUsecaseInterface
	tokenUsecase      usecases.TokenUsecaseInterface
}

// NewServer initializes dependencies from config and returns a ready server.
func NewServer(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Connect to database
	dsn := config.GenerateDSN(cfg)
	db, err := config.ConnectDB(dsn, cfg.AppConfig.Name, cfg.AppConfig.Version, cfg.AppConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := seederpkg.SeedDummyData(db); err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	// Initialize repositories
	userRepository := repositories.NewUserRepository(db)
	clientRepository := repositories.NewClientRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	permissionRepository := repositories.NewPermissionRepository(db)
	tokenRepository := repositories.NewTokenRepository(db)

	// Initialize shared services
	securityService := security.NewSecurityService(cfg.JWTConfig)
	logSvc := loggerpkg.NewLogger(cfg.AppConfig.Name, cfg.AppConfig.Version, cfg.AppConfig.Environment)

	// Initialize usecases
	userUsecase := usecases.NewUserUsecase(userRepository, securityService)
	clientUsecase := usecases.NewClientUsecase(clientRepository, securityService, logSvc)
	roleUsecase := usecases.NewRoleUsecase(roleRepository)
	permissionUsecase := usecases.NewPermissionUsecase(permissionRepository)
	tokenUsecase := usecases.NewTokenUsecase(tokenRepository)

	if cfg.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
	gin.DebugPrintFunc = func(format string, args ...interface{}) {}
	engine := gin.New()
	engine.Use(gin.Recovery())
	if err := engine.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}
	host := fmt.Sprintf(":%s", cfg.AppConfig.Port)

	return &Server{
		engine:               engine,
		host:                 host,
		cfg:                  cfg,
		logger:               logSvc,
		securityService:      securityService,
		userRepository:       userRepository,
		clientRepository:     clientRepository,
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		tokenRepository:      tokenRepository,
		userUsecase:          userUsecase,
		clientUsecase:        clientUsecase,
		roleUsecase:          roleUsecase,
		permissionUsecase:    permissionUsecase,
		tokenUsecase:         tokenUsecase,
	}, nil
}

// initRoutes initializes all routes and controllers
func (s *Server) initRoutes() {
	s.engine.Use(gin.Recovery(), middleware_http.RequestLoggerMiddleware(s.logger))
	rg := s.engine.Group("/api/v1")

	// Authentication endpoints
	authUsecase := usecases.NewAuthUsecase(s.userUsecase, s.clientUsecase, s.permissionUsecase, s.tokenUsecase, s.securityService, s.logger, s.cfg.AppConfig.Name, 15, 24)
	authCtrl := controllersAPI.NewAuthController(authUsecase, s.securityService, s.logger)
	authCtrl.Route(rg)

	// Protected resources
	protected := rg.Group("")
	protected.Use(middleware_http.AuthMiddleware(s.securityService, s.logger))
	userCtrl := controllersAPI.NewUserController(s.userUsecase, s.logger)
	userCtrl.Route(protected)

	roleCtrl := controllersAPI.NewRoleController(s.roleUsecase, s.logger)
	roleCtrl.Route(protected)

	permCtrl := controllersAPI.NewPermissionController(s.permissionUsecase, s.logger)
	permCtrl.Route(protected)

	clientCtrl := controllersAPI.NewClientController(s.clientUsecase, s.logger)
	clientCtrl.Route(protected)

	// OAuth endpoints remain open for token exchange
	tokenCtrlConfig := controllersAPI.TokenControllerConfig{
		AppName:               s.cfg.AppConfig.Name,
		AccessTokenTTLMinutes: 15,
	}
	tokenCtrl := controllersAPI.NewTokenController(
		s.userUsecase,
		s.clientUsecase,
		s.tokenUsecase,
		s.securityService,
		tokenCtrlConfig,
		s.logger,
	)
	tokenCtrl.Route(rg)

	s.engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "resource not found"})
	})
}

// Run starts the server
func (s *Server) RegisterSwaggerRoutes() {
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Run() {
	s.initRoutes()
	s.logger.Info("", "", "", "service started", slog.String("port", s.cfg.AppConfig.Port))
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server failed to run on %s: %v", s.host, err))
	}
}
