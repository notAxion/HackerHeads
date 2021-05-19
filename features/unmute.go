package features

import (
	"fmt"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/notAxion/HackerHeads/db"
)

type muteTime struct {
	GID, UserID string
	UnmuteTime  time.Time
}

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
	if err = removeRole(s, m.GuildID, user.ID, muteRoleID); err != nil {
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

func removeRole(s *dg.Session, GuildID, userID, muteRoleID string) error {
	return s.GuildMemberRoleRemove(GuildID, userID, muteRoleID)
}

// SetAllMutedTimer is only for init which will set the timer
// of all muted user which were in the db
func SetAllMutedTimer() {
	users, err := db.GetMutedUsers()
	if err != nil {
		// don't know what to with it
		fmt.Println(err)
	}
	for _, user := range users {
		fmt.Println(user) // check this
		// afterFunc will be here
	}
}
