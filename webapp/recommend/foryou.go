package recommend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

func (fy ForYou) Recommend(userID string, categories, groups bool, limit int) (ForYouResponse, error) {
	blueprint, err := fy.blueprint.fetch(userID)
	if err != nil {
		return ForYouResponse{}, err
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
		return ForYouResponse{}, err
	}
	return result, nil

}

func (fy ForYou) call(req forYouRequest) (ForYouResponse, error) {
	req, err := fy.prepareRequest(reqParams)
	if err != nil {
		return ForYouResponse{}, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return ForYouResponse{}, err
	}
	defer resp.Body.Close()
	body, err := fy.parseResponse(resp)
	if err != nil {
		return ForYouResponse{}, err
	}
	return body, nil
}

func (fy ForYou) prepareRequest(reqParams forYouRequest) (*http.Request, error) {
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

func (fy ForYou) parseResponse(resp *http.Response) (ForYouResponse, error) {
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ForYouResponse{}, err
	}
	var body ForYouResponse
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return ForYouResponse{}, err
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
	algo, err := json.Marshal(r.Algo)
	if err != nil {
		return nil, err
	}
	body = fmt.Sprintf(body, r.UserID, r.Limit, r.Blueprint)
	return []byte(body), nil
}
