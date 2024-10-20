package resources

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

func MethodRequest(httpMethod string, serviceEndpoint *url.URL, path string, query url.Values, body map[string]interface{}) (*http.Request, error) {
	reqUrl := serviceEndpoint.JoinPath(path)
	if query != nil {
		reqUrl.RawQuery = query.Encode()
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(httpMethod, reqUrl.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Helper function to set values of common query params
func CommonQueryParams(size, token, filter string) url.Values {
	values := url.Values{}
	if size != "" {
		values.Set("pageSize", size)
	}
	if token != "" {
		values.Set("pageToken", token)
	}
	if filter != "" {
		values.Set("filter", filter)
	}
	return values
}
