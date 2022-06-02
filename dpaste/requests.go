package dpaste

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CreateRequest is a Dpaste request to create a new paste
type CreateRequest struct {
	Content    io.Reader
	Title      string
	ExpiryDays int
	Syntax     string
}

// Parse CreateRequest to a valid query which won't break things
func (r CreateRequest) toQuery(client Dpaste) (url.Values, error) {
	data := url.Values{}
	var err error

	s, err := ioutil.ReadAll(r.Content)
	if err != nil {
		return nil, err
	}

	data.Set("content", string(s))

	if len(r.Title) > 0 {
		data.Set("title", r.Title)
	}
	if r.ExpiryDays > 0 {
		data.Set("expiry_days", strconv.Itoa(r.ExpiryDays))
	}
	// @todo To do anything with "syntax" we need to validate it here first
	if ok, err := client.isValidSyntax(r.Syntax); err == nil && ok { // add also syntax len check, emit error if invalid syntax
		data.Set("syntax", r.Syntax)
	}

	return data, err
}

// send CreateRequest (though, probably should be something else...)
func (r CreateRequest) send(client Dpaste) (response CreateResponse, err error) {

	// Setup request body
	pData, err := r.toQuery(client)
	if err != nil {
		return
	}

	hr, err := http.NewRequest("POST", apiEndpoint, strings.NewReader(pData.Encode()))
	if err != nil {
		return
	}

	// Add request headers
	hr.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpResponse, err := client.DoRequest(hr)
	if err != nil {
		return
	}

	success := httpResponse.StatusCode == CreateSuccessCode
	response = CreateResponse{httpResponse, success, httpResponse.StatusCode, httpResponse.Header.Get("Location"), httpResponse.Header.Get("Expires")}

	return
}
