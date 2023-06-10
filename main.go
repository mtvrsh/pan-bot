// simple discord bot https://github.com/m3tav3rse/pan-bot
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/m3tav3rse/pan-bot/pkg/bot"
)

var version = "unknown"

func main() {
	var token, configPath string
	var versionFlag, debug bool

	flag.StringVar(&token, "t", "", "discord app token")
	flag.StringVar(&configPath, "c", "", "config file path")
	flag.BoolVar(&versionFlag, "v", false, "print version")
	flag.BoolVar(&debug, "d", false, "enable debugging")
	flag.Parse()

	if versionFlag {
		fmt.Println("pan-bot: " + version)
		return
	}

	bot.Version = version

	var err error
	bot.Cfg, err = bot.LoadConfig(configPath)
	if err != nil {
		log.Println(err)
	}

	// no/empty flag, use environment variable
	if token == "" {
		token = bot.Cfg.Token
	}

	if token == "" {
		log.Fatalln("No token provided, check \"pan-bot -h\"")
	}

	if bot.Cfg.WpcChannelID == "" {
		log.Println("wpcChannelID is empty")
	}

	if bot.Cfg.AstroChannelID == "" {
		log.Println("astroChannelID is empty")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("Error creating Discord session: ", err)
	}

	if debug {
		dg.LogLevel = 9
	}

	dg.AddHandler(bot.MessageCreate)
	dg.AddHandler(bot.MessageReactionAdd)

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
