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
	token          string
	astroChannelID string
	wpcChannelID   string
}

func loadConfig(path string) (config, error) {
	var c map[string]string
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

	return config{
		token: c["Token"], astroChannelID: c["astroChannelID"],
		wpcChannelID: c["wpcChannelID"],
	}, nil
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

func pudzianConverter(s string) string {
	words := strings.Fields(s)

	if len(words) == 1 {
		switch {
		case strings.HasSuffix(words[0], "HP"):
			words[0] = strings.TrimSuffix(words[0], "HP")
			words = append(words, "HP")
		case strings.HasSuffix(words[0], "KM"):
			words[0] = strings.TrimSuffix(words[0], "KM")
			words = append(words, "KM")
		case strings.HasSuffix(words[0], "PM"):
			words[0] = strings.TrimSuffix(words[0], "PM")
			words = append(words, "PM")
		case strings.HasSuffix(words[0], "W"):
			words[0] = strings.TrimSuffix(words[0], "W")
			words = append(words, "W")
		}
	}

	if len(words) == 2 {
		number, err := strconv.ParseFloat(words[0], 64)
		if err != nil {
			log.Println(err) // TODO remove?
			return "Błąd 404: nie znaleziono Pudziana"
		}
		switch words[1] {
		case "HP":
			fallthrough
		case "KM":
			return fmt.Sprintf("%v %v to %.2f PM", words[0], words[1],
				pudzianmechaniczny.HPToPM(number))
		case "PM":
			return fmt.Sprintf("%v %v to %.2f KM", words[0], words[1],
				pudzianmechaniczny.PMToHP(number))
		case "W":
			return fmt.Sprintf("%v %v to %.2f PM", words[0], words[1],
				pudzianmechaniczny.WToPM(number))
		}
	}

	return "Błędne dane: podaj liczbę i jednostkę (KM,PM,W)\n" +
		"Na przykład: 1337 KM"
}
