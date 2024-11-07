package blueprint

import (
	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/mock_data"
)

func HandleContent() templ.Component {
	blueprintCourses := mock_data.GetBlueprintCourses()
	years := mock_data.GetCoursesByYears()
	return Content(&blueprintCourses, &years)
}

func HandlePage() templ.Component {
	blueprintCourses := mock_data.GetBlueprintCourses()
	years := mock_data.GetCoursesByYears()
	return Page(&blueprintCourses, &years)
}
