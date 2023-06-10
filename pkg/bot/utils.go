package bot

import "github.com/bwmarrin/discordgo"

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
