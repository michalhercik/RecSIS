package courses

import (
	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/mock_data"
)

func HandleContent() templ.Component {
	data := GetListOfCourses()
	return Content(&data)
}

func HandlePage() templ.Component {
	data := GetListOfCourses()
	return Page(&data)
}
