package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/codinja1188/go-ai-interviewer/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles user authentication
type AuthService struct {
	db *sql.DB
}

// NewAuthService creates a new AuthService
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(username, password, role string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user
	query := "INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)"
	result, err := s.db.Exec(query, username, string(hashedPassword), role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	user := &models.User{
		ID:           int(id),
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	return user, nil
}

// AuthenticateUser authenticates a user by username and password
func (s *AuthService) AuthenticateUser(username, password string) (*models.User, error) {
	query := "SELECT id, username, password_hash, role, created_at FROM users WHERE username = ?"
	row := s.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	query := "SELECT id, username, password_hash, role, created_at FROM users WHERE id = ?"
	row := s.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *AuthService) GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT id, username, password_hash, role, created_at FROM users WHERE username = ?"
	row := s.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return &user, nil
}
