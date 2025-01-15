// main_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockGameStorage struct{}

func (m *MockGameStorage) SaveGame(yourSelection, computerSelection, winner string) error {
	return nil
}

func TestApiPlay(t *testing.T) {
	gameStorage = &MockGameStorage{}
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/play", apiPlay)

	tests := []struct {
		name           string
		yourSelection  string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid ROCK input",
			yourSelection:  "ROCK",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				response := w.Body.String()
				// Extrahera bara vinnarnamnet (texten före JSON)
				winner := strings.Split(response, "{")[0]
				assert.Contains(t, []string{"You", "Computer", "Tie"}, winner)
			},
		},
		{
			name:           "Valid PAPER input",
			yourSelection:  "PAPER",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				response := w.Body.String()
				// Extrahera bara vinnarnamnet (texten före JSON)
				winner := strings.Split(response, "{")[0]
				assert.Contains(t, []string{"You", "Computer", "Tie"}, winner)
			},
		},
		{
			name:           "Valid SCISSOR input",
			yourSelection:  "SCISSOR",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				response := w.Body.String()
				// Extrahera bara vinnarnamnet (texten före JSON)
				winner := strings.Split(response, "{")[0]
				assert.Contains(t, []string{"You", "Computer", "Tie"}, winner)
			},
		},
		{
			name:           "Lowercase rock input should work",
			yourSelection:  "rock",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				response := w.Body.String()
				winner := strings.Split(response, "{")[0]
				assert.Contains(t, []string{"You", "Computer", "Tie"}, winner)
			},
		},
		{
			name:           "Empty selection",
			yourSelection:  "",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "Invalid selection", w.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/play?yourSelection="+tt.yourSelection, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
