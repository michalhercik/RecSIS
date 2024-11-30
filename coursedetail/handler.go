package coursedetail

import (
	"net/http"

	"github.com/a-h/templ"
)

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	code := r.PathValue("code")
	course, err := db.Course(code)
	if err != nil {
		return ContentNotFound(code)
	}
	return Content(course)
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	code := r.PathValue("code")
	course, err := db.Course(code)
	if err != nil {
		return PageNotFound(code)
	}
	return Page(course)
}
