package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	watchGames()
}

// fancy logic to avoid too many calls?
func watchGames() {
	for {
		games, err := getGames()
		if err != nil {
			log.Printf("error retrieving games")
		}

		for _, g := range games.Games {
			if g.Active {
				activateGameUrl := "http://localhost:8084/activate"
				_, err := http.Get(fmt.Sprintf("%s/%s", activateGameUrl, g.GameCode())) // don't rely on idempotence?
				if err != nil {
					log.Printf("something went wrong activating %s: %s\n", g.GameCode(), err)
				}
				fmt.Printf("%s activated\n", g.GameCode())
			}
		}
		time.Sleep(1 * time.Minute) // better way to do this for sure
	}
}

type Games struct {
	Games []Game `json:"games"`
}

type Game struct {
	Id           string    `json:"gameId"`
	StartTime    time.Time `json:"startTimeUTC"`
	VisitingTeam Team      `json:"vTeam"`
	HomeTeam     Team      `json:"hTeam"`
	Active       bool      `json:"isGameActivated"`
}

type Team struct {
	Id      string `json:"teamId"`
	TriCode string `json:"triCode"`
}

func (g Game) GameCode() string {
	return fmt.Sprintf("%s%s", g.VisitingTeam.TriCode, g.HomeTeam.TriCode)
}

func getGames() (Games, error) {
	resp, err := http.Get("http://localhost:8080/games")
	if err != nil {
		log.Printf("error retrieving games %s\n", err)
		return Games{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var games Games
	for dec.More() {
		err := dec.Decode(&games)
		if err != nil {
			log.Printf("error decoding games %s\n", err)
			return Games{}, err
		}
	}

	return games, nil
}
