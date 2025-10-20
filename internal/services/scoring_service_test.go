package services

import (
	"testing"
)

func TestScoringService_CalculateScore(t *testing.T) {
	service := NewScoringService()

	tests := []struct {
		name       string
		evalResult *EvaluationResult
		code       string
		wantTotal  float64 // Approximate expected total
	}{
		{
			name: "Perfect score",
			evalResult: &EvaluationResult{
				Success:       true,
				TestsPassed:   5,
				TestsTotal:    5,
				ExecutionTime: 0.5,
			},
			code: `
// This is a well-commented solution
func solution() {
    result := 0
    // Calculate result
    for i := 0; i < 10; i++ {
        result += i
    }
    return result
}`,
			wantTotal: 90.0, // Should be close to max
		},
		{
			name: "Partial correctness",
			evalResult: &EvaluationResult{
				Success:       false,
				TestsPassed:   3,
				TestsTotal:    5,
				ExecutionTime: 2.0,
			},
			code:      "func solution() { return 0 }",
			wantTotal: 40.0, // Lower score
		},
		{
			name: "Slow execution",
			evalResult: &EvaluationResult{
				Success:       true,
				TestsPassed:   5,
				TestsTotal:    5,
				ExecutionTime: 15.0,
			},
			code:      "func solution() { return 0 }",
			wantTotal: 50.0, // Penalized for slow execution
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.CalculateScore(tt.evalResult, tt.code)

			if score.Total < 0 || score.Total > 100 {
				t.Errorf("Total score should be between 0 and 100, got %.2f", score.Total)
			}

			if score.Correctness < 0 || score.Correctness > 40 {
				t.Errorf("Correctness score should be between 0 and 40, got %.2f", score.Correctness)
			}

			if score.Efficiency < 0 || score.Efficiency > 30 {
				t.Errorf("Efficiency score should be between 0 and 30, got %.2f", score.Efficiency)
			}

			if score.CodeQuality < 0 || score.CodeQuality > 30 {
				t.Errorf("Code quality score should be between 0 and 30, got %.2f", score.CodeQuality)
			}

			// Check if total is sum of components
			expectedTotal := score.Correctness + score.Efficiency + score.CodeQuality
			if score.Total != expectedTotal {
				t.Errorf("Total score %.2f should equal sum of components %.2f", score.Total, expectedTotal)
			}
		})
	}
}

func TestScoringService_GenerateFeedback(t *testing.T) {
	service := NewScoringService()

	evalResult := &EvaluationResult{
		Success:       true,
		TestsPassed:   5,
		TestsTotal:    5,
		ExecutionTime: 0.5,
	}

	code := `
// Well-documented code
func solution() int {
    result := 0
    return result
}`

	score := service.CalculateScore(evalResult, code)
	feedback := service.GenerateFeedback(score, evalResult, code)

	if feedback.Comments == "" {
		t.Error("Feedback comments should not be empty")
	}

	if feedback.Suggestions == "" {
		t.Error("Feedback suggestions should not be empty")
	}

	if feedback.Score.Total != score.Total {
		t.Error("Feedback score should match calculated score")
	}
}

func TestScoringService_EvaluateCodeQuality(t *testing.T) {
	service := NewScoringService()

	tests := []struct {
		name         string
		code         string
		expectHigher bool // true if this code should score higher than a simple "x"
	}{
		{
			name: "Well-formatted code with comments",
			code: `
// This function calculates the sum
func calculateSum(arr []int) int {
    sum := 0
    for _, num := range arr {
        sum += num
    }
    return sum
}`,
			expectHigher: true,
		},
		{
			name:         "Minimal code",
			code:         "x",
			expectHigher: false,
		},
		{
			name: "Code with descriptive variables",
			code: `
func process() {
    count := 0
    sum := 0
    index := 0
    result := count + sum + index
}`,
			expectHigher: true,
		},
	}

	baseline := service.evaluateCodeQuality("x")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quality := service.evaluateCodeQuality(tt.code)

			if quality < 0 || quality > 30 {
				t.Errorf("Code quality should be between 0 and 30, got %.2f", quality)
			}

			if tt.expectHigher && quality <= baseline {
				t.Errorf("Expected code quality %.2f to be higher than baseline %.2f", quality, baseline)
			}
		})
	}
}
