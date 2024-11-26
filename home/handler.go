package home

import (
	"github.com/a-h/templ"
	"net/http"
)

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	return Content()
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	return Page()
}
