package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/arpsch/players/client"
	"github.com/arpsch/players/players"
)

func main() {
	var server *string
	var teams *string

	server = flag.String("server", "https://api-origin.onefootball.com", "upstream server url")
	teams = flag.String("teams", "Germany,France",
		"comma separated team names to list players from")

	flag.Parse()

	teamsSlice := strings.Split(*teams, ",")
	doMain(*server, teamsSlice)
}

func doMain(server string, teams []string) {

	url, err := url.Parse(server)
	if err != nil {
		log.Fatalf("Failed to parse the base url: %v\n", err)
	}

	cPlayer, err := client.New(url, &http.Client{})
	if err != nil {
		log.Fatalf("Failed to create client instance: %v\n", err)
	}

	playerApp, err := players.NewPlayerApp(teams, cPlayer)
	if err != nil {
		log.Fatalf("Failed to create player app instance: %v\n", err)
	}

	i := 0
	for _, p := range playerApp.GetPlayersInformation() {
		i++
		fmt.Printf("%d. %s; %s; %s; %s\n", i, p.ID, p.Name, p.Age, strings.Join(p.Teams, ", "))
	}
}
