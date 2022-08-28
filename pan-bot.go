// simple discord bot https://github.com/m3tav3rse/pan-bot
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const helpMessage = `Dostępne komendy:
!echo <tekst>   Wypisuje tekst
!h|!help        Na pewno nie wyświetli tej listy komend
!sjp <wyraz>    Znaczenie wyrazu z sjp.pl
!v|!version     Wyświetla wersję bota
`

var version = "unknown"

func main() {
	var token string
	var showVersion bool

	flag.StringVar(&token, "t", "", "bot token")
	flag.BoolVar(&showVersion, "v", false, "print version")
	flag.Parse()

	if showVersion {
		fmt.Println("pan-bot: " + version)
		return
	}

	if token == "" {
		log.Fatalln("No token provided, check \"pan-bot -h\"")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("Error creating Discord session: ", err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildEmojis | discordgo.IntentsGuildMessageReactions

	err = dg.Open()
	if err != nil {
		log.Fatalln("Error opening Discord session: ", err)
	}
	defer func() {
		err = dg.Close()
		if err != nil {
			log.Println(err)
		}
	}()

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

	// better way to do it instead of this gigaswitch?
	// map[string]func() maybe
	switch {
	case strings.HasPrefix(m.Content, "!"):
		switch msg := strings.TrimSpace(m.Content)[1:]; {
		case strings.HasPrefix(msg, "echo"):
			_, err := s.ChannelMessageSend(m.ChannelID, strings.TrimSpace(msg[4:]))
			if err != nil {
				log.Println(err)
			}
		case msg == "h" || msg == "help":
			_, err := s.ChannelMessageSend(m.ChannelID, asMdCode(helpMessage))
			if err != nil {
				log.Println(err)
			}
		case msg == "v" || msg == "version":
			_, err := s.ChannelMessageSend(m.ChannelID, "pan-bot: "+version)
			if err != nil {
				log.Println(err)
			}
		case strings.HasPrefix(msg, "sjp"):
			message, err := sjpQuery(msg[3:])
			if err != nil {
				log.Println(err)
				_, err := s.ChannelMessageSend(m.ChannelID, "sjpQuery: błąd: "+err.Error())
				if err != nil {
					log.Println(err)
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Println(err)
				}
			}
		case strings.HasPrefix(msg, "s"):
			log.Println("Not implemented")
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
	case userListContains(m.Mentions, s.State.User.ID):
		err := reactWithGuildEmoji(s, m.Message, "🖕")
		if err != nil {
			log.Println(err)
		}
	}
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Println("Couldn't retrieve message to react to")
		return
	} else if msg.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case m.Emoji.Name == "KEKW":
		err := s.MessageReactionAdd(m.ChannelID, m.MessageID, m.Emoji.APIName())
		if err != nil {
			log.Println(err)
		}
	}
}

func getEmoji(s *discordgo.Session, emoji, guildID string) (string, error) {
	list, err := s.GuildEmojis(guildID)
	if err != nil {
		// different error/info
		return "", err
	}
	for _, e := range list {
		if e.Name == emoji {
			return e.APIName(), nil
		}
	}
	return emoji, nil
}

func reactWithGuildEmoji(s *discordgo.Session, m *discordgo.Message, name string) error {
	e, err := getEmoji(s, name, m.GuildID)
	if err != nil {
		return err
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, e)
	if err != nil {
		return err
	}
	return nil
}

func sjpQuery(word string) (string, error) {
	w := strings.TrimSpace(word)
	if w == "" {
		return "Puste zapytanie, pusta odpowiedź :)", nil
	}

	resp, err := http.Get("https://sjp.pl/" + w)
	if err != nil {
		log.Println(err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if strings.Contains(string(body), "<p>nie występuje w słowniku</p>") {
		return "Nie występuje w słowniku", nil
	}

	parAll := regexp.MustCompile("<p style=\"margin: .5em 0; font: medium/1.4 sans-serif; max-width: 34em; \">.*</p>")
	parStart := regexp.MustCompile("<p style=\"margin: .5em 0; font: medium/1.4 sans-serif; max-width: 34em; \">")
	parEnd := regexp.MustCompile("</p>")
	linebreak := regexp.MustCompile("<br />")
	cleanEnd := regexp.MustCompile("[;,\\. ]+\n")

	text := parAll.FindAllString(string(body), -1)

	var result string

	for _, x := range text {
		//s := html.UnescapeString(x) //check if needed
		s := linebreak.ReplaceAllString(x, "\n")
		s = parStart.ReplaceAllString(s, "")
		s = parEnd.ReplaceAllString(s, "")
		s = cleanEnd.ReplaceAllString(s, "\n")
		result += s + "\n"
	}

	result = strings.Trim(result, " \t\n")

	if result == "" {
		return w + " nie ma opisu", nil
	}
	return result, nil
}

func asMdCode(s string) string {
	return "```text\n" + s + "\n```"
}

func userListContains(a []*discordgo.User, b string) bool {
	for _, v := range a {
		if v.ID == b {
			return true
		}
	}
	return false
}
