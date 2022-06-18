package features

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

func (r *Mux) DC(s *dg.Session, m *dg.MessageCreate) {
	args := fieldsN(m.Content[1:], -1)
	if len(args) < 2 {
		r.helpDC(s, m.ChannelID)
		return
	}
	var user *dg.User
	if args[1] == "@me" {
		user = m.Author
	} else {
		var valid bool
		user, valid = r.validUserID(s, m, args[1])
		if !valid {
			r.helpDC(s, m.ChannelID)
			msg := fmt.Sprintf("I can't find that user, %s", args[1])
			s.ChannelMessageSend(m.ChannelID, InCodeBlock(msg))
			return
		}
	}
	s.GuildMemberMove(m.GuildID, user.ID, nil)
}

func (r *Mux) helpDC(s *dg.Session, chnID string) {
	desc := fmt.Sprintf(`
		**Description**: does what it sounds like it dc/ disconnects people.
		 doing @me will just dc you
		**Usage**: %sdc <@mention | @me>
		**Example**: 
		%sdc @me
		%sdc @raider
		`, r.botPrefix, r.botPrefix, r.botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       "\n**Command**: dc",
		Description: desc,
		Color:       0x00ff00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
	}
}
