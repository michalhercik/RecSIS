package recommend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

type RestCallWithAlgoSwitch struct {
	DB           *sqlx.DB
	Client       *http.Client
	Endpoint     string
	AlgoEndpoint string
}

func (c RestCallWithAlgoSwitch) Algorithms() ([]string, error) {
	resp, err := c.Client.Get(c.AlgoEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rawBody, _ := io.ReadAll(resp.Body)
	body := struct {
		Algorithms []string `json:"algorithms"`
	}{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, err
	}
	// TODO
	return body.Algorithms, nil
}

func (c RestCallWithAlgoSwitch) Recommend(userID, algo string) ([]string, error) {
	blueprint, err := c.blueprint(userID)
	if err != nil {
		return nil, err
	}
	studyInfo, err := c.studyInfo(userID)
	if err != nil {
		return nil, err
	}
	result, err := c.call(algo, userID, blueprint, studyInfo)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c RestCallWithAlgoSwitch) blueprint(userID string) (string, error) {
	query := `--sql
		WITH course_per_semester AS (
			-- 1️⃣ Gather courses for each academic year + semester
			SELECT 
				by.academic_year,
				bs.semester,
				jsonb_agg(bc.course_code) AS courses
			FROM blueprint_years by
			LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = by.id
			LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id = bs.id
			WHERE by.user_id = $1
			GROUP BY by.academic_year, bs.semester
		),
		semester_pivot AS (
			-- 2️⃣ Pivot semesters into JSON object keys
			SELECT
				academic_year,
				jsonb_object_agg(
					semester,
					COALESCE(
						NULLIF(courses, '[null]'::jsonb),  -- treat [null] as empty
						'[]'::jsonb
					)
				) AS semester_data
			FROM course_per_semester
			GROUP BY academic_year
		)

		-- 3️⃣ Final JSON array with {year, <semester>: [...]}
		SELECT jsonb_agg(
				jsonb_build_object('year', academic_year) || semester_data
				ORDER BY academic_year
			) AS result
		FROM semester_pivot;
	`
	var bp string
	err := c.DB.QueryRow(query, userID).Scan(&bp)
	if err != nil {
		return bp, err
	}
	return bp, nil
}

func (c RestCallWithAlgoSwitch) studyInfo(userID string) (string, error) {
	query := `--sql
		SELECT 
			jsonb_build_object(
				'degree_plan_code', degree_plan_code, 
				'start_year', start_year
			) 
		FROM studies
		WHERE user_id = $1
	`
	var dpCodeYear string
	err := c.DB.QueryRow(query, userID).Scan(&dpCodeYear)
	if err != nil {
		return dpCodeYear, err
	}
	return dpCodeYear, nil
}

func (c RestCallWithAlgoSwitch) call(algo, userID string, blueprint, studyInfo string) ([]string, error) {
	req, err := c.buildRequest(algo, userID, blueprint, studyInfo)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rawBody, _ := io.ReadAll(resp.Body)
	body := struct {
		Recommended []string `json:"recommended"`
	}{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, err
	}
	// TODO
	return body.Recommended, nil
}

func (c RestCallWithAlgoSwitch) buildRequest(algo, userID, blueprint, studyInfo string) (*http.Request, error) {
	body := `{
		"algo": 	"%s",
		"user_id":    "%s",
		"blueprint":  %s,
		"study_info": %s
	}`
	body = fmt.Sprintf(body, algo, userID, blueprint, studyInfo)
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
