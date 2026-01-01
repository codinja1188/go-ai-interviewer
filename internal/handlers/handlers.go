package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/codinja1188/go-ai-interviewer/internal/middleware"
	"github.com/codinja1188/go-ai-interviewer/internal/models"
	"github.com/codinja1188/go-ai-interviewer/internal/services"
)

// Handler holds all the services and dependencies
type Handler struct {
	db                *sql.DB
	authService       *services.AuthService
	questionService   *services.QuestionService
	evaluationService *services.CodeEvaluationService
	scoringService    *services.ScoringService
	interviewService  *services.InterviewService
}

// NewHandler creates a new Handler
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db:                db,
		authService:       services.NewAuthService(db),
		questionService:   services.NewQuestionService(db),
		evaluationService: services.NewCodeEvaluationService(30 * time.Second),
		scoringService:    services.NewScoringService(),
		interviewService:  services.NewInterviewService(db),
	}
}

// renderTemplate renders a specific template
func (h *Handler) renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/"+templateName)
	if err != nil {
		log.Printf("Template parse error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Home handles the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.GetSession(r)
	userID, ok := session.Values["userID"].(int)

	data := map[string]interface{}{
		"IsAuthenticated": ok && userID > 0,
	}

	if ok && userID > 0 {
		if username, ok := session.Values["username"].(string); ok {
			data["Username"] = username
		}
	}

	h.renderTemplate(w, "home.html", data)
}

// Login handles the login page
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.renderTemplate(w, "login.html", nil)
		return
	}

	// Handle POST request
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.authService.AuthenticateUser(username, password)
	if err != nil {
		h.renderTemplate(w, "login.html", map[string]interface{}{
			"Error": "Invalid username or password",
		})
		return
	}

	// Create session
	session, _ := middleware.GetSession(r)
	session.Values["userID"] = user.ID
	session.Values["username"] = user.Username
	session.Values["role"] = user.Role
	middleware.SaveSession(r, w, session)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.renderTemplate(w, "register.html", nil)
		return
	}

	// Handle POST request
	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	if role != "interviewer" && role != "candidate" {
		role = "candidate"
	}

	_, err := h.authService.CreateUser(username, password, role)
	if err != nil {
		h.renderTemplate(w, "register.html", map[string]interface{}{
			"Error": "Failed to create user. Username may already exist.",
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Logout handles user logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.GetSession(r)
	session.Values["userID"] = 0
	session.Values["username"] = ""
	session.Values["role"] = ""
	middleware.SaveSession(r, w, session)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Dashboard handles the dashboard page
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.GetSession(r)
	username, _ := session.Values["username"].(string)
	role, _ := session.Values["role"].(string)

	interviews, err := h.interviewService.GetAllInterviews()
	if err != nil {
		log.Printf("Failed to get interviews: %v", err)
		interviews = []models.Interview{}
	}

	data := map[string]interface{}{
		"Username":        username,
		"Role":            role,
		"Interviews":      interviews,
		"IsAuthenticated": true,
	}

	h.renderTemplate(w, "dashboard.html", data)
}

// Questions handles the questions page
func (h *Handler) Questions(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	difficulty := r.URL.Query().Get("difficulty")

	var questions []models.Question
	var err error

	if category != "" && difficulty != "" {
		questions, err = h.questionService.GetQuestionsByCategoryAndDifficulty(category, difficulty)
	} else if category != "" {
		questions, err = h.questionService.GetQuestionsByCategory(category)
	} else if difficulty != "" {
		questions, err = h.questionService.GetQuestionsByDifficulty(difficulty)
	} else {
		questions, err = h.questionService.GetAllQuestions()
	}

	if err != nil {
		log.Printf("Failed to get questions: %v", err)
		questions = []models.Question{}
	}

	session, _ := middleware.GetSession(r)
	username, _ := session.Values["username"].(string)

	data := map[string]interface{}{
		"Username":        username,
		"Questions":       questions,
		"Category":        category,
		"Difficulty":      difficulty,
		"IsAuthenticated": true,
	}

	h.renderTemplate(w, "questions.html", data)
}

// StartInterview handles starting a new interview
func (h *Handler) StartInterview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := middleware.GetSession(r)
	userID, _ := session.Values["userID"].(int)

	// For simplicity, candidate interviews themselves (candidateID = userID)
	interview, err := h.interviewService.CreateInterview(userID, userID)
	if err != nil {
		http.Error(w, "Failed to create interview", http.StatusInternalServerError)
		return
	}

	// Store interview ID in session
	session.Values["currentInterviewID"] = interview.ID
	middleware.SaveSession(r, w, session)

	http.Redirect(w, r, "/interview", http.StatusSeeOther)
}

// Interview handles the interview coding interface
func (h *Handler) Interview(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.GetSession(r)
	username, _ := session.Values["username"].(string)
	interviewID, ok := session.Values["currentInterviewID"].(int)

	if !ok || interviewID == 0 {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Get a random question (simplified: get first available question)
	questions, err := h.questionService.GetAllQuestions()
	if err != nil || len(questions) == 0 {
		http.Error(w, "No questions available", http.StatusInternalServerError)
		return
	}

	question := questions[0]

	data := map[string]interface{}{
		"Username":        username,
		"InterviewID":     interviewID,
		"Question":        question,
		"IsAuthenticated": true,
	}

	h.renderTemplate(w, "interview.html", data)
}

// SubmitCode handles code submission
func (h *Handler) SubmitCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	interviewIDStr := r.FormValue("interview_id")
	questionIDStr := r.FormValue("question_id")
	code := r.FormValue("code")
	language := r.FormValue("language")

	interviewID, _ := strconv.Atoi(interviewIDStr)
	questionID, _ := strconv.Atoi(questionIDStr)

	// Get question for test cases
	question, err := h.questionService.GetQuestionByID(questionID)
	if err != nil {
		http.Error(w, "Failed to get question", http.StatusInternalServerError)
		return
	}

	// Parse test cases
	testCases, err := services.ParseTestCases(question.TestCases)
	if err != nil {
		log.Printf("Failed to parse test cases: %v", err)
		testCases = []services.TestCase{}
	}

	// Evaluate code
	evalResult, err := h.evaluationService.EvaluateCode(code, language, testCases)
	if err != nil {
		http.Error(w, "Failed to evaluate code", http.StatusInternalServerError)
		return
	}

	// Calculate score
	score := h.scoringService.CalculateScore(evalResult, code)

	// Generate feedback
	feedback := h.scoringService.GenerateFeedback(score, evalResult, code)

	// Save response
	response := &models.Response{
		InterviewID: interviewID,
		QuestionID:  questionID,
		Code:        code,
		Language:    language,
		Score:       score.Total,
		Feedback:    feedback.Comments + " " + feedback.Suggestions,
		Correctness: score.Correctness,
		Efficiency:  score.Efficiency,
		CodeQuality: score.CodeQuality,
		SubmittedAt: time.Now(),
	}

	if err := h.interviewService.SaveResponse(response); err != nil {
		http.Error(w, "Failed to save response", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"score":    score,
		"feedback": feedback,
		"result":   evalResult,
	})
}

// CompleteInterview handles completing an interview
func (h *Handler) CompleteInterview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := middleware.GetSession(r)
	interviewID, ok := session.Values["currentInterviewID"].(int)

	if !ok || interviewID == 0 {
		http.Error(w, "No active interview", http.StatusBadRequest)
		return
	}

	if err := h.interviewService.CompleteInterview(interviewID); err != nil {
		http.Error(w, "Failed to complete interview", http.StatusInternalServerError)
		return
	}

	// Clear current interview from session
	session.Values["currentInterviewID"] = 0
	middleware.SaveSession(r, w, session)

	http.Redirect(w, r, "/transcript/"+strconv.Itoa(interviewID), http.StatusSeeOther)
}

// Transcript handles viewing interview transcripts
func (h *Handler) Transcript(w http.ResponseWriter, r *http.Request) {
	// Extract interview ID from URL
	idStr := r.URL.Path[len("/transcript/"):]
	interviewID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid interview ID", http.StatusBadRequest)
		return
	}

	transcript, err := h.interviewService.GetTranscript(interviewID)
	if err != nil {
		http.Error(w, "Failed to get transcript", http.StatusInternalServerError)
		return
	}

	session, _ := middleware.GetSession(r)
	username, _ := session.Values["username"].(string)

	data := map[string]interface{}{
		"Username":        username,
		"Transcript":      transcript,
		"IsAuthenticated": true,
	}

	h.renderTemplate(w, "transcript.html", data)
}
