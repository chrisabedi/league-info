package leagueapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Client struct {
	BaseUrl    string
	HTTPClient *http.Client
	Headers    map[string]string
	Token      string
}

type Participant struct {
	RiotIdGameName    string `json:"riotIdGameName"`
	ChampionName      string `json:"championName"`
	Win               bool   `json:"win"`
	Puuid             string `json:"puuid"`
	DangerPings       int    `json:"dangerPings"`
	GetBackPings      int    `json:"getBackPings"`
	CommandPings      int    `json:"CommandPings"`
	HoldPings         int    `json:"holdPings"`
	EnemyMissingPings int    `json:"enemyMissingPings"`
	EnemyVisionPings  int    `json:"enemyVisionPings"`
	OnMyWayPings      int    `json:"onMyWayPings"`
}

func NewParticipant(riotIdGameName string, championName string, win bool, puuid string, dangerPings int, getBackPings int, commandPings int, holdPings int, enemyMissingPings int, enemyVisionPings int, onMyWayPings int) *Participant {
	return &Participant{
		RiotIdGameName:    riotIdGameName,
		ChampionName:      championName,
		Win:               win,
		Puuid:             puuid,
		DangerPings:       dangerPings,
		GetBackPings:      getBackPings,
		CommandPings:      commandPings,
		HoldPings:         holdPings,
		EnemyMissingPings: enemyMissingPings,
		EnemyVisionPings:  enemyVisionPings,
		OnMyWayPings:      onMyWayPings,
	}
}

func (p *Participant) OutputMarkDown() string {

	output := fmt.Sprintf("| %s | value |\n", p.RiotIdGameName)

	v := reflect.ValueOf(*p)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if field.Name != "Puuid" && field.Name != "RiotIdGameName" {
			output += fmt.Sprintf("| %s | %v |\n", field.Name, value)
		}
	}

	return output
}

type LastMatchInfo struct {
	Info Info `json:"info"`
}

type Info struct {
	Participants []Participant `json:"participants"`
}

type PUUIDResponse struct {
	Puuid    string `json:puuid`
	GameName string `json:gameName`
	TagLine  string `json:tagLine`
}

func NewClient(baseUrl string, timeout time.Duration, token string, headers map[string]string) *Client {
	return &Client{
		BaseUrl: baseUrl,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		Headers: headers,
	}
}

func (c *Client) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var requestBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling body: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseUrl+path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set default headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	// Set Content-Type if body is present
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}

// Get sends a GET request
func (c *Client) Get(path string) (*http.Response, error) {
	return c.makeRequest(http.MethodGet, path, nil)
}

func (c *Client) GetPUUID(gameName string, tagLine string) (string, error) {

	uri := fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s?api_key=%s", gameName, tagLine, c.Token)

	resp, err := c.Get(uri)
	if err != nil {
		panic(err)
	}

	var newPuuidResp PUUIDResponse
	jerr := json.NewDecoder(resp.Body).Decode(&newPuuidResp)
	if jerr != nil {
		return "fail", jerr
	}

	return newPuuidResp.Puuid, err
}

func (c *Client) GetLastRankedMatchId(gameName string, tagLine string) ([2]string, error) {

	puuid, _ := c.GetPUUID(gameName, tagLine)

	uri := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids?start=0&count=1&api_key=%s", puuid, c.Token)
	resp, err := c.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return [2]string{string(body), puuid}, err
}

func (c *Client) GetLastRankedMatchInfo(gameName string, tagLine string) (*Participant, error) {

	values, _ := c.GetLastRankedMatchId(gameName, tagLine)

	matchId := strings.Trim(string(values[0]), "[]\"")
	uri := fmt.Sprintf("/lol/match/v5/matches/%s?api_key=%s", matchId, c.Token)

	resp, err := c.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	jsonData, _ := io.ReadAll(resp.Body)
	var lastMatchInfo LastMatchInfo

	jerr := json.Unmarshal(jsonData, &lastMatchInfo)
	if jerr != nil {
		fmt.Println("Error unmarshalling:", err)
		return nil, jerr
	}

	for _, participant := range lastMatchInfo.Info.Participants {

		if participant.RiotIdGameName == gameName {
			return &participant, nil
		}
	}

	return nil, nil
}
