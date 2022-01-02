package dpaste

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	apiEndpoint       = "https:/dpaste.com/api/v2"
	CreateSuccessCode = 201
)

// Dpaste is an instance of th dpaste.com service.
type Dpaste struct {
	Token      string
	HttpClient http.Client
}

// Send request from the dpaste
func (d Dpaste) Send(r request) (response *http.Response, err error) {
	// @todo change the request interface (since now HttpClient is part of the Dpaste instance).
	response, err = r.Send(d)

	return
}

// interface for Dpaste client compatible requests
type request interface {
	Send(client Dpaste) (response *http.Response, err error)
}

// CreateRequest is a Dpaste request to create a new paste
type CreateRequest struct {
	content    string
	title      string
	expiryDays int
	syntax     string
}

// Send CreateRequest (though, probably should be something else...)
func (r CreateRequest) Send(client Dpaste) (response *http.Response, err error) {
	pData := url.Values{}

	pData.Set("content", r.content)
	pData.Set("title", r.title)
	pData.Set("expiry_days", strconv.Itoa(r.expiryDays))

	// @todo: We need validation for this against the API enumeration, also probably "auto" resolution with the undocumented
	// 		API.
	pData.Set("syntax", r.syntax)

	hr, err := http.NewRequest("POST", apiEndpoint, strings.NewReader(pData.Encode()))

	hr.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if client.Token != "" {
		hr.Header.Add("Authorization", "Bearer"+client.Token)
	}

	// Parse response or what?
	response, err = client.HttpClient.Do(hr)

	return
}
