package resources

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

type resourcePathElement struct {
	Name   string
	Value  string
	parent *resourcePathElement
}

type ResourcePath struct {
	resource *resourcePathElement
}

func NewResourcePath(project, location, instance string, elements ...string) ResourcePath {
	if project == "" || location == "" || instance == "" {
		return ResourcePath{}
	}
	resourcePathSlice := []string{
		"projects",
		project,
		"locations",
		location,
		"instances",
		instance,
	}
	resourcePathSlice = append(resourcePathSlice, elements...)
	var currentResource *resourcePathElement = nil
	for i, element := range resourcePathSlice {
		// even indexes in a resource path are resource names
		if i%2 == 0 {
			currentResource = &resourcePathElement{
				Name:   element,
				parent: currentResource,
			}
		} else {
			currentResource.Value = element
		}
	}
	return ResourcePath{
		resource: currentResource,
	}
}

func (p *ResourcePath) Resource() *resourcePathElement {
	return p.resource
}

func (p *ResourcePath) UnmarshalJSON(data []byte) error {
	var pathString string
	err := json.Unmarshal(data, &pathString)
	if err != nil {
		return err
	}
	resourcePathSlice := strings.Split(pathString, "/")
	project := resourcePathSlice[1]
	location := resourcePathSlice[3]
	instance := resourcePathSlice[5]
	theRest := resourcePathSlice[6:]
	newResource := NewResourcePath(project, location, instance, theRest...)
	if newResource.resource == nil {
		return fmt.Errorf("invalid resource path %s", string(data))
	}
	*p = newResource
	return nil
}

func (p *ResourcePath) String() string {
	var resourceSlice []string
	currentResource := p.resource
	for currentResource != nil {
		resourceSlice = append(
			[]string{
				currentResource.Name,
				currentResource.Value,
			},
			resourceSlice...,
		)
		currentResource = currentResource.parent
	}
	return path.Join(resourceSlice...)
}

func (p *ResourcePath) Map() map[string]string {
	resourceMap := map[string]string{}
	currentResource := p.resource
	for currentResource.parent != nil {
		resourceMap[currentResource.Name] = currentResource.Value
	}
	return resourceMap
}

// Checks the last element in a URL path. If the value is non-empty, then return the path string with the last element removed. Otherwise return the unchanged path string.
func (p *ResourcePath) StripLastElement() string {
	if p.resource.Value != "" {
		newResource := *p
		newResource.resource.Value = ""
		return newResource.String()
	} else {
		return p.String()
	}
}

// Returns true if the resource value is not empty, false otherwise.
func (p *ResourcePath) HasValue() bool {
	if p.resource.Value != "" {
		return true
	} else {
		return false
	}
}
