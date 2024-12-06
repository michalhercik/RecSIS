package coursedetail

import (
	"net/http"

	"github.com/a-h/templ"
)

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	code := r.PathValue("code")
	course, err := db.Course(code)
	if err != nil {
		return ContentNotFound(code)
	}
	return Content(course)
}

func HandlePage(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	course, err := db.Course(code)
	if err != nil {
		PageNotFound(code).Render(r.Context(), w)
	} else {
		Page(course).Render(r.Context(), w)
	}
}

func HandleCommentAddition(w http.ResponseWriter, r *http.Request) {
	// get the course code and comment content from the request
	code := r.PathValue("code")
	commentContent := r.FormValue("comment")

	// sanitize the comment
	// TODO maybe not change the content, but return an error if the content is not valid
	sanitizedComment := sanitize(commentContent)

	// add the comment to the database
	err := db.AddComment(code, sanitizedComment)
	if err != nil {
		http.Error(w, "Unable to add comment", http.StatusInternalServerError)
		return
	}

	// return the page with the updated comments
	newComments, err := db.GetComments(code)
	if err != nil {
		http.Error(w, "Unable to retrieve comments", http.StatusInternalServerError)
		return
	}
	// TODO please check if this is the correct way to render the page
	Comments(newComments, code).Render(r.Context(), w)
}

func HandleLike(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	code := r.PathValue("code")

	// TODO implement the like functionality

	// TODO please check if this is the correct way to render the page
	Ratings([]Rating{}, code).Render(r.Context(), w)
}

func HandleDislike(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	code := r.PathValue("code")

	// TODO implement the dislike functionality

	// TODO please check if this is the correct way to render the page
	Ratings([]Rating{}, code).Render(r.Context(), w)
}

func HandleBlueprintAddition(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from URL
	code := r.PathValue("code")

	// Make data changes
	assignments, err := db.AddCourseToBlueprint(user, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Render the result
	BlueprintAssignment(assignments, code).Render(r.Context(), w)
}

func HandleBlueprintRemoval(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from URL
	code := r.PathValue("code")

	// Make data changes
	err := db.RemoveCourseFromBlueprint(user, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Render the result
	BlueprintAssignment([]Assignment{}, code).Render(r.Context(), w)
}
