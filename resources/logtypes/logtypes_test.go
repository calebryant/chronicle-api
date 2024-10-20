package logtypes_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/calebryant/chronicle-api/resources/logtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogTypeResource(t *testing.T) {
	testproject := "testproject"
	testlocation := "us"
	testinstance := "testinstance"
	testlogtype := "WINEVTLOG"
	tt := []struct {
		name      string
		value     *logtypes.LogTypeResource
		expectnil bool
		expected  string
	}{
		{
			name:     "test1",
			value:    logtypes.NewLogTypeResource(testproject, testlocation, testinstance, testlogtype),
			expected: fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes/%s", testproject, testlocation, testinstance, testlogtype),
		},
		{
			name:      "test2",
			value:     logtypes.NewLogTypeResource(testproject, testlocation, "", ""),
			expectnil: true,
			expected:  "",
		},
		{
			name:      "test3",
			value:     logtypes.NewLogTypeResource(testproject, "", "", ""),
			expectnil: true,
			expected:  "",
		},
		{
			name:      "test4",
			value:     logtypes.NewLogTypeResource("", "", "", ""),
			expectnil: true,
			expected:  "",
		},
		{
			name:     "No logtype",
			value:    logtypes.NewLogTypeResource(testproject, testlocation, testinstance, ""),
			expected: fmt.Sprintf("projects/%s/locations/%s/instances/%s/logTypes", testproject, testlocation, testinstance),
		},
	}
	for _, tt := range tt {
		if tt.expectnil {
			assert.Nil(t, tt.value)
			continue
		}
		assert.Equal(t, tt.expected, tt.value.Name.String())
	}
}

func TestLogTypeMethods(t *testing.T) {
	testproject := "testproject"
	testlocation := "us"
	testinstance := "testinstance"
	testlogtype := "WINEVTLOG"
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
			name:               "Test Get Method",
			value:              createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, testlogtype), "get", tu),
			expectedHttpMethod: "GET",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes/%s", testproject, testlocation, testinstance, testlogtype),
			expectedQuery:      "",
			expectedBody:       nil,
		},
		{
			name:       "Test Get Method (missing log type)",
			expectFail: true,
			value:      createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, ""), "get", tu),
		},
		{
			name:               "Test List Method (no query)",
			value:              createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, ""), "list", tu, "", ""),
			expectedHttpMethod: "GET",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes", testproject, testlocation, testinstance),
			expectedQuery:      "",
			expectedBody:       nil,
		},
		{
			name:               "Test List Method (with query)",
			value:              createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, ""), "list", tu, "100", "abcdefg"),
			expectedHttpMethod: "GET",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes", testproject, testlocation, testinstance),
			expectedQuery:      "pageSize=100&pageToken=abcdefg",
			expectedBody:       nil,
		},
		{
			name:               "Test List Method (with log type)",
			value:              createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, testlogtype), "list", tu, "100", "abcdefg"),
			expectedHttpMethod: "GET",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes", testproject, testlocation, testinstance),
			expectedQuery:      "pageSize=100&pageToken=abcdefg",
			expectedBody:       nil,
		},
		{
			name:       "Test RunParser Method (no cbn)",
			value:      createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, testlogtype), "runParser", tu, []byte(""), []byte(""), [][]byte{[]byte("")}, true),
			expectFail: true,
		},
		{
			name:               "Test RunParser Method (with cbn)",
			value:              createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, testlogtype), "runParser", tu, []byte("test cbn parser"), []byte(""), [][]byte{[]byte("")}, true),
			expectedHttpMethod: "POST",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes/%s:runParser", testproject, testlocation, testinstance, testlogtype),
			expectedQuery:      "",
			expectedBody: map[string]interface{}{
				"parser": map[string]interface{}{
					"cbn": "dGVzdCBjYm4gcGFyc2Vy",
				},
				"log":              []interface{}{""},
				"statedumpAllowed": true,
			},
		},
		{
			name:               "Test RunParser Method (with log)",
			value:              createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, testlogtype), "runParser", tu, []byte("test cbn parser"), []byte(""), [][]byte{[]byte("log message")}, true),
			expectedHttpMethod: "POST",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes/%s:runParser", testproject, testlocation, testinstance, testlogtype),
			expectedQuery:      "",
			expectedBody: map[string]interface{}{
				"parser": map[string]interface{}{
					"cbn": "dGVzdCBjYm4gcGFyc2Vy",
				},
				"log":              []interface{}{"bG9nIG1lc3NhZ2U="},
				"statedumpAllowed": true,
			},
		},
		{
			name:               "Test RunParser Method (with extension)",
			value:              createRequest(logtypes.NewLogTypeResource(testproject, testlocation, testinstance, testlogtype), "runParser", tu, []byte("test cbn parser"), []byte("test cbn snippet"), [][]byte{[]byte("log message")}, true),
			expectedHttpMethod: "POST",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s/logTypes/%s:runParser", testproject, testlocation, testinstance, testlogtype),
			expectedQuery:      "",
			expectedBody: map[string]interface{}{
				"parser": map[string]interface{}{
					"cbn": "dGVzdCBjYm4gcGFyc2Vy",
				},
				"parserExtension": map[string]interface{}{
					"cbnSnippet": "dGVzdCBjYm4gc25pcHBldA==",
				},
				"log":              []interface{}{"bG9nIG1lc3NhZ2U="},
				"statedumpAllowed": true,
			},
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

func createRequest(resource *logtypes.LogTypeResource, methodType string, u *url.URL, options ...interface{}) *http.Request {
	if resource == nil {
		return nil
	}
	var err error
	var req *http.Request
	switch methodType {
	case "get":
		req, err = resource.Get(u)
	case "list":
		pageSize := options[0].(string)
		pageToken := options[1].(string)
		req, err = resource.List(u, pageSize, pageToken)
	case "runParser":
		cbn := options[0].([]byte)
		cbnSnippet := options[1].([]byte)
		logs := options[2].([][]byte)
		statedumpAllowed := options[3].(bool)
		req, err = resource.RunParser(u, cbn, cbnSnippet, logs, statedumpAllowed)
	default:
		return nil
	}
	if err != nil {
		return nil
	}
	return req
}
