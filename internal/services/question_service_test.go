package services

import (
	"database/sql"
	"testing"

	"github.com/codinja1188/go-ai-interviewer/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables
	queries := []string{
		`CREATE TABLE questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL,
			difficulty TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			test_cases TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to create test table: %v", err)
		}
	}

	return db
}

func TestQuestionService_CreateQuestion(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewQuestionService(db)

	question := &models.Question{
		Category:    models.CategoryAlgorithms,
		Difficulty:  models.DifficultyBeginner,
		Title:       "Test Question",
		Description: "Test Description",
		TestCases:   `[{"input": "test", "output": "test"}]`,
	}

	err := service.CreateQuestion(question)
	if err != nil {
		t.Fatalf("Failed to create question: %v", err)
	}

	if question.ID == 0 {
		t.Error("Question ID should be set after creation")
	}
}

func TestQuestionService_GetQuestionByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewQuestionService(db)

	// Create a question first
	question := &models.Question{
		Category:    models.CategoryAlgorithms,
		Difficulty:  models.DifficultyBeginner,
		Title:       "Test Question",
		Description: "Test Description",
		TestCases:   `[{"input": "test", "output": "test"}]`,
	}

	err := service.CreateQuestion(question)
	if err != nil {
		t.Fatalf("Failed to create question: %v", err)
	}

	// Retrieve the question
	retrieved, err := service.GetQuestionByID(question.ID)
	if err != nil {
		t.Fatalf("Failed to get question: %v", err)
	}

	if retrieved.Title != question.Title {
		t.Errorf("Expected title %s, got %s", question.Title, retrieved.Title)
	}
}

func TestQuestionService_GetQuestionsByCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewQuestionService(db)

	// Create questions in different categories
	questions := []*models.Question{
		{
			Category:    models.CategoryAlgorithms,
			Difficulty:  models.DifficultyBeginner,
			Title:       "Algorithm Question",
			Description: "Test",
			TestCases:   `[]`,
		},
		{
			Category:    models.CategoryDataStructures,
			Difficulty:  models.DifficultyBeginner,
			Title:       "Data Structure Question",
			Description: "Test",
			TestCases:   `[]`,
		},
	}

	for _, q := range questions {
		if err := service.CreateQuestion(q); err != nil {
			t.Fatalf("Failed to create question: %v", err)
		}
	}

	// Get questions by category
	algoQuestions, err := service.GetQuestionsByCategory(models.CategoryAlgorithms)
	if err != nil {
		t.Fatalf("Failed to get questions by category: %v", err)
	}

	if len(algoQuestions) != 1 {
		t.Errorf("Expected 1 algorithm question, got %d", len(algoQuestions))
	}

	if algoQuestions[0].Category != models.CategoryAlgorithms {
		t.Errorf("Expected category %s, got %s", models.CategoryAlgorithms, algoQuestions[0].Category)
	}
}

func TestQuestionService_GetQuestionsByDifficulty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewQuestionService(db)

	// Create questions with different difficulties
	questions := []*models.Question{
		{
			Category:    models.CategoryAlgorithms,
			Difficulty:  models.DifficultyBeginner,
			Title:       "Easy Question",
			Description: "Test",
			TestCases:   `[]`,
		},
		{
			Category:    models.CategoryAlgorithms,
			Difficulty:  models.DifficultyAdvanced,
			Title:       "Hard Question",
			Description: "Test",
			TestCases:   `[]`,
		},
	}

	for _, q := range questions {
		if err := service.CreateQuestion(q); err != nil {
			t.Fatalf("Failed to create question: %v", err)
		}
	}

	// Get questions by difficulty
	beginnerQuestions, err := service.GetQuestionsByDifficulty(models.DifficultyBeginner)
	if err != nil {
		t.Fatalf("Failed to get questions by difficulty: %v", err)
	}

	if len(beginnerQuestions) != 1 {
		t.Errorf("Expected 1 beginner question, got %d", len(beginnerQuestions))
	}

	if beginnerQuestions[0].Difficulty != models.DifficultyBeginner {
		t.Errorf("Expected difficulty %s, got %s", models.DifficultyBeginner, beginnerQuestions[0].Difficulty)
	}
}
