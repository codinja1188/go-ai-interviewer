package models

import "time"

// QuestionCategory represents the category of a question
type QuestionCategory string

const (
	Algorithms     QuestionCategory = "algorithms"
	DataStructures QuestionCategory = "data_structures"
	SystemDesign   QuestionCategory = "system_design"
)

// DifficultyLevel represents the difficulty of a question
type DifficultyLevel string

const (
	Beginner     DifficultyLevel = "beginner"
	Intermediate DifficultyLevel = "intermediate"
	Advanced     DifficultyLevel = "advanced"
)

// ProgrammingLanguage represents supported programming languages
type ProgrammingLanguage string

const (
	LanguageGo     ProgrammingLanguage = "go"
	LanguagePython ProgrammingLanguage = "python"
	LanguageJava   ProgrammingLanguage = "java"
)

// Question represents an interview question
type Question struct {
	ID          string           `json:"id"`
	Category    QuestionCategory `json:"category"`
	Difficulty  DifficultyLevel  `json:"difficulty"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	TestCases   []TestCase       `json:"test_cases"`
	TimeLimit   int              `json:"time_limit"` // in seconds
	CreatedAt   time.Time        `json:"created_at"`
}

// TestCase represents a test case for a question
type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	IsHidden       bool   `json:"is_hidden"` // Hidden test cases not shown to candidate
}

// CodeSubmission represents a candidate's code submission
type CodeSubmission struct {
	ID            string              `json:"id"`
	InterviewID   string              `json:"interview_id"`
	QuestionID    string              `json:"question_id"`
	Code          string              `json:"code"`
	Language      ProgrammingLanguage `json:"language"`
	SubmittedAt   time.Time           `json:"submitted_at"`
	ExecutionTime int                 `json:"execution_time"` // in milliseconds
}

// EvaluationResult represents the result of code evaluation
type EvaluationResult struct {
	SubmissionID    string        `json:"submission_id"`
	Success         bool          `json:"success"`
	PassedTestCases int           `json:"passed_test_cases"`
	TotalTestCases  int           `json:"total_test_cases"`
	Score           float64       `json:"score"`
	Feedback        string        `json:"feedback"`
	TestResults     []TestResult  `json:"test_results"`
	ExecutionTime   int           `json:"execution_time"` // in milliseconds
	MemoryUsage     int           `json:"memory_usage"`   // in KB
	CodeQuality     CodeQuality   `json:"code_quality"`
	EvaluatedAt     time.Time     `json:"evaluated_at"`
}

// TestResult represents the result of a single test case
type TestResult struct {
	TestCase TestCase `json:"test_case"`
	Passed   bool     `json:"passed"`
	Output   string   `json:"output"`
	Error    string   `json:"error,omitempty"`
}

// CodeQuality represents code quality metrics
type CodeQuality struct {
	Readability  int    `json:"readability"`  // 0-100
	Efficiency   int    `json:"efficiency"`   // 0-100
	BestPractice int    `json:"best_practice"` // 0-100
	Comments     string `json:"comments"`
}

// InterviewSession represents an interview session
type InterviewSession struct {
	ID            string             `json:"id"`
	CandidateName string             `json:"candidate_name"`
	CandidateEmail string            `json:"candidate_email"`
	StartTime     time.Time          `json:"start_time"`
	EndTime       *time.Time         `json:"end_time,omitempty"`
	Status        string             `json:"status"` // "active", "completed", "cancelled"
	Questions     []Question         `json:"questions"`
	Submissions   []CodeSubmission   `json:"submissions"`
	Evaluations   []EvaluationResult `json:"evaluations"`
	TotalScore    float64            `json:"total_score"`
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never expose in JSON
	Role      string    `json:"role"` // "interviewer" or "candidate"
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// QuestionRequest represents a request to generate a question
type QuestionRequest struct {
	Category   QuestionCategory `json:"category" binding:"required"`
	Difficulty DifficultyLevel  `json:"difficulty" binding:"required"`
}

// SubmitCodeRequest represents a code submission request
type SubmitCodeRequest struct {
	InterviewID string              `json:"interview_id" binding:"required"`
	QuestionID  string              `json:"question_id" binding:"required"`
	Code        string              `json:"code" binding:"required"`
	Language    ProgrammingLanguage `json:"language" binding:"required"`
}
