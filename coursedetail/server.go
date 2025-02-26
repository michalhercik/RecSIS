package coursedetail

import (
	"fmt"
	"log"
	"net/http"
	"path"
)

type Server struct {
	Data DataManager
}

func (s Server) Register(router *http.ServeMux, prefix string) {
	router.HandleFunc(fmt.Sprintf("GET %s/{code}", prefix), s.page)
	router.HandleFunc(fmt.Sprintf("GET %s/{code}/en", prefix), s.page)
	router.HandleFunc(fmt.Sprintf("POST %s/{code}/comment", prefix), s.commentAddition)
	router.HandleFunc(fmt.Sprintf("POST %s/{code}/like", prefix), s.like)
	router.HandleFunc(fmt.Sprintf("POST %s/{code}/dislike", prefix), s.dislike)
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	lang := DBLang(path.Base(r.URL.Path))
	if lang != "en" {
		lang = cs
	}
	course, err := s.Data.Course(code, lang)
	log.Printf("Course: %v", course)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
		PageNotFound(code).Render(r.Context(), w)
	} else {
		Page(course).Render(r.Context(), w)
	}
}

func (s Server) commentAddition(w http.ResponseWriter, r *http.Request) {
	// get the course code and comment content from the request
	code := r.PathValue("code")
	commentContent := r.FormValue("comment")

	// sanitize the comment
	// TODO maybe not change the content, but return an error if the content is not valid
	sanitizedComment := sanitize(commentContent)

	// add the comment to the database
	err := s.Data.AddComment(code, sanitizedComment)
	if err != nil {
		http.Error(w, "Unable to add comment", http.StatusInternalServerError)
		return
	}

	// return the page with the updated comments
	newComments, err := s.Data.GetComments(code)
	if err != nil {
		http.Error(w, "Unable to retrieve comments", http.StatusInternalServerError)
		return
	}
	// TODO please check if this is the correct way to render the page
	Comments(newComments, code).Render(r.Context(), w)
}

func (s Server) like(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	code := r.PathValue("code")

	// TODO implement the like functionality

	// TODO please check if this is the correct way to render the page
	Ratings([]Rating{}, code).Render(r.Context(), w)
}

func (s Server) dislike(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	code := r.PathValue("code")

	// TODO implement the dislike functionality

	// TODO please check if this is the correct way to render the page
	Ratings([]Rating{}, code).Render(r.Context(), w)
}
