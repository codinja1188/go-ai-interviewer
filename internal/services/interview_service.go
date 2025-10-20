package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codinja1188/go-ai-interviewer/internal/models"
)

// InterviewService handles interview sessions and transcripts
type InterviewService struct {
	db *sql.DB
}

// NewInterviewService creates a new InterviewService
func NewInterviewService(db *sql.DB) *InterviewService {
	return &InterviewService{db: db}
}

// CreateInterview creates a new interview session
func (s *InterviewService) CreateInterview(candidateID, createdBy int) (*models.Interview, error) {
	query := "INSERT INTO interviews (candidate_id, created_by, status) VALUES (?, ?, ?)"
	result, err := s.db.Exec(query, candidateID, createdBy, "in_progress")
	if err != nil {
		return nil, fmt.Errorf("failed to create interview: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	interview := &models.Interview{
		ID:          int(id),
		CandidateID: candidateID,
		CreatedBy:   createdBy,
		Status:      "in_progress",
		TotalScore:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return interview, nil
}

// GetInterview retrieves an interview by ID
func (s *InterviewService) GetInterview(id int) (*models.Interview, error) {
	query := "SELECT id, candidate_id, created_by, status, total_score, created_at, updated_at FROM interviews WHERE id = ?"
	row := s.db.QueryRow(query, id)

	var interview models.Interview
	err := row.Scan(&interview.ID, &interview.CandidateID, &interview.CreatedBy, &interview.Status,
		&interview.TotalScore, &interview.CreatedAt, &interview.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("interview not found")
		}
		return nil, fmt.Errorf("failed to scan interview: %w", err)
	}

	return &interview, nil
}

// SaveResponse saves a candidate's response to a question
func (s *InterviewService) SaveResponse(response *models.Response) error {
	query := `INSERT INTO responses (interview_id, question_id, code, language, score, feedback, correctness, efficiency, code_quality)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, response.InterviewID, response.QuestionID, response.Code,
		response.Language, response.Score, response.Feedback, response.Correctness,
		response.Efficiency, response.CodeQuality)
	if err != nil {
		return fmt.Errorf("failed to save response: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	response.ID = int(id)
	return nil
}

// GetResponses retrieves all responses for an interview
func (s *InterviewService) GetResponses(interviewID int) ([]models.Response, error) {
	query := `SELECT id, interview_id, question_id, code, language, score, feedback, 
		correctness, efficiency, code_quality, submitted_at 
		FROM responses WHERE interview_id = ? ORDER BY submitted_at ASC`

	rows, err := s.db.Query(query, interviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to query responses: %w", err)
	}
	defer rows.Close()

	var responses []models.Response
	for rows.Next() {
		var r models.Response
		err := rows.Scan(&r.ID, &r.InterviewID, &r.QuestionID, &r.Code, &r.Language,
			&r.Score, &r.Feedback, &r.Correctness, &r.Efficiency, &r.CodeQuality, &r.SubmittedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan response: %w", err)
		}
		responses = append(responses, r)
	}

	return responses, nil
}

// UpdateInterviewScore updates the total score for an interview
func (s *InterviewService) UpdateInterviewScore(interviewID int, totalScore float64) error {
	query := "UPDATE interviews SET total_score = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := s.db.Exec(query, totalScore, interviewID)
	if err != nil {
		return fmt.Errorf("failed to update interview score: %w", err)
	}
	return nil
}

// CompleteInterview marks an interview as completed
func (s *InterviewService) CompleteInterview(interviewID int) error {
	// Calculate total score from all responses
	responses, err := s.GetResponses(interviewID)
	if err != nil {
		return fmt.Errorf("failed to get responses: %w", err)
	}

	totalScore := 0.0
	for _, r := range responses {
		totalScore += r.Score
	}

	// Update interview status and score
	query := "UPDATE interviews SET status = ?, total_score = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err = s.db.Exec(query, "completed", totalScore, interviewID)
	if err != nil {
		return fmt.Errorf("failed to complete interview: %w", err)
	}

	return nil
}

// GetTranscript generates a complete transcript for an interview
func (s *InterviewService) GetTranscript(interviewID int) (*models.Transcript, error) {
	// Get interview
	interview, err := s.GetInterview(interviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to get interview: %w", err)
	}

	// Get responses
	responses, err := s.GetResponses(interviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}

	// Get questions for each response
	questionService := NewQuestionService(s.db)
	questions := make([]models.Question, 0, len(responses))
	for _, r := range responses {
		q, err := questionService.GetQuestionByID(r.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get question: %w", err)
		}
		questions = append(questions, *q)
	}

	// Generate overall feedback
	feedback := s.generateOverallFeedback(interview.TotalScore, len(responses))

	transcript := &models.Transcript{
		ID:          interviewID,
		InterviewID: interviewID,
		Questions:   questions,
		Responses:   responses,
		TotalScore:  interview.TotalScore,
		Feedback:    feedback,
	}

	return transcript, nil
}

// generateOverallFeedback generates overall feedback for the interview
func (s *InterviewService) generateOverallFeedback(totalScore float64, numQuestions int) string {
	avgScore := 0.0
	if numQuestions > 0 {
		avgScore = totalScore / float64(numQuestions)
	}

	feedback := fmt.Sprintf("Interview completed. Total score: %.2f/100. ", totalScore)

	if avgScore >= 80 {
		feedback += "Excellent performance! You demonstrated strong problem-solving skills and code quality."
	} else if avgScore >= 60 {
		feedback += "Good performance. You showed solid understanding with room for improvement."
	} else if avgScore >= 40 {
		feedback += "Fair performance. Focus on practicing more problems and improving code efficiency."
	} else {
		feedback += "Needs improvement. Consider reviewing fundamental concepts and practicing more problems."
	}

	return feedback
}

// SaveTranscriptAsJSON saves the transcript as a JSON file
func (s *InterviewService) SaveTranscriptAsJSON(transcript *models.Transcript, filepath string) error {
	data, err := json.MarshalIndent(transcript, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transcript: %w", err)
	}

	// In a real application, you might save this to a file system or cloud storage
	// For now, we'll just return the data as an example
	_ = data
	return nil
}

// GetAllInterviews retrieves all interviews
func (s *InterviewService) GetAllInterviews() ([]models.Interview, error) {
	query := "SELECT id, candidate_id, created_by, status, total_score, created_at, updated_at FROM interviews ORDER BY created_at DESC"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query interviews: %w", err)
	}
	defer rows.Close()

	var interviews []models.Interview
	for rows.Next() {
		var i models.Interview
		err := rows.Scan(&i.ID, &i.CandidateID, &i.CreatedBy, &i.Status, &i.TotalScore, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interview: %w", err)
		}
		interviews = append(interviews, i)
	}

	return interviews, nil
}
