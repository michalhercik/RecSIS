package receval

import (
	"fmt"
	"net/http"

	"github.com/michalhercik/RecSIS/cas"
	"github.com/michalhercik/RecSIS/components/page"
	"github.com/michalhercik/RecSIS/recommend"

	"github.com/jmoiron/sqlx"
)

func NewServer(db *sqlx.DB, errorHandler Error, pageTempl page.Page, host string, port int ) http.Handler {
	result := &Server{
		Auth: cas.UserIDFromContext{},
		Error: errorHandler,
		Page: page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
		Experiment: recommend.RestCallWithAlgoSwitch{
			Client:       &http.Client{},
			DB:           db,
			Endpoint:     fmt.Sprintf("http://%s:%d/recommended", host, port),
			AlgoEndpoint: fmt.Sprintf("http://%s:%d/algorithms", host, port),
			FitEndpoint: fmt.Sprintf("http://%s:%d/fit", host, port),
		},
		Data: DBManager{DB: db},
	}
	result.Init()
	return result.Router()
}
