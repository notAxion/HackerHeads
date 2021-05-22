package features

import (
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func (r *Mux) MessageCreate(s *dg.Session, m *dg.MessageCreate) {

	if strings.HasPrefix(m.Content, r.botPrefix) {

		if m.Author.ID == r.botID || m.Author.Bot {
			return
		}
		cmd := command(m.Content)
		switch cmd {
		case "event":
			r.EventStart(s, m)
		case "mute":
			r.Mute(s, m)
		case "ping":
			r.Ping(s, m)
		case "remind":
			r.Remind(s, m)
		case "warn":
			r.Warn(s, m)
		case "unmute":
			r.Unmute(s, m)
		default:
			s.ChannelMessageSend(m.ChannelID, "**WHAT !?** :thinking::thinking:")
		}
	}

	if m.Content == "hi bot" || m.Content == "Hi bot" {
		s.ChannelMessageSendReply(m.ChannelID, "```Hey```", m.MessageReference)
	}

	// if m.ChannelID == "765909517879476224" {
	// fmt.Println(len(m.Attachments))
	// }
}
