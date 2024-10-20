package instances_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/calebryant/chronicle-api/resources/instances"
	"github.com/stretchr/testify/assert"
)

func TestNewInstanceResource(t *testing.T) {
	testproject := "testproject"
	testlocation := "us"
	testinstance := "testinstance"
	tt := []struct {
		name      string
		expectnil bool
		value     *instances.InstanceResource
		expected  string
	}{
		{
			name:     "test1",
			value:    instances.NewInstanceResource(testproject, testlocation, testinstance),
			expected: fmt.Sprintf("projects/%s/locations/%s/instances/%s", testproject, testlocation, testinstance),
		},
		{
			name:      "test2",
			expectnil: true,
			value:     instances.NewInstanceResource(testproject, testlocation, ""),
			expected:  "",
		},
		{
			name:      "test3",
			expectnil: true,
			value:     instances.NewInstanceResource(testproject, "", ""),
			expected:  "",
		},
		{
			name:      "test$",
			expectnil: true,
			value:     instances.NewInstanceResource("", "", ""),
			expected:  "",
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

func TestGetInstance(t *testing.T) {
	testproject := "testproject"
	testlocation := "us"
	testinstance := "testinstance"
	tu, _ := url.Parse("https://test.local")
	tt := []struct {
		name               string
		value              *http.Request
		expectnil          bool
		expectedHttpMethod string
		expectedUrlPath    string
		expectedQuery      string
		expectedBody       map[string]interface{}
	}{
		{
			name:               "Test Get Method",
			value:              createGetRequest(instances.NewInstanceResource(testproject, testlocation, testinstance), tu),
			expectedHttpMethod: "GET",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s", testproject, testlocation, testinstance),
			expectedQuery:      "",
			expectedBody:       map[string]interface{}{},
		},
		{
			name:               "Test Get Method (no instance)",
			value:              createGetRequest(instances.NewInstanceResource(testproject, testlocation, ""), tu),
			expectnil:          true,
			expectedHttpMethod: "GET",
			expectedUrlPath:    fmt.Sprintf("/projects/%s/locations/%s/instances/%s", testproject, testlocation, testinstance),
			expectedQuery:      "",
			expectedBody:       map[string]interface{}{},
		},
	}
	for _, tt := range tt {
		var bodyBytes []byte
		if tt.expectnil == true {
			assert.Nil(t, tt.value)
			continue
		}
		tt.value.Body.Read(bodyBytes)
		parsedBody := make(map[string]interface{})
		json.Marshal(parsedBody)
		assert.Equal(t, tt.expectedUrlPath, tt.value.URL.Path)
		assert.Equal(t, tt.expectedBody, parsedBody)
		assert.Equal(t, tt.expectedQuery, tt.value.URL.Query().Encode())
		assert.Equal(t, tt.expectedHttpMethod, tt.value.Method)
	}
}

func createGetRequest(i *instances.InstanceResource, u *url.URL) *http.Request {
	if i == nil {
		return nil
	}
	req, err := i.Get(u)
	if err != nil {
		return nil
	}
	return req
}
