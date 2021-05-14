package features

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

// 												***		U N M U T E 	***

//
func Unmute(s *dg.Session, m *dg.MessageCreate) {
	args := fieldsN(m.Content, -1)
	if len(args) < 2 {
		helpMute(s, m.ChannelID) //!valid or just checking help mute
		return
	}
	user, valid := validUserID(s, m, args[1])
	if !valid {
		helpMute(s, m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```I can't find that user, %s```", args[1]))
		return
	}
	muteRoleID, err := muteRole(s, m)
	if err != nil {
		fmt.Println("unmute mute role error ", err)
		return
	}
	if err = s.GuildMemberRoleRemove(m.GuildID, user.ID, muteRoleID); err != nil {
		helpMute(s, m.ChannelID)
		msg := fmt.Sprintf("```can't remove role of %s#%s```", user.Username, user.Discriminator)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	// sending final message of success to the channel
	unmuteEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: fmt.Sprintf(":white_check_mark: *%s#%s has been unmuted. \u200e*", user.Username, user.Discriminator),
		Color: 0x00fa00,
	}
	if _, err = s.ChannelMessageSendEmbed(m.ChannelID, unmuteEmbed); err != nil {
		fmt.Println(err)
	}

}
