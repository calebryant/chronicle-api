package parsers

import (
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/calebryant/chronicle-api/resources"
	"github.com/calebryant/chronicle-api/resources/instances"
)

// A parsers API resource object
//
// https://cloud.google.com/chronicle/docs/reference/rest/v1alpha/projects.locations.instances.logTypes.parsers
type ParserResource struct {
	Name                 resources.ResourcePath `json:"name,omitempty"`
	Creator              *creator               `json:"creator,omitempty"`
	CreateTime           string                 `json:"createTime,omitempty"`
	Changelogs           *changelogs            `json:"changelogs,omitempty"`
	ParserExtension      string                 `json:"parserExtension,omitempty"`
	Type                 string                 `json:"type,omitempty"`
	State                string                 `json:"state,omitempty"`
	ValidationReport     string                 `json:"validationReport,omitempty"`
	ValidatedOnEmptyLogs bool                   `json:"validatedOnEmptyLogs,omitempty"`
	Cbn                  []byte                 `json:"cbn,omitempty"`
	ReleaseStage         string                 `json:"releaseStage,omitempty"`
	ValidationStage      string                 `json:"validationStage,omitempty"`
}

func NewParserResource(project, location, instance, logtype, parserId string) *ParserResource {
	if !instances.ValidInstance(project, location, instance) || logtype == "" {
		return nil
	}
	return &ParserResource{
		Name: resources.NewResourcePath(
			project,
			location,
			instance,
			resources.LogtypesResourceName,
			logtype,
			resources.ParsersResourceName,
			parserId,
		),
	}
}

func (p *ParserResource) Activate(serviceEndpoint *url.URL) (*http.Request, error) {
	return resources.CreateActivateRequest(serviceEndpoint, p.Name)
}

func (p *ParserResource) Deactivate(serviceEndpoint *url.URL) (*http.Request, error) {
	return resources.CreateDeactivateRequest(serviceEndpoint, p.Name)
}

func (p *ParserResource) Create(serviceEndpoint *url.URL) (*http.Request, error) {
	body := map[string]interface{}{
		"cbn":                  base64.StdEncoding.EncodeToString(p.Cbn),
		"validatedOnEmptyLogs": p.ValidatedOnEmptyLogs,
	}
	return resources.MethodRequest(
		http.MethodPost,
		serviceEndpoint,
		p.Name.StripLastElement(),
		nil,
		body,
	)
}

type creator struct {
	Customer string `json:"customer,omitempty"`
	Author   string `json:"author,omitempty"`
	Source   string `json:"source,omitempty"`
}

type changelogs struct {
	Entries []changeEntries
}

type changeEntries struct {
	CreateTime    string
	ChangeMessage string
}
