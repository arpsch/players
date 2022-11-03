package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/arpsch/players/model"
)

// ClientRunner is an interface of teams api client
type ClientRunner interface {
	GetTeamInformationByID(ctx context.Context, teamID string) (string, []model.Player, error)
}

const (
	teamUri        = "/score-one-proxy/api/teams/en/:teamID.json"
	defaultTimeout = time.Second * 10
)

// Client represents the HTTP client object
type Client struct {
	httpRunner *http.Client

	BaseURL *url.URL
}

const (
	MAX_IDLE_CONN          = 100
	MAX_CONN_PER_HOST      = 100
	MAX_IDLE_CONN_PER_HOST = 100
)

// New retrns the new client instance
func New(baseURL *url.URL, runner *http.Client) (*Client, error) {

	if baseURL == nil || runner == nil {
		return nil, errors.New("invalid paramters")
	}

	// set default  time out to 10 s
	runner.Timeout = defaultTimeout
	// don't use the default transport, instead custom
	// with 100 max idle connections
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = MAX_IDLE_CONN
	// allow max conn and idle conn per host
	t.MaxConnsPerHost = MAX_CONN_PER_HOST
	t.MaxIdleConnsPerHost = MAX_IDLE_CONN_PER_HOST

	lrt, err := NewLogRoundTripper(t, os.Stdout)
	if err != nil {
		return nil, err
	}

	runner.Transport = lrt

	return &Client{
		BaseURL:    baseURL,
		httpRunner: runner,
	}, nil
}

// extractPlayerInfo extracts the players information from the
// rest response
func extractPlayerInfo(res Response) []model.Player {

	players := []model.Player{}

	for _, p := range res.Data.Team.Players {
		// few entries didn't carry neither name nor first, last name
		if p.Name == "" && (p.FirstName == "" && p.LastName == "") {
			continue
		}
		// prepare name from fn and ln when present
		if p.Name == "" && (p.FirstName != "" && p.LastName != "") {
			p.Name = p.FirstName + " " + p.LastName
		}

		player := model.Player{
			ID:    p.ID,
			Name:  p.Name,
			Age:   p.Age,
			Teams: []string{res.Data.Team.Name},
		}
		players = append(players, player)
	}
	return players
}

// GetTeamInformationByID returns the playes information along with the team name
func (c *Client) GetTeamInformationByID(ctx context.Context, teamID string) (string, []model.Player, error) {
	repl := strings.NewReplacer(":teamID", teamID)
	uri := repl.Replace(teamUri)

	rel := &url.URL{Path: uri}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := c.httpRunner.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", nil, err
	}
	players := extractPlayerInfo(response)
	return response.Data.Team.Name, players, nil
}
