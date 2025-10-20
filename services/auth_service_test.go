package services

import (
	"os"
	"testing"

	"github.com/codinja1188/go-ai-interviewer/database"
)

func TestAuthService_Register(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_auth.db"
	defer os.Remove(dbPath)

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	authService := NewAuthService(db)

	tests := []struct {
		name     string
		username string
		password string
		role     string
		wantErr  bool
	}{
		{
			name:     "Register interviewer",
			username: "interviewer1",
			password: "password123",
			role:     "interviewer",
			wantErr:  false,
		},
		{
			name:     "Register candidate",
			username: "candidate1",
			password: "password123",
			role:     "candidate",
			wantErr:  false,
		},
		{
			name:     "Duplicate username",
			username: "interviewer1",
			password: "password456",
			role:     "interviewer",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Register(tt.username, tt.password, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify user was created
			user, err := db.GetUserByUsername(tt.username)
			if err != nil {
				t.Errorf("Failed to get user: %v", err)
			}
			if user == nil {
				t.Error("User should exist after registration")
			}
			if user.Username != tt.username {
				t.Errorf("Expected username %s, got %s", tt.username, user.Username)
			}
			if user.Role != tt.role {
				t.Errorf("Expected role %s, got %s", tt.role, user.Role)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_login.db"
	defer os.Remove(dbPath)

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	authService := NewAuthService(db)

	// Register a test user
	err = authService.Register("testuser", "testpass123", "interviewer")
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid login",
			username: "testuser",
			password: "testpass123",
			wantErr:  false,
		},
		{
			name:     "Invalid password",
			username: "testuser",
			password: "wrongpass",
			wantErr:  true,
		},
		{
			name:     "Non-existent user",
			username: "nonexistent",
			password: "password",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.username, tt.password)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == nil {
				t.Fatal("Expected login response")
			}

			if response.Token == "" {
				t.Error("Token should not be empty")
			}

			if response.Username != tt.username {
				t.Errorf("Expected username %s, got %s", tt.username, response.Username)
			}

			if response.Role == "" {
				t.Error("Role should not be empty")
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_token.db"
	defer os.Remove(dbPath)

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	authService := NewAuthService(db)

	// Register and login to get a valid token
	err = authService.Register("tokenuser", "password123", "interviewer")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	loginResp, err := authService.Login("tokenuser", "password123")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid token",
			token:   loginResp.Token,
			wantErr: false,
		},
		{
			name:    "Invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateToken(tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if claims == nil {
				t.Fatal("Expected claims")
			}

			if claims["username"] != "tokenuser" {
				t.Errorf("Expected username 'tokenuser', got %v", claims["username"])
			}

			if claims["role"] != "interviewer" {
				t.Errorf("Expected role 'interviewer', got %v", claims["role"])
			}
		})
	}
}
