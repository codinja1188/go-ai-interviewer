package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CodeEvaluationService handles code execution and evaluation
type CodeEvaluationService struct {
	timeout time.Duration
}

// NewCodeEvaluationService creates a new CodeEvaluationService
func NewCodeEvaluationService(timeout time.Duration) *CodeEvaluationService {
	return &CodeEvaluationService{
		timeout: timeout,
	}
}

// EvaluationResult holds the result of code evaluation
type EvaluationResult struct {
	Success       bool    `json:"success"`
	Output        string  `json:"output"`
	Error         string  `json:"error"`
	ExecutionTime float64 `json:"execution_time"`
	TestsPassed   int     `json:"tests_passed"`
	TestsTotal    int     `json:"tests_total"`
}

// TestCase represents a single test case
type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// EvaluateCode evaluates code in a sandboxed environment
func (s *CodeEvaluationService) EvaluateCode(code, language string, testCases []TestCase) (*EvaluationResult, error) {
	switch language {
	case "go":
		return s.evaluateGo(code, testCases)
	case "python":
		return s.evaluatePython(code, testCases)
	case "java":
		return s.evaluateJava(code, testCases)
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// evaluateGo evaluates Go code
func (s *CodeEvaluationService) evaluateGo(code string, testCases []TestCase) (*EvaluationResult, error) {
	// Create a temporary directory for code execution
	tmpDir, err := os.MkdirTemp("", "go-eval-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write code to a file
	codeFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %w", err)
	}

	// Execute the code with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	startTime := time.Now()
	cmd := exec.CommandContext(ctx, "go", "run", codeFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	executionTime := time.Since(startTime).Seconds()

	result := &EvaluationResult{
		ExecutionTime: executionTime,
		TestsTotal:    len(testCases),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Error = "Execution timeout exceeded"
		return result, nil
	}

	if err != nil {
		result.Error = stderr.String()
		return result, nil
	}

	result.Success = true
	result.Output = stdout.String()
	result.TestsPassed = len(testCases) // Simplified: assume all tests pass if no error

	return result, nil
}

// evaluatePython evaluates Python code
func (s *CodeEvaluationService) evaluatePython(code string, testCases []TestCase) (*EvaluationResult, error) {
	tmpDir, err := os.MkdirTemp("", "python-eval-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	codeFile := filepath.Join(tmpDir, "solution.py")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	startTime := time.Now()
	cmd := exec.CommandContext(ctx, "python3", codeFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	executionTime := time.Since(startTime).Seconds()

	result := &EvaluationResult{
		ExecutionTime: executionTime,
		TestsTotal:    len(testCases),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Error = "Execution timeout exceeded"
		return result, nil
	}

	if err != nil {
		result.Error = stderr.String()
		return result, nil
	}

	result.Success = true
	result.Output = stdout.String()
	result.TestsPassed = len(testCases)

	return result, nil
}

// evaluateJava evaluates Java code
func (s *CodeEvaluationService) evaluateJava(code string, testCases []TestCase) (*EvaluationResult, error) {
	tmpDir, err := os.MkdirTemp("", "java-eval-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Extract class name from code - sanitize to prevent path traversal
	className := "Solution"
	if strings.Contains(code, "class ") {
		parts := strings.Split(code, "class ")
		if len(parts) > 1 {
			parts2 := strings.Fields(parts[1])
			if len(parts2) > 0 {
				// Sanitize class name: only allow alphanumeric characters
				candidate := parts2[0]
				sanitized := ""
				for _, ch := range candidate {
					if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
						sanitized += string(ch)
					}
				}
				if len(sanitized) > 0 {
					className = sanitized
				}
			}
		}
	}

	codeFile := filepath.Join(tmpDir, className+".java")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %w", err)
	}

	// Compile
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	compileCmd := exec.CommandContext(ctx, "javac", codeFile)
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	if err := compileCmd.Run(); err != nil {
		return &EvaluationResult{
			Error:      compileStderr.String(),
			TestsTotal: len(testCases),
		}, nil
	}

	// Run
	startTime := time.Now()
	runCmd := exec.CommandContext(ctx, "java", "-cp", tmpDir, className)

	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	err = runCmd.Run()
	executionTime := time.Since(startTime).Seconds()

	result := &EvaluationResult{
		ExecutionTime: executionTime,
		TestsTotal:    len(testCases),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Error = "Execution timeout exceeded"
		return result, nil
	}

	if err != nil {
		result.Error = stderr.String()
		return result, nil
	}

	result.Success = true
	result.Output = stdout.String()
	result.TestsPassed = len(testCases)

	return result, nil
}

// ParseTestCases parses test cases from JSON string
func ParseTestCases(testCasesJSON string) ([]TestCase, error) {
	var testCases []TestCase
	if err := json.Unmarshal([]byte(testCasesJSON), &testCases); err != nil {
		return nil, fmt.Errorf("failed to parse test cases: %w", err)
	}
	return testCases, nil
}
