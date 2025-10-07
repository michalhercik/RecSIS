package cas

/** PACKAGE DESCRIPTION

The cas package provides authentication middleware and utilities for integrating Central Authentication Service (CAS) single sign-on into this application. Its main purpose is to manage user sessions, handle login and logout flows, and securely associate requests with authenticated users. The package abstracts the details of CAS protocol, ticket validation, and session cookie management, so developers can easily add authentication without dealing with low-level CAS API calls.

It has two main structs: UserIDFromContext, which extracts the user ID from the request context, and is injected into servers to access the authenticated user. The second struct is Authentication, which handles authentication (see main -> authenticationHandler) for access to protected parts of the application. It manages the CAS login flow, session management, and logout processes. It is used as a handler middleware to ensure that only authenticated users can access certain routes.

*/

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
		sessionID, err := r.Cookie(sessionCookieName)
		if err != nil {
			a.loginPage(w, r)
			return
		}
		lang := language.FromContext(r.Context())
		userID, err := a.Data.authenticate(sessionID.Value, lang)
		if err != nil {
			a.loginPage(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userIDKey{}, userID))
		next.ServeHTTP(w, r)
	})
}

func (a Authentication) loginPage(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	model := loginModel{
		lang:     lang,
		text:     texts[lang],
		loginURL: a.CAS.loginURLToCAS(a.loginURL(r)),
	}
	err := Login(model).Render(r.Context(), w)
	if err != nil {
		err = errorx.AddContext(err)
		a.Error.Log(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
		HttpOnly: true,
		// TODO: define subdomains in docker and then use http.SameSiteStrictMode
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}

func (a Authentication) loginURL(r *http.Request) string {
	return "https://" + r.Host + a.loginPath
}

type userIDKey struct{}

type Error interface {
	// Logs the provided error.
	Log(err error)

	// Renders an error message to the user as a floating window, with a status code and localized message.
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)

	// Renders a full error page, including title and user ID, for major errors or page-level failures.
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)

	// Renders a fallback error page when a regular page cannot be rendered due to an error.
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)

	// Renders a floating window with error when any component cannot be rendered due to an error.
	CannotRenderComponent(w http.ResponseWriter, r *http.Request, err error, lang language.Language)
}
