package errorx

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)

type ErrorHandler struct {
	Page Page
	// TODO: add log targets (file, console, etc.)
	// TODO: add possibilities for configuration
}

// implementing Error interface
func (eh ErrorHandler) Log(err error) {
	log.Println(err)
}

func (eh ErrorHandler) Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language) {
	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Retarget", "#error-content")
		w.Header().Set("HX-Reswap", "innerHTML")
	}
	ErrorMessageTopOfPage(code, userMsg, texts[lang]).Render(r.Context(), w)
}

func (eh ErrorHandler) RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language) {
	main := ErrorMessageContent(code, userMsg, texts[lang])
	eh.Page.View(main, lang, title, userID).Render(r.Context(), w)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

// global error struct for HTTP errors
type HTTPError struct {
	Err     error
	Code    int
	UserMsg string
}

func (he HTTPError) Error() string {
	return he.Err.Error()
}

func NewHTTPErr(err error, code int, userMsg string) HTTPError {
	return HTTPError{
		Err:     err,
		Code:    code,
		UserMsg: userMsg,
	}
}

func (he HTTPError) StatusCode() int {
	return he.Code
}

func (he HTTPError) UserMessage() string {
	return he.UserMsg
}

type AppError interface {
	error
	StatusCode() int
	UserMessage() string
}

type Param struct {
	Name  string
	Value any
}

func P(name string, value any) Param {
	return Param{Name: name, Value: value}
}

// AddContext wraps an error with package.structure.method and parameters context.
func AddContext(err error, params ...Param) error {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("AddContext: failed to get caller info -> original error: %w", err)
	}

	funcFull := runtime.FuncForPC(pc).Name()
	funcParts := strings.Split(funcFull, "/")
	funcName := funcParts[len(funcParts)-1]

	// Build param string: p1=v1, p2=v2, ...
	paramStr := ""
	for i, p := range params {
		if i > 0 {
			paramStr += ", "
		}
		paramStr += fmt.Sprintf("%s=%v", p.Name, p.Value)
	}

	return fmt.Errorf("%s(%s): %w", funcName, paramStr, err)
}

// tries to unwraps user message and HTTP status code from an error
// if it is not an AppError, it returns http.StatusInternalServerError and a generic user message
func UnwrapError(err error, lang language.Language) (int, string) {
	if err == nil {
		return http.StatusOK, texts[lang].errOk
	}

	var appErr AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode(), appErr.UserMessage()
	}

	// if not an AppError, return generic error message
	return http.StatusInternalServerError, texts[lang].errGeneric
}

// transform int slice to string slice
func ItoaSlice(ints []int) []string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return strs
}
