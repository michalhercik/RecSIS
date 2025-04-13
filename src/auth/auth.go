package auth

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type UserIDFromContext struct{}

func (UserIDFromContext) UserID(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(userIDKey{}).(string)
	if !ok {
		return "", sql.ErrNoRows
	}
	return userID, nil
}

type Authentication struct {
	Authenticate func(string) (string, error)
}

func (a Authentication) AuthenticateHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://cas.cuni.cz/cas/login", http.StatusFound)
		}
		sessionID, err := r.Cookie("recsis_session_key")
		if err != nil {
			login(w, r)
			return
		}
		userID, err := a.Authenticate(sessionID.Value)
		if err != nil {
			login(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userIDKey{}, userID))
		next.ServeHTTP(w, r)
	})
}

// func UserIDFromContext(ctx context.Context) string {
// 	userID, ok := ctx.Value(userIDKey{}).(string)
// 	if !ok {
// 		return ""
// 	}
// 	return userID
// }

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) Authenticate(sessionID string) (string, error) {
	var userID sql.NullString
	err := m.DB.Get(&userID, "SELECT user_id FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		return "", err
	}
	if !userID.Valid {
		return "", sql.ErrNoRows
	}
	return userID.String, nil
}

type userIDKey struct{}

func NoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), userIDKey{}, "81411247"))
		next.ServeHTTP(w, r)
	})
}
