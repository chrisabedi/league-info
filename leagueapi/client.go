package leagueapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseUrl    string
	HTTPClient *http.Client
	Headers    map[string]string
	Token      string
}

type LeagueQueue struct {
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

type PUUIDResponse struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
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

func (c *Client) GetRankedTierInfo(gameName string, tagLine string, puuid string) (*LeagueQueue, error) {

	uri := fmt.Sprintf("/lol/league/v4/entries/by-puuid/%s?api_key=%s", puuid, c.Token)
	resp, err := c.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	jsonData, _ := io.ReadAll(resp.Body)
	var leagueQueues []LeagueQueue

	jerr := json.Unmarshal(jsonData, &leagueQueues)
	if jerr != nil {
		fmt.Println("Error unmarshalling:", err)
		return nil, jerr
	}
	for _, leagueQueue := range leagueQueues {

		if leagueQueue.QueueType == "RANKED_SOLO_5x5" {
			return &leagueQueue, nil
		}
	}

	return nil, nil
}
