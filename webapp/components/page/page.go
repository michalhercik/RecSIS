package page

/** PACKAGE DESCRIPTION

The page package provides a reusable framework for rendering consistent page layouts across the application, including navigation, search, error handling, and main content injection. Its main purpose is to standardize how pages are constructed and displayed, ensuring a uniform user experience and simplifying the integration of features like quick search, navigation bar, and error reporting.

Typical usage involves creating a Page instance (as in main.go) and configuring it with navigation items, search parameters, endpoints, and error handling. Typically, you would need only one such instance in the application. The View method is used to render a page by injecting the main content component, language, title, potentially search query, and user ID. This method automatically adds headers, navigation, and optionally search/filter controls. For pages that should not display filters or should forget the search query on refresh, the PageWithNoFiltersAndForgetsSearchQueryOnRefresh struct can be used.

The package also provides routing for quick search functionality (quickSearch handler), which allows users to search for courses directly from the navigation bar. Error handling is integrated via the Error interface, ensuring that any rendering or data-fetching issues are reported to the user in a consistent, localized manner. Developers should use the Page type to wrap their main content and rely on its methods to handle layout, navigation, and search, making it easy to maintain and extend the application's UI.

*/

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
	courses, err := p.Search.quickSearchResult(query, lang)
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
	// Logs the provided error.
	Log(err error)

	// Renders an error message to the user as a floating window, with a status code and localized message.
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)

	// Renders a floating window with error when any component cannot be rendered due to an error.
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
