package logs_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/calebryant/chronicle-api/resources/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogResource(t *testing.T) {
	testproject := "testproject"
	testlocation := "us"
	testinstance := "testinstance"
	testlogtype := "WINEVTLOG"
	testlogval := "1234567890"
	tt := []struct {
		name       string
		value      *logs.LogResource
		expectfail bool
		expected   string
	}{
		{
			name:     "Valid test",
			value:    logs.NewLogResource(testproject, testlocation, testinstance, testlogtype, testlogval),
			expected: fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s/logs/%s", testproject, testlocation, testinstance, testlogtype, testlogval),
		},
		{
			name:     "Valid test (no log value)",
			value:    logs.NewLogResource(testproject, testlocation, testinstance, testlogtype, ""),
			expected: fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s/logs", testproject, testlocation, testinstance, testlogtype),
		},
		{
			name:       "test2",
			value:      logs.NewLogResource(testproject, testlocation, "", "", ""),
			expectfail: true,
			expected:   "",
		},
		{
			name:       "test3",
			value:      logs.NewLogResource(testproject, "", "", "", ""),
			expectfail: true,
			expected:   "",
		},
		{
			name:       "test4",
			value:      logs.NewLogResource("", "", "", "", ""),
			expectfail: true,
			expected:   "",
		},
		{
			name:       "No logtype1",
			expectfail: true,
			value:      logs.NewLogResource(testproject, testlocation, testinstance, "", ""),
			expected:   fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s/logs", testproject, testlocation, testinstance, testlogval),
		},
		{
			name:       "No logtype2",
			expectfail: true,
			value:      logs.NewLogResource(testproject, testlocation, testinstance, "", testlogval),
			expected:   fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s/logs", testproject, testlocation, testinstance, testlogval),
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

func TestLogMethods(t *testing.T) {
	testproject := "testproject"
	testlocation := "us"
	testinstance := "testinstance"
	testlogtype := "WINEVTLOG"
	testlogval := "1234567890"
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
			name:               "Valid Test List Method (no query)",
			value:              createRequest(logs.NewLogResource(testproject, testlocation, testinstance, testlogtype, ""), "list", tu, "", "", ""),
			expectedHttpMethod: "GET",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes/%s/logs", testproject, testlocation, testinstance, testlogtype),
			expectedQuery:      "",
			expectedBody:       nil,
		},
		{
			name:               "Test List Method (with query)",
			value:              createRequest(logs.NewLogResource(testproject, testlocation, testinstance, testlogtype, testlogval), "list", tu, "100", "abcdefg", "filterquery"),
			expectedHttpMethod: "GET",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes/%s/logs", testproject, testlocation, testinstance, testlogtype),
			expectedQuery:      "filter=filterquery&pageSize=100&pageToken=abcdefg",
			expectedBody:       nil,
		},
	}
	for _, tt := range tt {
		if tt.expectFail {
			require.Nil(t, tt.value)
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

func createRequest(resource *logs.LogResource, methodType string, u *url.URL, options ...interface{}) *http.Request {
	if resource == nil {
		return nil
	}
	var err error
	var req *http.Request
	switch methodType {
	case "list":
		pageSize := options[0].(string)
		pageToken := options[1].(string)
		filter := options[2].(string)
		req, err = resource.List(u, pageSize, pageToken, filter)
	default:
		return nil
	}
	if err != nil {
		return nil
	}
	return req
}
