package services

import (
	"database/sql"
	"fmt"
	"github.com/codinja1188/go-ai-interviewer/internal/models"
)

// QuestionService handles question generation and retrieval
type QuestionService struct {
	db *sql.DB
}

// NewQuestionService creates a new QuestionService
func NewQuestionService(db *sql.DB) *QuestionService {
	return &QuestionService{db: db}
}

// GetQuestionsByCategory retrieves questions by category
func (s *QuestionService) GetQuestionsByCategory(category string) ([]models.Question, error) {
	query := "SELECT id, category, difficulty, title, description, test_cases, created_at FROM questions WHERE category = ?"
	rows, err := s.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	return s.scanQuestions(rows)
}

// GetQuestionsByDifficulty retrieves questions by difficulty level
func (s *QuestionService) GetQuestionsByDifficulty(difficulty string) ([]models.Question, error) {
	query := "SELECT id, category, difficulty, title, description, test_cases, created_at FROM questions WHERE difficulty = ?"
	rows, err := s.db.Query(query, difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	return s.scanQuestions(rows)
}

// GetQuestionsByCategoryAndDifficulty retrieves questions by category and difficulty
func (s *QuestionService) GetQuestionsByCategoryAndDifficulty(category, difficulty string) ([]models.Question, error) {
	query := "SELECT id, category, difficulty, title, description, test_cases, created_at FROM questions WHERE category = ? AND difficulty = ?"
	rows, err := s.db.Query(query, category, difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	return s.scanQuestions(rows)
}

// GetQuestionByID retrieves a question by ID
func (s *QuestionService) GetQuestionByID(id int) (*models.Question, error) {
	query := "SELECT id, category, difficulty, title, description, test_cases, created_at FROM questions WHERE id = ?"
	row := s.db.QueryRow(query, id)

	var q models.Question
	err := row.Scan(&q.ID, &q.Category, &q.Difficulty, &q.Title, &q.Description, &q.TestCases, &q.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("question not found")
		}
		return nil, fmt.Errorf("failed to scan question: %w", err)
	}

	return &q, nil
}

// GetAllQuestions retrieves all questions
func (s *QuestionService) GetAllQuestions() ([]models.Question, error) {
	query := "SELECT id, category, difficulty, title, description, test_cases, created_at FROM questions"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	return s.scanQuestions(rows)
}

// CreateQuestion creates a new question
func (s *QuestionService) CreateQuestion(q *models.Question) error {
	query := "INSERT INTO questions (category, difficulty, title, description, test_cases) VALUES (?, ?, ?, ?, ?)"
	result, err := s.db.Exec(query, q.Category, q.Difficulty, q.Title, q.Description, q.TestCases)
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	q.ID = int(id)
	return nil
}

// scanQuestions is a helper function to scan multiple questions
func (s *QuestionService) scanQuestions(rows *sql.Rows) ([]models.Question, error) {
	var questions []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(&q.ID, &q.Category, &q.Difficulty, &q.Title, &q.Description, &q.TestCases, &q.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return questions, nil
}
