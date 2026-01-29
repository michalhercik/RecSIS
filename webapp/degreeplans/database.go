package degreeplans

import (
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

type dbDegreePlanRecord struct {
	Code      string `db:"plan_code"`
	Title     string `db:"title"`
	StudyType string `db:"study_type"`
	ValidFrom dpYear `db:"valid_from"`
	ValidTo   dpYear `db:"valid_to"`
}

func (m DBManager) degreePlanMetadata(dpCodes []string, lang language.Language) ([]degreePlanSearchResult, error) {
	var records []dbDegreePlanRecord
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

func intoDPSearchResults(records []dbDegreePlanRecord) []degreePlanSearchResult {
	results := make([]degreePlanSearchResult, len(records))
	for i, record := range records {
		results[i] = degreePlanSearchResult{
			code:      record.Code,
			title:     record.Title,
			studyType: record.StudyType,
			validFrom: record.ValidFrom,
			validTo:   record.ValidTo,
		}
	}
	return results
}
