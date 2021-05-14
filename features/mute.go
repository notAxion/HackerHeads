package features

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/notAxion/HackerHeads/db"
)

//											***		M U T E		***

// Mute command that will mute the user so that they can't talk
// or chat in any channel however they an join the VC
// and will be able to see the message history by default
func Mute(s *dg.Session, m *dg.MessageCreate) {
	//  Getting the args
	args := fieldsN(m.Content[1:], -1)
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
	switch len(args) {
	case 1:
		helpMute(s, m.ChannelID)
	case 2:
		muteComplex(s, m, user, -1, "")
	case 3:
		mute3Arg(s, m, user, args)
	case 4:
		muteAllArgs(s, m, user, args)
	}

}

func mute3Arg(s *dg.Session, m *dg.MessageCreate, user *dg.User, args []string) {
	if len(args) != 3 {
		return
	}

	muteDur, err := time.ParseDuration(args[2])
	if err != nil {
		muteComplex(s, m, user, -1, args[2])
		return
	}
	muteComplex(s, m, user, muteDur, "")
}

func muteAllArgs(s *dg.Session, m *dg.MessageCreate, user *dg.User, args []string) {
	if len(args) < 4 {
		return
	}

	muteDur, err := time.ParseDuration(args[2])
	if err == nil {
		reason := strings.Join(args[3:], " ")
		muteComplex(s, m, user, muteDur, reason)
		return
	}
	lastIndx := len(args) - 1
	muteDur, err = time.ParseDuration(args[lastIndx])
	if err == nil {
		reason := strings.Join(args[2:lastIndx], " ")
		muteComplex(s, m, user, muteDur, reason)
		return
	}
	// no mute duration
	reason := strings.Join(args[2:], " ")
	muteComplex(s, m, user, -1, reason)

}

func muteComplex(s *dg.Session, m *dg.MessageCreate, user *dg.User, dur time.Duration, reason string) {
	// if dur < 0 { no duration}
	muteRoleID, err := muteRole(s, m)
	if err != nil {
		fmt.Println("mute role error ", err)
		return
	}
	if err = s.GuildMemberRoleAdd(m.GuildID, user.ID, muteRoleID); err != nil {
		msg := fmt.Sprintf("```sorry i cant update the role of %s#%s```", user.Username, user.Discriminator)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	// sending final message of success to the channel
	muteEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: fmt.Sprintf(":white_check_mark: *%s#%s has been muted. \u200e*", user.Username, user.Discriminator),
		Color: 0x00fa00,
	}
	if _, err = s.ChannelMessageSendEmbed(m.ChannelID, muteEmbed); err != nil {
		fmt.Println(err)
	}

	var dmOpen bool
	muteDMChan, err := s.UserChannelCreate(user.ID)
	if err == nil {
		dmOpen = true
	}

	guild, _ := s.Guild(m.GuildID)
	if dur > 0 { // with mute duration
		if dmOpen {
			s.ChannelMessageSend(muteDMChan.ID, fmt.Sprintf("you were muted from %s | %s.", guild.Name, reason))
		}
		time.Sleep(dur)
		Unmute(s, m)

	} else { // no mute duration
		if dmOpen {
			s.ChannelMessageSend(muteDMChan.ID, fmt.Sprintf("you were muted from %s | %s.", guild.Name, reason))
		}
	}
}

//												*** 	helpMute	***

func helpMute(s *dg.Session, chnID string) {
	desc := fmt.Sprintf(`
**Description**: muting a member from a server will revoke them from chatting or talking from a Channel
however then can see the message history and will be able to connect in the channels by default.
**Usage**: %smute [@user] <limit> [reason]  
**Example**:
%smute @raider 3d now cry in a corner
	`, botPrefix, botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       "\n**Command**: mute",
		Description: desc,
		Color:       0x00ff00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
	}
}

//-----------------------------------------------------------||-----------------------------------------------------------

//												***		createMuteRole		***

