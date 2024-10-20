package logtypes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/calebryant/chronicle-api/resources"
	"github.com/calebryant/chronicle-api/resources/instances"
)

// A logTypes API resource object
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances.logTypes
type LogTypeResource struct {
	Name        resources.ResourcePath `json:"name,omitempty"`
	DisplayName string                 `json:"displayName,omitempty"`
	Golden      bool                   `json:"golden"`
}

func NewLogTypeResource(project, location, instance, logtype string) *LogTypeResource {
	if !instances.ValidInstance(project, location, instance) {
		return nil
	}
	return &LogTypeResource{
		Name: resources.NewResourcePath(
			project,
			location,
			instance,
			resources.LogtypesResourceName,
			logtype,
		),
	}
}

// creates a get log type resource method http request
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances.logTypes/get
func (l *LogTypeResource) Get(serviceEndpoint *url.URL) (*http.Request, error) {
	return resources.CreateGetRequest(serviceEndpoint, l.Name)
}

// creates a list log types resource method http request
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances.logTypes/list
func (l *LogTypeResource) List(serviceEndpoint *url.URL, pageSize, pageToken string) (*http.Request, error) {
	return resources.CreateListRequest(
		serviceEndpoint,
		l.Name,
		resources.CommonQueryParams(pageSize, pageToken, ""),
	)
}

// creates a run parser resource method request
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances.logTypes/runParser
func (l *LogTypeResource) RunParser(serviceEndpoint *url.URL, cbnBytes, cbnSnippetBytes []byte, logs [][]byte, statedumpAllowed bool) (*http.Request, error) {
	if len(cbnBytes) == 0 {
		return nil, fmt.Errorf("no cbn provided")
	}
	// set body params
	cbn := map[string]string{
		"cbn": base64.StdEncoding.EncodeToString(cbnBytes),
	}
	reqBody := map[string]interface{}{
		"parser":           cbn,
		"log":              encodeLogs(logs),
		"statedumpAllowed": statedumpAllowed,
	}
	if len(cbnSnippetBytes) != 0 {
		cbnSnippet := map[string]string{
			"cbnSnippet": base64.StdEncoding.EncodeToString(cbnSnippetBytes),
		}
		reqBody["parserExtension"] = cbnSnippet
	}
	return resources.MethodRequest(
		http.MethodPost,
		serviceEndpoint,
		l.Name.String()+":runParser",
		nil,
		reqBody,
	)
}

func encodeLogs(logs [][]byte) []string {
	encodedLogs := []string{}
	for _, logmsg := range logs {
		encodedLogs = append(encodedLogs, base64.StdEncoding.EncodeToString(logmsg))
	}
	return encodedLogs
}
