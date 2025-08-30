# Repo structure

- [Repo structure](#repo-structure)
  - [Webapp](#webapp)
  - [Bert](#bert)
  - [Cert](#cert)
  - [Docs](#docs)
  - [ELT](#elt)
  - [Init\_db](#init_db)
  - [Mock\_cas](#mock_cas)
  - [Scripts](#scripts)

## Webapp

We start with the most important part, which is the webapp. It consists of several packages that work together to provide the desired functionality. They are all connected through `main.go` file. Here we provide the overview and API of all the packages.

### `stringx`

This package is really simple and serves as an extension to the standard library's string functions. If you need to define any custom string manipulation functions, this is the place to do it.

Functions:

- `Capitalize(string) string`
  > Capitalizes the first letter of the input string. Supports UTF-8 strings.

### `language`

This package provides utilities for handling language localization in the web application. It defines supported languages, manages language-specific strings, parses language codes from URLs, and integrates language preferences into HTTP request contexts.

Types and methods:

- `Language` 
  > Is a string, with three constants used across the application: `Default` (=`CS`), `CS`, and `EN`. This type is used for multilingual texts in all `texts.go` files. It is also used for identification requested language in the request and for language-specific database calls.
- `(Language) LocalizeURL(string) string`
  > Takes a URL and returns a localized version of it. Adds `cs` or `en` prefix to the url path.
- `LangString`
  > Is used to store language variants of the same string. It is a struct that currently contains two strings, one for the Czech language and one for the English language.
- `(LangString) String(Language) string`
  > Returns the language variant of a string based on the provided language.

Functions:

- `FromString(string) (Language, bool)`
  > Take a string representing a language code and returns the corresponding `Language` value and a boolean indicating success.
- `FromContext(context.Context) Language`
  > Gets the `Language` from the context, which was set there by the middleware `language.SetAndStripLanguageHandler`, or returns `Default` if not found. This function is mainly used in `server.go` files for getting the language from `http.Request.Context()`.
- `MakeLangString(string, string) LangString`
  > Constructor for `LangString` instance with the provided Czech and English strings.
- `SetAndStripLanguageHandler(http.Handler) http.Handler`
  > Middleware that strips the language from the URL path and sets it in the request context. This function is used in `main.go` for handling language-specific routes. Its counterpart is `FromContext(context.Context)`, which is used to retrieve the language from the request context.

### `recommend`

This package provides strategies for recommendations. It is used by home page to recommend courses.

Types and methods:

- `MeiliSearchSimilarToBlueprint` 
  > Recommendation strategy that takes courses from user's blueprint and finds similar courses using MeiliSearch's [hybrid search](https://www.meilisearch.com/docs/reference/api/search#hybrid-search) which is configured to uses bert service for embeddings. It also filters results to only include informatics courses, filters out courses that are already in user's blueprint and picks 10 random courses from top 30 results.
- `(m MeiliSearchSimilarToBlueprint) Recommend(userID string) ([]string, error)`
  > Does the recommendation and returns course codes of recommended courses.
- `NewCourses` 
  > Recommendation strategy that returns courses with newest *valid_from* year. It also filters out courses that are in user's blueprint and courses that are not informatics courses. Lastly it picks 10 random courses from the top 30 courses.
- `(m NewCourses) Recommend(userID string) ([]string, error)`
  > Does the recommendation and returns course codes of recommended courses.

### `cas`

This package provides authentication middleware and utilities for integrating Central Authentication Service (CAS) single sign-on into this application. Its main purpose is to manage user sessions, handle login and logout flows, and securely associate requests with authenticated users. It authenticates user using session key and sets user ID to request context. If session key is not present or authentication fails then it redirects to login page.

Types and methods:

- `UserIDFromContext`
  > A struct that extracts the user ID from the request context. It is used to access the authenticated user in `server.go` files.
- `(UserIDFromContext) UserID(*http.Request) string`
  > Extracts the user ID from the context of the provided request. Used in `server.go` files to get the authenticated user ID.
- `Authentication`
  > A struct that handles the authentication process, including login, logout, and session management. It is used as middleware to protect routes that require authentication.
- `(Authentication) AuthenticateHTTP(http.Handler) http.Handler`
  > Middleware that protects routes requiring authentication. It takes care of login and logout, checking for a valid session and user ID and redirecting to the login page if not authenticated. This method is used in `main.go` for handling routes that require authentication.
- `CAS`
  > A struct that contains link to the Central Authentication Service (CAS) server and provides internal methods for interaction with the server. Documentation can be found at <https://apereo.github.io/cas/>. Can be set to mock CAS server for development or testing purposes.

This package works as follows:
1. In `main.go`, an instance of `cas.Authentication` is created and used as middleware for the router. Parameters are set to configure the authentication behavior. `Data` is initialized SQL database, `Error` is an instance of an error handler (package [`errorx`](#errorx)), `CAS` is an instance of a `CAS` structure with URL to the CAS server. For development purposes, we use [`mock_cas`](#mock_cas), but it is replaced with a real CAS server in production.
2. In `main.go` all servers are configured to use `cas.UserIDFromContext` structure which implements their `Auth` interface.
3. The middleware checks for a valid session and user ID for each incoming request. This is done in `authentication.go` file.
4. If the user is not authenticated, they are redirected to the login page. Logging in using (mock-)CAS and logging out (which does not communicate with the CAS server) is done in `cas.go` file. Database interactions are handled in the `database.go` file.
5.  Once authenticated, servers can call `Auth.UserID(r)` to get the authenticated user ID from the request context.

Logout and login pages (their templated HTML) can be seen in `view.templ` file. This file uses HTML templates to render the HTML for these pages. Login uses `loginModel` structure defined in `model.go`. `texts.go` contains multi-language texts used in the HTML and in error messages.

### `dbds`

This package provides structures that maps to database tables. The mapping is not necessary one to one, it can be used to combine multiple tables into one structure or to split one table into multiple structures. The main idea is to have only one existing mapping of column to structure field so that when any change to the name happens it can be easily resolved. This package should be used whenever relevant data are fetched from a database. After data are fetched from the database into `dbds` structures they should be remapped to local structures to depend on common structures as little as possible so that any change to `dbds` package can be easily resolved.

This package creates mapping for the most common structures, which are structures for course and teacher. It is therefore divided into two files, where each file contains definitions relevant to its specific structure. The files are:
- `course.go`: contains struct `Course` that represents a course in the system. It should be used whenever data related to courses is fetched from the database (as can be seen in pages-related `database.go` files). The file also contains another struct definitions related to course and if you want to learn more about them, please explore the file and its usage yourself, and explore the [Data Model](./data-model.md).
- `teacher.go`: contains struct `Teacher` that represents a teacher in the system. It also defines alias for slice of teachers `TeacherSlice`. The slice is used in `Course` struct, but also in some `database.go` files. For better understanding, please explore the structs usage yourself, and explore the [Data Model](./data-model.md).

Typical usage of this package can be demonstrated on `database.go` file from `courses` package. You can see that the `dbds` package is imported
```go
import (
    "github.com/michalhercik/RecSIS/dbds"
)
```
and the `dbds.Course` struct is used in internal database `courses` struct slice, which represent courses fetched from the database. You can see that `dbds` is used for the most common structures, and you can embed them anywhere you need relevant data from the database.
```go
type courses []struct {
	dbds.Course
	BlueprintSemesters pq.BoolArray `db:"semesters"`
	InDegreePlan       bool         `db:"in_degree_plan"`
}
```
`courses()` method then fetch courses from the database which are automatically mapped to the `courses` struct. However, the Search page (`courses` package) uses its own internal `course` representation which can be seen in `model.go` file (`courses` package) and the fetched courses are then manually mapped to this representation using `intoCourses()` and other functions as seen in the file.

### `errorx`

The errorx package provides a centralized and extensible error handling solution for this application, with a focus on HTTP error reporting, user-friendly messaging, and localization. Its main purpose is to standardize how errors are logged, rendered to users, and propagated through the application, making it easier for developers to maintain consistent error handling across different components.

Types and methods:

- `ErrorHandler`
  > A struct that encapsulates error handling logic, including logging, rendering error messages, and managing fallback scenarios. It is used throughout the application to handle HTTP errors in a consistent manner.
- `(ErrorHandler) Log(error)`
  > Logs the error using the standard `log` package.
- `(ErrorHandler) Render(http.ResponseWriter, *http.Request, int, string, language.Language)`
  > Writes error code to the response writer, set correct `htmx` headers, and renders a small floating window with error code and message. If the error cannot be rendered, it falls back to logging the error and creating a http error.
- `(ErrorHandler) RenderPage(http.ResponseWriter, *http.Request, int, string, string, string, language.Language)`
  > Renders a full error page to the user using the [`page`](#page) package, with the provided error code and message. Title and user ID are necessary for rendering using the page package. If the error page cannot be rendered, it falls back to logging the error and creating a http error.
- `(ErrorHandler) CannotRenderPage(http.ResponseWriter, *http.Request, string, string, error, language.Language)`
  > Handles the scenario where any full page cannot be rendered, logging the error and displaying a cannot-render-page message. This method should be used anywhere, where rendering a full page fails.
- `(ErrorHandler) CannotRenderComponent(http.ResponseWriter, *http.Request, error, language.Language)`
  > Handles the scenario where any non-page component cannot be rendered, logging the error and displaying a cannot-render-component message. This method should be used anywhere, where rendering a non-page component fails.
- `HTTPError`
  > A custom error type that on top of the standard error type includes an HTTP status code and a user-friendly error message. It is used to represent errors that occur during HTTP requests and can be easily rendered by the errorx package.
- `(HTTPError) Error() string`
  > Returns the `error` as a string.
- `(HTTPError) StatusCode() int`
  > Returns the HTTP status code associated with the error.
- `(HTTPError) UserMessage() string`
  > Returns the user-friendly error message.
- `Param`
  > A struct representing a function/method parameter for error context, with a name and value. Is used in `AddContext()` function representing parameters that caused the error.

Functions:

- `NewHTTPError(statusCode int, userMessage string) HTTPError`
  > Constructor for a `HTTPError` with the given error, status code and user message.
- `P(string, any) Param`
  > Constructor for a `Param` with the given name and value.
- `AddContext(error, ...Param) error`
  > Wraps an error with additional context information, such as package, structure and method name and parameters names and values.
- `UnwrapError(error, language.Language) (int, string)`
  > Tries to represent an error as appError (interface with `Error`, `StatusCode`, and `UserMessage` methods), if successful it returns the status code and user message. If not, it returns http.StatusInternalServerError and a generic user message.

This package works as follows:
1. In `main.go`, an `ErrorHandler` instance is created.
2. The instance is injected into an instance of the page server, but as it has not yet been injected itself with an implementation of `Page` interface, it cannot render any pages from [`page`](#page) package.
3. The instance is injected with `Page` interface implementation for rendering error pages.
4. The `ErrorHandler` instance is now injected into all servers for pages, implementing their `Error` interface.

Typical usage of this package can be demonstrated on `server.go` and `database.go` files from `courses` package. You can see the `errorx` package imported in both files.
```go
import (
	"github.com/michalhercik/RecSIS/errorx"
)
```
In the `server.go` file, there is a `Server` struct that expect an implementation of the `Error` interface, which is satisfied by the `ErrorHandler` instance injected into it in `main.go`.

Now, when a new HTTP error occurs, as seen in `database.go` file,
```go
if err := ...; err != nil {
	return nil, errorx.NewHTTPErr(
		errorx.AddContext(
			fmt.Errorf("sqlquery.Courses: %w", err),
			errorx.P("courseCodes", strings.Join(courseCodes, ",")),
			errorx.P("lang", lang),
		),
		http.StatusInternalServerError,
		texts[lang].errCannotLoadCourses,
	)
}
```
it is wrapped with `errorx.NewHTTPErr()` function. This function takes an error, the HTTP status code, which you must provide, and a user-friendly error message, typically stored in `texts.go` file (which must be, for now, correctly localized), and returns a new `HTTPError` instance. The occurred error can be wrapped in context using `errorx.AddContext()` function which takes the original error and any relevant public parameters (we consider `userID` to be a secret). The parameters are created using `errorx.P()` function.

This error then bubbles through the application. When it is rereturned it should be again wrapped in context using `errorx.AddContext()` function, as seen in `search()` method in `server.go` file,
```go
if err != nil {
	return result, errorx.AddContext(err)
}
```
for adding the function/method context, which is however optional. Finally, when the error reaches any rendering method, it must be unwrapped using `errorx.UnwrapError()` function, which will extract the HTTP status code, and user message. If the error was not an `appError` it will be treated as a generic internal server error.

You will then have an error, error code and user error message. Error should be logged using `Error`s log method `s.Error.Log()`, again wrapped in context, and could (but mostly should) be rendered as a floating window using `s.Error.Render()` method or as a full page using `s.Error.RenderPage()` method. You should use the second method when the original rendering method renders full page.

Finally, when rendering any component or page results in error, you should use `s.Error.CannotRenderComponent()` respectively `s.Error.CannotRenderPage()` method to handle the error gracefully.

Creating a `HTTPError` from an error can happen anywhere on the error's journey through the application. You only need to provide HTTP status code and user-friendly error message. You can redefine the error using `errorx.NewHTTPErr()` function again. `errorx.AddContext()` function edits only `error`.

### `filters`

Parse filters for search from user request and convert them into filter query for search engine. To be able to parse user request the package needs to be used also for generating the filtering options so that they use proper IDs. 

Types and methods:

- `Filters`
  > Object that represents a collection of filter categories and their values. It is used to display filtering options with proper references to search engine.
- `(*filters) Init() error`
  > Initializes the filters by fetching the corresponding filter categories and their values from the database.
- `(filters) Facets() []string`
  > Returns list of fields for which facets should be generated. The fields are used in MeiliSearch search request.
- `(filters) ParseURLQuery(url.Values, language.Language) (expression, error)`
  > Parses The URL query parameters to create a filter expression for MeiliSearch. Takes URL values and language as input and returns a filter expression and an error if parsing fails.
- `(Filters) IterFiltersWithFacets(Facets, url.Values, language.Language) iter.Seq[FacetIterator]`
  > Iterates over the filter categories, returning an iterator of `FacetIterator` for each category. Takes facets returned by MeiliSearch, URL values, and language as input.
- `Facets`
  > alias for `map[string]map[string]int`. Represents Category>Value>Count mapping of facets returned by MeiliSearch.
- `FacetValue`
  > Structure representing a single value in a filter category, including its ID, title, description, count of items matching this value, and whether it is currently selected.
- `FacetIterator`
  > Structure for filter category with info about it. For detail its see methods and usage.
- `(FacetIterator) IterWithFacets() iter.Seq2[int, FacetValue]`
  > Returns an iterator over the values in the filter category, yielding index and `FacetValue` for each value.
- `(FacetIterator) Size() int`
  > Returns the number of all possible values in the filter category.
- `(FacetIterator) Count() int`
  > Returns the number of currently non-zero values in the filter category.
- `(FacetIterator) ID() string`
  > Returns the unique identifier of the filter category.
- `(FacetIterator) Title() string`
  > Returns the title of the filter category..
- `(FacetIterator) Desc() string`
  > Returns the description of the filter category.
- `(FacetIterator) DisplayedValueLimit() int`
  > Returns the count of values that should be displayed.
- `(FacetIterator) Active() bool`
  > Returns whether the filter category is active. That means that at least one value is selected.
- `expression`
  > Alias for `[]condition`. Represents a filter expression for MeiliSearch.
- `(*expression) Append(param string, values ...string)`
  > Appends a new condition to the expression based on the provided parameter and values. 
- `(expression) String() string`
  > Converts the expression to a string representation suitable for MeiliSearch filter expressions.
- `(expression) ConditionsCount() int`
  > Returns the number of conditions in the expression.
- `(expression) Except() func(func(string, string) bool)`
  > Return iterator which returns all variants of the expression without one condition. It is used for Meilisearch multi-search request to get disjunctive facets ([see discussion](https://github.com/orgs/meilisearch/discussions/187)) - used in courses package.
- `condition`
  > Represents category condition. It stores category ID and selected values IDs.
- `(condition) String() string`
  > Converts the condition to a string representation suitable for MeiliSearch filter expressions.

Functions:

- `MakeFilters(*sqlx.DB, string) Filters`
  > Constructor for `Filters` that initializes it with the provided database connection and filter ID.
- `SkipEmptyFacet(iter.Seq2[int, FacetValue]) iter.Seq2[int, FacetValue]`
  > If the value has zero count given the search query, it is skipped. Used for displaying survey filters.

<!-- my notes - delete -->
in DB there are filter categories and their values

we currently have filters for courses and surveys

filtering is done using [MeiliSearch](https://www.meilisearch.com/docs/learn/filtering_and_sorting/filter_expression_reference)
<!--  -->

1. Inject filters into a server using `filters.MakeFilters` in `main.go`. As `source`, database with filter categories and their values should be used. `id` should be the `filter_id` of the filter category. That can be found in the database, see the [Data Model](./data-model.md). As we currently have filters for courses and surveys, you can see in `main.go` that we use `courses` and `course-surveys` as `id`.
2. In `server.go`, the injected filters must be initialized using `s.Filters.Init()` method. This will fetch the corresponding filter categories and their values from the database.
3. Then you would make search request to MeiliSearch. In the request, you must specify for which fields you want facets to be generated. You can get the list of fields using `s.Filters.Facets()` method. It is also possible to filter the search by parsing URL query parameters using `s.Filters.ParseURLQuery(r.URL.Query(), lang)` method. It will return a filter expression that can be used in MeiliSearch search request.
4. Part of the MeiliSearch search response are `FacetsDistribution` which is a mapping of Category>Value>Count in other words it is `Facets` type. 
5. You can then display the filters on the page using `s.Filters.IterFiltersWithFacets()` method. It takes `Facets` from MeiliSearch response, URL values (to know which values are selected), and language (for displaying titles and descriptions in the correct language). The method returns an iterator of `FacetIterator` which represents a filter category with its values. You can then use its methods to get information about the category and iterate over its values.

### `bpbtn`

The `bpbtn` package provides reusable components and logic for adding courses to a user's blueprint (own study plan) in the application. Its main purpose is to encapsulate the UI and backend logic for the *add to blueprint* button, including request parsing, validation, error handling, and database operations. This package is designed to be injected into servers (such as courses or degree plan) so that the add button can be rendered and its actions handled consistently across different parts of the app.

Types and methods:

- `AddWithTwoTemplComponents`
  > This type extends the basic `Add` functionality by providing a possibility to render a second component alongside the main add button.
- `(AddWithTwoTemplComponents) PartialComponentSecond(language.Language) func(string, string, string, []bool, string) templ.Component`
  > Creates a partial component for the second template. Partial meaning, that you need to provide HTMX attributes `hx-swap`, `hx-target`, `hx-include`, semester flags denoting the available semesters for a course in blueprint, and a course code.
- `Add`
  > This type provides the basic functionality for the *add to blueprint* button, including request parsing, validation, and database operations. It stores the database where the blueprint data is kept, template for rendering the button, and HTMX configuration.
- `(Add) Endpoint() string`
  > Returns the HTTP method and endpoint which is called when the add button is clicked. Servers using this package should register this endpoint to handle the add requests.
- `(Add) PartialComponent(language.Language) func(string, string, string, []bool, string) templ.Component`
  > Creates a partial component for the template. Partial meaning, that you need to provide HTMX attributes `hx-swap`, `hx-target`, `hx-include`, semester flags denoting the available semesters for a course in blueprint, and a course code.
- `(Add) ParseRequest(*http.Request, []string) ([]string, int, int, error)`
  > Parses the HTTP request to extract course code, year, and semester information. Takes to HTTP request to parse and extra course codes to add to the blueprint.
- `(Add) Action(string, int, int, language.Language, ...string) ([]int, error)`
  > Adds the specified course(s) to the user's blueprint for the given year and semester. Takes the user's ID, year and semester number, language and course(s) to add. Makes a simple database call.
- `ViewModel`
  > Represents the data model for the add button component, including course information, semester flags, and HTMX attributes.

This package works as follows:
1. In `main.go`, any server that wants to add course(s) to blueprint is injected with an instance of `bpbtn.Add` or `bpbtn.AddWithTwoTemplComponents`. The instance is initialized with the database connection, template(s) for rendering the button (which are defined in this package), and the base for HTMX POST request.
2. Servers register a handler for the add button's HTTP endpoint using `s.BpBtn.Endpoint()`.
3. When a user clicks the add button, the server can implement it own logic for handling such request but it should use the built-in methods from the `bpbtn` package, which take care of parsing the request, making database calls, and creating a template component for the response.
4. For rendering the add button, the server should call `s.BpBtn.PartialComponent()` method, which result should be passed to the model for the page.
5. The partial component can be then rendered in the page template. The template requires HTMX attributes `hx-swap`, `hx-target`, `hx-include`, semester flags denoting the available semesters for the course in blueprint, and the course code itself.

The usage of this package can be demonstrated on the `courses` package. In `main.go`, the courses server is injected with an instance of `bpbtn.Add`.
```go
courses := courses.Server{
	BpBtn: bpbtn.Add{
		DB:    db,
		Templ: bpbtn.AddBtn,
		Options: bpbtn.Options{
			HxPostBase: coursesRoot,
		},
	},
  ...
}
```
You can see that the template for the add button was chosen as `bpbtn.AddBtn`. Also, the base for HTMX POST request is `/courses/` (`coursesRoot` variable). The courses server then registers a handler for the add button's endpoint.
```go
router.HandleFunc(s.BpBtn.Endpoint(), s.addCourseToBlueprint)
```
with the following handler function:
```go
func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	courseCodes, year, semester, err := s.BpBtn.ParseRequest(r, nil)
	if err != nil {
		...
	}
	if len(courseCodes) != 1 {
		...
	}
	courseCode := courseCodes[0]
	_, err = s.BpBtn.Action(userID, year, semester, lang, courseCode)
	if err != nil {
		...
	}
	course, err := s.Data.courses(userID, []string{courseCode}, lang)
	if err != nil {
		...
	}
	btn := s.BpBtn.PartialComponent(lang)
	err = CourseCard(&course[0], t, btn).Render(r.Context(), w)
	if err != nil {
		...
	}
}
```
The `s.BpBtn.ParseRequest()` method parses the request to extract course code, year and semester. There are no extra courses to add (`nil` parameter for the method) as on the courses page, each button corresponds to a single course. Then the `s.BpBtn.Action()` method is called to add the course to the user's blueprint. Finally, the partial component for the button is created using the `s.BpBtn.PartialComponent()` method which is then rendered in the response. When the partial component is passed to the model, it is used to render the button in the appropriate place within the page template.
```go
templ CourseCard(..., addBtn PartialBlueprintAdd) {
  <div>
    ...
    <div>
      ...
      // add to blueprint button
      <div class="d-flex justify-content-end">
        @addBtn("outerHTML", "#course-card-" + course.code, "", course.blueprintSemesters, course.code)
      </div>
  </div>
}
```
Here, the partial component is given the necessary context for the button to work correctly. We can see that the `hx-target` here is the course card element. That means that the handler for the add button's endpoint should render this component.

### Pages

Now for the packages that are actually *seen*. That means that their are rendered on the FE. We start with the `page` package.

#### `page`

This package provides the page layout and structure for the application. Specifically, it defines the page header with navigation bar and a footer. All the contents of our pages are inserted into this layout which creates the final look.

The `page` API is not interesting, and so, the package will be documented through its usage. If you would like to create a new page using this package, please refer to [Add new page](./how-to-extend.md#add-new-page) part.

In `main.go`, a `Page` instance is created and configured with the necessary parameters.
```go
pageTempl := page.Page{
	Error: errorHandler,
	Home:  homeRoot,
	NavItems: []page.NavItem{
		{Title: language.MakeLangString("Domů", "Home"), Path: homeRoot, Skeleton: home.Skeleton, Indicator: "#home-skeleton"},
		{Title: language.MakeLangString("Hledání", "Search"), Path: coursesRoot, Skeleton: courses.Skeleton, Indicator: "#courses-skeleton"},
		{Title: language.MakeLangString("Blueprint", "Blueprint"), Path: blueprintRoot, Skeleton: blueprint.Skeleton, Indicator: "#blueprint-skeleton"},
		{Title: language.MakeLangString("Studijní plán", "Degree plan"), Path: degreePlanRoot, Skeleton: degreeplan.Skeleton, Indicator: "#degreeplan-skeleton"},
	},
	Search: page.MeiliSearch{
		Client: meiliClient,
		Index:  "courses",
		Limit:  5,
	},
	Param:          "search",
	SearchEndpoint: coursesRoot,
	ResultsDetailEndpoint: func(code string) string {
		return courseDetailRoot + code
	},
}
pageTempl.Init()
```
`Error` expects an error handler which implements `page`'s `Error` interface. We use our `errorx` package to provide this functionality. `Home` is path to home page - `/`. `NavItems` are links to parts of the application which will be seen in the navigation bar. `Search` is used for quick searching of courses using `MeiliSearch` in search bar which is also in the navigation bar. Other parameters are also used for quick search. Next the `Page` instance is initialized.

The result is than injected into every page server. Servers either use the `Page` instance directly or wrap it in `PageWithNoFiltersAndForgetsSearchQueryOnRefresh` which does exactly what its name suggests. Both `Page` and `PageWithNoFiltersAndForgetsSearchQueryOnRefresh` have `View()` method, which is responsible for rendering the page and its content. Difference between the two methods can be seen in `page.go` file directly. They both use the same private method which created a page model from provided parameters and returns template for the page which can be rendered. The most important parameter is template `templ.Component` for the content of the page.

The `View()` method is used in `server.go` files for rendering pages. For more information, please refer to the `server.go` file. The server always provide content for the page.

#### Specific pages

In our application, we currently have five specific pages:

1. Home page
2. Blueprint page
3. Courses page
4. Course detail page - there is a course detail page for each course, but they all share the same template
5. Degree plan page

Their structure is similar, but each page has its own specific content and functionality. An overview of how to pages work can be seen in the following diagram:

![Pages structure diagram](pages.svg)

Every page has a server structure `Server` located in `server.go` file. These structures consists of interfaces which are implemented by the packages that are described above through dependency injection and some internal structures that are also injected with dependencies. The structure and interfaces of each server can be seen in the `server.go` files.

Each `server.go` file also contains router, specific for each page. Every path has defined its own handler.

A page typically works by this flow:

1. The user makes a request to a specific path.
2. The router in the corresponding `server.go` file receives the request and calls the appropriate handler.
3. The handler prepares the data needed for the page, typically from database, using `database.go` file located in the package.
4. The data are then parsed into a model of the page or some of its component, which can be seen in every package `model.go` file.
5. The page, or some component of it, is rendered using the specified template with the model data. Templates are located in the `view.templ` files.

SQL queries needed for database called are usually stored in `internal/sqlquery` package for each page. Bi-lingual texts used on the pages are stored in `texts.go` files.

Some packages have some extra files, specific for their functionality:
- `sanitizer.go` - sanitize and transform texts seen on the course detail page. For more information, please refer to the file itself.
- `search.go` - implements search functionality using Meilisearch client. For more information, please refer to the file itself or [Meilisearch API documentation](https://www.meilisearch.com/docs/reference/api).

If you want to learn how to pages and servers work in greater detail, please read the [Add new page](./how-to-extend.md#add-new-page) part.

### Connecting the packages

How the packages are connected can be seen from imports in each file. Typically, they all go through `main.go` which is the most important file for understanding the flow of the application.

Apart from packages, the `webapp` directory contains:
- `static` folder with static assets: icons, CSS, JavaScript and user guide (in czech language)
- config file which is loaded in `main.go` for application configuration
- `Dockerfile` for building the application container
- `go.mod` and `go.sum` files for Go module management
- `main.go` file
- `main_test.go` file containing integration tests, for more see [Testing](./testing.md)

## Bert

Simple service that provides BERT embeddings for given texts. It is used by
MeiliSearch to embed courses. The embeddings are then used for simple
recommendations. The service implements single endpoint:
- `POST /embedding`
The endpoint expects JSON body with a single field `text` which is a text to be
embedded. Response then contains single field `embedding` which is an array of
float32 numbers representing the embedding.

## Cert

`cert` directory contains `server.crt` and `server.key` files used self-signed SSL certificate for the application.

## Docs

`docs` directory contains developer documentation for the application.

## ELT

We would also like to give you an high level overview of how the ELT process
works.  We decided to implement it using Go and SQL. The entire process is
simple and we didn't feel the need to use any sophisticated tools. Most of the
logic is implemented in SQL. Go serves mainly as orchestrator of the process.
Therefore the source code or at least the main file serves as a high level
overview of the process. We decided to KISS (Keep It Simple, Stupid) and
therefore running ELT deletes all the data in the database and repopulates it
from scratch. The entire ELT process takes few minutes and doing anything more
sophisticated would be overkill at this time.

As the name suggest ELT consists of Extract, Load and Transform steps - we also
added fourth step which is migration. The first two steps (extract and load) are
pretty straightforward and each table is extracted in parallel into local
database. The only caveat was related to bulk insert and you can read more about
it in [this
article](https://klotzandrew.com/blog/postgres-passing-65535-parameter-limit/).
It also worth noting that before loading course descriptions into database we
had to remove null bytes as the PostgreSQL doesn't support it. The extract
and load process for each table is defined in structure (one structure for each
source table) implementing `operation` interface (see below). The last tricky
extraction was related to degree plans. We don't have access to list of degree
plans (see [Data Model](./data-model.md)) with appropriate years and to fetch degree
plan from SIS database we need degree plan code and year. We did a dirty
workaround by taking studies in ten year window, drop duplicates and for each
degree plan code extracted variant for every year. We don't know how good or bad
the solution is because we don't know much about the degree plans but we believe
this solution works.

```go
type operation interface {
    name() string // for logging purposes
    selectData(from *sqlx.DB, to *sqlx.DB) error // to parameter is there from historical reasons and can be used for filtering purposes
    insertData(to *sqlx.DB) error
}
```

The third step (transform) can be a bit harder to follow. Some transformations
can run in parallel but not all since some of them depend on the previous ones.
Each transformation takes table and produces another table. Each transformation
is defined as instance of `transformation` structure (see below).

```go
type transformation struct {
    name  string // for logging purposes
    query string
}
```

After the transformation the forth step (migration) is run. The migration takes
transformed data and migrates it into the final tables in different schema and
into search engine. It is possible that the ELT process will result in
inconsistent state as failure of between tables migration does not rollback
migration of data into search engine. This should be addressed in the future.

## Init_db

Scripts that are used for initializing the local database. They are run when the
postgres container is started for the first time. The scripts are run in
lexicographical order.

## Mock_cas

Very simple service that mimics the Central Authentication Service (CAS) server.
It naively implements three endpoints:
- `POST /cas/login`
- `GET /cas/login`
- `GET /cas/serviceValidate`

To get better understanding of how CAS works please refer to [CAS protocol
documentation](https://apereo.github.io/cas/).

## Scripts

This directory contains various scripts for managing the application. Their usage is explained in [How to run RecSIS](./how-to-run.md).
