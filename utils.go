package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/m3tav3rse/pan-bot/pkg/pudzianmechaniczny"
)

const defaultConfigPath = "/run/secrets/pan-bot-config"

const helpMessage = `Dostępne komendy:
!echo <tekst>   Wypisuje tekst
!h|!help        Na pewno nie wyświetli tej listy komend
!sjp <wyraz>    Znaczenie wyrazu z sjp.pl
!v|!version     Wyświetla wersję bota
!pm             Konwerter Pudzianów Mechanicznyh (PM) do watów i koni mechanicznych (W/KM)
`

type config struct {
	Token          string
	AstroChannelID string
	WpcChannelID   string
}

func loadConfig(path string) (config, error) {
	var c config
	var err error

	if path == "" {
		path, err = os.UserConfigDir()
		if err != nil {
			return config{}, err
		}
		path = filepath.Join(path, "pan-bot", "config.json")

		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			path = defaultConfigPath
		}
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return config{}, err
	}
	if err = json.Unmarshal(f, &c); err != nil {
		return config{}, err
	}

	log.Println("config loaded from: " + path)

	return c, nil
}

func messageReactionAdd(session *discordgo.Session, m *discordgo.MessageReactionAdd) {
	msg, err := session.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Println("Couldn't retrieve message to react to")
		return
	} else if msg.Author.ID == session.State.User.ID {
		return
	}

	switch {
	case m.Emoji.Name == "KEKW":
		err := session.MessageReactionAdd(m.ChannelID, m.MessageID, m.Emoji.APIName())
		if err != nil {
			log.Println(err)
		}
	}
}

func getGuildEmoji(s *discordgo.Session, emoji, guildID string) (string, error) {
	list, err := s.GuildEmojis(guildID)
	if err != nil {
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
	e, err := getGuildEmoji(s, name, m.GuildID)
	if err != nil {
		return err
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, e)
	if err != nil {
		return err
	}

	return nil
}

func addMdCode(s string) string {
	return "```text\n" + s + "\n```"
}

func containsUser(a []*discordgo.User, b string) bool {
	for _, v := range a {
		if v.ID == b {
			return true
		}
	}

	return false
}

func pudzianConverter(s string) string { // refactor
	var unit string
	failMsg := "Błędne dane: podaj liczbę i jednostkę (KM,PM,W)\n" +
		"Na przykład: 1337 KM"
	words := strings.Fields(s)

	// TODO rm HP?
	if len(words) == 1 {
		switch {
		case strings.HasSuffix(words[0], "HP"):
			words[0] = strings.TrimSuffix(words[0], "HP")
			fallthrough
		case strings.HasSuffix(words[0], "KM"):
			unit = "KM"
		case strings.HasSuffix(words[0], "PM"):
			unit = "PM"
		case strings.HasSuffix(words[0], "W"):
			unit = "W"
		}
		words[0] = strings.TrimSuffix(words[0], unit)
		words = append(words, unit)
	}

	if len(words) == 2 {
		var result float64
		number, err := strconv.ParseFloat(words[0], 64)
		if err != nil {
			return failMsg
		}
		switch words[1] {
		case "HP":
			fallthrough
		case "KM":
			unit = "PM"
			result = pudzianmechaniczny.HPToPM(number)
		case "PM":
			unit = "KM"
			result = pudzianmechaniczny.PMToHP(number)
		case "W":
			unit = "PM"
			result = pudzianmechaniczny.WToPM(number)
		default:
			return failMsg
		}
		return fmt.Sprintf("%v %v to %.2f %v", words[0], words[1], result, unit)
	}

	return failMsg
}
