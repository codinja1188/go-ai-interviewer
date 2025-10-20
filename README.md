# Go AI Interviewer

A comprehensive AI-powered technical interview platform built with Go. This application enables interviewers to conduct coding interviews with real-time question generation, code evaluation, and detailed feedback.

## 🚀 Features

### 1. **Real-time Question Generation**
- Dynamic question generation across multiple categories:
  - **Algorithms**: Two Sum, Reverse String, Longest Substring, Merge Intervals, Median of Arrays
  - **Data Structures**: Valid Parentheses, LRU Cache, In-Memory File System
  - **System Design**: URL Shortener, Rate Limiter, Distributed Cache
- Three difficulty levels: Beginner, Intermediate, Advanced
- Includes test cases and time limits for each question

### 2. **Code Evaluation & Execution**
- Multi-language support: Go, Python, Java
- Secure sandboxed code execution
- Test case validation with hidden test cases
- Execution time tracking
- Support for multiple test cases per question

### 3. **Intelligent Scoring & Feedback**
- Automated scoring based on test case pass rate
- Code quality analysis:
  - Readability metrics
  - Efficiency assessment
  - Best practices evaluation
- Detailed feedback with improvement suggestions
- Performance metrics (execution time, memory usage)

### 4. **Interview Session Management**
- Create and manage multiple interview sessions
- Save complete interview transcripts to SQLite database
- Track candidate responses and evaluations
- Session status management (active, completed, cancelled)
- Historical interview data retrieval

### 5. **Modern Web Interface**
- Clean, responsive UI with gradient design
- Separate dashboards for interviewers and candidates
- Real-time coding interface with syntax highlighting
- Login/logout authentication with JWT tokens
- Interactive interview session management
- Live evaluation results display

### 6. **Security & Authentication**
- JWT-based authentication
- Bcrypt password hashing
- Role-based access control (Interviewer/Candidate)
- Secure database operations

## 📋 Prerequisites

- Go 1.21 or higher
- Git

## 🔧 Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/codinja1188/go-ai-interviewer.git
cd go-ai-interviewer
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Build the Application

```bash
go build -o go-ai-interviewer
```

### 4. Run the Application

```bash
./go-ai-interviewer
```

The application will start on `http://localhost:8080`

## 🎯 Quick Start

### Default Demo Accounts

The application comes with two pre-configured demo accounts:

**Interviewer Account:**
- Username: `interviewer`
- Password: `password123`

**Candidate Account:**
- Username: `candidate`
- Password: `password123`

### Access the Application

1. Open your browser and navigate to `http://localhost:8080`
2. Login with one of the demo accounts
3. Interviewers can create new interview sessions and add questions
4. Candidates can join sessions and submit code

## 📚 API Documentation

### Authentication Endpoints

#### Register a New User
```http
POST /api/register
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "role": "interviewer"  // or "candidate"
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "username": "interviewer",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "interviewer",
  "role": "interviewer"
}
```

### Interview Management Endpoints

All these endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your-token>
```

#### Create Interview Session
```http
POST /api/interviews
Content-Type: application/json

{
  "candidate_name": "John Doe",
  "candidate_email": "john@example.com"
}

Response:
{
  "id": "uuid-string",
  "candidate_name": "John Doe",
  "candidate_email": "john@example.com",
  "start_time": "2025-10-20T09:00:00Z",
  "status": "active",
  "total_score": 0
}
```

#### List All Interviews
```http
GET /api/interviews

Response:
{
  "sessions": [
    {
      "id": "uuid-string",
      "candidate_name": "John Doe",
      "candidate_email": "john@example.com",
      "status": "active",
      "total_score": 85.5
    }
  ]
}
```

#### Get Interview Details
```http
GET /api/interviews/:id

Response:
{
  "id": "uuid-string",
  "candidate_name": "John Doe",
  "candidate_email": "john@example.com",
  "start_time": "2025-10-20T09:00:00Z",
  "status": "active",
  "questions": [...],
  "submissions": [...],
  "evaluations": [...],
  "total_score": 85.5
}
```

#### Add Question to Interview
```http
POST /api/interviews/:id/questions
Content-Type: application/json

{
  "category": "algorithms",  // or "data_structures", "system_design"
  "difficulty": "beginner"   // or "intermediate", "advanced"
}

