package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/onlinejudge/backend/internal/handlers"
	"github.com/onlinejudge/backend/internal/middleware"
	"github.com/onlinejudge/backend/internal/services"
	"github.com/onlinejudge/backend/pkg/broker"
	"github.com/onlinejudge/backend/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize NATS client
	natsClient, err := broker.NewNATSClient()
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	// Initialize evaluator
	evaluator, err := services.NewEvaluator()
	if err != nil {
		log.Fatalf("Failed to initialize evaluator: %v", err)
	}

	// Subscribe to submission evaluations
	if err := natsClient.SubscribeToSubmissions(evaluator); err != nil {
		log.Fatalf("Failed to subscribe to submissions: %v", err)
	}

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	problemHandler := handlers.NewProblemHandler(db)
	contestHandler := handlers.NewContestHandler(db)
	submissionHandler := handlers.NewSubmissionHandler(db, natsClient)

	// Initialize router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/auth/register", userHandler.Register)
		public.POST("/auth/login", userHandler.Login)
		public.GET("/problems", problemHandler.ListProblems)
		public.GET("/problems/:id", problemHandler.GetProblem)
		public.GET("/contests", contestHandler.ListContests)
		public.GET("/contests/:id", contestHandler.GetContest)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.Auth())
	{
		// User routes
		protected.GET("/users/me", userHandler.GetProfile)
		protected.PUT("/users/me", userHandler.UpdateProfile)

		// Problem routes
		protected.POST("/problems", problemHandler.CreateProblem)
		protected.PUT("/problems/:id", problemHandler.UpdateProblem)
		protected.DELETE("/problems/:id", problemHandler.DeleteProblem)

		// Contest routes
		protected.POST("/contests", contestHandler.CreateContest)
		protected.PUT("/contests/:id", contestHandler.UpdateContest)
		protected.DELETE("/contests/:id", contestHandler.DeleteContest)
		protected.POST("/contests/:id/register", contestHandler.RegisterForContest)

		// Submission routes
		protected.POST("/submissions", submissionHandler.Submit)
		protected.GET("/submissions/:id", submissionHandler.GetSubmission)
		protected.GET("/submissions", submissionHandler.ListSubmissions)
		protected.GET("/submissions/:id/results", submissionHandler.GetSubmissionResults)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 