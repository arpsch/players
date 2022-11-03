package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/arpsch/players/model"
	"github.com/stretchr/testify/assert"
)

// return mock http server returning status code 'status'
func NewMockServer(status int, body []byte) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.WriteHeader(status)
		w.Write(body)
	}))
	return srv
}

func TestNewClient(t *testing.T) {

	c, err := New(nil, nil)
	assert.Nil(t, c)
	assert.NotNil(t, err)

	c, err = New(&url.URL{Path: "/test"}, &http.Client{})
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestExtractPlayerInfo(t *testing.T) {

	cases := []struct {
		name   string
		input  Response
		output []model.Player
		length int
	}{
		{
			name: "valid case",
			input: Response{
				Data: Data{
					Team: Team{
						Name: "Germany",
						ID:   0,
						Players: []Player{
							{
								ID:        "1244",
								Name:      "A",
								FirstName: "fn",
								LastName:  "ln",
							},
							{
								ID:        "1245",
								Name:      "B",
								FirstName: "fn1",
								LastName:  "ln1",
							},
						},
					},
				},
			},
			length: 2,
			output: []model.Player{
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
		},
		{
			name: "missing name",
			input: Response{
				Data: Data{
					Team: Team{
						Name: "Germany",
						ID:   0,
						Players: []Player{
							{
								ID:        "1244",
								Name:      "",
								FirstName: "fn",
								LastName:  "ln",
							},
							{
								ID:        "1245",
								Name:      "",
								FirstName: "fn1",
								LastName:  "ln1",
							},
						},
					},
				},
			},
			length: 2,
			output: []model.Player{
				{
					ID:    "1244",
					Name:  "fn ln",
					Teams: []string{"Germany"},
				},
				{
					ID:    "1245",
					Name:  "fn1 ln1",
					Teams: []string{"Germany"},
				},
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {

			players := extractPlayerInfo(tc.input)
			assert.NotNil(t, players)

			assert.Equal(t, tc.length, len(players))
			assert.Equal(t, tc.output, players)
		})
	}
}

func TestGetTeamInformationByID(t *testing.T) {
	cases := []struct {
		name   string
		input  Response
		output []model.Player
		status int
		length int
	}{
		{
			name:   "valid case",
			status: 200,
			input: Response{
				Data: Data{
					Team: Team{
						Name: "Germany",
						ID:   0,
						Players: []Player{
							{
								ID:        "1244",
								Name:      "A",
								FirstName: "fn",
								LastName:  "ln",
							},
							{
								ID:        "1245",
								Name:      "B",
								FirstName: "fn1",
								LastName:  "ln1",
							},
						},
					},
				},
			},
			length: 2,
			output: []model.Player{
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
		},
		{
			name:   "missing name",
			status: 200,
			input: Response{
				Data: Data{
					Team: Team{
						Name: "Germany",
						ID:   0,
						Players: []Player{
							{
								ID:        "1244",
								Name:      "",
								FirstName: "fn",
								LastName:  "ln",
							},
							{
								ID:        "1245",
								Name:      "",
								FirstName: "fn1",
								LastName:  "ln1",
							},
						},
					},
				},
			},
			length: 2,
			output: []model.Player{
				{
					ID:    "1244",
					Name:  "fn ln",
					Teams: []string{"Germany"},
				},
				{
					ID:    "1245",
					Name:  "fn1 ln1",
					Teams: []string{"Germany"},
				},
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {

			inpJson, err := json.Marshal(tc.input)
			assert.Nil(t, err)

			s := NewMockServer(tc.status, inpJson)

			url, err := url.Parse(s.URL)
			assert.Nil(t, err)

			c, err := New(url, s.Client())
			assert.Nil(t, err)
			assert.NotNil(t, c)

			team, players, err := c.GetTeamInformationByID(context.Background(), strconv.Itoa(tc.input.Data.Team.ID))
			assert.Nil(t, err)
			assert.Equal(t, team, tc.input.Data.Team.Name)
			assert.Equal(t, players, tc.output)

			s.CloseClientConnections()
			s.Close()
		})
	}
}
