package dpaste

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	apiEndpoint                   = "https://dpaste.com/api/v2/"
	syntaxChoicesEndpointRelative = "syntax-choices/"
	httpPacerDuration             = time.Second
	CreateSuccessCode             = 201
)

var (
	// A JSON object with the syntax choices, populated from the Dpaste APIs syntax-choices endpoint
	syntaxJson map[string]interface{}

	// We need to make sure we time the requests 1s or more apart, for this reason every time we send a request we need to mark the time
	lastRequest time.Time
)

// New creates and returns an instance of the Dpaste client
func New(token string, httpClient *http.Client) *Dpaste {
	return &Dpaste{token, httpClient}
}

// Dpaste is an instance of the dpaste.com service.
type Dpaste struct {
	token      string
	httpClient *http.Client
}

// DoRequest is a practical override to http.Client.Do, allowing for pacing between requests.
//
// The Authorization header is populated if it's set in the dpaste client
func (d Dpaste) DoRequest(r *http.Request) (response *http.Response, err error) {
	// We need a pacer here
	if lastRequest.IsZero() || time.Since(lastRequest) >= httpPacerDuration {
		if d.token != "" {
			r.Header.Add("Authorization", "Bearer "+d.token)
		}
		response, err = d.httpClient.Do(r)
		lastRequest = time.Now()
		return
	}

	time.Sleep(httpPacerDuration - time.Since(lastRequest))
	return d.DoRequest(r)
}

// Get an url from the dpaste API. The url must be relative to the apiEndpoint
func (d Dpaste) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", apiEndpoint+url, nil)
	if err != nil {
		return nil, err
	}
	return d.DoRequest(req)
}

// Post to a dpaste API endpoint. The default apiEndpoint is sufficient for creating pastes
// and incidentally, the only POSTable endpoint on the public dpaste API. Url must be relative to apiEndpoint
func (d Dpaste) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", apiEndpoint+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return d.DoRequest(req)

}

// Create request from the dpaste
func (d Dpaste) Create(r CreateRequest) (CreateResponse, error) {
	return r.send(d)
}

// isValidSyntax() returns true iff syntax is a valid syntax choice on dpaste.com
func (d Dpaste) isValidSyntax(syntax string) (bool, error) {
	if len(syntaxJson) == 0 {
		// Populate
		// We will go with a fetch-each-time approach right now
		response, err := d.Get(syntaxChoicesEndpointRelative)
		if err != nil {
			return false, err
		}
		jsonBody, err := ioutil.ReadAll(response.Body)
		var decodedData interface{}
		err = json.Unmarshal(jsonBody, &decodedData)
		if err != nil {
			return false, errors.New("failed to unmarshal json")
		}
		syntaxJson = decodedData.(map[string]interface{})
		if len(syntaxJson) == 0 {
			return false, errors.New("failed to decode response from " + syntaxChoicesEndpointRelative)
		}
	}
	if _, ok := syntaxJson[syntax]; ok {
		return true, nil
	}

	return false, errors.New("invalid syntax choice")
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

// Parse CreateRequest to a valid query which won't break things
func (r CreateRequest) toQuery(client Dpaste) (url.Values, error) {
	data := url.Values{}
	var err error

	// Content is required
	if len(r.Content) == 0 {
		return nil, errors.New("invalid request") // content needed
	}
	data.Set("content", r.Content)

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
