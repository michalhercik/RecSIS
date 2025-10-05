package recommend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type recRequest struct {
	Algo           string
	Limit          int
	UserID         string
	DegreePlan     string
	EnrollmentYear int
	Blueprint      string
}

func (r recRequest) MarshalJSON() ([]byte, error) {
	body := `{
		"algo": 	"%s",
		"limit":      %d,
		"user_id":    "%s",
		"blueprint":  %s,
		"degree_plan": "%s",
		"enrollment_year": %d
	}`
	body = fmt.Sprintf(body, r.Algo, r.Limit, r.UserID, r.Blueprint, r.DegreePlan, r.EnrollmentYear)
	return []byte(body), nil
}

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

func (c RestCallWithAlgoSwitch) Recommend(userID, algo string, limit int) ([]string, error) {
	req := recRequest{
		Algo:   algo,
		Limit:  limit,
		UserID: userID,
	}
	var err error
	req.Blueprint, err = c.blueprint(userID)
	if err != nil {
		return nil, err
	}
	req.DegreePlan, req.EnrollmentYear, err = c.studyInfo(userID)
	if err != nil {
		return nil, err
	}
	result, err := c.call(req)
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
					CASE
						WHEN semester = 0 THEN 'unassigned'
						WHEN semester = 1 THEN 'winter'
						WHEN semester = 2 THEN 'summer'
						ELSE semester::text
					END,
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

func (c RestCallWithAlgoSwitch) studyInfo(userID string) (string, int, error) {
	query := `--sql
		SELECT 
			degree_plan_code, start_year
		FROM studies
		WHERE user_id = $1
	`
	var degreePlan string
	var enrollmentYear int
	err := c.DB.QueryRow(query, userID).Scan(&degreePlan, &enrollmentYear)
	if err != nil {
		return degreePlan, enrollmentYear, err
	}
	return degreePlan, enrollmentYear, nil
}

func (c RestCallWithAlgoSwitch) call(reqParams recRequest) ([]string, error) {
	req, err := c.buildRequest(reqParams)
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

func (c RestCallWithAlgoSwitch) buildRequest(reqParams recRequest) (*http.Request, error) {
	payload, err := reqParams.MarshalJSON()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
