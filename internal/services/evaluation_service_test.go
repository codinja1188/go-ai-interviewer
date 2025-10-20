package services

import (
	"testing"
	"time"
)

func TestCodeEvaluationService_EvaluateGo(t *testing.T) {
	service := NewCodeEvaluationService(10 * time.Second)

	tests := []struct {
		name      string
		code      string
		testCases []TestCase
		wantErr   bool
	}{
		{
			name: "Valid Go code",
			code: `
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
			testCases: []TestCase{},
			wantErr:   false,
		},
		{
			name: "Invalid Go code",
			code: `
package main

func main() {
	// Missing closing brace
`,
			testCases: []TestCase{},
			wantErr:   false, // Service handles error gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.EvaluateCode(tt.code, "go", tt.testCases)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != nil {
				if result.TestsTotal != len(tt.testCases) {
					t.Errorf("Expected %d tests, got %d", len(tt.testCases), result.TestsTotal)
				}

				if result.ExecutionTime < 0 {
					t.Error("Execution time should be non-negative")
				}
			}
		})
	}
}

func TestCodeEvaluationService_EvaluatePython(t *testing.T) {
	service := NewCodeEvaluationService(10 * time.Second)

	code := `print("Hello, Python!")`
	testCases := []TestCase{}

	result, err := service.EvaluateCode(code, "python", testCases)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.TestsTotal != 0 {
		t.Errorf("Expected 0 tests, got %d", result.TestsTotal)
	}
}

func TestCodeEvaluationService_UnsupportedLanguage(t *testing.T) {
	service := NewCodeEvaluationService(10 * time.Second)

	code := `console.log("Hello");`
	testCases := []TestCase{}

	_, err := service.EvaluateCode(code, "javascript", testCases)
	if err == nil {
		t.Error("Expected error for unsupported language")
	}
}

func TestParseTestCases(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantLen int
		wantErr bool
	}{
		{
			name:    "Valid JSON",
			json:    `[{"input": "test1", "output": "out1"}, {"input": "test2", "output": "out2"}]`,
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "Empty array",
			json:    `[]`,
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			json:    `{invalid}`,
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCases, err := ParseTestCases(tt.json)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.wantErr && len(testCases) != tt.wantLen {
				t.Errorf("Expected %d test cases, got %d", tt.wantLen, len(testCases))
			}
		})
	}
}
