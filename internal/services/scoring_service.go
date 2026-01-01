package services

import (
	"strings"
)

// ScoringService handles scoring and feedback generation
type ScoringService struct{}

// NewScoringService creates a new ScoringService
func NewScoringService() *ScoringService {
	return &ScoringService{}
}

// Score represents a detailed score breakdown
type Score struct {
	Correctness float64 `json:"correctness"`
	Efficiency  float64 `json:"efficiency"`
	CodeQuality float64 `json:"code_quality"`
	Total       float64 `json:"total"`
}

// Feedback contains detailed feedback for a submission
type Feedback struct {
	Score       Score  `json:"score"`
	Comments    string `json:"comments"`
	Suggestions string `json:"suggestions"`
}

// CalculateScore calculates the score based on evaluation results and code quality
func (s *ScoringService) CalculateScore(evalResult *EvaluationResult, code string) *Score {
	score := &Score{}

	// Calculate correctness (40% of total)
	if evalResult.TestsTotal > 0 {
		score.Correctness = (float64(evalResult.TestsPassed) / float64(evalResult.TestsTotal)) * 40
	}

	// Calculate efficiency (30% of total) - based on execution time
	// Lower execution time = higher score
	if evalResult.ExecutionTime > 0 {
		if evalResult.ExecutionTime < 1.0 {
			score.Efficiency = 30
		} else if evalResult.ExecutionTime < 5.0 {
			score.Efficiency = 25
		} else if evalResult.ExecutionTime < 10.0 {
			score.Efficiency = 20
		} else {
			score.Efficiency = 15
		}
	}

	// Calculate code quality (30% of total)
	score.CodeQuality = s.evaluateCodeQuality(code)

	// Calculate total score
	score.Total = score.Correctness + score.Efficiency + score.CodeQuality

	return score
}

// evaluateCodeQuality evaluates code quality based on simple heuristics
func (s *ScoringService) evaluateCodeQuality(code string) float64 {
	quality := 30.0 // Start with max score

	// Check for comments
	hasComments := strings.Contains(code, "//") || strings.Contains(code, "/*") || strings.Contains(code, "#")
	if !hasComments {
		quality -= 5
	}

	// Check for proper indentation (simple check)
	lines := strings.Split(code, "\n")
	properIndentation := true
	for _, line := range lines {
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			// Has some indentation
			continue
		}
		if strings.TrimSpace(line) != "" && strings.Contains(line, "{") {
			properIndentation = true
			break
		}
	}
	if !properIndentation {
		quality -= 5
	}

	// Check code length - penalize very short or very long solutions
	codeLength := len(strings.TrimSpace(code))
	if codeLength < 50 {
		quality -= 5
	} else if codeLength > 1000 {
		quality -= 3
	}

	// Check for variable naming (simple check)
	hasDescriptiveNames := strings.Contains(code, "result") ||
		strings.Contains(code, "count") ||
		strings.Contains(code, "sum") ||
		strings.Contains(code, "index")
	if !hasDescriptiveNames {
		quality -= 2
	}

	// Ensure quality is not negative
	if quality < 0 {
		quality = 0
	}

	return quality
}

// GenerateFeedback generates detailed feedback for a submission
func (s *ScoringService) GenerateFeedback(score *Score, evalResult *EvaluationResult, code string) *Feedback {
	feedback := &Feedback{
		Score: *score,
	}

	// Generate comments based on scores
	var comments []string

	if score.Correctness >= 35 {
		comments = append(comments, "Excellent! Your solution passes all test cases.")
	} else if score.Correctness >= 25 {
		comments = append(comments, "Good effort. Your solution passes most test cases.")
	} else if score.Correctness >= 15 {
		comments = append(comments, "Your solution passes some test cases but needs improvement.")
	} else {
		comments = append(comments, "Your solution has significant correctness issues.")
	}

	if score.Efficiency >= 25 {
		comments = append(comments, "Your solution is efficient with good execution time.")
	} else if score.Efficiency >= 20 {
		comments = append(comments, "Your solution's efficiency is acceptable but could be improved.")
	} else {
		comments = append(comments, "Consider optimizing your solution for better performance.")
	}

	if score.CodeQuality >= 25 {
		comments = append(comments, "Your code quality is excellent with good practices.")
	} else if score.CodeQuality >= 20 {
		comments = append(comments, "Your code quality is good but has room for improvement.")
	} else {
		comments = append(comments, "Focus on improving code readability and structure.")
	}

	feedback.Comments = strings.Join(comments, " ")

	// Generate suggestions
	var suggestions []string

	if score.Correctness < 35 {
		suggestions = append(suggestions, "Review the problem requirements and edge cases carefully.")
		suggestions = append(suggestions, "Test your solution with different inputs before submitting.")
	}

	if score.Efficiency < 25 {
		suggestions = append(suggestions, "Consider using more efficient algorithms or data structures.")
		suggestions = append(suggestions, "Analyze the time complexity of your solution.")
	}

	if score.CodeQuality < 25 {
		suggestions = append(suggestions, "Add comments to explain complex logic.")
		suggestions = append(suggestions, "Use meaningful variable names.")
		suggestions = append(suggestions, "Follow proper indentation and formatting conventions.")
	}

	if len(suggestions) > 0 {
		feedback.Suggestions = "Improvement suggestions: " + strings.Join(suggestions, " ")
	} else {
		feedback.Suggestions = "Great work! Keep up the excellent coding practices."
	}

	return feedback
}
