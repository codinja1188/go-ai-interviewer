package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/codinja1188/go-ai-interviewer/models"
	"github.com/google/uuid"
)

// QuestionService handles question generation
type QuestionService struct {
	questionTemplates map[models.QuestionCategory]map[models.DifficultyLevel][]models.Question
}

// NewQuestionService creates a new question service
func NewQuestionService() *QuestionService {
	qs := &QuestionService{
		questionTemplates: make(map[models.QuestionCategory]map[models.DifficultyLevel][]models.Question),
	}
	qs.initializeTemplates()
	return qs
}

// initializeTemplates initializes the question templates
func (qs *QuestionService) initializeTemplates() {
	// Initialize maps
	qs.questionTemplates[models.Algorithms] = make(map[models.DifficultyLevel][]models.Question)
	qs.questionTemplates[models.DataStructures] = make(map[models.DifficultyLevel][]models.Question)
	qs.questionTemplates[models.SystemDesign] = make(map[models.DifficultyLevel][]models.Question)

	// Algorithms - Beginner
	qs.questionTemplates[models.Algorithms][models.Beginner] = []models.Question{
		{
			ID:          "algo-beginner-1",
			Category:    models.Algorithms,
			Difficulty:  models.Beginner,
			Title:       "Two Sum",
			Description: "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\n\nYou may assume that each input would have exactly one solution, and you may not use the same element twice.\n\nExample:\nInput: nums = [2,7,11,15], target = 9\nOutput: [0,1]\nExplanation: Because nums[0] + nums[1] == 9, we return [0, 1].",
			TestCases: []models.TestCase{
				{Input: "[2,7,11,15]\n9", ExpectedOutput: "[0,1]", IsHidden: false},
				{Input: "[3,2,4]\n6", ExpectedOutput: "[1,2]", IsHidden: false},
				{Input: "[3,3]\n6", ExpectedOutput: "[0,1]", IsHidden: true},
			},
			TimeLimit: 300,
		},
		{
			ID:          "algo-beginner-2",
			Category:    models.Algorithms,
			Difficulty:  models.Beginner,
			Title:       "Reverse String",
			Description: "Write a function that reverses a string. The input string is given as an array of characters.\n\nExample:\nInput: ['h','e','l','l','o']\nOutput: ['o','l','l','e','h']",
			TestCases: []models.TestCase{
				{Input: "['h','e','l','l','o']", ExpectedOutput: "['o','l','l','e','h']", IsHidden: false},
				{Input: "['H','a','n','n','a','h']", ExpectedOutput: "['h','a','n','n','a','H']", IsHidden: true},
			},
			TimeLimit: 300,
		},
	}

	// Algorithms - Intermediate
	qs.questionTemplates[models.Algorithms][models.Intermediate] = []models.Question{
		{
			ID:          "algo-intermediate-1",
			Category:    models.Algorithms,
			Difficulty:  models.Intermediate,
			Title:       "Longest Substring Without Repeating Characters",
			Description: "Given a string s, find the length of the longest substring without repeating characters.\n\nExample:\nInput: s = \"abcabcbb\"\nOutput: 3\nExplanation: The answer is \"abc\", with the length of 3.",
			TestCases: []models.TestCase{
				{Input: "abcabcbb", ExpectedOutput: "3", IsHidden: false},
				{Input: "bbbbb", ExpectedOutput: "1", IsHidden: false},
				{Input: "pwwkew", ExpectedOutput: "3", IsHidden: true},
			},
			TimeLimit: 600,
		},
		{
			ID:          "algo-intermediate-2",
			Category:    models.Algorithms,
			Difficulty:  models.Intermediate,
			Title:       "Merge Intervals",
			Description: "Given an array of intervals where intervals[i] = [starti, endi], merge all overlapping intervals.\n\nExample:\nInput: intervals = [[1,3],[2,6],[8,10],[15,18]]\nOutput: [[1,6],[8,10],[15,18]]",
			TestCases: []models.TestCase{
				{Input: "[[1,3],[2,6],[8,10],[15,18]]", ExpectedOutput: "[[1,6],[8,10],[15,18]]", IsHidden: false},
				{Input: "[[1,4],[4,5]]", ExpectedOutput: "[[1,5]]", IsHidden: true},
			},
			TimeLimit: 600,
		},
	}

	// Algorithms - Advanced
	qs.questionTemplates[models.Algorithms][models.Advanced] = []models.Question{
		{
			ID:          "algo-advanced-1",
			Category:    models.Algorithms,
			Difficulty:  models.Advanced,
			Title:       "Median of Two Sorted Arrays",
			Description: "Given two sorted arrays nums1 and nums2 of size m and n respectively, return the median of the two sorted arrays.\n\nThe overall run time complexity should be O(log (m+n)).\n\nExample:\nInput: nums1 = [1,3], nums2 = [2]\nOutput: 2.0",
			TestCases: []models.TestCase{
				{Input: "[1,3]\n[2]", ExpectedOutput: "2.0", IsHidden: false},
				{Input: "[1,2]\n[3,4]", ExpectedOutput: "2.5", IsHidden: true},
			},
			TimeLimit: 900,
		},
	}

	// Data Structures - Beginner
	qs.questionTemplates[models.DataStructures][models.Beginner] = []models.Question{
		{
			ID:          "ds-beginner-1",
			Category:    models.DataStructures,
			Difficulty:  models.Beginner,
			Title:       "Valid Parentheses",
			Description: "Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.\n\nAn input string is valid if:\n1. Open brackets must be closed by the same type of brackets.\n2. Open brackets must be closed in the correct order.\n\nExample:\nInput: s = \"()\"\nOutput: true",
			TestCases: []models.TestCase{
				{Input: "()", ExpectedOutput: "true", IsHidden: false},
				{Input: "()[]{}", ExpectedOutput: "true", IsHidden: false},
				{Input: "(]", ExpectedOutput: "false", IsHidden: true},
			},
			TimeLimit: 300,
		},
	}

	// Data Structures - Intermediate
	qs.questionTemplates[models.DataStructures][models.Intermediate] = []models.Question{
		{
			ID:          "ds-intermediate-1",
			Category:    models.DataStructures,
			Difficulty:  models.Intermediate,
			Title:       "Implement LRU Cache",
			Description: "Design a data structure that follows the constraints of a Least Recently Used (LRU) cache.\n\nImplement the LRUCache class with get(key) and put(key, value) methods.\n\nExample:\nInput: [\"LRUCache\", \"put\", \"put\", \"get\", \"put\", \"get\"]\n[[2], [1, 1], [2, 2], [1], [3, 3], [2]]\nOutput: [null, null, null, 1, null, -1]",
			TestCases: []models.TestCase{
				{Input: "2\nput 1 1\nput 2 2\nget 1\nput 3 3\nget 2", ExpectedOutput: "1\n-1", IsHidden: false},
			},
			TimeLimit: 600,
		},
	}

	// Data Structures - Advanced
	qs.questionTemplates[models.DataStructures][models.Advanced] = []models.Question{
		{
			ID:          "ds-advanced-1",
			Category:    models.DataStructures,
			Difficulty:  models.Advanced,
			Title:       "Design In-Memory File System",
			Description: "Design a data structure that simulates an in-memory file system.\n\nImplement the FileSystem class with operations like ls, mkdir, addContentToFile, and readContentFromFile.",
			TestCases: []models.TestCase{
				{Input: "mkdir /a\naddContent /a/b.txt hello\nreadContent /a/b.txt", ExpectedOutput: "hello", IsHidden: false},
			},
			TimeLimit: 900,
		},
	}

	// System Design - Beginner
	qs.questionTemplates[models.SystemDesign][models.Beginner] = []models.Question{
		{
			ID:          "sd-beginner-1",
			Category:    models.SystemDesign,
			Difficulty:  models.Beginner,
			Title:       "Design a URL Shortener",
			Description: "Design a URL shortening service like bit.ly. Your design should include:\n1. A function to encode a URL to a shortened URL\n2. A function to decode a shortened URL to its original URL\n\nExample:\nInput: encode(\"https://leetcode.com/problems/design-tinyurl\")\nOutput: \"http://tinyurl.com/4e9iAk\"",
			TestCases: []models.TestCase{
				{Input: "encode https://example.com/very/long/url", ExpectedOutput: "short_code", IsHidden: false},
			},
			TimeLimit: 600,
		},
	}

	// System Design - Intermediate
	qs.questionTemplates[models.SystemDesign][models.Intermediate] = []models.Question{
		{
			ID:          "sd-intermediate-1",
			Category:    models.SystemDesign,
			Difficulty:  models.Intermediate,
			Title:       "Design a Rate Limiter",
			Description: "Design a rate limiter that limits the number of requests a user can make within a time window.\n\nImplement the RateLimiter class with allowRequest(userId) method that returns true if the request is allowed, false otherwise.",
			TestCases: []models.TestCase{
				{Input: "limit=5 window=60\nallowRequest user1\nallowRequest user1\nallowRequest user1\nallowRequest user1\nallowRequest user1\nallowRequest user1", ExpectedOutput: "true\ntrue\ntrue\ntrue\ntrue\nfalse", IsHidden: false},
			},
			TimeLimit: 600,
		},
	}

	// System Design - Advanced
	qs.questionTemplates[models.SystemDesign][models.Advanced] = []models.Question{
		{
			ID:          "sd-advanced-1",
			Category:    models.SystemDesign,
			Difficulty:  models.Advanced,
			Title:       "Design a Distributed Cache",
			Description: "Design a distributed caching system like Memcached or Redis. Consider:\n1. Data distribution across multiple nodes\n2. Consistency mechanisms\n3. Fault tolerance\n4. Eviction policies",
			TestCases: []models.TestCase{
				{Input: "put key1 value1\nget key1\ndelete key1\nget key1", ExpectedOutput: "value1\nnull", IsHidden: false},
			},
			TimeLimit: 1200,
		},
	}
}

