package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/codinja1188/go-ai-interviewer/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

func setupInterviewTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	queries := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL,
			difficulty TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			test_cases TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE interviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			candidate_id INTEGER NOT NULL,
			created_by INTEGER NOT NULL,
			status TEXT NOT NULL,
			total_score REAL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (candidate_id) REFERENCES users(id),
			FOREIGN KEY (created_by) REFERENCES users(id)
		)`,
		`CREATE TABLE responses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			interview_id INTEGER NOT NULL,
			question_id INTEGER NOT NULL,
			code TEXT NOT NULL,
			language TEXT NOT NULL,
			score REAL DEFAULT 0,
			feedback TEXT,
			correctness REAL DEFAULT 0,
			efficiency REAL DEFAULT 0,
			code_quality REAL DEFAULT 0,
			submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (interview_id) REFERENCES interviews(id),
			FOREIGN KEY (question_id) REFERENCES questions(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to create test table: %v", err)
		}
	}

	// Insert test user
	_, err = db.Exec("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		"testuser", "hash", "candidate")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Insert test question
	_, err = db.Exec("INSERT INTO questions (category, difficulty, title, description, test_cases) VALUES (?, ?, ?, ?, ?)",
		"Algorithms", "beginner", "Test", "Description", "[]")
	if err != nil {
		t.Fatalf("Failed to insert test question: %v", err)
	}

	return db
}

func TestInterviewService_CreateInterview(t *testing.T) {
	db := setupInterviewTestDB(t)
	defer db.Close()

	service := NewInterviewService(db)

	interview, err := service.CreateInterview(1, 1)
	if err != nil {
		t.Fatalf("Failed to create interview: %v", err)
	}

	if interview.ID == 0 {
		t.Error("Interview ID should be set")
	}

	if interview.Status != "in_progress" {
		t.Errorf("Expected status 'in_progress', got '%s'", interview.Status)
	}

	if interview.CandidateID != 1 {
		t.Errorf("Expected candidate ID 1, got %d", interview.CandidateID)
	}
}

func TestInterviewService_SaveResponse(t *testing.T) {
	db := setupInterviewTestDB(t)
	defer db.Close()

	service := NewInterviewService(db)

	// Create interview first
	interview, err := service.CreateInterview(1, 1)
	if err != nil {
		t.Fatalf("Failed to create interview: %v", err)
	}

	response := &models.Response{
		InterviewID: interview.ID,
		QuestionID:  1,
		Code:        "test code",
		Language:    "go",
		Score:       85.5,
		Feedback:    "Good job",
		Correctness: 40.0,
		Efficiency:  25.5,
		CodeQuality: 20.0,
		SubmittedAt: time.Now(),
	}

	err = service.SaveResponse(response)
	if err != nil {
		t.Fatalf("Failed to save response: %v", err)
	}

	if response.ID == 0 {
		t.Error("Response ID should be set")
	}
}

func TestInterviewService_GetResponses(t *testing.T) {
	db := setupInterviewTestDB(t)
	defer db.Close()

	service := NewInterviewService(db)

	interview, err := service.CreateInterview(1, 1)
	if err != nil {
		t.Fatalf("Failed to create interview: %v", err)
	}

	// Save multiple responses
	for i := 0; i < 3; i++ {
		response := &models.Response{
			InterviewID: interview.ID,
			QuestionID:  1,
			Code:        "test code",
			Language:    "go",
			Score:       80.0,
			Feedback:    "Good",
			Correctness: 35.0,
			Efficiency:  25.0,
			CodeQuality: 20.0,
			SubmittedAt: time.Now(),
		}
		if err := service.SaveResponse(response); err != nil {
			t.Fatalf("Failed to save response: %v", err)
		}
	}

	responses, err := service.GetResponses(interview.ID)
	if err != nil {
		t.Fatalf("Failed to get responses: %v", err)
	}

	if len(responses) != 3 {
		t.Errorf("Expected 3 responses, got %d", len(responses))
	}
}

func TestInterviewService_CompleteInterview(t *testing.T) {
	db := setupInterviewTestDB(t)
	defer db.Close()

	service := NewInterviewService(db)

	interview, err := service.CreateInterview(1, 1)
	if err != nil {
		t.Fatalf("Failed to create interview: %v", err)
	}

	// Add a response
	response := &models.Response{
		InterviewID: interview.ID,
		QuestionID:  1,
		Code:        "test",
		Language:    "go",
		Score:       75.0,
		Feedback:    "Good",
		Correctness: 35.0,
		Efficiency:  20.0,
		CodeQuality: 20.0,
		SubmittedAt: time.Now(),
	}
	if err := service.SaveResponse(response); err != nil {
		t.Fatalf("Failed to save response: %v", err)
	}

	// Complete the interview
	err = service.CompleteInterview(interview.ID)
	if err != nil {
		t.Fatalf("Failed to complete interview: %v", err)
	}

	// Verify interview is completed
	completed, err := service.GetInterview(interview.ID)
	if err != nil {
		t.Fatalf("Failed to get interview: %v", err)
	}

	if completed.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", completed.Status)
	}

	if completed.TotalScore != 75.0 {
		t.Errorf("Expected total score 75.0, got %.2f", completed.TotalScore)
	}
}

func TestInterviewService_GetTranscript(t *testing.T) {
	db := setupInterviewTestDB(t)
	defer db.Close()

	service := NewInterviewService(db)

	interview, err := service.CreateInterview(1, 1)
	if err != nil {
		t.Fatalf("Failed to create interview: %v", err)
	}

	response := &models.Response{
		InterviewID: interview.ID,
		QuestionID:  1,
		Code:        "test",
		Language:    "go",
		Score:       80.0,
		Feedback:    "Good",
		Correctness: 35.0,
		Efficiency:  25.0,
		CodeQuality: 20.0,
		SubmittedAt: time.Now(),
	}
	if err := service.SaveResponse(response); err != nil {
		t.Fatalf("Failed to save response: %v", err)
	}

	if err := service.CompleteInterview(interview.ID); err != nil {
		t.Fatalf("Failed to complete interview: %v", err)
	}

	transcript, err := service.GetTranscript(interview.ID)
	if err != nil {
		t.Fatalf("Failed to get transcript: %v", err)
	}

	if transcript.InterviewID != interview.ID {
		t.Errorf("Expected interview ID %d, got %d", interview.ID, transcript.InterviewID)
	}

	if len(transcript.Responses) != 1 {
		t.Errorf("Expected 1 response, got %d", len(transcript.Responses))
	}

	if len(transcript.Questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(transcript.Questions))
	}

	if transcript.Feedback == "" {
		t.Error("Transcript feedback should not be empty")
	}
}
