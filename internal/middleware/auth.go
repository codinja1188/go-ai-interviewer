package middleware

import (
	"context"
	"net/http"
)

// ContextKey is the type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "userID"
	// UsernameKey is the context key for username
	UsernameKey ContextKey = "username"
	// UserRoleKey is the context key for user role
	UserRoleKey ContextKey = "userRole"
)

// AuthMiddleware checks if user is authenticated
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session
		session, err := GetSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Check if user is authenticated
		userID, ok := session.Values["userID"].(int)
		if !ok || userID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user info to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, userID)
		if username, ok := session.Values["username"].(string); ok {
			ctx = context.WithValue(ctx, UsernameKey, username)
		}
		if role, ok := session.Values["role"].(string); ok {
			ctx = context.WithValue(ctx, UserRoleKey, role)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// GetUsernameFromContext retrieves username from context
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

// GetUserRoleFromContext retrieves user role from context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
