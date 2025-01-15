package data

import "time"

type Game struct {
	Id                int
	CreatedAt         time.Time
	Winner            string `json:"winner"`
	YourSelection     string `json:"yourSelection"`
	ComputerSelection string `json:"computerSelection"`
}
