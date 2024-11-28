package courses

import (
	"github.com/a-h/templ"
	"net/http"
)

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	data := db.GetData()
	return Content(&data)
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	data := db.GetData()
	return Page(&data)
}
