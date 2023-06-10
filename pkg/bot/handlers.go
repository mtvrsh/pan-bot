package bot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/m3tav3rse/pan-bot/pkg/crappers"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			_, err := s.ChannelMessageSend(m.ChannelID, "pan-bot: "+Version)
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
			_, err = s.ChannelEdit(Cfg.WpcChannelID, &discordgo.ChannelEdit{Name: newName})
			if err != nil {
				log.Println(err)
			} else {
				log.Println("WPC exchange rate updated successfully")
			}
		case strings.HasPrefix(msg, "pm"):
			_, err := s.ChannelMessageSend(m.ChannelID, pudzianConverter(msg[2:]))
			if err != nil {
				log.Println(err)
			}
		case strings.HasPrefix(msg, "test"):
			log.Printf("%+v", Cfg)
			log.Printf("%#v", Cfg)
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

func MessageReactionAdd(session *discordgo.Session, m *discordgo.MessageReactionAdd) {
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
