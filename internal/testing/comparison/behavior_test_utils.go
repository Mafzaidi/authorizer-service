package comparison

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// HTTPRequest represents a test HTTP request
type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    interface{}
}

// HTTPResponse represents a test HTTP response
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       map[string]interface{}
}

// ExecuteRequest executes an HTTP request against an Echo instance
func ExecuteRequest(e *echo.Echo, req HTTPRequest) (*HTTPResponse, error) {
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	httpReq := httptest.NewRequest(req.Method, req.Path, bodyReader)
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httpReq)

	response := &HTTPResponse{
		StatusCode: rec.Code,
		Headers:    make(map[string]string),
	}

	for key := range rec.Header() {
		response.Headers[key] = rec.Header().Get(key)
	}

	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &response.Body); err != nil {
			return nil, err
		}
	}

	return response, nil
}

// CompareResponses compares two HTTP responses for equality
func CompareResponses(t *testing.T, expected, actual *HTTPResponse) bool {
	t.Helper()

	if expected.StatusCode != actual.StatusCode {
		t.Errorf("Status code mismatch: expected %d, got %d", expected.StatusCode, actual.StatusCode)
		return false
	}

	// Compare critical headers (Content-Type, Authorization, etc.)
	criticalHeaders := []string{"Content-Type", "Authorization"}
	for _, header := range criticalHeaders {
		if expected.Headers[header] != actual.Headers[header] {
			t.Errorf("Header %s mismatch: expected %s, got %s", 
				header, expected.Headers[header], actual.Headers[header])
			return false
		}
	}

	// Compare body structure
	if !compareJSONStructure(t, expected.Body, actual.Body) {
		return false
	}

	return true
}

// compareJSONStructure compares the structure of two JSON objects
func compareJSONStructure(t *testing.T, expected, actual map[string]interface{}) bool {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Body structure mismatch: expected %d fields, got %d fields", 
			len(expected), len(actual))
		return false
	}

	for key := range expected {
		if _, exists := actual[key]; !exists {
			t.Errorf("Missing field in response body: %s", key)
			return false
		}
	}

	return true
}

// ResponseMatcher is a function that checks if a response matches expected criteria
type ResponseMatcher func(*HTTPResponse) bool

// StatusCodeMatcher creates a matcher for status code
func StatusCodeMatcher(expectedCode int) ResponseMatcher {
	return func(resp *HTTPResponse) bool {
		return resp.StatusCode == expectedCode
	}
}

// BodyFieldMatcher creates a matcher for a specific body field
func BodyFieldMatcher(field string, expectedValue interface{}) ResponseMatcher {
	return func(resp *HTTPResponse) bool {
		if resp.Body == nil {
			return false
		}
		actualValue, exists := resp.Body[field]
		if !exists {
			return false
		}
		return actualValue == expectedValue
	}
}

// AllMatchers combines multiple matchers with AND logic
func AllMatchers(matchers ...ResponseMatcher) ResponseMatcher {
	return func(resp *HTTPResponse) bool {
		for _, matcher := range matchers {
			if !matcher(resp) {
				return false
			}
		}
		return true
	}
}
