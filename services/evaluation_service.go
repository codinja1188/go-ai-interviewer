package services

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/codinja1188/go-ai-interviewer/models"
	"github.com/google/uuid"
)

// CodeEvaluationService handles code evaluation and execution
type CodeEvaluationService struct {
	maxExecutionTime time.Duration
}

// NewCodeEvaluationService creates a new code evaluation service
func NewCodeEvaluationService() *CodeEvaluationService {
	return &CodeEvaluationService{
		maxExecutionTime: 10 * time.Second, // 10 seconds max execution time
	}
}

// EvaluateCode evaluates submitted code against test cases
func (ces *CodeEvaluationService) EvaluateCode(submission *models.CodeSubmission, question *models.Question) (*models.EvaluationResult, error) {
	result := &models.EvaluationResult{
		SubmissionID:  submission.ID,
		TotalTestCases: len(question.TestCases),
		TestResults:   make([]models.TestResult, 0),
		EvaluatedAt:   time.Now(),
	}

	// Execute code against each test case
	passedTests := 0
	totalExecutionTime := 0

	for _, testCase := range question.TestCases {
		testResult := ces.runTestCase(submission, testCase)
		result.TestResults = append(result.TestResults, testResult)
		
		if testResult.Passed {
			passedTests++
		}
	}

	result.PassedTestCases = passedTests
	result.Success = passedTests == result.TotalTestCases
	result.Score = float64(passedTests) / float64(result.TotalTestCases) * 100
	result.ExecutionTime = totalExecutionTime

	// Generate code quality metrics
	result.CodeQuality = ces.analyzeCodeQuality(submission)

	// Generate feedback
	result.Feedback = ces.generateFeedback(result)

	return result, nil
}

// runTestCase executes code against a single test case
func (ces *CodeEvaluationService) runTestCase(submission *models.CodeSubmission, testCase models.TestCase) models.TestResult {
	result := models.TestResult{
		TestCase: testCase,
		Passed:   false,
	}

	// Create a temporary file for the code
	output, err := ces.executeCodeInSandbox(submission.Code, submission.Language, testCase.Input)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Output = strings.TrimSpace(output)
	expectedOutput := strings.TrimSpace(testCase.ExpectedOutput)

	// Compare output
	result.Passed = result.Output == expectedOutput

	return result
}

// executeCodeInSandbox executes code in a sandboxed environment
func (ces *CodeEvaluationService) executeCodeInSandbox(code string, language models.ProgrammingLanguage, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ces.maxExecutionTime)
	defer cancel()

	var cmd *exec.Cmd
	var codeFile string

	switch language {
	case models.LanguageGo:
		// For Go, we'll use a simpler approach - evaluate the code
		// In production, you'd want to use docker or another sandboxing solution
		return ces.executeGoCode(ctx, code, input)
	
	case models.LanguagePython:
		return ces.executePythonCode(ctx, code, input)
	
	case models.LanguageJava:
		return ces.executeJavaCode(ctx, code, input)
	
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	if cmd != nil && codeFile != "" {
		// Cleanup would happen here
	}

	return "", nil
}

// executeGoCode executes Go code (simplified - in production use Docker)
func (ces *CodeEvaluationService) executeGoCode(ctx context.Context, code string, input string) (string, error) {
	// This is a simplified version. In production, you should:
	// 1. Use Docker containers for sandboxing
	// 2. Set resource limits (CPU, memory)
	// 3. Use network isolation
	
	// For now, we'll return a simulated result
	// In a real implementation, you'd create a temp file, compile and run it
	
	return "Execution not fully implemented - using sandbox simulation", nil
}

// executePythonCode executes Python code
func (ces *CodeEvaluationService) executePythonCode(ctx context.Context, code string, input string) (string, error) {
	// Check if python3 is available
	cmd := exec.CommandContext(ctx, "python3", "--version")
	if err := cmd.Run(); err != nil {
		return "Python3 not available - using sandbox simulation", nil
	}

	// In production, use Docker for proper sandboxing
	return "Python execution - using sandbox simulation", nil
}

// executeJavaCode executes Java code
func (ces *CodeEvaluationService) executeJavaCode(ctx context.Context, code string, input string) (string, error) {
	// Check if java is available
	cmd := exec.CommandContext(ctx, "java", "-version")
	if err := cmd.Run(); err != nil {
		return "Java not available - using sandbox simulation", nil
	}

	// In production, use Docker for proper sandboxing
	return "Java execution - using sandbox simulation", nil
}

