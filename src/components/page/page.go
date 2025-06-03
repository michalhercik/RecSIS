package page

import (
	"net/http"

	"github.com/a-h/templ"
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
	Home            string
	NavItems        []NavItem
	QuickSearchPath string
	SearchBar       SearchBar
	router          *http.ServeMux
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.Render(r.Context(), w)
}

type pageModel struct {
	title    string
	main     templ.Component
	lang     language.Language
	text     Text
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
