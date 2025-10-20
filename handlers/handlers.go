package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/codinja1188/go-ai-interviewer/database"
	"github.com/codinja1188/go-ai-interviewer/models"
	"github.com/codinja1188/go-ai-interviewer/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler holds all the services needed for HTTP handlers
type Handler struct {
	authService       *services.AuthService
	questionService   *services.QuestionService
	evaluationService *services.CodeEvaluationService
	db                *database.Database
}

// NewHandler creates a new handler
func NewHandler(
	authService *services.AuthService,
	questionService *services.QuestionService,
	evaluationService *services.CodeEvaluationService,
	db *database.Database,
) *Handler {
	return &Handler{
		authService:       authService,
		questionService:   questionService,
		evaluationService: evaluationService,
		db:                db,
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role != "interviewer" && req.Role != "candidate" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be 'interviewer' or 'candidate'"})
		return
	}

	err := h.authService.Register(req.Username, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AuthMiddleware validates JWT tokens
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := h.authService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

// GenerateQuestion generates a new interview question
func (h *Handler) GenerateQuestion(c *gin.Context) {
	var req models.QuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.questionService.GenerateQuestion(req.Category, req.Difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

// CreateInterview creates a new interview session
func (h *Handler) CreateInterview(c *gin.Context) {
	var req struct {
		CandidateName  string `json:"candidate_name" binding:"required"`
		CandidateEmail string `json:"candidate_email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := &models.InterviewSession{
		ID:             uuid.New().String(),
		CandidateName:  req.CandidateName,
		CandidateEmail: req.CandidateEmail,
		StartTime:      time.Now(),
		Status:         "active",
		Questions:      []models.Question{},
		Submissions:    []models.CodeSubmission{},
		Evaluations:    []models.EvaluationResult{},
	}

	err := h.db.CreateInterviewSession(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GetInterview retrieves an interview session
func (h *Handler) GetInterview(c *gin.Context) {
	id := c.Param("id")

	session, err := h.db.GetInterviewSession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "interview not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// ListInterviews lists all interview sessions
func (h *Handler) ListInterviews(c *gin.Context) {
	sessions, err := h.db.ListInterviewSessions(100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// AddQuestionToInterview adds a question to an interview
func (h *Handler) AddQuestionToInterview(c *gin.Context) {
	interviewID := c.Param("id")

	var req models.QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.questionService.GenerateQuestion(req.Category, req.Difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.db.SaveQuestion(interviewID, question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}

// SubmitCode handles code submission
func (h *Handler) SubmitCode(c *gin.Context) {
	var req models.SubmitCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create submission
	submission := h.evaluationService.CreateSubmission(
		req.InterviewID,
		req.QuestionID,
		req.Code,
		req.Language,
	)

	// Save submission
	if err := h.db.SaveSubmission(submission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get question for evaluation
	questions, err := h.db.GetQuestionsByInterviewID(req.InterviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var question *models.Question
	for _, q := range questions {
		if q.ID == req.QuestionID {
			question = &q
			break
		}
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
		return
	}

	// Evaluate code
	evaluation, err := h.evaluationService.EvaluateCode(submission, question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save evaluation
	if err := h.db.SaveEvaluation(evaluation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update interview total score
	session, _ := h.db.GetInterviewSession(req.InterviewID)
	if session != nil {
		totalScore := 0.0
		for _, eval := range session.Evaluations {
			totalScore += eval.Score
		}
		if len(session.Evaluations) > 0 {
			session.TotalScore = totalScore / float64(len(session.Evaluations))
		}
		h.db.UpdateInterviewSession(session)
	}

	c.JSON(http.StatusOK, evaluation)
}

// EndInterview ends an interview session
func (h *Handler) EndInterview(c *gin.Context) {
	interviewID := c.Param("id")

	session, err := h.db.GetInterviewSession(interviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "interview not found"})
		return
	}

	now := time.Now()
	session.EndTime = &now
	session.Status = "completed"

	err = h.db.UpdateInterviewSession(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}
