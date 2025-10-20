package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/codinja1188/go-ai-interviewer/models"
	_ "modernc.org/sqlite"
)

// Database handles all database operations
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}
	
	// Initialize schema
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

// initSchema initializes the database schema
func (d *Database) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS interview_sessions (
		id TEXT PRIMARY KEY,
		candidate_name TEXT NOT NULL,
		candidate_email TEXT NOT NULL,
		start_time DATETIME NOT NULL,
		end_time DATETIME,
		status TEXT NOT NULL,
		total_score REAL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS questions (
		id TEXT PRIMARY KEY,
		interview_id TEXT NOT NULL,
		category TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		test_cases TEXT NOT NULL,
		time_limit INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (interview_id) REFERENCES interview_sessions(id)
	);

	CREATE TABLE IF NOT EXISTS submissions (
		id TEXT PRIMARY KEY,
		interview_id TEXT NOT NULL,
		question_id TEXT NOT NULL,
		code TEXT NOT NULL,
		language TEXT NOT NULL,
		submitted_at DATETIME NOT NULL,
		execution_time INTEGER DEFAULT 0,
		FOREIGN KEY (interview_id) REFERENCES interview_sessions(id),
		FOREIGN KEY (question_id) REFERENCES questions(id)
	);

	CREATE TABLE IF NOT EXISTS evaluations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		submission_id TEXT NOT NULL,
		success BOOLEAN NOT NULL,
		passed_test_cases INTEGER NOT NULL,
		total_test_cases INTEGER NOT NULL,
		score REAL NOT NULL,
		feedback TEXT NOT NULL,
		test_results TEXT NOT NULL,
		execution_time INTEGER DEFAULT 0,
		memory_usage INTEGER DEFAULT 0,
		code_quality TEXT NOT NULL,
		evaluated_at DATETIME NOT NULL,
		FOREIGN KEY (submission_id) REFERENCES submissions(id)
	);

	CREATE INDEX IF NOT EXISTS idx_interview_sessions_status ON interview_sessions(status);
	CREATE INDEX IF NOT EXISTS idx_questions_interview ON questions(interview_id);
	CREATE INDEX IF NOT EXISTS idx_submissions_interview ON submissions(interview_id);
	CREATE INDEX IF NOT EXISTS idx_evaluations_submission ON evaluations(submission_id);
	`

	_, err := d.db.Exec(schema)
	return err
}

// CreateUser creates a new user
func (d *Database) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, username, password, role, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := d.db.Exec(query, user.ID, user.Username, user.Password, user.Role, user.CreatedAt)
	return err
}

// GetUserByUsername retrieves a user by username
func (d *Database) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, role, created_at FROM users WHERE username = ?`
	
	user := &models.User{}
	err := d.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// CreateInterviewSession creates a new interview session
