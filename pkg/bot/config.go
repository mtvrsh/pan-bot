package bot

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

const defaultConfigPath = "/run/secrets/pan-bot-config"

type config struct {
	Token          string
	AstroChannelID string
	WpcChannelID   string
}

func LoadConfig(path string) (config, error) {
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
