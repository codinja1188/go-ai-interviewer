package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// InitDB initializes the database connection and creates tables
func InitDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return database, nil
}

// createTables creates the necessary database tables
func (db *DB) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL,
			difficulty TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			test_cases TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS interviews (
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
		`CREATE TABLE IF NOT EXISTS responses (
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
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// SeedQuestions adds sample questions to the database
func (db *DB) SeedQuestions() error {
	questions := []struct {
		category    string
		difficulty  string
		title       string
		description string
		testCases   string
	}{
		{
			"Algorithms",
			"beginner",
			"Two Sum",
			"Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
			`[{"input": "[2,7,11,15], 9", "output": "[0,1]"}, {"input": "[3,2,4], 6", "output": "[1,2]"}]`,
		},
		{
			"Data Structures",
			"intermediate",
			"Valid Parentheses",
			"Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.",
			`[{"input": "()", "output": "true"}, {"input": "()[]{}", "output": "true"}, {"input": "(]", "output": "false"}]`,
		},
		{
			"System Design",
			"advanced",
			"Design a URL Shortener",
			"Design a URL shortening service like bit.ly. Explain your design choices, data structures, and scalability considerations.",
			`[{"input": "Design discussion", "output": "Discussion"}]`,
		},
	}

	for _, q := range questions {
		_, err := db.Exec(
			"INSERT INTO questions (category, difficulty, title, description, test_cases) VALUES (?, ?, ?, ?, ?)",
			q.category, q.difficulty, q.title, q.description, q.testCases,
		)
		if err != nil {
			return fmt.Errorf("failed to seed question: %w", err)
		}
	}

	log.Println("Sample questions seeded successfully")
	return nil
}
