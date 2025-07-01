package leagueapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseUrl    string
	HTTPClient *http.Client
	Headers    map[string]string
	Token      string
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

type PUUIDResponse struct {
	Puuid    string `json:puuid`
	GameName string `json:gameName`
	TagLine  string `json:tagLine`
}

func (c *Client) GetPUUID(gameName string, tagLine string) (string, error) {

	uri := fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s?api_key=%s", gameName, tagLine, c.Token)

	resp, err := c.Get(uri)
	if err != nil {
		panic(err)
	}

	var newPuuid PUUIDResponse
	jerr := json.NewDecoder(resp.Body).Decode(&newPuuid)
	if jerr != nil {
		return "fail", jerr
	}

	return newPuuid.Puuid, err
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

func (c *Client) GetLastRankedMatchInfo(gameName string, tagLine string) {

	values, _ := c.GetLastRankedMatchId(gameName, tagLine)

	matchId := strings.Trim(string(values[0]), "[]\"")
	uri := fmt.Sprintf("/lol/match/v5/matches/%s?api_key=%s", matchId, c.Token)

	fmt.Println(uri)
	// resp, err := c.Get(uri)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(body))
	// return string(body), err
}
