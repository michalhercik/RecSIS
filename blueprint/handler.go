package blueprint

import (
	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/mock_data"
	"net/http"
	"strings"
	"strconv"
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

func HandleLastYearRemoval(w http.ResponseWriter, r *http.Request) {
	year := r.PathValue("year")

	//remove data from DB
	yearInt, _ := strconv.Atoi(year)
	mock_data.RemoveYear(yearInt)

	// remove data from UI
	var sb strings.Builder
	sb.WriteString(`<tr id="`)
	sb.WriteString("Year")
	sb.WriteString(year)
	sb.WriteString(`" hx-swap-oob="delete"></tr>`)

	sb.WriteString(`<div id="`)
	sb.WriteString("Year")
	sb.WriteString(year)
	sb.WriteString(`" hx-swap-oob="delete"></div>`)

	sb.WriteString(`<button id="RemoveYearButton" hw-swap-oob="true" hx-post="/blueprint/remove-year/`)
	sb.WriteString(strconv.Itoa(yearInt - 1))
	sb.WriteString(`" hx-target="#Year`)
	sb.WriteString(strconv.Itoa(yearInt - 1))
	sb.WriteString(`"> Remove last year </button>`)


	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(sb.String()))
}

func HandleBLueprintUnassignedRemoval(w http.ResponseWriter, r *http.Request) {
	//remove data from DB
	idInt, _ := strconv.Atoi(r.PathValue("id"))
	mock_data.RemoveFromBlueprint(idInt)
	w.WriteHeader(http.StatusOK)
}
