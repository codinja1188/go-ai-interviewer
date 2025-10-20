package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codinja1188/go-ai-interviewer/internal/database"
	"github.com/codinja1188/go-ai-interviewer/internal/handlers"
	"github.com/codinja1188/go-ai-interviewer/internal/middleware"
)

func main() {
	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./ai_interviewer.db"
	}

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Seed initial questions
	if err := db.SeedQuestions(); err != nil {
		log.Printf("Warning: Failed to seed questions: %v", err)
	}

	// Initialize session
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "your-secret-key-change-this-in-production"
		log.Println("Warning: Using default session secret. Set SESSION_SECRET environment variable in production.")
	}
	middleware.InitSession(sessionSecret)

	// Initialize handlers
	h := handlers.NewHandler(db.DB)

	// Setup routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/logout", h.Logout)

	// Protected routes
	mux.HandleFunc("/dashboard", middleware.AuthMiddleware(h.Dashboard))
	mux.HandleFunc("/questions", middleware.AuthMiddleware(h.Questions))
	mux.HandleFunc("/start-interview", middleware.AuthMiddleware(h.StartInterview))
	mux.HandleFunc("/interview", middleware.AuthMiddleware(h.Interview))
	mux.HandleFunc("/submit-code", middleware.AuthMiddleware(h.SubmitCode))
	mux.HandleFunc("/complete-interview", middleware.AuthMiddleware(h.CompleteInterview))
	mux.HandleFunc("/transcript/", middleware.AuthMiddleware(h.Transcript))

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
