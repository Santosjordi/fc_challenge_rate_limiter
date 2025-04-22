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