func createMuteRole(s *dg.Session, m *dg.MessageCreate) (muteRole *dg.Role, err error) {
	muteRole, err = s.GuildRoleCreate(m.GuildID)
	if err != nil {
		return
	}

	//(for  muted role)
	var perm int64 = 0x400 | 0x10000 | 0x100000

	_, err = s.GuildRoleEdit(
		m.GuildID, muteRole.ID, "Muted",
		0x6b6b6b, false, perm, false) // bools are hoist and
	if err != nil {
		return
	}
	return muteRole, nil
}

//												***		muteRole		***

// muteRole
func muteRole(s *dg.Session, m *dg.MessageCreate) (string, error) {
	gID := m.GuildID
	roleID, err := db.MuteRoleID(gID)
	if err != sql.ErrNoRows && err != nil {
		fmt.Println("features.muteRole error")
		return "", err
	}
	valid := validRoleID(s, m, roleID)

	if err == sql.ErrNoRows || !valid { // err = new guild, !valid = something wrong with role
		newRole, err := createMuteRole(s, m)
		if err != nil {
			fmt.Println("create role error")
			return "", err
		}

		if err = db.UpsertRole(gID, newRole.ID); err != nil {
			fmt.Println("Upsert in DB error")
			return "", err
		}

		if err = revokeChannelPerms(s, m, newRole.ID); err != nil {
			fmt.Println("revoke Channel perms error")
			return "", err
		}
	}
	return roleID, nil

}

//												***		revokeChannelPerms		***

// revokeChannelPerms will go on each channels of the guild when a mute command is called called for the first time in a guild
func revokeChannelPerms(s *dg.Session, m *dg.MessageCreate, muteRoleID string) error {
	chans, err := s.GuildChannels(m.GuildID)
	if err != nil {
		return err
	}
	textPerm := &dg.PermissionOverwrite{
		ID:   muteRoleID,
		Type: dg.PermissionOverwriteTypeRole,
		Deny: 0x800 | 0x40,
	}
	textEdit := &dg.ChannelEdit{
		PermissionOverwrites: []*dg.PermissionOverwrite{textPerm},
	}
	voicePerm := &dg.PermissionOverwrite{
		ID:   muteRoleID,
		Type: dg.PermissionOverwriteTypeRole,
		Deny: 0x200000,
	}
	voiceEdit := &dg.ChannelEdit{
		PermissionOverwrites: []*dg.PermissionOverwrite{voicePerm},
	}
	for i := range chans {
		if chans[i].Type == 0 {
			_, err = s.ChannelEditComplex(chans[i].ID, textEdit)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if chans[i].Type == 2 {
			_, err = s.ChannelEditComplex(chans[i].ID, voiceEdit)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return nil
}

// 											***		AddMuteRole		***

// AddMuteRole Will add the mute role to the channel called through bot.ManangeChannel
func AddMuteRole(s *dg.Session, chans *dg.ChannelCreate) { // ** todo move this to another file
	muteRoleID := "772777995025907732"
	textPerm := &dg.PermissionOverwrite{
		ID:   muteRoleID,
		Type: dg.PermissionOverwriteTypeRole,
		Deny: 0x800 | 0x40,
	}
	textEdit := &dg.ChannelEdit{
		PermissionOverwrites: []*dg.PermissionOverwrite{textPerm},
	}
	voicePerm := &dg.PermissionOverwrite{
		ID:   muteRoleID,
		Type: dg.PermissionOverwriteTypeRole,
		Deny: 0x200000,
	}
	voiceEdit := &dg.ChannelEdit{
		PermissionOverwrites: []*dg.PermissionOverwrite{voicePerm},
	}

	// This functions will add the role with the perms needed for the mute role

	if chans.Type == 0 {
		_, err := s.ChannelEditComplex(chans.ID, textEdit)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if chans.Type == 2 {
		_, err := s.ChannelEditComplex(chans.ID, voiceEdit)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
