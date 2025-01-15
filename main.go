package main

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin" // swagger embed files
	// gin-swagger middleware
	"systementor.se/cloudgolangapi/data"
	docs "systementor.se/cloudgolangapi/docs" // swagger docs
)

var config Config
var theRandom *rand.Rand

// GameStorage definierar ett interface för spellagring
type GameStorage interface {
	SaveGame(yourSelection, computerSelection, winner string) error
}

// Database interface
var gameStorage GameStorage

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	theRandom = rand.New(source)
	// In Prod use the real database implemation
	gameStorage = &data.DBGameStorage{}
}

// @Summary Get start
// @Description Get startpage
// @Success 200 {object} map[string]interface{}
// @Router /swagger/index.html [get]
func start(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Message": "Välkommen, /swagger/index.html#/"})
}

func enableCors(c *gin.Context) {
	(*c).Header("Access-Control-Allow-Origin", "*")
}

// @Summary Get stats
// @Description Get game statistics
// @Success 200 {object} map[string]interface{}
// @Router /api/stats [get]
func apiStats(c *gin.Context) {
	enableCors(c)
	totalGames, wins := data.Stats()
	c.JSON(http.StatusOK, gin.H{"totalGames": totalGames, "wins": wins})
}

// @BasePath /api/v1
// @Summary Play the game
// @Description Play a game of rock, paper, scissor
// @Param yourSelection query string true "Your selection"
// @Success 200 {string} string "Winner"
// @Router /api/play [get]
func apiPlay(c *gin.Context) {
	// Convert to uppercase
	yourSelection := strings.ToUpper(c.Query("yourSelection"))

	// Validat input not empty
	if yourSelection == "" {
		c.String(400, "Invalid selection")
		return
	}

	// Validate correct input
	validMoves := map[string]bool{
		"ROCK":     true,
		"PAPER":    true,
		"SCISSOR":  true,
		"SCISSORS": true, // allow both spellings
	}

	if !validMoves[yourSelection] {
		c.String(http.StatusBadRequest, "Invalid selection. Use ROCK, PAPER, or SCISSORS")
		return
	}

	// Convert SCISSORS --> SCISSOR
	if yourSelection == "SCISSORS" {
		yourSelection = "SCISSOR"
	}

	// Computer choice
	moves := []string{"ROCK", "PAPER", "SCISSOR"}
	computerSelection := moves[rand.Intn(len(moves))]

	// Game logic
	var winner string
	if yourSelection == computerSelection {
		winner = "Tie"
	} else if (yourSelection == "ROCK" && computerSelection == "SCISSOR") ||
		(yourSelection == "PAPER" && computerSelection == "ROCK") ||
		(yourSelection == "SCISSOR" && computerSelection == "PAPER") {
		winner = "You"
	} else {
		winner = "Computer"
	}

	c.String(200, winner)
	// Ignore storing for test
	_ = gameStorage.SaveGame(yourSelection, computerSelection, winner)
	c.JSON(http.StatusOK, gin.H{"winner": winner, "yourSelection": yourSelection, "computerSelection": computerSelection})

}

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	readConfig(&config)

	data.InitDatabase(config.Database.File,
		config.Database.Server,
		config.Database.Database,
		config.Database.Username,
		config.Database.Password,
		config.Database.Port)

	router := gin.Default()
	// Swagger docs
	docs.SwaggerInfo.BasePath = "/"

	// Serve the Swagger UI at the root endpoint
	//router.Get("/swagger/*", httpSwagger.Handler(
	//	httpSwagger.URL("http://golangsite1204.chickenkiller.com/swagger/doc.json"), //The url pointing to API definition
	//))
	router.GET("/", start)
	router.GET("/api/play", apiPlay)
	router.GET("/api/stats", apiStats)
	// Swagger setup
	router.Run(":8080")
}
