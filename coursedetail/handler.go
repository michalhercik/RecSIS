package coursedetail

import (
	"github.com/a-h/templ"
	"net/http"
	"strconv"
)

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	courseId, _ := strconv.Atoi(r.PathValue("id"))
	// TODO handle error
	data := db.GetData(courseId)
	return Content(&data)
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	courseId, _ := strconv.Atoi(r.PathValue("id"))
	// TODO handle error
	data := db.GetData(courseId)
	return Page(&data)
}