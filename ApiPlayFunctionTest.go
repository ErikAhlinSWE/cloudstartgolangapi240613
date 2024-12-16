package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock för randomizeSelection för att kontrollera testscenarier
func mockRandomizeSelection(selection string) func() string {
	return func() string {
		return selection
	}
}

func TestAPIPlay(t *testing.T) {
	// Sätt Gin till test-läge
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name               string
		yourSelection      string
		computerSelection  string
		expectedWinner     string
		expectedStatusCode int
	}{

		// Scenarion där spelaren vinner
		{
			name:               "Rock beats Scissors - Player Wins",
			yourSelection:      "ROCK",
			computerSelection:  "SCISSOR",
			expectedWinner:     "You",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Scissors beats Paper - Player Wins",
			yourSelection:      "SCISSOR",
			computerSelection:  "PAPER",
			expectedWinner:     "You",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Paper beats Rock - Player Wins",
			yourSelection:      "PAPER",
			computerSelection:  "ROCK",
			expectedWinner:     "You",
			expectedStatusCode: http.StatusOK,
		},

		// Scenarion där datorn vinner
		{
			name:               "Rock beats Scissors - Computer Wins",
			yourSelection:      "SCISSOR",
			computerSelection:  "ROCK",
			expectedWinner:     "Computer",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Paper beats Rock - Computer Wins",
			yourSelection:      "ROCK",
			computerSelection:  "PAPER",
			expectedWinner:     "Computer",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Scissors beats Paper - Computer Wins",
			yourSelection:      "PAPER",
			computerSelection:  "SCISSOR",
			expectedWinner:     "Computer",
			expectedStatusCode: http.StatusOK,
		},

		// Oavgjort-scenarion
		{
			name:               "Rock vs Rock - Tie",
			yourSelection:      "ROCK",
			computerSelection:  "ROCK",
			expectedWinner:     "Tie",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Paper vs Paper - Tie",
			yourSelection:      "PAPER",
			computerSelection:  "PAPER",
			expectedWinner:     "Tie",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Scissors vs Scissors - Tie",
			yourSelection:      "SCISSOR",
			computerSelection:  "SCISSOR",
			expectedWinner:     "Tie",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Tillfälligt ersätt RandomizeSelection med en mock
			originalRandomizer := RandomizeSelection

			// Tillfälligt ersätt
			RandomizeSelection = func() string {
				return tc.computerSelection
			}

			// Återställ efter testet
			defer func() {
				RandomizeSelection = originalRandomizer
			}()

			// Skapa en testsvarsinspelare
			w := httptest.NewRecorder()

			// Skapa en Gin-kontext
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodGet, "/?yourSelection="+tc.yourSelection, nil)

			// Anropa funktionen
			apiPlay(c)

			// Kontrollera svarsstatuskod
			assert.Equal(t, tc.expectedStatusCode, w.Code)

			// Kontrollera JSON-svar
			expectedResponse := gin.H{
				"winner":            tc.expectedWinner,
				"yourSelection":     tc.yourSelection,
				"computerSelection": tc.computerSelection,
			}

			// Parse JSON-svaret
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Verifiera JSON-svaret
			assert.Equal(t, expectedResponse["winner"], response["winner"])
			assert.Equal(t, expectedResponse["yourSelection"], response["yourSelection"])
			assert.Equal(t, expectedResponse["computerSelection"], response["computerSelection"])
		})
	}
}
