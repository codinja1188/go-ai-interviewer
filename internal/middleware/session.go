package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

// InitSession initializes the session store
func InitSession(secret string) {
	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// GetSession retrieves the session for a request
func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "ai-interviewer-session")
}

// SaveSession saves the session
func SaveSession(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return session.Save(r, w)
}
