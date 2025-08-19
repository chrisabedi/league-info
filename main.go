package main

import (
	"fmt"
	"league-info/leagueapi"
	"log"
	"os"
	"os/signal"
	"regexp"
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

	command := LpCommandRegex(m.Content)

	if len(command) == 3 {

		response, _ := GetRankedTierInfo(command[1], command[2])
		var output string
		if response != nil {
			output = fmt.Sprintf("`%+v\n`", *response)
		} else {
			output = ":warning:  usage: `!lp <summoner>#tagLine `\n(only set up for NA1)"
		}
		s.ChannelMessageSend(m.ChannelID, output)

	}

	if m.Content == "!ping" {

		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func LpCommandRegex(content string) []string {

	reg := regexp.MustCompile(`!lp ([a-zA-Z0-9 _]{3,16})#([a-zA-Z0-9]{2,5})`)

	matches := reg.FindStringSubmatch(content)

	return matches
}

func GetPUUID(gameName string, tagLine string) (string, error) {
	ApiToken := os.Getenv("LEAGUE_API_TOKEN")

	client := leagueapi.NewClient("https://americas.api.riotgames.com", 10*time.Second, ApiToken, map[string]string{})

	return client.GetPUUID(gameName, tagLine)
}

func GetRankedTierInfo(gameName string, tagLine string) (*leagueapi.LeagueQueue, error) {
	ApiToken := os.Getenv("LEAGUE_API_TOKEN")

	puuid, _ := GetPUUID(gameName, tagLine)
	client := leagueapi.NewClient("https://na1.api.riotgames.com", 10*time.Second, ApiToken, map[string]string{})

	information, _ := client.GetRankedTierInfo(gameName, tagLine, puuid)

	return information, nil
}