// GenerateQuestion generates a random question based on category and difficulty
func (qs *QuestionService) GenerateQuestion(category models.QuestionCategory, difficulty models.DifficultyLevel) (*models.Question, error) {
	categoryQuestions, ok := qs.questionTemplates[category]
	if !ok {
		return nil, errors.New("invalid category")
	}

	questions, ok := categoryQuestions[difficulty]
	if !ok || len(questions) == 0 {
		return nil, errors.New("no questions available for this difficulty")
	}

	// Select a random question
	rand.Seed(time.Now().UnixNano())
	selectedQuestion := questions[rand.Intn(len(questions))]

	// Create a new instance with unique ID and timestamp
	question := models.Question{
		ID:          fmt.Sprintf("%s-%s", selectedQuestion.ID, uuid.New().String()[:8]),
		Category:    selectedQuestion.Category,
		Difficulty:  selectedQuestion.Difficulty,
		Title:       selectedQuestion.Title,
		Description: selectedQuestion.Description,
		TestCases:   selectedQuestion.TestCases,
		TimeLimit:   selectedQuestion.TimeLimit,
		CreatedAt:   time.Now(),
	}

	return &question, nil
}

// GetAllCategories returns all available question categories
func (qs *QuestionService) GetAllCategories() []models.QuestionCategory {
	return []models.QuestionCategory{
		models.Algorithms,
		models.DataStructures,
		models.SystemDesign,
	}
}

// GetAllDifficulties returns all available difficulty levels
func (qs *QuestionService) GetAllDifficulties() []models.DifficultyLevel {
	return []models.DifficultyLevel{
		models.Beginner,
		models.Intermediate,
		models.Advanced,
	}
}
