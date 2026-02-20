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
	Student        string
	DegreePlan     string
	EnrollmentYear int
	Blueprint      string
}

func (r recRequest) MarshalJSON() ([]byte, error) {
	body := `{
		"algo": 	"%s",
		"limit":      %d,
		"user_id":    "%s",
		"student":    "%s",
		"blueprint":  %s,
		"degree_plan": "%s",
		"enrollment_year": %d
	}`
	body = fmt.Sprintf(body, r.Algo, r.Limit, r.UserID, r.Student, r.Blueprint, r.DegreePlan, r.EnrollmentYear)
	return []byte(body), nil
}

type RestCallWithAlgoSwitch struct {
	DB           *sqlx.DB
	Client       *http.Client
	Endpoint     string
	AlgoEndpoint string
	FitEndpoint  string
}

func (c RestCallWithAlgoSwitch) Algorithms() ([]string, []bool, error) {
	resp, err := c.Client.Get(c.AlgoEndpoint)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	rawBody, _ := io.ReadAll(resp.Body)
	body := struct {
		Algorithms []string `json:"algorithms"`
		Fit        []bool   `json:"fit"`
	}{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, nil, err
	}
	// TODO
	return body.Algorithms, body.Fit, nil
}

func (c RestCallWithAlgoSwitch) Recommend(userID, student, algo string, limit int) ([]string, []string, []string, []bool, error) {
	req := recRequest{
		Algo:   algo,
		Limit:  limit,
		UserID: userID,
		Student: student,
		Blueprint: "null",
	}
	var err error
	if student == "" {
		req.Blueprint, err = c.blueprint(userID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		req.DegreePlan, req.EnrollmentYear, err = c.studyInfo(userID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	finished, recommended, expected, target, err := c.call(req)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return finished, recommended, expected, target, nil
}

func (c RestCallWithAlgoSwitch) Fit(algo string) error {
	payload, err := json.Marshal(struct {
		Algo string `json:"algo"`
	}{
		Algo: algo,
	})
	if err != nil {
		return err
	}
	_, err = c.Client.Post(c.FitEndpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	return nil
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

func (c RestCallWithAlgoSwitch) call(reqParams recRequest) ([]string, []string, []string, []bool, error) {
	req, err := c.buildRequest(reqParams)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer resp.Body.Close()
	rawBody, _ := io.ReadAll(resp.Body)
	body := struct {
		Recommended []string `json:"recommended"`
		Finished []string `json:"finished"`
		Expected []string `json:"expected"`
		Target []bool `json:"target"`
	}{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, nil, nil, nil, err
	}
	// TODO
	return body.Finished, body.Recommended, body.Expected, body.Target, nil
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