func (d *Database) CreateInterviewSession(session *models.InterviewSession) error {
	query := `INSERT INTO interview_sessions (id, candidate_name, candidate_email, start_time, status, total_score) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := d.db.Exec(query, session.ID, session.CandidateName, session.CandidateEmail, 
		session.StartTime, session.Status, session.TotalScore)
	return err
}

// GetInterviewSession retrieves an interview session by ID
func (d *Database) GetInterviewSession(id string) (*models.InterviewSession, error) {
	query := `SELECT id, candidate_name, candidate_email, start_time, end_time, status, total_score 
	          FROM interview_sessions WHERE id = ?`
	
	session := &models.InterviewSession{}
	var endTime sql.NullTime
	
	err := d.db.QueryRow(query, id).Scan(
		&session.ID, &session.CandidateName, &session.CandidateEmail,
		&session.StartTime, &endTime, &session.Status, &session.TotalScore,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if endTime.Valid {
		session.EndTime = &endTime.Time
	}

	// Load associated data
	session.Questions, _ = d.GetQuestionsByInterviewID(id)
	session.Submissions, _ = d.GetSubmissionsByInterviewID(id)
	session.Evaluations, _ = d.GetEvaluationsByInterviewID(id)
	
	return session, nil
}

// UpdateInterviewSession updates an existing interview session
func (d *Database) UpdateInterviewSession(session *models.InterviewSession) error {
	query := `UPDATE interview_sessions SET end_time = ?, status = ?, total_score = ? WHERE id = ?`
	_, err := d.db.Exec(query, session.EndTime, session.Status, session.TotalScore, session.ID)
	return err
}

// ListInterviewSessions lists all interview sessions
func (d *Database) ListInterviewSessions(limit, offset int) ([]*models.InterviewSession, error) {
	query := `SELECT id, candidate_name, candidate_email, start_time, end_time, status, total_score 
	          FROM interview_sessions ORDER BY start_time DESC LIMIT ? OFFSET ?`
	
	rows, err := d.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []*models.InterviewSession{}
	for rows.Next() {
		session := &models.InterviewSession{}
		var endTime sql.NullTime
		
		err := rows.Scan(
			&session.ID, &session.CandidateName, &session.CandidateEmail,
			&session.StartTime, &endTime, &session.Status, &session.TotalScore,
		)
		if err != nil {
			return nil, err
		}
		
		if endTime.Valid {
			session.EndTime = &endTime.Time
		}
		
		sessions = append(sessions, session)
	}
	
	return sessions, nil
}

// SaveQuestion saves a question to the database
func (d *Database) SaveQuestion(interviewID string, question *models.Question) error {
	testCasesJSON, err := json.Marshal(question.TestCases)
	if err != nil {
		return err
	}

	query := `INSERT INTO questions (id, interview_id, category, difficulty, title, description, test_cases, time_limit, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = d.db.Exec(query, question.ID, interviewID, question.Category, question.Difficulty,
		question.Title, question.Description, string(testCasesJSON), question.TimeLimit, question.CreatedAt)
	return err
}

// GetQuestionsByInterviewID retrieves all questions for an interview
func (d *Database) GetQuestionsByInterviewID(interviewID string) ([]models.Question, error) {
	query := `SELECT id, category, difficulty, title, description, test_cases, time_limit, created_at 
	          FROM questions WHERE interview_id = ?`
	
	rows, err := d.db.Query(query, interviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := []models.Question{}
	for rows.Next() {
		question := models.Question{}
		var testCasesJSON string
		
		err := rows.Scan(
			&question.ID, &question.Category, &question.Difficulty, &question.Title,
			&question.Description, &testCasesJSON, &question.TimeLimit, &question.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		json.Unmarshal([]byte(testCasesJSON), &question.TestCases)
		questions = append(questions, question)
	}
	
	return questions, nil
}

// SaveSubmission saves a code submission
func (d *Database) SaveSubmission(submission *models.CodeSubmission) error {
	query := `INSERT INTO submissions (id, interview_id, question_id, code, language, submitted_at, execution_time) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := d.db.Exec(query, submission.ID, submission.InterviewID, submission.QuestionID,
		submission.Code, submission.Language, submission.SubmittedAt, submission.ExecutionTime)
	return err
}

// GetSubmissionsByInterviewID retrieves all submissions for an interview
func (d *Database) GetSubmissionsByInterviewID(interviewID string) ([]models.CodeSubmission, error) {
	query := `SELECT id, interview_id, question_id, code, language, submitted_at, execution_time 
	          FROM submissions WHERE interview_id = ?`
	
	rows, err := d.db.Query(query, interviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := []models.CodeSubmission{}
	for rows.Next() {
		submission := models.CodeSubmission{}
		
		err := rows.Scan(
			&submission.ID, &submission.InterviewID, &submission.QuestionID,
			&submission.Code, &submission.Language, &submission.SubmittedAt, &submission.ExecutionTime,
		)
		if err != nil {
			return nil, err
		}
		
		submissions = append(submissions, submission)
	}
	
	return submissions, nil
}

// SaveEvaluation saves an evaluation result
func (d *Database) SaveEvaluation(evaluation *models.EvaluationResult) error {
	testResultsJSON, err := json.Marshal(evaluation.TestResults)
	if err != nil {
		return err
	}

	codeQualityJSON, err := json.Marshal(evaluation.CodeQuality)
	if err != nil {
		return err
	}

	query := `INSERT INTO evaluations (submission_id, success, passed_test_cases, total_test_cases, 
	          score, feedback, test_results, execution_time, memory_usage, code_quality, evaluated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err = d.db.Exec(query, evaluation.SubmissionID, evaluation.Success, evaluation.PassedTestCases,
		evaluation.TotalTestCases, evaluation.Score, evaluation.Feedback, string(testResultsJSON),
		evaluation.ExecutionTime, evaluation.MemoryUsage, string(codeQualityJSON), evaluation.EvaluatedAt)
	return err
}

// GetEvaluationsByInterviewID retrieves all evaluations for an interview
func (d *Database) GetEvaluationsByInterviewID(interviewID string) ([]models.EvaluationResult, error) {
	query := `SELECT e.submission_id, e.success, e.passed_test_cases, e.total_test_cases, 
	          e.score, e.feedback, e.test_results, e.execution_time, e.memory_usage, 
	          e.code_quality, e.evaluated_at 
	          FROM evaluations e
	          JOIN submissions s ON e.submission_id = s.id
	          WHERE s.interview_id = ?`
	
	rows, err := d.db.Query(query, interviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	evaluations := []models.EvaluationResult{}
	for rows.Next() {
		evaluation := models.EvaluationResult{}
		var testResultsJSON, codeQualityJSON string
		
		err := rows.Scan(
			&evaluation.SubmissionID, &evaluation.Success, &evaluation.PassedTestCases,
			&evaluation.TotalTestCases, &evaluation.Score, &evaluation.Feedback,
			&testResultsJSON, &evaluation.ExecutionTime, &evaluation.MemoryUsage,
			&codeQualityJSON, &evaluation.EvaluatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		json.Unmarshal([]byte(testResultsJSON), &evaluation.TestResults)
		json.Unmarshal([]byte(codeQualityJSON), &evaluation.CodeQuality)
		evaluations = append(evaluations, evaluation)
	}
	
	return evaluations, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
