package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize NATS client
	natsClient, err := services.NewNATSClient(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	// Initialize code evaluator
	evaluator := services.NewEvaluator()

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Global middleware
	router.Use(middleware.ErrorHandler())

	// Public routes
	router.POST("/api/auth/register", handlers.Register)
	router.POST("/api/auth/login", handlers.Login)

	// Protected routes
	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// User routes
		auth.GET("/users/profile", handlers.GetProfile)
		auth.PUT("/users/profile", handlers.UpdateProfile)

		// Problem routes
		auth.GET("/problems", handlers.ListProblems)
		auth.GET("/problems/:id", handlers.GetProblem)
		auth.POST("/problems", handlers.CreateProblem)
		auth.PUT("/problems/:id", handlers.UpdateProblem)
		auth.DELETE("/problems/:id", handlers.DeleteProblem)

		// Contest routes
		auth.GET("/contests", handlers.ListContests)
		auth.GET("/contests/:id", handlers.GetContest)
		auth.POST("/contests", handlers.CreateContest)
		auth.PUT("/contests/:id", handlers.UpdateContest)
		auth.DELETE("/contests/:id", handlers.DeleteContest)
		auth.POST("/contests/:id/register", handlers.RegisterForContest)

		// Submission routes
		auth.POST("/submissions", handlers.CreateSubmission)
		auth.GET("/submissions", handlers.ListSubmissions)
		auth.GET("/submissions/:id", handlers.GetSubmission)
	}

	// Admin routes
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminMiddleware())
	{
		admin.GET("/users", handlers.ListUsers)
		admin.PUT("/users/:id", handlers.UpdateUser)
		admin.DELETE("/users/:id", handlers.DeleteUser)
	}

	// Start the server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 