package page

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type SearchBar interface {
	View(string, language.Language, bool) templ.Component
	QuickSearchResult(query string, lang language.Language) (templ.Component, error)
	QuickSearchEndpoint() string
	SearchParam() string
}

type PageWithNoFiltersAndForgetsSearchQueryOnRefresh struct {
	Page
}

func (p PageWithNoFiltersAndForgetsSearchQueryOnRefresh) View(main templ.Component, lang language.Language, title string, userID string) templ.Component {
	return p.view(main, lang, title, "", false, userID)
}

type Page struct {
	Error           Error
	Home            string
	NavItems        []NavItem
	QuickSearchPath string
	router          *http.ServeMux
	SearchBar       SearchBar
}

func (p Page) SearchParam() string {
	return p.SearchBar.SearchParam()
}

func (p *Page) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET "+p.QuickSearchPath, p.quickSearch)
	p.router = router
}

func (p Page) Router() *http.ServeMux {
	return p.router

}

func (p Page) View(main templ.Component, lang language.Language, title string, searchQuery string, userID string) templ.Component {
	return p.view(main, lang, title, searchQuery, true, userID)
}

func (p Page) view(main templ.Component, lang language.Language, title string, searchQuery string, includeFilters bool, userID string) templ.Component {
	searchBarView := p.SearchBar.View(searchQuery, lang, includeFilters)
	model := pageModel{
		title:    title,
		main:     main,
		lang:     lang,
		text:     texts[lang],
		search:   searchBarView,
		home:     p.Home,
		navItems: p.NavItems,
		userID:   userID,
	}
	return PageView(model)
}

func (p Page) quickSearch(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	query := r.FormValue(p.SearchBar.SearchParam())
	view, err := p.SearchBar.QuickSearchResult(query, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		p.Error.Log(errorx.AddContext(err))
		p.Error.Render(w, r, code, userMsg, lang)
		return
	}
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
	title    string
	main     templ.Component
	lang     language.Language
	text     text
	search   templ.Component
	home     string
	navItems []NavItem
	userID   string
}

type NavItem struct {
	Title     language.LangString
	Path      string
	Skeleton  func(language.Language) templ.Component
	Indicator string
}
