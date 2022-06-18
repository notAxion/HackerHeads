package features

import (
	"strings"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

func (r *Mux) MessageCreate(s *dg.Session, m *dg.MessageCreate) {

	if strings.HasPrefix(m.Content, r.botPrefix) {

		if m.Author.ID == r.botID || m.Author.Bot {
			return
		}
		cmd := command(m.Content[1:])
		if _, ok := r.Cooldown[cmd+m.Author.ID]; ok {
			s.ChannelMessageSend(m.ChannelID, "take A chil pill this command is in cooldown")
			return
		}
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
		case "dc":
			r.DC(s, m)
		// case "whois":
		// 	r.Whois(s, m)
		case "connect-server":
			r.ConnectServer(s, m)
		case "meme":
			r.meme(s, m)
		case "prefix":
			r.Prefix(s, m)
		case "make-meme": // don't commit
			r.CaptionMeme(s, m)
		case "makememe": // don't commit
			r.CaptionMeme(s, m)
		case "clean-mess": // don't commit
			r.CleanWebhook(s, m)
		default:
			if m.Author.ID == "663985628194668566" {
				return
			}
			s.ChannelMessageSendReply(m.ChannelID, "**WHAT !?** :thinking::thinking:", m.Reference())
			return
		}
		k := cmd + m.Author.ID
		r.Cooldown[k] = true
		go func() {
			time.Sleep(3 * time.Second)
			delete(r.Cooldown, k)
		}()
	} else {
		r.ConnectServerRelay(s, m)
	}

	if m.Content == "hi bot" || m.Content == "Hi bot" {
		s.ChannelMessageSendReply(m.ChannelID, "```Hey```", m.MessageReference)
	}

	// if m.ChannelID == "765909517879476224" {
	// fmt.Println(len(m.Attachments))
	// }
}
