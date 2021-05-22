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

// var muteDone map[string]chan struct{}
// 												***		U N M U T E 	***

//
func (r *Mux) Unmute(s *dg.Session, m *dg.MessageCreate) {
	args := fieldsN(m.Content, -1)
	if len(args) < 2 {
		r.helpMute(s, m.ChannelID) //!valid or just checking help mute
		return
	}
	user, valid := r.validUserID(s, m, args[1])
	if !valid {
		r.helpMute(s, m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```I can't find that user, %s```", args[1]))
		return
	}
	muteRoleID, err := r.muteRole(s, m)
	if err != nil {
		fmt.Println("unmute mute role error ", err)
		return
	}
	if err = r.removeRole(s, m.GuildID, user.ID, muteRoleID); err != nil {
		r.helpMute(s, m.ChannelID)
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

func (r *Mux) removeRole(s *dg.Session, GuildID, userID, muteRoleID string) error {
	return s.GuildMemberRoleRemove(GuildID, userID, muteRoleID)
}

// SetAllMutedTimer is only for init which will set the timer
// of all muted user which were in the db
func (r *Mux) SetAllMutedTimer() {
	users, err := db.TGetMutedUsers()
	if err != nil {
		// don't know what to with it
		fmt.Println(err)
	}
	for k, user := range users {
		fmt.Println(k)
		fmt.Println(user) // check this
		// afterFunc will be here
		// delete the key when unmute
	}
}
