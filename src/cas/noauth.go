package cas

import (
	"context"
	"net/http"
)

func NoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), userIDKey{}, "81411247"))
		next.ServeHTTP(w, r)
	})
}
