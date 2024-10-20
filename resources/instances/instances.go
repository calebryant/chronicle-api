package instances

import (
	"net/http"
	"net/url"

	"github.com/calebryant/chronicle-api/resources"
)

// An instance API resource object
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances#InstanceResource
type InstanceResource struct {
	Name resources.ResourcePath `json:"name,omitempty"`
}

func NewInstanceResource(project, location, instance string) *InstanceResource {
	if !ValidInstance(project, location, instance) {
		return nil
	}
	return &InstanceResource{
		Name: resources.NewResourcePath(
			project,
			location,
			instance,
		),
	}
}

// creates a get instance resource method http request
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances/get
func (i *InstanceResource) Get(serviceEndpoint *url.URL) (*http.Request, error) {
	return resources.MethodRequest(
		http.MethodGet,
		serviceEndpoint,
		i.Name.String(),
		nil,
		nil,
	)
}

func ValidInstance(project, location, instance string) bool {
	if project == "" || location == "" || instance == "" {
		return false
	} else {
		return true
	}
}
