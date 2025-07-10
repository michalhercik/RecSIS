package cas

import (
	"context"
	"net/http"

	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

const sessionCookieName = "recsis_session_key"

type UserIDFromContext struct{}

func (UserIDFromContext) UserID(r *http.Request) string {
	val := r.Context().Value(userIDKey{})
	userID, _ := val.(string)
	return userID
}

type Authentication struct {
	AfterLoginPath string
	CAS            CAS
	Data           DBManager
	Error          Error
	loginPath      string
}

func (a Authentication) AuthenticateHTTP(next http.Handler) http.Handler {
	a.loginPath = "/cas/login"
	router := http.NewServeMux()
	router.HandleFunc("/", a.authenticate(next))
	router.HandleFunc("GET "+a.loginPath, a.login)
	router.HandleFunc("POST "+a.loginPath, a.logoutFromCAS)
	router.HandleFunc("POST /logout", a.logoutFromUser)
	return router
}

func (a Authentication) authenticate(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := language.FromContext(r.Context())
		login := func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, a.CAS.loginURLToCAS(a.loginURL(r)), http.StatusFound)
		}
		sessionID, err := r.Cookie(sessionCookieName)
		if err != nil {
			login(w, r)
			return
		}
		userID, err := a.Data.authenticate(sessionID.Value, lang)
		if err != nil {
			login(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userIDKey{}, userID))
		next.ServeHTTP(w, r)
	})
}

func (a Authentication) login(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID, ticket, err := a.CAS.validateTicket(r, a.loginURL(r))
	if err != nil {
		http.Redirect(w, r, a.CAS.loginURLToCAS(a.loginURL(r)), http.StatusFound)
		a.Error.Log(errorx.AddContext(err))
		return
	}
	sessionID, err := a.Data.login(userID, ticket, lang)
	if err != nil {
		a.Error.Log(errorx.AddContext(err))
		return
	}
	a.setSessionCookie(w, sessionID)
	http.Redirect(w, r, a.AfterLoginPath, http.StatusFound)
}

func (a Authentication) logoutFromCAS(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID, ticket, err := a.CAS.userIDTicketFromCASLogoutRequest(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		a.Error.Log(errorx.AddContext(err))
		a.Error.Render(w, r, code, userMsg, lang)
		return
	}
	err = a.Data.logoutWithTicket(userID, ticket, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		a.Error.Log(errorx.AddContext(err))
		a.Error.Render(w, r, code, userMsg, lang)
		return
	}
}

func (a Authentication) logoutFromUser(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	sessionID, err := r.Cookie(sessionCookieName)
	if err != nil {
		a.Error.Log(errorx.AddContext(err))
		err = Logout(a.CAS.loginURLToCAS(a.loginURL(r)), t).Render(r.Context(), w)
		if err != nil {
			a.Error.CannotRenderComponent(w, r, err, lang)
		}
		return
	}
	userID, err := a.Data.authenticate(sessionID.Value, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		a.Error.Log(errorx.AddContext(err))
		a.Error.Render(w, r, code, userMsg, lang)
		return
	}
	err = a.Data.logoutWithSession(userID, sessionID.Value, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		a.Error.Log(errorx.AddContext(err))
		a.Error.Render(w, r, code, userMsg, lang)
		return
	}
	a.deleteSessionCookie(w)
	page := Logout(a.CAS.loginURLToCAS(a.loginURL(r)), t)
	err = page.Render(r.Context(), w)
	if err != nil {
		a.Error.CannotRenderComponent(w, r, err, lang)
	}
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

type userIDKey struct{}

type Error interface {
	Log(err error)
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)
	CannotRenderComponent(w http.ResponseWriter, r *http.Request, err error, lang language.Language)
}
