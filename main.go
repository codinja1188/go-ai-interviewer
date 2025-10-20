package main

import (
	"log"

	"github.com/codinja1188/go-ai-interviewer/database"
	"github.com/codinja1188/go-ai-interviewer/handlers"
	"github.com/codinja1188/go-ai-interviewer/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.NewDatabase("./interviewer.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize services
	authService := services.NewAuthService(db)
	questionService := services.NewQuestionService()
	evaluationService := services.NewCodeEvaluationService()

	// Create default users for demo
	_ = authService.Register("interviewer", "password123", "interviewer")
	_ = authService.Register("candidate", "password123", "candidate")

	// Initialize handlers
	handler := handlers.NewHandler(authService, questionService, evaluationService, db)

	// Setup router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// HTML templates
	r.LoadHTMLGlob("templates/*")

	// Public routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)

		// Protected routes
		protected := api.Group("")
		protected.Use(handler.AuthMiddleware())
		{
			protected.POST("/questions/generate", handler.GenerateQuestion)
			protected.POST("/interviews", handler.CreateInterview)
			protected.GET("/interviews", handler.ListInterviews)
			protected.GET("/interviews/:id", handler.GetInterview)
			protected.POST("/interviews/:id/questions", handler.AddQuestionToInterview)
			protected.POST("/interviews/:id/end", handler.EndInterview)
			protected.POST("/submit", handler.SubmitCode)
		}
	}

	// Protected web routes
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", nil)
	})
	r.GET("/candidate-dashboard", func(c *gin.Context) {
		c.HTML(200, "candidate-dashboard.html", nil)
	})
	r.GET("/interview/:id", func(c *gin.Context) {
		c.HTML(200, "interview.html", nil)
	})

	log.Println("🚀 Go AI Interviewer started on http://localhost:8080")
	r.Run(":8080")
}