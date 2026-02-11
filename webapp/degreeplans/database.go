package degreeplans

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/degreeplans/internal/sqlquery"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) degreePlanMetadata(dpCodes []string, lang language.Language) ([]degreePlanSearchResult, error) {
	var records []dbds.DegreePlan
	err := m.DB.Select(&records, sqlquery.DegreePlanMetadataForSearch, pq.Array(dpCodes), lang)
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DegreePlanMetadataForSearch: %w", err), errorx.P("dpCodes", dpCodes), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errFailedDPSearch,
		)
	}
	return intoDPSearchResults(records), nil
}

func intoDPSearchResults(records []dbds.DegreePlan) []degreePlanSearchResult {
	results := make([]degreePlanSearchResult, len(records))
	for i, record := range records {
		results[i] = degreePlanSearchResult{
			code:      record.Code,
			title:     record.Title,
			studyType: record.StudyType,
			validFrom: dpYear(record.ValidFrom),
			validTo:   dpYear(record.ValidTo),
		}
	}
	return results
}
