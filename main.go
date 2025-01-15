package main

import (
	// Import the generated docs

	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"systementor.se/cloudgolangapi/data"
	docs "systementor.se/cloudgolangapi/docs"
)

var config Config
var theRandom *rand.Rand

// GameStorage definierar ett interface för spellagring
type GameStorage interface {
	SaveGame(yourSelection, computerSelection, winner string) error
}

// Använd interface istället för direkt databasanrop
var gameStorage GameStorage

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	theRandom = rand.New(source)
	// I produktion, använd den riktiga databasimplementationen
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

//func RandomizeSelection() string {
//	selections := []string{"ROCK", "PAPER", "SCISSOR"}
//	return selections[theRandom.Intn(len(selections))]
//}

// @BasePath /api/v1
// @Summary Play the game
// @Description Play a game of rock, paper, scissor
// @Param yourSelection query string true "Your selection"
// @Success 200 {string} string "Winner"
// @Router /api/play [get]
func apiPlay(c *gin.Context) {
	// Konvertera input till uppercase
	yourSelection := strings.ToUpper(c.Query("yourSelection"))

	// Validera input först
	if yourSelection == "" {
		c.String(400, "Invalid selection")
		return
	}

	// Konvertera input till uppercase
	yourSelection = strings.ToUpper(yourSelection)

	// Validera input
	validMoves := map[string]bool{
		"ROCK":     true,
		"PAPER":    true,
		"SCISSOR":  true,
		"SCISSORS": true, // Tillåt båda stavningarna
	}

	if !validMoves[yourSelection] {
		c.String(http.StatusBadRequest, "Invalid selection. Use ROCK, PAPER, or SCISSORS")
		return
	}

	// Normalisera SCISSORS/SCISSOR till en konsekvent form
	if yourSelection == "SCISSORS" {
		yourSelection = "SCISSOR"
	}

	// Datorns val
	moves := []string{"ROCK", "PAPER", "SCISSOR"}
	computerSelection := moves[rand.Intn(len(moves))]

	// Bestäm vinnaren
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
	//data.SaveGame(yourSelection, mySelection, winner)
	// Ignorera eventuella fel vid lagring i testerna
	_ = gameStorage.SaveGame(yourSelection, computerSelection, winner)
	c.JSON(http.StatusOK, gin.H{"winner": winner, "yourSelection": yourSelection, "computerSelection": computerSelection})

}

//var RandomizeSelection = func() string {
//	val := theRandom.Intn(3) + 1
//	if val == 1 {
//		return "ROCK"
//	}
//	if val == 2 {
//		return "SCISSOR"
//	}
//	if val == 3 {
//		return "PAPER"
//	}
//	return "ERROR"
//}

func main() {
	//theRandom = rand.New(rand.NewSource(time.Now().UnixNano()))
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", start)
	router.GET("/api/play", apiPlay)
	router.GET("/api/stats", apiStats)
	// Swagger setup
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// router.GET("/api/employee/:id", apiEmployeeById)
	// router.PUT("/api/employee/:id", apiEmployeeUpdateById)
	// router.DELETE("/api/employee/:id", apiEmployeeDeleteById)
	// router.POST("/api/employee", apiEmployeeAdd)

	// router.GET("/api/employees", employeesJson)
	// router.GET("/api/addemployee", addEmployee)
	// router.GET("/api/addmanyemployees", addManyEmployees)
	router.Run(":8080")

}
