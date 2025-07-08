package main

import (
	"fmt"
	"league-info/leagueapi"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DiscordToken := os.Getenv("BOT_TOKEN")

	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	var command []string
	command = PuuidCommandRegex(m.Content)

	if len(command) == 3 {

		response, _ := GetPUUID(command[1], command[2])

		s.ChannelMessageSend(m.ChannelID, response)

	}

	command = LastMatchCommandRegex(m.Content)

	if len(command) == 3 {

		response, _ := GetLastRankedMatch(command[1], command[2])

		s.ChannelMessageSend(m.ChannelID, response)

	}

	command = LastMatchInfoCommandRegex(m.Content)

	if len(command) == 3 {

		response, _ := GetLastRankedMatchInfo(command[1], command[2])

		s.ChannelMessageSend(m.ChannelID, response.OutputMarkDown())

	}
	if m.Content == "!ping" {

		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func PuuidCommandRegex(content string) []string {

	reg := regexp.MustCompile(`!puuid (\w*)#(\w*)`)

	matches := reg.FindStringSubmatch(content)

	return matches
}
func LastMatchCommandRegex(content string) []string {

	reg := regexp.MustCompile(`!lm (\w*)#(\w*)`)

	matches := reg.FindStringSubmatch(content)

	return matches
}

func LastMatchInfoCommandRegex(content string) []string {

	reg := regexp.MustCompile(`!lminfo ([^#]*)#(\w*)`)

	matches := reg.FindStringSubmatch(content)

	return matches
}

func GetPUUID(gameName string, tagLine string) (string, error) {
	ApiToken := os.Getenv("LEAGUE_API_TOKEN")

	client := leagueapi.NewClient("https://americas.api.riotgames.com", 10*time.Second, ApiToken, map[string]string{})

	return client.GetPUUID(gameName, tagLine)
}

func GetLastRankedMatch(gameName string, tagLine string) (string, error) {
	ApiToken := os.Getenv("LEAGUE_API_TOKEN")

	client := leagueapi.NewClient("https://americas.api.riotgames.com", 10*time.Second, ApiToken, map[string]string{})

	values, err := client.GetLastRankedMatchId(gameName, tagLine)
	if err != nil {
		log.Fatal("Error getting last Ranked matchId")
	}

	return strings.Trim(string(values[0]), "[]\""), err
}

func GetLastRankedMatchInfo(gameName string, tagLine string) (*leagueapi.Participant, error) {
	ApiToken := os.Getenv("LEAGUE_API_TOKEN")

	client := leagueapi.NewClient("https://americas.api.riotgames.com", 10*time.Second, ApiToken, map[string]string{})

	information, _ := client.GetLastRankedMatchInfo(gameName, tagLine)

	return information, nil
}
