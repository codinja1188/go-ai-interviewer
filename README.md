# Go AI Interviewer

A comprehensive Go-based AI interviewer application that helps candidates practice coding interviews with real-time question generation, code evaluation, scoring, and detailed feedback.

## Features

### 1. **Real-time Question Generation**
- Dynamically generates interview questions based on predefined templates
- Categories: Algorithms, Data Structures, and System Design
- Adjustable difficulty levels: Beginner, Intermediate, Advanced
- Filter questions by category and difficulty

### 2. **Code Evaluation**
- Secure code execution in a sandboxed environment
- Supports multiple programming languages:
  - Go
  - Python
  - Java
- Automatic test case execution
- Timeout protection to prevent infinite loops

### 3. **Score and Feedback System**
- Comprehensive scoring based on three criteria:
  - **Correctness** (40%): Based on test cases passed
  - **Efficiency** (30%): Based on execution time
  - **Code Quality** (30%): Based on comments, formatting, and best practices
- Detailed feedback and improvement suggestions
- Total score out of 100

### 4. **Interview Transcript Saving**
- Complete interview session storage
- Includes questions, candidate responses, evaluations, and feedback
- Stored in SQLite database for easy retrieval
- View historical interview transcripts

### 5. **Web Interface**
- User-friendly web interface for interviewers and candidates
- Features include:
  - User registration and authentication
  - Interview dashboard
  - Real-time coding interface
  - Live feedback display
  - Transcript viewing

## Architecture

```
go-ai-interviewer/
├── cmd/
│   └── server/         # Main application server
├── internal/
│   ├── database/       # Database initialization and management
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # Authentication and session middleware
│   ├── models/         # Data models
│   └── services/       # Business logic services
│       ├── auth_service.go
│       ├── evaluation_service.go
│       ├── interview_service.go
│       ├── question_service.go
│       └── scoring_service.go
├── web/
│   ├── templates/      # HTML templates
│   └── static/         # Static assets (CSS, JS)
└── tests/              # Test files
```

## Prerequisites

- Go 1.24 or higher
- Python 3.x (for Python code evaluation)
- Java JDK (for Java code evaluation)
- SQLite3

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/codinja1188/go-ai-interviewer.git
   cd go-ai-interviewer
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build the application:**
   ```bash
   go build -o ai-interviewer ./cmd/server
   ```

## Running the Application

1. **Start the server:**
   ```bash
   ./ai-interviewer
   ```
   
   Or run directly with Go:
   ```bash
   go run ./cmd/server
   ```

2. **Access the application:**
   Open your browser and navigate to `http://localhost:8080`

3. **Environment Variables (Optional):**
   ```bash
   export PORT=8080                           # Server port (default: 8080)
   export DB_PATH=./ai_interviewer.db        # Database path (default: ./ai_interviewer.db)
   export SESSION_SECRET=your-secret-key     # Session secret for security
   ```

## Usage

### For Candidates

1. **Register an account:**
   - Go to the registration page
   - Create a username and password
   - Select "Candidate" as your role

2. **Start an interview:**
   - Login with your credentials
   - Click "Start New Interview" from the dashboard
   - You'll be presented with a coding question

3. **Submit your solution:**
   - Select your preferred programming language (Go, Python, or Java)
   - Write your solution in the code editor
   - Click "Submit Code" to evaluate your solution
   - Review the immediate feedback including score breakdown and suggestions

4. **Complete the interview:**
   - After solving questions, click "Complete Interview"
   - View your complete transcript with all questions and scores

### For Interviewers

1. **Register as an interviewer:**
   - Create an account with "Interviewer" role

2. **Browse questions:**
   - View available questions in the Questions section
   - Filter by category and difficulty

3. **View interview results:**
   - Access the dashboard to see all interviews
   - Review completed interview transcripts
   - Analyze candidate performance

## Testing

Run all tests:
```bash
go test -v ./...
```

Run specific test packages:
```bash
go test -v ./internal/services/...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## API Endpoints

### Public Endpoints
- `GET /` - Home page
- `GET /login` - Login page
- `POST /login` - Login authentication
- `GET /register` - Registration page
- `POST /register` - User registration

### Protected Endpoints (Require Authentication)
- `GET /dashboard` - User dashboard
- `GET /questions` - Browse questions
- `POST /start-interview` - Start new interview
- `GET /interview` - Interview coding interface
- `POST /submit-code` - Submit code for evaluation
- `POST /complete-interview` - Complete interview session
- `GET /transcript/:id` - View interview transcript
- `GET /logout` - Logout

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Questions Table
```sql
CREATE TABLE questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL,
    difficulty TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    test_cases TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Interviews Table
```sql
CREATE TABLE interviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    candidate_id INTEGER NOT NULL,
    created_by INTEGER NOT NULL,
    status TEXT NOT NULL,
    total_score REAL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (candidate_id) REFERENCES users(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);
```

### Responses Table
```sql
CREATE TABLE responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    interview_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    code TEXT NOT NULL,
    language TEXT NOT NULL,
    score REAL DEFAULT 0,
    feedback TEXT,
    correctness REAL DEFAULT 0,
    efficiency REAL DEFAULT 0,
    code_quality REAL DEFAULT 0,
    submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (interview_id) REFERENCES interviews(id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);
```

## Security Features

- Password hashing using bcrypt
- Session-based authentication with secure cookies
- Code execution in temporary sandboxed environments
- Timeout protection for code execution (30 seconds default)
- SQL injection prevention through prepared statements
- XSS protection through HTML template escaping

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing Guidelines

- Write unit tests for all new services and functions
- Ensure test coverage is above 80%
- Test both success and error cases
- Use table-driven tests for multiple scenarios

## License

This project is open source and available under the MIT License.

## Future Enhancements

- Add more programming languages (JavaScript, C++, Rust)
- Implement AI-powered question generation using LLMs
- Add video/audio recording for interviews
- Implement real-time collaboration features
- Add analytics dashboard for interviewers
- Support for custom test case creation
- Integration with GitHub for code submission
- Email notifications for completed interviews
- Export transcripts as PDF
- Multi-language support for the UI

## Troubleshooting

### Database Issues
If you encounter database errors, delete the database file and restart:
```bash
rm ai_interviewer.db
./ai-interviewer
```

### Port Already in Use
Change the port using the PORT environment variable:
```bash
PORT=8081 ./ai-interviewer
```

### Code Evaluation Timeout
If code evaluation times out, ensure the code doesn't have infinite loops and completes within 30 seconds.

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## Acknowledgments

- Built with Go standard library
- Uses SQLite for data persistence
- Gorilla Sessions for session management
- bcrypt for password hashing