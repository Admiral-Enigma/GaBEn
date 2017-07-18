package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"io"
	"github.com/bwmarrin/discordgo"
)

var (
	token        string
	catsBaseURL  string
)

func init() {
	flag.StringVar(&token, "t", "", "Discord token")
	flag.Parse()
	catsBaseURL = "http://thecatapi.com/api/images/get?format=src&type=gif"
}

func getCats() {
	res, err := http.Get(catsBaseURL)
	if err != nil {
		fmt.Println("Got error while downloading cat,", err)
		return
	}
	defer res.Body.Close()

	file, err := os.Create("./cats/cat.gif")
	if err != nil {
		fmt.Println("Got error while creating cat,", err)
		return
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println("Got error while creating cat,", err)
		return
	}
	file.Close()
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error while connecting to Discord api", err)
		return
	}

	dg.AddHandler(onMessage)
	err = dg.Open()
	if err != nil {
		fmt.Println("Error while connecting", err)
		return
	}

	fmt.Println("Bot is now running!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!cat" {
		getCats()
		cat, err := os.Open("./cats/cat.gif")
		if err != nil {
			fmt.Println(err)
			return
		}
		s.ChannelFileSendWithMessage(m.ChannelID, "CAT!", "cat.gif", cat)
	} else if m.Content == "!bingbingbong" {
		s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=kv5mPHGOp5E")
	} else if m.Content == "!fuckdig" {
		s.ChannelMessageSend(m.ChannelID, "Fuck dig jeg kommer og chopper dig " + m.Author.Username)
	} else if m.Content == "!megaman" {
		s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=DyDxh3fjUv8")
	}
}
