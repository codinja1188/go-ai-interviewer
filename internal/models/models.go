package models

import "time"

// Question represents an interview question
type Question struct {
	ID          int       `json:"id" db:"id"`
	Category    string    `json:"category" db:"category"`
	Difficulty  string    `json:"difficulty" db:"difficulty"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	TestCases   string    `json:"test_cases" db:"test_cases"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// User represents a user (interviewer or candidate)
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Interview represents an interview session
type Interview struct {
	ID          int       `json:"id" db:"id"`
	CandidateID int       `json:"candidate_id" db:"candidate_id"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
	Status      string    `json:"status" db:"status"`
	TotalScore  float64   `json:"total_score" db:"total_score"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Response represents a candidate's response to a question
type Response struct {
	ID          int       `json:"id" db:"id"`
	InterviewID int       `json:"interview_id" db:"interview_id"`
	QuestionID  int       `json:"question_id" db:"question_id"`
	Code        string    `json:"code" db:"code"`
	Language    string    `json:"language" db:"language"`
	Score       float64   `json:"score" db:"score"`
	Feedback    string    `json:"feedback" db:"feedback"`
	Correctness float64   `json:"correctness" db:"correctness"`
	Efficiency  float64   `json:"efficiency" db:"efficiency"`
	CodeQuality float64   `json:"code_quality" db:"code_quality"`
	SubmittedAt time.Time `json:"submitted_at" db:"submitted_at"`
}

// Transcript represents a complete interview transcript
type Transcript struct {
	ID          int        `json:"id"`
	InterviewID int        `json:"interview_id"`
	Questions   []Question `json:"questions"`
	Responses   []Response `json:"responses"`
	TotalScore  float64    `json:"total_score"`
	Feedback    string     `json:"feedback"`
}

// Category constants
const (
	CategoryAlgorithms     = "Algorithms"
	CategoryDataStructures = "Data Structures"
	CategorySystemDesign   = "System Design"
)

// Difficulty constants
const (
	DifficultyBeginner     = "beginner"
	DifficultyIntermediate = "intermediate"
	DifficultyAdvanced     = "advanced"
)

// Language constants
const (
	LanguageGo     = "go"
	LanguagePython = "python"
	LanguageJava   = "java"
)
