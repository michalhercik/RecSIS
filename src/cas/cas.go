package cas

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

type CAS struct {
	Host string
}

func (c CAS) loginURLToCAS(service string) string {
	url := url.URL{
		Scheme:   "https",
		Host:     c.Host,
		Path:     "cas/login",
		RawQuery: url.Values{"service": []string{service}}.Encode(),
	}
	return url.String()
}

func (c CAS) validateTicket(r *http.Request, service string) (string, string, error) {
	validateReq, ticket := c.validateTicketURLToCAS(r, service)
	res, err := http.Get(validateReq)
	if err != nil {
		return "", "", err
	}
	var valRes validationResponse
	err = json.NewDecoder(res.Body).Decode(&valRes)
	if err != nil {
		return "", "", err
	}
	return valRes.ServiceResponse.AuthenticationSuccess.User, ticket, nil
}

func (c CAS) validateTicketURLToCAS(r *http.Request, service string) (string, string) {
	ticket := r.FormValue("ticket")
	query := url.Values{}
	query.Add("service", service)
	query.Add("ticket", ticket)
	query.Add("format", "json")
	validateReq := url.URL{
		Scheme:   "https",
		Host:     c.Host,
		Path:     "/cas/serviceValidate",
		RawQuery: query.Encode(),
	}
	return validateReq.String(), ticket
}

func (c CAS) UserIDTicketFromCASLogoutRequest(r *http.Request) (string, string, error) {
	rawPayload := r.FormValue("logoutRequest")
	var payload logoutRequest
	err := xml.Unmarshal([]byte(rawPayload), &payload)
	if err != nil {
		return "", "", fmt.Errorf("UserIDTicketFromCASLogoutRequest: %w", err)
	}
	return payload.UserID, payload.Ticket, nil
}

type logoutRequest struct {
	UserID string `xml:"NameID"`
	Ticket string `xml:"SessionIndex"`
}

type validationResponse struct {
	ServiceResponse struct {
		AuthenticationSuccess struct {
			Attributes struct {
				// 	AuthenticationDate                     []float64 `json:"authenticationDate"`
				// 	AuthenticationMethod                   []string  `json:"authenticationMethod"`
				// 	ClientIPAddress                        []string  `json:"clientIpAddress"`
				// 	CN                                     []string  `json:"cn"`
				// 	CredentialType                         []string  `json:"credentialType"`
				// 	CuniPersonalID                         []string  `json:"cunipersonalid"`
				EduPersonScopedAffiliation []string `json:"edupersonscopedaffiliation"`
				// 	GivenName                              []string  `json:"givenname"`
				// 	IsFromNewLogin                         []bool    `json:"isFromNewLogin"`
				// 	LongTermAuthenticationRequestTokenUsed []bool    `json:"longTermAuthenticationRequestTokenUsed"`
				// 	Mail                                   []string  `json:"mail"`
				// 	SAMLAuthenticationStatementAuthMethod  []string  `json:"samlAuthenticationStatementAuthMethod"`
				// 	ServerIPAddress                        []string  `json:"serverIpAddress"`
				// 	SN                                     []string  `json:"sn"`
				// 	SuccessfulAuthenticationHandlers       []string  `json:"successfulAuthenticationHandlers"`
				// 	UID                                    []string  `json:"uid"`
				// 	UserAgent                              []string  `json:"userAgent"`
			} `json:"attributes"`
			User string `json:"user"`
		} `json:"authenticationSuccess"`
	} `json:"serviceResponse"`
}
