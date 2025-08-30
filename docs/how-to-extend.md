# How to extend this application

Before diving into implementing new features you should get even deeper
understanding by reading at least some of the package implementations. We
suggest you to start in the main file and then continue in any of the packages
implementing page handlers (e.g.  coursedetail, courses, ...).

- [How to extend this application](#how-to-extend-this-application)
  - [Add new page](#add-new-page)
  - [Add a filter](#add-a-filter)
  - [Add a recommender](#add-a-recommender)
  - [Add error configuration](#add-error-configuration)
  - [Use `LangString`s](#use-langstrings)

## Add new page

We will demonstrate this use-case by creating a new page with teachers.

First you should create a new directory for the teachers page. This directory should be placed in the `webapp` directory and should follow the naming conventions used in the existing pages. Let's call it `teachers`.

Inside it you should create the following files:
- `view.templ` - this file will contain the HTML template for the teachers page.
- `model.go` - this file will contain the data model for the teachers page.
- `server.go` - this file will contain the server-side logic for the teachers page.
- `database.go` - this file will contain the database access logic for the teachers page.

You can start by copying the structure from an existing page and then modify it to fit the teachers page requirements or you can create these files from scratch and implement the required functionality, which is described below.

In every file, the first line must define the package.
```go
package teachers

```
Next, in `server.go`, you should create a `Server` struct and define the necessary interfaces. The most basic implementation would look like this:
```go
type Server struct {
	Auth        Authentication
	Page        Page
	router      http.Handler
}
```
The interfaces can be copied from e.g. `server.go` from `home` package. As router uses `net/http` package and the interfaces use `templ` and `language` packages you must import them in your `server.go` file using
```go
import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)
```
Next you must create routing for the teachers page. You should define a getter for the router like:
```go
func (s *Server) Router() http.Handler {
	return s.router
}
```
and a method, where you will define the routes for the teachers page. For now, let's add a simple route that renders the whole page:
```go
func (s *Server) initRouter() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.page)
	s.router = router
}
```
You can see that this route is generic and without any `/teacher/` prefix. This will be taken care of in `main.go` by middleware, to which we will get.

You have to implement the `page()` method in `server.go` to render the teachers page. For now let's just log the request (need to import `log`):
```go
func (s *Server) page(w http.ResponseWriter, r *http.Request) {
	log.Println("Teachers page requested")
}
```
Finally, you must define `Init()` method for initializing the server.
```go
func (s *Server) Init() {
	s.initRouter()
}
```
Currently, only router need initialization but in the future we might need to initialize other components as well.

Second, we can connect the server to the main application. In `main.go`, you must import the `teachers` package:
```go
import (
	"github.com/michalhercik/RecSIS/teachers"
)
```
Then, you must create an instance of the `teachers.Server` struct and inject it with the necessary dependencies. You can copy the structure from any existing server. You are advised to create own function for initializing the server:
```go
func teachersServer(pageTempl page.Page) http.Handler {
  teachers := teachers.Server{ // inject dependencies
    Auth: cas.UserIDFromContext{},
    Page: page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
  }
  teachers.Init() // init router
  return teachers.Router() // we need only the router
}
```
From the server, we only need the router. The router must then be registered in the main application. You can do that in `setupHandler()` function. There, put newly created router into the `servers` struct:
```go
s := servers{
  ...
  teachersServer: teachersServer(pageTempl),
}

// you need to update the servers struct too
type servers struct {
  ...
  teachersServer http.Handler
}
```
This struct is passed to `protectedHandler()` function. That means that the teachers server is hidden behind authentication, and can only be accessed by authenticated users. If you want your page to be visible to every user, you should use the `unprotectedHandler()` function. Be beware, that you cannot expose any sensitive information through this route. In `protectedHandler()` function, you just need to add this row:
```go
handle(protectedRouter, teachersRoot, s.teachersServer)
```
`teachersRoot` should be string constant (those are defined at the end of the file) and should contain the root path for the teachers page, e.g. `/teachers/`. You can inspect the `handle` function to see how the middleware is applied.

Now you have running server with routing. You can test the logging but you might run into problems with `view_templ.go` file. You can solve them by putting the following template inside `view.templ` file. 

The sensible thing to do next is to add some HTML. This should be done in the `view.templ` file. Again, you can take inspiration from existing templates. We will create a very simple page, but if you want to make more complicated designs, you should get acquainted with go templates, htmx, bootstrap and alpine JS. For now let's this be our page:
```html
templ Content() {
  <div id="teachers-page" class="container pt-3">
    <h4>Teachers page</h4>
    <p>This is a page for teachers.</p>
    <div class="card">
      <div class="card-body">
        Teacher 1
      </div>
    </div>
    <div class="card">
      <div class="card-body">
        Teacher 2
      </div>
    </div>
    <div class="card">
      <div class="card-body">
        Teacher 3
      </div>
    </div>
  </div>
}
```
This is our content. Now to render it, we need to modify the `page()` method in `server.go` to use the new template. We will be injecting our content into the page template using the defined interface:
```go
View(main templ.Component, lang language.Language, title string, userID string) templ.Component
```
We have `main`, which is the `Content` component defined in `view.templ`, but we have to get language and user ID. That can easily be done from the context. Title of our page can be anything, for example "Teachers".
```go
func (s *Server) page(w http.ResponseWriter, r *http.Request) {
  main := Content() // from view.templ
  lang := language.FromContext(r.Context())
  title := "Teachers"
  userID := s.Auth.UserID(r)

  page := s.Page.View(main, lang, title, userID)
  page.Render(r.Context(), w)
}
```
Now on the `/teachers/` path, we should be able to see the teachers page.

Since we have the access to language, we can make our page multilingual. You should create a `texts.go` file for such string constants. You can take inspiration from such existing files but the content should look something like this:
```go
package teachers

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
  pageTitle       string
  headline        string
  description     string
}

var texts = map[language.Language]text{
  language.CS: {
    pageTitle:   "Učitelé",
    headline:    "Stránka učitelů",
    description: "Toto je stránka pro učitele.",
  },
  language.EN: {
    pageTitle:   "Teachers",
    headline:    "Teachers page",
    description: "This is a page for teachers.",
  },
}
```
Next, let's update our template:
```html
templ Content(t text) {
  <div id="teachers-page" class="container pt-3">
    <h4>{ t.headline }</h4>
    <p>{ t.description }</p>
    <div class="card">
      <div class="card-body">
        Teacher 1
      </div>
    </div>
    <div class="card">
      <div class="card-body">
        Teacher 2
      </div>
    </div>
    <div class="card">
      <div class="card-body">
        Teacher 3
      </div>
    </div>
  </div>
}
```
Our template now needs a `text` parameter which are the string constants. This parameter must be provided in the `page()` method:
```go
func (s *Server) page(w http.ResponseWriter, r *http.Request) {
  lang := language.FromContext(r.Context())
  text := texts[lang] // get the text for the current language from texts.go file
  title := text.pageTitle // get the page title
  userID := s.Auth.UserID(r)

  main := Content(text) // content now needs text

  page := s.Page.View(main, lang, title, userID)
  page.Render(r.Context(), w)
}
```

Finally, we would like to add real data to our page. For that we will use our database. But first, we must define what data we want to display. That shall be done in `model.go` file. Let's say that our data consists of the names of guarantors of `NPRG030 - Programming 1`. There are three of them, so the file will contain something like this:
```go
type teachersPage struct {
  guarantor1 teacher
  guarantor2 teacher
  guarantor3 teacher
}

type teacher struct {
  firstName   string
  lastName    string 
}

// you will need to import fmt package
func (t teacher) string() string {
  return fmt.Sprintf("%s %s", t.firstName, t.lastName)
}
```
The `teachersPage` struct is now ready to hold the data for our page. We will use it in the view:
```html
templ Content(tp *teachersPage, t text) {
  <div id="teachers-page" class="container pt-3">
    <h4>{ t.headline }</h4>
    <p>{ t.description }</p>
    <div class="card">
      <div class="card-body">
        { tp.guarantor1.string() }
      </div>
    </div>
    <div class="card">
      <div class="card-body">
        { tp.guarantor2.string() }
      </div>
    </div>
    <div class="card">
      <div class="card-body">
        { tp.guarantor3.string() }
      </div>
    </div>
  </div>
}
```
Of course, the struct must be passed to the template and populated with real data from the database. That we be again done in the `page()` method:
```go
func (s *Server) page(w http.ResponseWriter, r *http.Request) {
  lang := language.FromContext(r.Context())
  text := texts[lang]
  title := text.pageTitle
  userID := s.Auth.UserID(r)

  data := getTeachersPage() // but how??

  main := Content(data, text) // content now needs data and text

  page := s.Page.View(main, lang, title, userID)
  page.Render(r.Context(), w)
}
```
But what is this mysterious `getTeachersPage()` function? Also, do we even have a database connection? Well, no. Not now. Let's fix that.

First, let's dive into our `database.go` file. Here, we shall make all database calls, for which we need the database. Let's define a `DBManager` type that will hold our database connection:
```go
import "github.com/jmoiron/sqlx"

type DBManager struct {
  DB *sqlx.DB
}
```
This struct will be used to manage our database connection and provide methods for querying the database. We should put it in our `Server` struct, so it can be accessed from the server for database calls, and injected with a database connection when the server is created:
```go
type Server struct {
	Auth        Authentication
	Data        DBManager // now we have a database manager
	Page        Page
	router      http.Handler
}
```
We still need to inject the database connection into our `DBManager`. This is done in `main.go`. Let's update the `teachersServer()` function:
```go
func teachersServer(db *sqlx.DB, pageTempl page.Page) http.Handler {
	teachers := teachers.Server{
		Auth: cas.UserIDFromContext{},
		Data: teachers.DBManager{DB: db}, // injecting db connection
		Page: page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
	}
	teachers.Init()
	return teachers.Router()
}
```
Of course, we need to provide the database connection to the `teachersServer()` function:
```go
s := servers{
  ...
  teachersServer: teachersServer(db, pageTempl),
}
```
`db` is already defined and working in `setupHandler()` so we can use it directly.

Second, let's go back to `database.go` file. We need to define a method for getting the data. We also need a database model for the data, which will be transformed into our `teachersPage` struct:`
```go
type dbTeachersPage struct {
	Guarantors dbds.TeacherSlice `db:"guarantors"`
}

func (m *DBManager) teachersPage() (*teachersPage, error) {
	var dtp dbTeachersPage
	sql := "SELECT guarantors FROM courses WHERE code = 'NPRG030' AND lang = 'cs'"
	err := m.DB.Get(&dtp, sql)
	if err != nil {
		return nil, err
	}

	tp := teachersPage{
		guarantor1: intoTeacher(dtp.Guarantors[0]),
		guarantor2: intoTeacher(dtp.Guarantors[1]),
		guarantor3: intoTeacher(dtp.Guarantors[2]),
	}
	return &tp, nil
}

// data transformation from db model to page model
func intoTeacher(from dbds.Teacher) teacher {
	return teacher{
		firstName: from.FirstName,
		lastName:  from.LastName,
	}
}
```
We use `dbds` package for database models and as the teacher is defined there, we can reuse it. Do not forget to include it. SQL queries can get ugly, so they should be stored in separate files in `teachers/internal/sqlquery` directory in `sqlquery` package. Then you would have to include the appropriate SQL file in your `database.go` file and use it in your query. If you would like to use different data, refer to [Data Model](./data-model) section. It is of course possible to make some transformations, add new tables and views and use those. This is just for demonstration.

Third, we return to `server.go` file and out `page()` method. We no longer use some mysterious function but our `DBManager` method:
```go
func (s *Server) page(w http.ResponseWriter, r *http.Request) {
  lang := language.FromContext(r.Context())
  text := texts[lang]
  title := text.pageTitle
  userID := s.Auth.UserID(r)

  data, _ := s.Data.teachersPage() // this is clear now

  main := Content(data, text)

  page := s.Page.View(main, lang, title, userID)
  page.Render(r.Context(), w)
}
```

We have now successfully created a new page. There are many more things you can do with it, but we hope that this simple introduction has given you a good starting point. If you want a good practice, try extending this example by adding a search bar, which takes a course code as input and after clicking on a button, it fetches the corresponding guarantors/teachers of such course.

## Add a filter

Let's say we have in Meilisearch index *teachers*. The index contains documents with fields *firstName*, *lastName* and *department*. We would like to add a filter for *department* field which can have values *A*,*B* or *C*. The mapping of A,B,C letters to departments go like this A->KSI, B->UFAL, C->KAM. To do that, we need add a new filter, filter category and filter values to the database. This can be done by modifying ELT transformation responsible for creating filters (`initFilterTables`) or simply executing the following SQL queries in the database but the changes will not survive next ELT run:

```sql
INSERT INTO filters(id) VALUES ('teachers');
INSERT INTO filter_categories(id, filter_id, facet_id, title_cs, title_en, description_cs, description_en, displayed_value_limit, position) VALUES 
  (42, 'teachers', 'department', 'Katedra', 'Department', 'Katedra učitele', 'Teacher department', 3, 1);
INSERT INTO filter_values(id, category_id, facet_id, title_cs, title_en, description_cs, description_en, position) VALUES 
  (100, 'department', 'A', 'KSI', 'KSI', 'Katedra KSI', 'Department KSI', 1),
  (101, 'department', 'B', 'UFAL', 'UFAL', 'Katedra UFAL', 'Department UFAL', 2),
  (102, 'department', 'C', 'KAM', 'KAM', 'Katedra KAM', 'Department KAM', 3);
```

Now we can use this filter in our page. Let's say we want to add it to our teachers page. It could look roughly like this: 
```go
teacher.Server{
  Filters: filters.MakeFilters(db, "teachers"),
  //...
}
```
Lastly we would use the `Filters` as in an other server to make Search requests and display filters.

## Add a recommender

Let's say that the current for you Recommender on home page does not work for you and you would like to change it. Let's say you already implemented advanced search engine as a standalone service with REST API. The API implements `GET /{user_id}` endpoint which returns list of recommended course codes to a given user. Then you only need to create recommendation strategy that would utilize such endpoint and inject the strategy in `ForYou` field in `home.Server` struct. It would look something like this:

1. Create strategy in `recommend` package.
```go
package recommend

type MyAwesomeRecEngine struct {}

func (m MyAwesomeRecEngine) Recommend(userID string) ([]string, error) {
  res := m.makeRequestToAdvancedSearchEngineService(userID)
  return res.ListOfRecommendedCourseCodes
}
```
2. Inject `MyAwesomeRecEngine` into home page defined in `main.go` file. The expected type for `home.Server.ForYou` is type that implements interface with only single method `Recommend(userID string) ([]string, error)`. Our type `recommend.MyAwesomeRecEngine` implements the interface. and we can simply replace used type for home page. It should look something like this:
```go
home.Server{
  ForYou: recommend.MyAwesomeRecEngine{},
  //...
}
```

## Add error configuration

A simple thing to do is to add configuration options for error handling. This can include things like:
- logging to different outputs
- logging with different levels (e.g., info, warning, error)

Sensible way to do this, would be to read about [`log` package](https://pkg.go.dev/log). Then extent `config` struct in `main.go` file to contain some error handling configuration (output file, log level, ...) and update `config.toml` file. Then, pass this configuration to the error handler in `setupHandler()` function.
```go
errorHandler := errorx.ErrorHandler{
  // place for error handling configuration
}
```
Also, you would have to update the `ErrorHandler` struct in `error.go` file.
```go
type ErrorHandler struct {
	Page Page
	// place for error handling configuration
}
```

Finally, you can extent the `Log()` method using the new configuration options.
```go
func (eh ErrorHandler) Log(err error) {
  log.Println(fmt.Errorf("ERROR: %w", err))
  // log to more places
  // check log level
  // ...
}
```

## Use `LangString`s

Current application uses user related string constants without using `LangString` type. The example can be seen in all `texts.go` files. For standart texts, this does not create any complications, but for error messages, this means that every function/method that uses these strings has to be aware of the current language.

That is why we should use `LangString` type for, at least, all error messages. In `texts.go` all user error messages should be refactored to use `LangString` type. An example would be make this:
```go
var errPageNotFound := language.MakeLangString(
	"Stránka nenalezena",
	"Page not found",
)
```
instead of what it is now.

`HTTPError` struct from `errorx` package would then look like this:
```go
type HTTPError struct {
	Err     error
	Code    int
	UserMsg language.LangString
}
```
`UnwrapError()` function would no longer need language, and would return `(int,language.LangString)` and all render methods would use `userMsg language.LangString`. You would only need to use `userMsg.String(lang)` to get the appropriate string representation.
