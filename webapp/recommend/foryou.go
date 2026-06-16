package recommend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/michalhercik/RecSIS/errorx"
	"io"
	"net/http"
)

type ForYouRecommender interface {
	Recommend(userID string, limit int) (ForYouResponse, error)
}

type ForYouResponse struct {
	Courses    []string `json:"pred"`
	Groups     [][]int  `json:"groups"`
	Categories Category `json:"categories"`
}

type Category struct {
	Names  []string   `json:"names"`
	Values [][]string `json:"pred"`
	Groups [][][]int  `json:"groups"`
}

type ForYou struct {
	Client    *http.Client
	Endpoint  string
	blueprint BlueprintFetcher
}

func (fy ForYou) Recommend(userID string, limit int) (ForYouResponse, error) {
	return fy.RecommendWith(userID, limit, false, false)
}

func (fy ForYou) RecommendWith(userID string, limit int, categories, groups bool) (ForYouResponse, error) {
	blueprint, err := fy.blueprint.fetch(userID)
	if err != nil {
		return ForYouResponse{}, errorx.AddContext(
			err,
			errorx.P("limit", limit),
			errorx.P("categories", categories),
			errorx.P("groups", groups),
		)
	}
	req := forYouRequest{
		UserID:     userID,
		Blueprint:  blueprint,
		Limit:      limit,
		Categories: categories,
		Groups:     groups,
	}
	result, err := fy.call(req)
	if err != nil {
		return ForYouResponse{}, errorx.AddContext(err, errorx.P("limit", limit))
	}
	return result, nil
}

func (fy ForYou) call(reqParams forYouRequest) (ForYouResponse, error) {
	req, err := fy.prepareRequest(reqParams)
	if err != nil {
		return ForYouResponse{}, errorx.AddContext(err)
	}
	resp, err := fy.Client.Do(req)
	if err != nil {
		return ForYouResponse{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("do request: %w", err)),
			http.StatusInternalServerError,
			"",
		)
	}
	defer resp.Body.Close()
	body, err := fy.parseResponse(resp)
	if err != nil {
		return ForYouResponse{}, errorx.AddContext(err)
	}
	return body, nil
}

func (fy ForYou) prepareRequest(reqParams forYouRequest) (*http.Request, error) {
	payload, err := reqParams.MarshalJSON()
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("marshal json: %w", err)),
			http.StatusInternalServerError,
			"",
		)
	}
	req, err := http.NewRequest(http.MethodPost, fy.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("create HTTP request: %w", err)),
			http.StatusInternalServerError,
			"",
		)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (fy ForYou) parseResponse(resp *http.Response) (ForYouResponse, error) {
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ForYouResponse{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("read io: %w", err)),
			http.StatusInternalServerError,
			"",
		)

	}
	var body ForYouResponse
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return ForYouResponse{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("unmarshal: %w", err)),
			http.StatusInternalServerError,
			"",
		)
	}
	return body, nil
}

type forYouRequest struct {
	UserID     string
	Blueprint  string
	Limit      int
	Categories bool
	Groups     bool
}

func (r forYouRequest) MarshalJSON() ([]byte, error) {
	body := `{
		"user_id":    "%s",
		"limit":      %d,
		"blueprint":  %s,
	}`
	body = fmt.Sprintf(body, r.UserID, r.Limit, r.Blueprint)
	return []byte(body), nil
}
