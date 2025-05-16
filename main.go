package main

import (
	"learn/docs"
	"learn/internal/config"
	"learn/internal/handlers"
	"learn/internal/middleware"
	"learn/internal/repository"
	"learn/internal/service"
	"log"
	"net/http"
	_ "net/http/pprof" // Import pprof for profiling

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Learn API
// @version         1.0
// @description     A simple authentication API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create repositories
	userRepo := repository.NewInMemoryUserRepository()
	tokenRepo := repository.NewInMemoryTokenRepository()

	// Create services
	authService := service.NewAuthService(userRepo, tokenRepo, cfg)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)

	// Create router
	router := gin.Default()

	// Create JWT middleware
	jwtMiddleware := middleware.JWTMiddleware(cfg)

	// Register routes
	authHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router, jwtMiddleware)

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start pprof server on a separate port (6060 is the conventional port for pprof)
	go func() {
		log.Println("Starting pprof server on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
