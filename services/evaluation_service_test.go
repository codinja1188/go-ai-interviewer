package services

import (
	"strings"
	"testing"

	"github.com/codinja1188/go-ai-interviewer/models"
)

func TestCodeEvaluationService_CreateSubmission(t *testing.T) {
	ces := NewCodeEvaluationService()

	submission := ces.CreateSubmission(
		"interview-123",
		"question-456",
		"func main() {}",
		models.LanguageGo,
	)

	if submission == nil {
		t.Fatal("Expected submission to be created")
	}

	if submission.ID == "" {
		t.Error("Submission ID should not be empty")
	}

	if submission.InterviewID != "interview-123" {
		t.Errorf("Expected interview ID 'interview-123', got %s", submission.InterviewID)
	}

	if submission.QuestionID != "question-456" {
		t.Errorf("Expected question ID 'question-456', got %s", submission.QuestionID)
	}

	if submission.Code != "func main() {}" {
		t.Errorf("Expected code 'func main() {}', got %s", submission.Code)
	}

	if submission.Language != models.LanguageGo {
		t.Errorf("Expected language Go, got %s", submission.Language)
	}
}

func TestCodeEvaluationService_AnalyzeCodeQuality(t *testing.T) {
	ces := NewCodeEvaluationService()

	tests := []struct {
		name     string
		code     string
		language models.ProgrammingLanguage
	}{
		{
			name: "Go code with error handling",
			code: `func main() {
				// This is a comment
				if err != nil {
					log.Fatal(err)
				}
			}`,
			language: models.LanguageGo,
		},
		{
			name: "Python code with try-except",
			code: `def main():
				# This is a comment
				try:
					result = some_function()
				except Exception as e:
					print(e)
			`,
			language: models.LanguagePython,
		},
		{
			name: "Java code with try-catch",
			code: `public class Main {
				// This is a comment
				public static void main(String[] args) {
					try {
						doSomething();
					} catch (Exception e) {
						e.printStackTrace();
					}
				}
			}`,
			language: models.LanguageJava,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			submission := &models.CodeSubmission{
				ID:       "test-123",
				Code:     tt.code,
				Language: tt.language,
			}

			quality := ces.analyzeCodeQuality(submission)

			if quality.Readability < 0 || quality.Readability > 100 {
				t.Errorf("Readability should be between 0 and 100, got %d", quality.Readability)
			}

			if quality.Efficiency < 0 || quality.Efficiency > 100 {
				t.Errorf("Efficiency should be between 0 and 100, got %d", quality.Efficiency)
			}

			if quality.BestPractice < 0 || quality.BestPractice > 100 {
				t.Errorf("BestPractice should be between 0 and 100, got %d", quality.BestPractice)
			}

			if quality.Comments == "" {
				t.Error("Quality comments should not be empty")
			}
		})
	}
}

func TestCodeEvaluationService_EvaluateCode(t *testing.T) {
	ces := NewCodeEvaluationService()

	question := &models.Question{
		ID:          "q1",
		Title:       "Test Question",
		Description: "Test",
		TestCases: []models.TestCase{
			{
				Input:          "5",
				ExpectedOutput: "10",
				IsHidden:       false,
			},
		},
	}

	submission := &models.CodeSubmission{
		ID:       "s1",
		Code:     "print(10)",
		Language: models.LanguagePython,
	}

	result, err := ces.EvaluateCode(submission, question)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected evaluation result")
	}

	if result.SubmissionID != submission.ID {
		t.Errorf("Expected submission ID %s, got %s", submission.ID, result.SubmissionID)
	}

	if result.TotalTestCases != len(question.TestCases) {
		t.Errorf("Expected %d total test cases, got %d", len(question.TestCases), result.TotalTestCases)
	}

	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Score should be between 0 and 100, got %f", result.Score)
	}

	if result.Feedback == "" {
		t.Error("Feedback should not be empty")
	}

	if len(result.TestResults) != len(question.TestCases) {
		t.Errorf("Expected %d test results, got %d", len(question.TestCases), len(result.TestResults))
	}
}

func TestCodeEvaluationService_GenerateFeedback(t *testing.T) {
	ces := NewCodeEvaluationService()

	tests := []struct {
		name         string
		result       *models.EvaluationResult
		wantContains []string
	}{
		{
			name: "All tests passed",
			result: &models.EvaluationResult{
				Success:         true,
				PassedTestCases: 3,
				TotalTestCases:  3,
				Score:           100,
				CodeQuality: models.CodeQuality{
					Readability:  90,
					Efficiency:   90,
					BestPractice: 90,
					Comments:     "Great code",
				},
			},
			wantContains: []string{"Congratulations", "passed"},
		},
		{
			name: "Some tests failed",
			result: &models.EvaluationResult{
				Success:         false,
				PassedTestCases: 1,
				TotalTestCases:  3,
				Score:           33.33,
				CodeQuality: models.CodeQuality{
					Readability:  70,
					Efficiency:   70,
					BestPractice: 70,
					Comments:     "Good code",
				},
			},
			wantContains: []string{"1 out of 3", "failed", "Suggestions"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feedback := ces.generateFeedback(tt.result)

			if feedback == "" {
				t.Error("Feedback should not be empty")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(feedback, want) {
					t.Errorf("Expected feedback to contain '%s', but it didn't. Feedback: %s", want, feedback)
				}
			}
		})
	}
}
