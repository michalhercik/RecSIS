package blueprint

import (
	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/mock_data"
	"net/http"
	"strconv"
	"fmt"
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

    // Remove data from DB
    yearInt, _ := strconv.Atoi(year)
    mock_data.RemoveYear(yearInt)

    // Generate updated button HTML
    buttonHTML := fmt.Sprintf(`<button
        hx-swap="outerHTML"
        hx-post="/blueprint/remove-year/%d"
        hx-target="this">
        Remove last year
    </button>`, yearInt - 1)

    // Generate removal HTML for the div and tr elements
    divRemovalHTML := fmt.Sprintf(`<div id="BlueprintYear%d" hx-swap-oob="delete"></div>`, yearInt)
    trWinterRemovalHTML := fmt.Sprintf(`<template><tr id="SumWinterYear%d" hx-swap-oob="delete"></tr></template>`, yearInt)
    trSummerRemovalHTML := fmt.Sprintf(`<template><tr id="SumSummerYear%d" hx-swap-oob="delete"></tr></template>`, yearInt)

    // Write updated button HTML and removal HTML to response
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(buttonHTML + divRemovalHTML + trWinterRemovalHTML + trSummerRemovalHTML))
}

func HandleBLueprintUnassignedRemoval(w http.ResponseWriter, r *http.Request) {
	//remove data from DB
	idInt, _ := strconv.Atoi(r.PathValue("id"))
	mock_data.RemoveFromBlueprint(idInt)
	w.WriteHeader(http.StatusOK)
}
