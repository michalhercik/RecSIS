package coursedetail

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	Data DataManager
}

func (s Server) Register(router *http.ServeMux, prefix string) {
	//router.HandleFunc(fmt.Sprintf("GET %s/{code}", prefix), s.page) // TODO get language from http header
	router.HandleFunc(fmt.Sprintf("GET /cs%s/{code}", prefix), s.csPage)
	router.HandleFunc(fmt.Sprintf("GET /en%s/{code}", prefix), s.enPage)
	// TODO: should we differentiate between languages for POSTs?
	router.HandleFunc(fmt.Sprintf("PUT %s/rating/{code}", prefix), s.rate)
	router.HandleFunc(fmt.Sprintf("DELETE %s/rating/{code}", prefix), s.deleteRating)
	router.HandleFunc(fmt.Sprintf("PUT %s/rating/{code}/{category}", prefix), s.rateCategory)
	router.HandleFunc(fmt.Sprintf("DELETE %s/rating/{code}/{category}", prefix), s.deleteCategoryRating)
}

func (s Server) csPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, texts["cs"], cs)
}

func (s Server) enPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, texts["en"], en)
}

func (s Server) page(w http.ResponseWriter, r *http.Request, t text, lang DBLang) {
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("courseDetail error: %v", err)
		return
	}
	code := r.PathValue("code")
	course, err := s.Data.Course(sessionCookie.Value, code, lang)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
		PageNotFound(code, t).Render(r.Context(), w)
	} else {
		Page(course, t).Render(r.Context(), w)
	}
}

func (s Server) rate(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("rate error: %v", err)
		return
	}
	code := r.PathValue("code")
	rating, err := strconv.Atoi(r.FormValue("value"))
	if err != nil && (rating == 1 || rating == 0) {
		log.Printf("rate error: %v", err)
		return
	}
	if err = s.Data.Rate(sessionCookie.Value, code, rating); err != nil {
		log.Printf("rate error: %v", err)
	}

}

func (s Server) deleteRating(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("deleteRating error: %v", err)
		return
	}
	code := r.PathValue("code")
	if err = s.Data.DeleteRating(sessionCookie.Value, code); err != nil {
		log.Printf("deleteRating error: %v", err)
	}
}

func (s Server) rateCategory(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("rateCategory error: %v", err)
		return
	}
	code := r.PathValue("code")
	category := r.PathValue("category")
	rating, err := strconv.Atoi(r.FormValue("value"))
	if err != nil {
		log.Printf("rateCategory error: %v", err)
		return
	}
	if err = s.Data.RateCategory(sessionCookie.Value, code, category, rating); err != nil {
		log.Printf("rateCategory error: %v", err)
	}
}

func (s Server) deleteCategoryRating(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("deleteCategoryRating error: %v", err)
		return
	}
	code := r.PathValue("code")
	category := r.PathValue("category")
	if err = s.Data.DeleteCategoryRating(sessionCookie.Value, code, category); err != nil {
		log.Printf("deleteCategoryRating error: %v", err)
	}
}
