package database

import (
	"os"
	"testing"
	"time"

	"github.com/codinja1188/go-ai-interviewer/models"
	"github.com/google/uuid"
)

func TestDatabase_CreateAndGetUser(t *testing.T) {
	dbPath := "/tmp/test_user.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	user := &models.User{
		ID:        uuid.New().String(),
		Username:  "testuser",
		Password:  "hashedpassword",
		Role:      "interviewer",
		CreatedAt: time.Now(),
	}

	err = db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := db.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected user but got nil")
	}

	if retrieved.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved.Username)
	}

	if retrieved.Role != user.Role {
		t.Errorf("Expected role %s, got %s", user.Role, retrieved.Role)
	}
}

func TestDatabase_InterviewSession(t *testing.T) {
	dbPath := "/tmp/test_interview.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create interview session
	session := &models.InterviewSession{
		ID:             uuid.New().String(),
		CandidateName:  "John Doe",
		CandidateEmail: "john@example.com",
		StartTime:      time.Now(),
		Status:         "active",
		TotalScore:     0,
	}

	err = db.CreateInterviewSession(session)
	if err != nil {
		t.Fatalf("Failed to create interview session: %v", err)
	}

	// Get interview session
	retrieved, err := db.GetInterviewSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get interview session: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected interview session but got nil")
	}

	if retrieved.CandidateName != session.CandidateName {
		t.Errorf("Expected name %s, got %s", session.CandidateName, retrieved.CandidateName)
	}

	if retrieved.Status != session.Status {
		t.Errorf("Expected status %s, got %s", session.Status, retrieved.Status)
	}

	// Update interview session
	endTime := time.Now()
	retrieved.EndTime = &endTime
	retrieved.Status = "completed"
	retrieved.TotalScore = 85.5

	err = db.UpdateInterviewSession(retrieved)
	if err != nil {
		t.Fatalf("Failed to update interview session: %v", err)
	}

	// Get updated session
	updated, err := db.GetInterviewSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get updated session: %v", err)
	}

	if updated.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", updated.Status)
	}

	if updated.TotalScore != 85.5 {
		t.Errorf("Expected score 85.5, got %f", updated.TotalScore)
	}

	if updated.EndTime == nil {
		t.Error("Expected end time to be set")
	}
}

func TestDatabase_SaveAndGetQuestion(t *testing.T) {
	dbPath := "/tmp/test_question.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create interview first
	interviewID := uuid.New().String()
	session := &models.InterviewSession{
		ID:             interviewID,
		CandidateName:  "Test Candidate",
		CandidateEmail: "test@example.com",
		StartTime:      time.Now(),
		Status:         "active",
	}
	db.CreateInterviewSession(session)

	// Create question
	question := &models.Question{
		ID:          uuid.New().String(),
		Category:    models.Algorithms,
		Difficulty:  models.Beginner,
		Title:       "Two Sum",
		Description: "Find two numbers that add up to target",
		TestCases: []models.TestCase{
			{
				Input:          "[2,7,11,15]\n9",
				ExpectedOutput: "[0,1]",
				IsHidden:       false,
			},
		},
		TimeLimit: 300,
		CreatedAt: time.Now(),
	}

	err = db.SaveQuestion(interviewID, question)
	if err != nil {
		t.Fatalf("Failed to save question: %v", err)
	}

	// Get questions
	questions, err := db.GetQuestionsByInterviewID(interviewID)
	if err != nil {
		t.Fatalf("Failed to get questions: %v", err)
	}

	if len(questions) != 1 {
		t.Fatalf("Expected 1 question, got %d", len(questions))
	}

	q := questions[0]
	if q.Title != question.Title {
		t.Errorf("Expected title %s, got %s", question.Title, q.Title)
	}

	if q.Category != question.Category {
		t.Errorf("Expected category %s, got %s", question.Category, q.Category)
	}

	if len(q.TestCases) != len(question.TestCases) {
		t.Errorf("Expected %d test cases, got %d", len(question.TestCases), len(q.TestCases))
	}
}

func TestDatabase_SaveAndGetSubmission(t *testing.T) {
	dbPath := "/tmp/test_submission.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create interview first
	interviewID := uuid.New().String()
	session := &models.InterviewSession{
		ID:             interviewID,
		CandidateName:  "Test Candidate",
		CandidateEmail: "test@example.com",
		StartTime:      time.Now(),
		Status:         "active",
	}
	db.CreateInterviewSession(session)

	// Create submission
	submission := &models.CodeSubmission{
		ID:            uuid.New().String(),
		InterviewID:   interviewID,
		QuestionID:    "q1",
		Code:          "func main() {}",
		Language:      models.LanguageGo,
		SubmittedAt:   time.Now(),
		ExecutionTime: 100,
	}

	err = db.SaveSubmission(submission)
	if err != nil {
		t.Fatalf("Failed to save submission: %v", err)
	}

	// Get submissions
	submissions, err := db.GetSubmissionsByInterviewID(interviewID)
	if err != nil {
		t.Fatalf("Failed to get submissions: %v", err)
	}

	if len(submissions) != 1 {
		t.Fatalf("Expected 1 submission, got %d", len(submissions))
	}

	s := submissions[0]
	if s.Code != submission.Code {
		t.Errorf("Expected code %s, got %s", submission.Code, s.Code)
	}

	if s.Language != submission.Language {
		t.Errorf("Expected language %s, got %s", submission.Language, s.Language)
	}
}

func TestDatabase_ListInterviewSessions(t *testing.T) {
	dbPath := "/tmp/test_list.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create multiple sessions
	for i := 0; i < 5; i++ {
		session := &models.InterviewSession{
			ID:             uuid.New().String(),
			CandidateName:  "Candidate " + string(rune(i+'A')),
			CandidateEmail: "candidate@example.com",
			StartTime:      time.Now(),
			Status:         "active",
		}
		db.CreateInterviewSession(session)
	}

	// List sessions
	sessions, err := db.ListInterviewSessions(10, 0)
	if err != nil {
		t.Fatalf("Failed to list sessions: %v", err)
	}

	if len(sessions) != 5 {
		t.Errorf("Expected 5 sessions, got %d", len(sessions))
	}
}
