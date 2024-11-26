package coursedetail

import (
	"github.com/a-h/templ"
    "github.com/michalhercik/RecSIS/database"
	"net/http"
	"strconv"
)

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	courseId, er := strconv.Atoi(r.PathValue("id"))
	if er != nil {
		panic(er)
	}
	data := database.GetCourseData(courseId)
	return Content(&data)
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	courseId, er := strconv.Atoi(r.PathValue("id"))
	if er != nil {
		panic(er)
	}
	data := database.GetCourseData(courseId)
	return Page(&data)
}