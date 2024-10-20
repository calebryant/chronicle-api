package parsers_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/calebryant/chronicle-api/resources/parsers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParserResource(t *testing.T) {
	testproject := "testproject"
	testlocation := "us"
	testinstance := "testinstance"
	testlogtype := "WINEVTLOG"
	testparserval := "1234567890"
	tt := []struct {
		name       string
		value      *parsers.ParserResource
		expectfail bool
		expected   string
	}{
		{
			name:     "Valid test",
			value:    parsers.NewParserResource(testproject, testlocation, testinstance, testlogtype, testparserval),
			expected: fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s/parsers/%s", testproject, testlocation, testinstance, testlogtype, testparserval),
		},
		{
			name:     "Valid test (no parser value)",
			value:    parsers.NewParserResource(testproject, testlocation, testinstance, testlogtype, ""),
			expected: fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s/parsers", testproject, testlocation, testinstance, testlogtype),
		},
		{
			name:       "test2",
			value:      parsers.NewParserResource(testproject, testlocation, "", "", ""),
			expectfail: true,
			expected:   "",
		},
		{
			name:       "test3",
			value:      parsers.NewParserResource(testproject, "", "", "", ""),
			expectfail: true,
			expected:   "",
		},
		{
			name:       "test4",
			value:      parsers.NewParserResource("", "", "", "", ""),
			expectfail: true,
			expected:   "",
		},
		{
			name:       "No logtype1",
			expectfail: true,
			value:      parsers.NewParserResource(testproject, testlocation, testinstance, "", ""),
			expected:   fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s/parsers", testproject, testlocation, testinstance, testparserval),
		},
		{
			name:       "No logtype2",
			expectfail: true,
			value:      parsers.NewParserResource(testproject, testlocation, testinstance, "", testparserval),
			expected:   fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s/parsers", testproject, testlocation, testinstance, testparserval),
		},
	}
	for _, tt := range tt {
		if tt.expectfail {
			assert.Nil(t, tt.value)
			continue
		}
		assert.Equal(t, tt.expected, tt.value.Name.String())
	}
}

func TestParserMethods(t *testing.T) {
	testproject := "testproject"
	testlocation := "us"
	testinstance := "testinstance"
	testlogtype := "WINEVTLOG"
	testparserval := "1234567890"
	tu, _ := url.Parse("https://test.local")
	tt := []struct {
		name               string
		expectFail         bool
		value              *http.Request
		expectedHttpMethod string
		expectedUrlPath    string
		expectedQuery      string
		expectedBody       map[string]interface{}
	}{
		{
			name:               "Test Valid Activate Method",
			value:              createRequest(parsers.NewParserResource(testproject, testlocation, testinstance, testlogtype, testparserval), "activate", tu),
			expectedHttpMethod: "POST",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes/%s/parsers/%s:activate", testproject, testlocation, testinstance, testlogtype, testparserval),
			expectedQuery:      "",
			expectedBody:       nil,
		},
		{
			name:       "Test Invalid Activate Method",
			expectFail: true,
			value:      createRequest(parsers.NewParserResource(testproject, testlocation, testinstance, testlogtype, ""), "activate", tu),
		},
		{
			name:               "Test Valid Deactivate Method",
			value:              createRequest(parsers.NewParserResource(testproject, testlocation, testinstance, testlogtype, testparserval), "deactivate", tu),
			expectedHttpMethod: "POST",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes/%s/parsers/%s:activate", testproject, testlocation, testinstance, testlogtype, testparserval),
			expectedQuery:      "",
			expectedBody:       nil,
		},
		{
			name:       "Test Invalid Deactivate Method",
			expectFail: true,
			value:      createRequest(parsers.NewParserResource(testproject, testlocation, testinstance, testlogtype, ""), "activate", tu),
		},
	}
	for _, tt := range tt {
		if tt.value == nil {
			require.True(t, tt.expectFail, tt.name)
			continue
		}
		bodyBytes, _ := io.ReadAll(tt.value.Body)
		parsedBody := make(map[string]interface{})
		json.Unmarshal(bodyBytes, &parsedBody)
		assert.Equal(t, tt.expectedUrlPath, tt.value.URL.Path)
		assert.Equal(t, tt.expectedBody, parsedBody)
		assert.Equal(t, tt.expectedQuery, tt.value.URL.Query().Encode())
		assert.Equal(t, tt.expectedHttpMethod, tt.value.Method)
	}
}

func createRequest(resource *parsers.ParserResource, methodType string, u *url.URL, options ...interface{}) *http.Request {
	if resource == nil {
		return nil
	}
	var err error
	var req *http.Request
	switch methodType {
	case "activate":
		req, err = resource.Activate(u)
	default:
		return nil
	}
	if err != nil {
		return nil
	}
	return req
}
