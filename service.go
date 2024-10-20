package chronicleapi

import (
	"net/url"
)

const (
	baseServiceEndpoint = "chronicle.googleapis.com"
)

// Builds a new Chronicle API service endpoint URL
//
// ex. https://us-chronicle.googleapis.com/v1alpha
func NewServiceEndpoint(region, version string) *url.URL {
	baseurl, _ := url.Parse("https://" + region + "-" + baseServiceEndpoint)
	return baseurl.JoinPath(version)
}
