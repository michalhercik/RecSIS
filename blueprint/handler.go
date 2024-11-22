package blueprint

import (
	"github.com/a-h/templ"
    "github.com/michalhercik/RecSIS/database"
	"net/http"
	"strconv"
)

const user = 42

func HandleContent() templ.Component {
	data := database.GetBlueprintData(user)
	return Content(&data)
}

func HandlePage() templ.Component {
	data := database.GetBlueprintData(user)
	return Page(&data)
}

// func HandleLastYearRemoval(w http.ResponseWriter, r *http.Request) {
//     year := r.PathValue("year")

//     // Remove data from DB
//     yearInt, _ := strconv.Atoi(year)
//     database.RemoveYear(yearInt)

//     // Send refresh header
//     w.Header().Set("HX-Refresh", "true") // htmx will trigger a full page reload
// 	w.WriteHeader(http.StatusOK)
// }

func HandleBLueprintUnassignedRemoval(w http.ResponseWriter, r *http.Request) {
	// Remove data from DB
	courseID, _ := strconv.Atoi(r.PathValue("id"))
	database.RemoveFromBlueprint(user, courseID)
    // Send http response
	w.WriteHeader(http.StatusOK)
}
