package degreeplans

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/degreeplans/internal/sqlquery"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) degreePlanMetadata(dpCodes []string, lang language.Language) ([]degreePlanSearchResult, error) {
	var records []degreePlanSearchResult
	err := m.DB.Select(&records, sqlquery.DegreePlanMetadataForSearch, pq.Array(dpCodes), lang)
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DegreePlanMetadataForSearch: %w", err), errorx.P("dpCodes", dpCodes), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errFailedDPSearch,
		)
	}
	return records, nil
}

func (m DBManager) userHasSelectedDegreePlan(uid string) bool {
	var userPlan sql.NullString
	err := m.DB.Get(&userPlan, sqlquery.UserDegreePlanCode, uid)
	if err != nil {
		return false
	}
	return userPlan.Valid
}
