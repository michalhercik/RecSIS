package language

import (
	"context"
	"net/http"
)

type Language string

const langLen int = 2
const (
	Default Language = CS
	CS      Language = "cs"
	EN      Language = "en"
)

func FromString(lang string) (Language, bool) {
	switch lang {
	case string(CS):
		return CS, true
	case string(EN):
		return EN, true
	default:
		return "", false
	}
}

func FromContext(ctx context.Context) Language {
	lang, ok := ctx.Value(key{}).(Language)
	if ok {
		return lang
	}
	return Default
}

func SetAndStripLanguage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalReq := r
		originalPath := r.URL.Path
		lang, ok := parseLanguage(r)
		if ok {
			r = requestWithLanguage(r, lang)
		}
		next.ServeHTTP(w, r)
		r = originalReq
		r.URL.Path = originalPath
	})
}

func parseLanguage(r *http.Request) (Language, bool) {
	if len(r.URL.Path) >= langLen+1 {
		prefix := r.URL.Path[1 : langLen+1]
		lang, ok := FromString(prefix)
		return lang, ok
	}
	return Default, false
}

func requestWithLanguage(r *http.Request, lang Language) *http.Request {
	r.URL.Path = r.URL.Path[langLen+1:]
	if len(r.URL.Path) == 0 {
		r.URL.Path = "/"
	}
	ctx := context.WithValue(r.Context(), key{}, lang)
	r = r.WithContext(ctx)
	return r
}

type key struct{}
