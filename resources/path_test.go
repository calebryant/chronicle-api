package resources_test

import (
	"encoding/json"
	"testing"

	"github.com/calebryant/chronicle-api/resources"
	"github.com/stretchr/testify/assert"
)

func TestResourcePathUnmarshaJSON(t *testing.T) {
	ts := struct {
		Name *resources.ResourcePath `json:"name"`
	}{}
	testJson := []byte(`{"name": "projects/testproject/locations/us/instances/12345/logTypes/WINEVTLOG/parsers/12345"}`)
	json.Unmarshal(testJson, &ts)
	assert.Equal(t, "parsers", ts.Name.Resource().Name)
}

func TestResourcePathString(t *testing.T) {
	tt := []struct {
		name     string
		expected string
	}{
		{
			name:     "test1",
			expected: "projects/testproject/locations/us/instances/123456789",
		},
		{
			name:     "test2",
			expected: "projects/testproject/locations/us/instances/123456789/logTypes",
		},
		{
			name:     "test3",
			expected: "projects/testproject/locations/us/instances/123456789/logTypes/WINEVTLOG/parsers",
		},
		{
			name:     "test4",
			expected: "projects/testproject/locations/us/instances/123456789/logTypes/WINEVTLOG/parsers/12345",
		},
	}
	for _, testCase := range tt {
		path := &resources.ResourcePath{}
		path.UnmarshalJSON([]byte(testCase.expected))
		assert.Equal(t, testCase.expected, path.String())
	}
}
