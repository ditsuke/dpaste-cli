package dpaste

import "net/http"

// CreateResponse is the response from the Dpaste API to a CreateRequest
type CreateResponse struct {
	Response *http.Response
	Success  bool
	Code     int
	Location string
	Expiry   string
}
