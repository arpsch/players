package players

import (
	"testing"

	"github.com/arpsch/players/client/mocks"
	"github.com/arpsch/players/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPlayerApp(t *testing.T) {

	c, err := NewPlayerApp(nil, nil)
	assert.Nil(t, c)
	assert.NotNil(t, err)

	c, err = NewPlayerApp([]string{"test"}, &mocks.ClientRunner{})
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestGetPlayersInformation(t *testing.T) {
	cases := []struct {
		name     string
		teamName string
		players  []model.Player
		output   map[string]model.Player
	}{
		{
			name:     "valid case",
			teamName: "Germany",
			players: []model.Player{
				{
					ID:    "1244",
					Name:  "A",
					Teams: []string{"Germany"},
				},
				{
					ID:    "1245",
					Name:  "B",
					Teams: []string{"Germany"},
				},
			},
			output: map[string]model.Player{
				"1244": {
					ID:    "1244",
					Name:  "A",
					Teams: []string{"Germany"},
				},
				"1245": {
					ID:    "1245",
					Name:  "B",
					Teams: []string{"Germany"},
				},
			},
		},
		{
			name:     "missing name",
			teamName: "France",
			players: []model.Player{
				{
					ID:    "1244",
					Name:  "fn ln",
					Teams: []string{"France"},
				},
				{
					ID:    "1245",
					Name:  "fn1 ln1",
					Teams: []string{"France", "Germany"},
				},
			},
			output: map[string]model.Player{
				"1244": {
					ID:    "1244",
					Name:  "fn ln",
					Teams: []string{"France"},
				},
				"1245": {
					ID:    "1245",
					Name:  "fn1 ln1",
					Teams: []string{"France", "Germany"},
				},
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mocks.ClientRunner{}
			playerApp, err := NewPlayerApp([]string{"Germany", "France"}, mockClient)
			assert.Nil(t, err)
			assert.NotNil(t, playerApp)

			mockClient.On("GetTeamInformationByID", mock.AnythingOfType("*context.emptyCtx"),
				mock.AnythingOfType("string")).Return(tc.teamName, tc.players, nil)

			playersMap := playerApp.GetPlayersInformation()
			assert.NotNil(t, playersMap)
			assert.Equal(t, playersMap, tc.output)
		})
	}
}
