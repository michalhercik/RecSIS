package cas

import (
	"encoding/json"
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

func (c CAS) validateTicket(r *http.Request, service string) (string, error) {
	validateReq := c.validateTicketURLToCAS(r, service)
	res, err := http.Get(validateReq)
	if err != nil {
		return "", err
	}
	var valRes validationResponse
	err = json.NewDecoder(res.Body).Decode(&valRes)
	if err != nil {
		return "", err
	}
	return valRes.ServiceResponse.AuthenticationSuccess.User, nil
}

func (c CAS) validateTicketURLToCAS(r *http.Request, service string) string {
	query := url.Values{}
	query.Add("service", service)
	query.Add("ticket", r.FormValue("ticket"))
	query.Add("format", "json")
	validateReq := url.URL{
		Scheme:   "https",
		Host:     c.Host,
		Path:     "/cas/serviceValidate",
		RawQuery: query.Encode(),
	}
	return validateReq.String()
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
