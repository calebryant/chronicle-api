package resources

import (
	"fmt"
	"net/http"
	"net/url"
)

// shared function for all REST API Activate methods
func CreateActivateRequest(endpoint *url.URL, path ResourcePath) (*http.Request, error) {
	if !path.HasValue() {
		return nil, fmt.Errorf("missing resource value")
	}
	return MethodRequest(
		http.MethodPost,
		endpoint,
		path.String()+":activate",
		nil,
		nil,
	)
}

// shared function for all REST API Activate methods
func CreateDeactivateRequest(endpoint *url.URL, path ResourcePath) (*http.Request, error) {
	if !path.HasValue() {
		return nil, fmt.Errorf("missing resource value")
	}
	return MethodRequest(
		http.MethodPost,
		endpoint,
		path.String()+":deactivate",
		nil,
		nil,
	)
}

// shared function for all REST API Get methods
func CreateGetRequest(endpoint *url.URL, path ResourcePath) (*http.Request, error) {
	if !path.HasValue() {
		return nil, fmt.Errorf("missing resource value")
	}
	return MethodRequest(
		http.MethodGet,
		endpoint,
		path.String(),
		nil,
		nil,
	)
}

// shared function for all REST API List methods
func CreateListRequest(endpoint *url.URL, path ResourcePath, query url.Values) (*http.Request, error) {
	return MethodRequest(
		http.MethodGet,
		endpoint,
		path.StripLastElement(),
		query,
		nil,
	)
}