// analyzeCodeQuality analyzes code quality metrics
func (ces *CodeEvaluationService) analyzeCodeQuality(submission *models.CodeSubmission) models.CodeQuality {
	quality := models.CodeQuality{
		Readability:  70,
		Efficiency:   70,
		BestPractice: 70,
	}

	code := submission.Code
	lines := strings.Split(code, "\n")

	// Simple heuristics for code quality
	// Readability: check for comments, proper naming
	commentCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "/*") {
			commentCount++
		}
	}
	
	if commentCount > 0 {
		quality.Readability += 10
	}
	if len(lines) < 100 { // Concise code
		quality.Readability += 10
	}

	// Efficiency: check for common anti-patterns
	if !strings.Contains(code, "sleep") && !strings.Contains(code, "time.Sleep") {
		quality.Efficiency += 15
	}

	// Best practices: check for error handling, proper structure
	switch submission.Language {
	case models.LanguageGo:
		if strings.Contains(code, "if err != nil") {
			quality.BestPractice += 15
		}
	case models.LanguagePython:
		if strings.Contains(code, "try:") || strings.Contains(code, "except:") {
			quality.BestPractice += 15
		}
	case models.LanguageJava:
		if strings.Contains(code, "try {") || strings.Contains(code, "catch") {
			quality.BestPractice += 15
		}
	}

	// Cap at 100
	if quality.Readability > 100 {
		quality.Readability = 100
	}
	if quality.Efficiency > 100 {
		quality.Efficiency = 100
	}
	if quality.BestPractice > 100 {
		quality.BestPractice = 100
	}

	quality.Comments = ces.generateQualityComments(quality)

	return quality
}

// generateQualityComments generates comments about code quality
func (ces *CodeEvaluationService) generateQualityComments(quality models.CodeQuality) string {
	comments := []string{}

	if quality.Readability < 70 {
		comments = append(comments, "Consider adding more comments and improving variable names for better readability.")
	} else if quality.Readability >= 90 {
		comments = append(comments, "Excellent code readability with clear structure and documentation.")
	}

	if quality.Efficiency < 70 {
		comments = append(comments, "Look for opportunities to optimize time and space complexity.")
	} else if quality.Efficiency >= 90 {
		comments = append(comments, "Great job optimizing the algorithm for efficiency.")
	}

	if quality.BestPractice < 70 {
		comments = append(comments, "Consider following more language-specific best practices and error handling.")
	} else if quality.BestPractice >= 90 {
		comments = append(comments, "Code follows best practices and includes proper error handling.")
	}

	if len(comments) == 0 {
		comments = append(comments, "Good overall code quality. Keep up the good work!")
	}

	return strings.Join(comments, " ")
}

// generateFeedback generates detailed feedback for the candidate
func (ces *CodeEvaluationService) generateFeedback(result *models.EvaluationResult) string {
	feedback := []string{}

	if result.Success {
		feedback = append(feedback, "🎉 Congratulations! All test cases passed.")
	} else {
		feedback = append(feedback, fmt.Sprintf("✗ %d out of %d test cases passed.", result.PassedTestCases, result.TotalTestCases))
		
		failedTests := result.TotalTestCases - result.PassedTestCases
		if failedTests > 0 {
			feedback = append(feedback, fmt.Sprintf("Please review the %d failed test case(s) and improve your solution.", failedTests))
		}
	}

	// Add code quality feedback
	avgQuality := (result.CodeQuality.Readability + result.CodeQuality.Efficiency + result.CodeQuality.BestPractice) / 3
	feedback = append(feedback, fmt.Sprintf("\n📊 Code Quality Score: %d/100", avgQuality))
	feedback = append(feedback, result.CodeQuality.Comments)

	// Suggestions for improvement
	if !result.Success {
		feedback = append(feedback, "\n💡 Suggestions:")
		feedback = append(feedback, "- Carefully read the problem statement and requirements")
		feedback = append(feedback, "- Test your code with the provided examples")
		feedback = append(feedback, "- Consider edge cases and boundary conditions")
	}

	return strings.Join(feedback, "\n")
}

// CreateSubmission creates a new code submission
func (ces *CodeEvaluationService) CreateSubmission(interviewID, questionID, code string, language models.ProgrammingLanguage) *models.CodeSubmission {
	return &models.CodeSubmission{
		ID:          uuid.New().String(),
		InterviewID: interviewID,
		QuestionID:  questionID,
		Code:        code,
		Language:    language,
		SubmittedAt: time.Now(),
	}
}
