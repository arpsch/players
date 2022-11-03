package players

import (
	"context"
	"errors"
	"strconv"

	"log"

	"github.com/arpsch/players/client"
	"github.com/arpsch/players/model"
	"github.com/arpsch/players/utils"
)

// PlayersApp represents the playersApp object
// central object of the application
type PlayersApp struct {
	teamsMap map[string]struct{}
	cPlayer  client.ClientRunner
}

// NewPlayerApp constructs the PlayerApp provided initial values
func NewPlayerApp(teams []string, cPlayer client.ClientRunner) (*PlayersApp, error) {

	if len(teams) == 0 || cPlayer == nil {

		return nil, errors.New("invalid paramters")
	}

	teamsMap := make(map[string]struct{})

	for _, c := range teams {
		teamsMap[c] = struct{}{}
	}

	return &PlayersApp{
		teamsMap: teamsMap,
		cPlayer:  cPlayer,
	}, nil
}

// GetPlayersInformation returns the player informats for every team
// that are configured
func (pa *PlayersApp) GetPlayersInformation() map[string]model.Player {

	playersMap := make(map[string]model.Player)
	teamCount := 0

	for i := 0; ; i++ {
		c, players, err := pa.cPlayer.GetTeamInformationByID(context.Background(), strconv.Itoa(i))
		if err != nil {
			log.Printf("Failed to fetch team information for team id %d\n", i)
			// there could be errors for the given team id, but continue with other
			continue
		}

		if _, ok := pa.teamsMap[c]; ok {
			teamCount++
			for _, p := range players {
				if v, ok := playersMap[p.ID]; !ok {
					playersMap[p.ID] = p
				} else {
					if utils.ContainsString(c, v.Teams) != true {
						v.Teams = append(v.Teams, c)
						playersMap[v.ID] = v
					}
				}
			}
		}

		if len(pa.teamsMap) == teamCount {
			break
		}
	}
	return playersMap
}
