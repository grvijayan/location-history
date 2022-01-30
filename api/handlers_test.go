package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddLocationData(t *testing.T) {
	values := map[string]float64{"lat": 12.34, "lng": 56.78}
	reqBody, _ := json.Marshal(values)
	request, err := http.NewRequest("POST", "/location/", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"order_id": "abc123",
	}

	request = mux.SetURLVars(request, vars)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AddLocationData)
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "", responseRecorder.Body.String())
}

func TestAddLocationDataInvalidId(t *testing.T) {
	values := map[string]float64{"lat": 12.34, "lng": 56.78}
	reqBody, _ := json.Marshal(values)
	request, err := http.NewRequest("POST", "/location/", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AddLocationData)
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "Bad request: Invalid order_id", responseRecorder.Body.String())
}

// Further test cases
// Test delete functionality - do a get after deleting an order_id - for valid and invalid IDs
// Test get functionality - with/without max - for valid and invalid IDs
