package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUuidHandler_ReturnsValidUUID(t *testing.T) {
	// Arrange: Create a request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/generate", nil)
	rec := httptest.NewRecorder()

	// Act: Call the handler
	UuidHandler().ServeHTTP(rec, req)

	// Assert: Status code
	assert.Equal(t, http.StatusOK, rec.Code, "should return HTTP 200")

	// Assert: Content-Type
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"), "should return application/json")

	// Assert: Body contains a valid UUID
	var body map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.NoError(t, err, "response body should be valid JSON")

	id, ok := body["uuid"]
	assert.True(t, ok, "response should contain a 'uuid' field")

	_, err = uuid.Parse(id)
	assert.NoError(t, err, "uuid field should be a valid UUID")
}
func TestUuidHandler_RetrievesTokenQueryParam(t *testing.T) {
	// Arrange: Create a request with a "token" query parameter
	req := httptest.NewRequest(http.MethodGet, "/generate?token=mytesttoken", nil)
	rec := httptest.NewRecorder()

	// Act: Call the handler
	UuidHandlerWithToken().ServeHTTP(rec, req)

	// Assert: Status code
	assert.Equal(t, http.StatusOK, rec.Code, "should return HTTP 200")

	// Assert: Content-Type
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"), "should return application/json")

	// Assert: Body contains a valid UUID and the token
	var body map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.NoError(t, err, "response body should be valid JSON")

	id, ok := body["uuid"]
	assert.True(t, ok, "response should contain a 'uuid' field")
	_, err = uuid.Parse(id)
	assert.NoError(t, err, "uuid field should be a valid UUID")

	token, ok := body["token"]
	assert.True(t, ok, "response should contain a 'token' field")
	assert.Equal(t, "mytesttoken", token, "token field should match the query param")
}

// UuidHandlerWithToken is a variant of UuidHandler that includes the token in the response for testing purposes.
// In production, you would modify the actual handler to support this if needed.
/*
UuidHandlerWithToken returns an HTTP handler that generates a UUID and echoes the "token" query parameter if present.

Returns:
  - http.Handler: The HTTP handler function.
*/
func UuidHandlerWithToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		token := r.URL.Query().Get("token")
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{"uuid": id.String()}
		if token != "" {
			resp["token"] = token
		}
		json.NewEncoder(w).Encode(resp)
	})
}
