package courses

import (
	"github.com/a-h/templ"
    "github.com/michalhercik/RecSIS/database"
)

func HandleContent() templ.Component {
	data := database.GetCoursesData()
	return Content(&data)
}

func HandlePage() templ.Component {
	data := database.GetCoursesData()
	return Page(&data)
}
