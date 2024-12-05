package degreeplan

import (
	"net/http"
)

func HandlePage(w http.ResponseWriter, r *http.Request) {
	Page().Render(r.Context(), w)
}
