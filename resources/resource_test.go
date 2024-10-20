package resources_test

import (
	"net/url"
	"testing"

	"github.com/calebryant/chronicle-api/resources/logtypes"
)

func TestStripLastElement(t *testing.T) {
	testUrl, _ := url.Parse("https://example.local")
	tt := []struct {
		name     string
		input    *logtypes.LogTypeResource
		expected string
	}{
		{
			name:     "test1",
			input:    logtypes.NewLogTypeResource("testproject", "us", "12345", "WINEVTLOG"),
			expected: "logTypes",
		},
		{
			name:     "test2",
			input:    logtypes.NewLogTypeResource("testproject", "us", "12345", ""),
			expected: "logTypes",
		},
	}
	for _, testCase := range tt {
		testCase.input.List(testUrl, "", "")
	}
}
