package logs

import (
	"net/http"
	"net/url"

	"github.com/calebryant/chronicle-api/resources"
	"github.com/calebryant/chronicle-api/resources/instances"
)

// A logs API resource object
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances.logTypes.logs
type LogResource struct {
	Name                 resources.ResourcePath `json:"name,omitempty"`
	Data                 string                 `json:"data,omitempty"`
	LogEntryTime         string                 `json:"logEntryTime,omitempty"`
	CollectionTime       string                 `json:"collectionTime,omitempty"`
	EnvironmentNamespace string                 `json:"environmentNamespace,omitempty"`
	Labels               *logLabel              `json:"labels,omitempty"`
	Additionals          map[string]interface{} `json:"additionals,omitempty"`
}

func NewLogResource(project, location, instance, logtype, logVal string) *LogResource {
	if !instances.ValidInstance(project, location, instance) || logtype == "" {
		return nil
	}
	return &LogResource{
		Name: resources.NewResourcePath(
			project,
			location,
			instance,
			resources.LogtypesResourceName,
			logtype,
			resources.LogsResourceName,
			logVal,
		),
	}
}

// creates a list logs resource method http request
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances.logs/list
func (l *LogResource) List(serviceEndpoint *url.URL, pageSize, pageToken, filter string) (*http.Request, error) {
	return resources.MethodRequest(
		http.MethodGet,
		serviceEndpoint,
		l.Name.StripLastElement(),
		resources.CommonQueryParams(pageSize, pageToken, filter),
		nil,
	)
}

type logLabel struct {
	Value       string `json:"value,omitempty"`
	RbacEnabled bool   `json:"rbacEnabled,omitempty"`
}
