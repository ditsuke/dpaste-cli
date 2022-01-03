package dpaste

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	apiEndpoint       = "https://dpaste.com/api/"
	CreateSuccessCode = 201
)

// Dpaste is an instance of the dpaste.com service.
type Dpaste struct {
	Token      string
	HttpClient http.Client
}

// Create request from the dpaste
func (d Dpaste) Create(r CreateRequest) (CreateResponse, error) {
	return r.send(d)
}

// interface for Dpaste client compatible requests
type request interface {
	Send(client Dpaste) (response *http.Response, err error)
}

// CreateRequest is a Dpaste request to create a new paste
type CreateRequest struct {
	Content    string
	Title      string
	ExpiryDays int
	Syntax     string
}

type CreateResponse struct {
	Response *http.Response
	Success  bool
	Code     int
	Location string
	Expiry   string
}

func (r CreateRequest) toQuery() (url.Values, error) {
	data := url.Values{}

	// Content is required
	if len(r.Content) == 0 {
		return nil, errors.New("invalid request")
	}
	data.Set("content", r.Content)

	if len(r.Title) > 0 {
		data.Set("title", r.Title)
	}
	if r.ExpiryDays > 0 {
		data.Set("expiry_days", strconv.Itoa(r.ExpiryDays))
	}
	// @todo To do anything with "syntax" we need to validate it here first

	return data, nil
}

// send CreateRequest (though, probably should be something else...)
func (r CreateRequest) send(client Dpaste) (response CreateResponse, err error) {

	// Setup request body
	pData, err := r.toQuery()
	if err != nil {
		return
	}

	hr, err := http.NewRequest("POST", apiEndpoint, strings.NewReader(pData.Encode()))
	if err != nil {
		return
	}

	// Add request headers
	hr.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if client.Token != "" {
		hr.Header.Add("Authorization", "Bearer "+client.Token)
	}

	httpResponse, err := client.HttpClient.Do(hr)
	if err != nil {
		return
	}

	success := httpResponse.StatusCode == CreateSuccessCode
	response = CreateResponse{httpResponse, success, httpResponse.StatusCode, httpResponse.Header.Get("Location"), httpResponse.Header.Get("Expires")}

	return
}
