package page

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type PageWithNoFiltersAndForgetsSearchQueryOnRefresh struct {
	Page
}

func (p PageWithNoFiltersAndForgetsSearchQueryOnRefresh) View(main templ.Component, lang language.Language, title string, userID string) templ.Component {
	return p.view(main, lang, title, "", false, userID)
}

type Page struct {
	Error                 Error
	Home                  string
	NavItems              []NavItem
	Search                MeiliSearch
	Param                 string
	SearchEndpoint        string
	ResultsDetailEndpoint func(code string) string
	router                *http.ServeMux
}

func (p Page) SearchParam() string {
	return p.Param
}

func (p *Page) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /quicksearch", p.quickSearch)
	p.router = router
}

func (p Page) Router() *http.ServeMux {
	return p.router

}

func (p Page) View(main templ.Component, lang language.Language, title string, searchQuery string, userID string) templ.Component {
	return p.view(main, lang, title, searchQuery, true, userID)
}

func (p Page) view(main templ.Component, lang language.Language, title string, searchInput string, includeFilters bool, userID string) templ.Component {
	model := pageModel{
		title:          title,
		main:           main,
		lang:           lang,
		text:           texts[lang],
		home:           p.Home,
		navItems:       p.NavItems,
		userID:         userID,
		searchInput:    searchInput,
		includeFilters: includeFilters,
		searchParam:    p.Param,
		searchEndpoint: p.SearchEndpoint,
	}
	return PageView(model)
}

func (p Page) quickSearch(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	query := r.FormValue(p.Param)
	courses, err := p.Search.QuickSearchResult(query, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		p.Error.Log(errorx.AddContext(err))
		p.Error.Render(w, r, code, userMsg, lang)
		return
	}
	model := quickResultsModel{
		t:                    t,
		lang:                 lang,
		courses:              courses,
		resultDetailEndpoint: p.ResultsDetailEndpoint,
	}
	view := QuickResults(model)
	err = view.Render(r.Context(), w)
	if err != nil {
		p.Error.CannotRenderComponent(w, r, err, lang)
	}
}

type Error interface {
	Log(err error)
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)
	CannotRenderComponent(w http.ResponseWriter, r *http.Request, err error, lang language.Language)
}

type pageModel struct {
	title          string
	main           templ.Component
	lang           language.Language
	text           text
	home           string
	navItems       []NavItem
	userID         string
	includeFilters bool
	searchParam    string
	searchInput    string
	searchEndpoint string
}

type quickResultsModel struct {
	t                    text
	lang                 language.Language
	courses              []quickCourse
	resultDetailEndpoint func(code string) string
}

type NavItem struct {
	Title     language.LangString
	Path      string
	Skeleton  func(language.Language) templ.Component
	Indicator string
}
