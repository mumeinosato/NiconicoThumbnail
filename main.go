package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	token := getToken()
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	session.AddHandler(Handler)

	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	_ = session.Close()
}

func getToken() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return ""
	}

	token := os.Getenv("TOKEN")
	return token
}

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*
		https://www.nicovideo.jp/watch/aaaabbbbb
		↓
		https://www.nicovideo.gay/watch/aaaabbbbb

		delete: embed
		send: replaced url

		regular expression: https://www\.nicovideo\.jp/watch/([a-zA-Z0-9]+)
	*/

	if m.Author.ID == s.State.User.ID {
		return
	}

	re := regexp.MustCompile("https://www\\.nicovideo\\.jp/watch/([a-zA-Z0-9]+)")

	if re.MatchString(m.Content) {
		fmt.Println(m.Content)

		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: m.ChannelID,
			ID:      m.ID,
			Flags:   discordgo.MessageFlagsSupressEmbeds,
		})

		if err != nil {
			fmt.Println("Error editing message:", err)
		}

		replaced := re.ReplaceAllString(m.Content, "https://www.nicovideo.gay/watch/$1")

		_, err = s.ChannelMessageSend(m.ChannelID, replaced)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}

}
