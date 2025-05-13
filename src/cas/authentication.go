package cas

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/michalhercik/RecSIS/language"
)

const sessionCookieName = "recsis_session_key"

type UserIDFromContext struct{}

func (UserIDFromContext) UserID(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(userIDKey{}).(string)
	if !ok {
		return "", sql.ErrNoRows
	}
	return userID, nil
}

type userIDKey struct{}

type Authentication struct {
	Data           DBManager
	CAS            CAS
	AfterLoginPath string
	loginPath      string
}

func (a Authentication) AuthenticateHTTP(next http.Handler) http.Handler {
	a.loginPath = "/cas/login"
	router := http.NewServeMux()
	router.HandleFunc("/", a.authenticate(next))
	router.HandleFunc(a.loginPath, a.login)
	router.HandleFunc("/logout", a.logout)
	return router
}

func (a Authentication) authenticate(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, a.CAS.loginURLToCAS(a.loginURL(r)), http.StatusFound)
		}
		sessionID, err := r.Cookie(sessionCookieName)
		if err != nil {
			login(w, r)
			return
		}
		userID, err := a.Data.Authenticate(sessionID.Value)
		if err != nil {
			login(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userIDKey{}, userID))
		next.ServeHTTP(w, r)
	})
}

func (a Authentication) login(w http.ResponseWriter, r *http.Request) {
	userID, err := a.CAS.validateTicket(r, a.loginURL(r))
	if err != nil {
		http.Redirect(w, r, a.CAS.loginURLToCAS(a.loginURL(r)), http.StatusFound)
		log.Println(err)
		return
	}
	sessionID, err := a.Data.Login(userID)
	if err != nil {
		log.Println(err)
		return
	}
	a.setSessionCookie(w, sessionID)
	http.Redirect(w, r, a.AfterLoginPath, http.StatusFound)
}

func (a Authentication) logout(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	sessionID, err := r.Cookie(sessionCookieName)
	if err != nil {
		log.Println(err)
		Logout(a.CAS.loginURLToCAS(a.loginURL(r)), t).Render(r.Context(), w)
		return
	}
	userID, err := a.Data.Authenticate(sessionID.Value)
	if err != nil {
		log.Println(err)
		return
	}
	err = a.Data.Logout(userID, sessionID.Value)
	if err != nil {
		log.Println(err)
		return
	}
	a.deleteSessionCookie(w)
	Logout(a.CAS.loginURLToCAS(a.loginURL(r)), t).Render(r.Context(), w)
}

func (a Authentication) setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := a.sessionCookie()
	cookie.Value = sessionID
	http.SetCookie(w, &cookie)
}

func (a Authentication) deleteSessionCookie(w http.ResponseWriter) {
	cookie := a.sessionCookie()
	cookie.MaxAge = -1
	http.SetCookie(w, &cookie)
}

func (a Authentication) sessionCookie() http.Cookie {
	cookie := http.Cookie{
		Name:     sessionCookieName,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	return cookie
}

func (a Authentication) loginURL(r *http.Request) string {
	return "https://" + r.Host + a.loginPath
}
