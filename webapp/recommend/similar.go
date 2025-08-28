package recommend

import (
	"encoding/json"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearchSimilarToBlueprint struct {
	Search      meilisearch.ServiceManager
	SearchIndex meilisearch.IndexConfig
	QueryPrefix string
	Embedder    string
	DB          *sqlx.DB
}

func (m MeiliSearchSimilarToBlueprint) Recommend(userID string) ([]string, error) {
	blueprintCourses, err := m.blueprintCourses(userID)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	query := m.buildQuery(blueprintCourses)
	filter := m.buildFilter(blueprintCourses)
	similarCourses, err := m.similarCourses(query, filter, 30)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	selected := chooseRandom(similarCourses, 10)
	return selected, nil
}

func (m MeiliSearchSimilarToBlueprint) blueprintCourses(userID string) ([]string, error) {
	var courses []string
	query := `--sql
		SELECT c.title FROM blueprint_years by
		INNER JOIN blueprint_semesters bs ON by.id = bs.blueprint_year_id
		INNER JOIN blueprint_courses bc ON bs.id = bc.blueprint_semester_id
		INNER JOIN courses c ON bc.course_code = c.code 
		WHERE by.user_id = $1
		AND c.lang = 'en'
	`
	err := m.DB.Select(&courses, query, userID)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	return courses, nil
}

func (m MeiliSearchSimilarToBlueprint) buildQuery(blueprintCourses []string) string {
	query := m.QueryPrefix + strings.Join(blueprintCourses, ",")
	return query
}

func (m MeiliSearchSimilarToBlueprint) buildFilter(blueprintCourses []string) string {
	filter := "code NOT IN ['" + strings.Join(blueprintCourses, "','") + "']"
	filter += " AND section=NI"
	return filter
}

func (m MeiliSearchSimilarToBlueprint) similarCourses(query, filter string, limit int64) ([]string, error) {
	const SemanticOnlyRatio = 1.0
	req := &meilisearch.SearchRequest{
		Hybrid: &meilisearch.SearchRequestHybrid{
			SemanticRatio: SemanticOnlyRatio,
			Embedder:      m.Embedder,
		},
		AttributesToRetrieve: []string{"code"},
		Limit:                limit,
		Filter:               filter,
	}
	rawRes, err := m.Search.Index(m.SearchIndex.Uid).SearchRaw(query, req)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	var res response
	rawResByte, err := rawRes.MarshalJSON()
	if err != nil {
		// TODO: add context
		return nil, err
	}
	if err := json.Unmarshal(rawResByte, &res); err != nil {
		// TODO: add context
		return nil, err
	}
	courses := make([]string, len(res.Hits))
	for i, hit := range res.Hits {
		courses[i] = hit.Code
	}
	return courses, nil
}

type course struct {
	Code string `json:"code"`
}
type response struct {
	Hits []course `json:"hits"`
}