Response:
{
  "id": "question-uuid",
  "category": "algorithms",
  "difficulty": "beginner",
  "title": "Two Sum",
  "description": "Given an array...",
  "test_cases": [...],
  "time_limit": 300
}
```

#### Submit Code for Evaluation
```http
POST /api/submit
Content-Type: application/json

{
  "interview_id": "interview-uuid",
  "question_id": "question-uuid",
  "code": "def solution():\n    return result",
  "language": "python"  // or "go", "java"
}

Response:
{
  "submission_id": "submission-uuid",
  "success": true,
  "passed_test_cases": 3,
  "total_test_cases": 3,
  "score": 100.0,
  "feedback": "🎉 Congratulations! All test cases passed...",
  "test_results": [...],
  "code_quality": {
    "readability": 85,
    "efficiency": 90,
    "best_practice": 80,
    "comments": "Good code quality..."
  }
}
```

#### End Interview Session
```http
POST /api/interviews/:id/end

Response:
{
  "id": "interview-uuid",
  "status": "completed",
  "end_time": "2025-10-20T10:00:00Z",
  "total_score": 85.5
}
```

## 🧪 Testing

### Run All Tests
```bash
go test ./... -v
```

### Run Tests with Coverage
```bash
go test ./... -cover
```

### Run Specific Package Tests
```bash
# Test services
go test ./services/... -v

# Test database
go test ./database/... -v
```

### Test Coverage Report
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 🏗️ Project Structure

```
go-ai-interviewer/
├── main.go                 # Application entry point
├── models/                 # Data models
│   └── models.go          # All data structures
├── services/              # Business logic
│   ├── question_service.go      # Question generation
│   ├── evaluation_service.go    # Code evaluation
│   ├── auth_service.go          # Authentication
│   └── *_test.go               # Unit tests
├── handlers/              # HTTP handlers
│   └── handlers.go        # API endpoints
├── database/              # Database layer
│   ├── database.go        # SQLite operations
│   └── database_test.go   # Database tests
├── templates/             # HTML templates
│   ├── login.html
│   ├── register.html
│   ├── dashboard.html
│   ├── candidate-dashboard.html
│   └── interview.html
├── static/                # Static assets
│   └── css/
│       └── style.css      # Styling
└── README.md             # This file
```

## 🔒 Security Features

- **Password Hashing**: All passwords are hashed using bcrypt
- **JWT Authentication**: Secure token-based authentication
- **SQL Injection Prevention**: Parameterized queries
- **Role-Based Access**: Separate permissions for interviewers and candidates
- **Code Sandboxing**: Isolated execution environment (production-ready implementation requires Docker)

## 🚧 Production Deployment Considerations

For production deployment, consider:

1. **Environment Variables**: Store sensitive data in environment variables
   ```go
   jwtSecret := os.Getenv("JWT_SECRET")
   ```

2. **Database**: Consider PostgreSQL or MySQL for production
   
3. **Code Sandboxing**: Implement Docker-based sandboxing for secure code execution
   
4. **HTTPS**: Use TLS certificates for secure communication

5. **Rate Limiting**: Implement rate limiting on API endpoints

6. **Logging**: Add structured logging with log levels

7. **Monitoring**: Integrate with monitoring tools (Prometheus, Grafana)

## 🔄 Development

### Adding New Questions

Edit `services/question_service.go` and add questions to the `initializeTemplates()` method:

```go
qs.questionTemplates[models.Algorithms][models.Beginner] = append(
    qs.questionTemplates[models.Algorithms][models.Beginner],
    models.Question{
        ID:          "algo-beginner-3",
        Category:    models.Algorithms,
        Difficulty:  models.Beginner,
        Title:       "Your Question Title",
        Description: "Your question description...",
        TestCases: []models.TestCase{
            {
                Input:          "input data",
                ExpectedOutput: "expected output",
                IsHidden:       false,
            },
        },
        TimeLimit: 300,
    },
)
```

### Running in Development Mode

```bash
# With hot reload (requires air)
air

# Manual restart after changes
go run main.go
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- Built with [Gin Web Framework](https://github.com/gin-gonic/gin)
- Database powered by [SQLite](https://www.sqlite.org/)
- JWT authentication using [golang-jwt](https://github.com/golang-jwt/jwt)
- Password hashing with [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## 📧 Support

For issues and questions, please open an issue on GitHub.

---

**Note**: This is a demonstration application. For production use, implement additional security measures, proper code sandboxing with Docker, and comprehensive error handling.