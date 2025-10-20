package services

import (
	"testing"

	"github.com/codinja1188/go-ai-interviewer/models"
)

func TestQuestionService_GenerateQuestion(t *testing.T) {
	qs := NewQuestionService()

	tests := []struct {
		name       string
		category   models.QuestionCategory
		difficulty models.DifficultyLevel
		wantErr    bool
	}{
		{
			name:       "Generate Algorithms Beginner Question",
			category:   models.Algorithms,
			difficulty: models.Beginner,
			wantErr:    false,
		},
		{
			name:       "Generate Data Structures Intermediate Question",
			category:   models.DataStructures,
			difficulty: models.Intermediate,
			wantErr:    false,
		},
		{
			name:       "Generate System Design Advanced Question",
			category:   models.SystemDesign,
			difficulty: models.Advanced,
			wantErr:    false,
		},
		{
			name:       "Invalid Category",
			category:   "invalid",
			difficulty: models.Beginner,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			question, err := qs.GenerateQuestion(tt.category, tt.difficulty)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if question == nil {
				t.Error("Expected question but got nil")
				return
			}

			if question.Category != tt.category {
				t.Errorf("Expected category %s, got %s", tt.category, question.Category)
			}

			if question.Difficulty != tt.difficulty {
				t.Errorf("Expected difficulty %s, got %s", tt.difficulty, question.Difficulty)
			}

			if question.Title == "" {
				t.Error("Question title should not be empty")
			}

			if question.Description == "" {
				t.Error("Question description should not be empty")
			}

			if len(question.TestCases) == 0 {
				t.Error("Question should have at least one test case")
			}

			if question.TimeLimit <= 0 {
				t.Error("Question time limit should be positive")
			}
		})
	}
}

func TestQuestionService_GetAllCategories(t *testing.T) {
	qs := NewQuestionService()
	categories := qs.GetAllCategories()

	if len(categories) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(categories))
	}

	expectedCategories := map[models.QuestionCategory]bool{
		models.Algorithms:     true,
		models.DataStructures: true,
		models.SystemDesign:   true,
	}

	for _, cat := range categories {
		if !expectedCategories[cat] {
			t.Errorf("Unexpected category: %s", cat)
		}
	}
}

func TestQuestionService_GetAllDifficulties(t *testing.T) {
	qs := NewQuestionService()
	difficulties := qs.GetAllDifficulties()

	if len(difficulties) != 3 {
		t.Errorf("Expected 3 difficulties, got %d", len(difficulties))
	}

	expectedDifficulties := map[models.DifficultyLevel]bool{
		models.Beginner:     true,
		models.Intermediate: true,
		models.Advanced:     true,
	}

	for _, diff := range difficulties {
		if !expectedDifficulties[diff] {
			t.Errorf("Unexpected difficulty: %s", diff)
		}
	}
}
