// simple discord bot https://github.com/m3tav3rse/pan-bot
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/m3tav3rse/pan-bot/pkg/crappers"

	"github.com/bwmarrin/discordgo"
)

var (
	version = "unknown"
	cfg     config
)

func main() {
	var token, configPath string
	var printVersion, debug bool

	flag.StringVar(&token, "t", "", "discord app token")
	flag.StringVar(&configPath, "c", "", "config file path")
	flag.BoolVar(&printVersion, "v", false, "print version")
	flag.BoolVar(&debug, "d", false, "enable debugging")
	flag.Parse()

	if printVersion {
		fmt.Println("pan-bot: " + version)
		return
	}

	var err error
	cfg, err = loadConfig(configPath)
	if err != nil {
		log.Println(err)
	}

	// no/empty flag, use environment variable
	if token == "" {
		token = cfg.Token
	}

	if token == "" {
		log.Fatalln("No token provided, check \"pan-bot -h\"")
	}

	if cfg.WpcChannelID == "" {
		log.Println("wpcChannelID is empty")
	}

	if cfg.AstroChannelID == "" {
		log.Println("astroChannelID is empty")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("Error creating Discord session: ", err)
	}

	if debug {
		dg.LogLevel = 9
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)

	dg.Identify.Intents |= discordgo.IntentsGuilds |
		discordgo.IntentsGuildEmojis |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions

	err = dg.Open()
	if err != nil {
		log.Fatalln("Error opening Discord session: ", err)
	}
	defer dg.Close()

	log.Println("Pan Bot started")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("Pan Bot stopped")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case strings.HasPrefix(m.Content, "!"):
		switch msg := strings.TrimSpace(m.Content)[1:]; {
		case strings.HasPrefix(msg, "echo"):
			_, err := s.ChannelMessageSend(m.ChannelID, strings.TrimSpace(msg[4:]))
			if err != nil {
				log.Println(err)
			}
		case msg == "h" || msg == "help":
			_, err := s.ChannelMessageSend(m.ChannelID, addMdCode(helpMessage))
			if err != nil {
				log.Println(err)
			}
		case msg == "v" || msg == "version":
			_, err := s.ChannelMessageSend(m.ChannelID, "pan-bot: "+version)
			if err != nil {
				log.Println(err)
			}
		case strings.HasPrefix(msg, "sjp"):
			message, err := crappers.QuerySjp(msg[3:])
			if err != nil {
				log.Println(err)

				_, err := s.ChannelMessageSend(m.ChannelID, "sjp: błąd: "+err.Error())
				if err != nil {
					log.Println(err)
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Println(err)
				}
			}
		case strings.HasPrefix(msg, "wpc"):
			wpcPrice, err := crappers.GetWpcPrice()
			if err != nil {
				log.Println(err)
				_, err := s.ChannelMessageSend(m.ChannelID, "wpc.GetWpcPrice: błąd: "+err.Error())
				if err != nil {
					log.Println(err)
				}
			}
			newName := "700g WPC ↹ " + wpcPrice + "PLN"
			_, err = s.ChannelEdit(cfg.WpcChannelID, &discordgo.ChannelEdit{Name: newName})
			if err != nil {
				log.Println(err)
			} else {
				log.Println("WPC exchange rate updated successfully")
			}
		case strings.HasPrefix(msg, "pm"):
			response := pudzianConverter(msg[2:])
			_, err := s.ChannelMessageSend(m.ChannelID, response)
			if err != nil {
				log.Println(err)
			}
		case strings.HasPrefix(msg, "test"):
			log.Printf("%+v", cfg)
			log.Printf("%#v", cfg)
			log.Println(m.Content)
		}
	case strings.ToLower(m.Content) == "kek",
		strings.ToLower(m.Content) == "lol",
		strings.Contains(strings.ToLower(m.Content), "kekw"):
		err := reactWithGuildEmoji(s, m.Message, "KEKW")
		if err != nil {
			log.Println(err)
		}
	case strings.Contains(strings.ToLower(m.Content), "based"):
		err := reactWithGuildEmoji(s, m.Message, "gigachad")
		if err != nil {
			log.Println(err)
		}
	case strings.Contains(strings.ToLower(m.Content), "ayy lmao"):
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "👽")
		if err != nil {
			log.Println(err)
		}
	case strings.Contains(strings.ToLower(m.Content), "idź sobie"):
		_, err := s.ChannelMessageSendReply(m.ChannelID, "Nie", m.Reference())
		if err != nil {
			log.Println(err)
		}
	case containsUser(m.Mentions, s.State.User.ID):
		err := reactWithGuildEmoji(s, m.Message, "🖕")
		if err != nil {
			log.Println(err)
		}
	}
}
