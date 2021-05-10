package features

import (
	"database/sql"
	"fmt"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/notAxion/HackerHeads/db"
)

//											***		M U T E		***

// Mute command that will mute the user so that they can't talk
// or chat in any channel however they an join the VC
// and will be able to see the message history by default
func Mute(s *dg.Session, m *dg.MessageCreate) {
	muteRoleID := "772777995025907732"
	var dmOpen bool

	//if !(already created) :
	/*
		err := createMuteRole(s, m)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	*/
	/*
		err := revokeChannelPerms(s, m)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Done")
	*/

	// Help mute
	if m.Content[1:] == "mute" || m.Content[1:] == "mute " {
		helpMute(s, m.ChannelID)
		return
	}
	//  Getting the args
	args := fieldsN(m.Content[1:], 3)
	if len(args) == 0 {
		helpMute(s, m.ChannelID) //!valid
		return
	}

	muteID, valid := validUserID(s, m, args[1])
	if !valid {
		idError := &dg.MessageEmbed{
			Type:  "rich",
			Title: fmt.Sprintf(":x: **I can't find that user, %s**", args[1]),
			Color: 0xff0000,
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, idError)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
	user, err := s.User(muteID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	guild, _ := s.Guild(m.GuildID) // ** todo try some other way to get the guild name cause i have to get all values just to get the guild name
	limitArgs := fieldsN(args[2], 2)
	// if len(limitArgs) == 0 { // if the length is 0 then there was just the reason there was reason entered no time limit
	//goto noLimit
	// }
	limit, errLimit := time.ParseDuration(limitArgs[0])
	err = s.GuildMemberRoleAdd(m.GuildID, user.ID, muteRoleID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("sorry i cant update the role  of %s#%s", user.Username, user.Discriminator))
		return
	}
	muteEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: fmt.Sprintf(":white_check_mark: *%s#%s has been muted. \u200e*", user.Username, user.Discriminator),
		Color: 0x00fa00,
	}
	// Sending the embed text
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, muteEmbed)
	if err != nil {
		fmt.Println(err.Error())
	}

	//  Creates dm channel for the the muted user
	muteDMChan, err := s.UserChannelCreate(user.ID)
	if err == nil {
		dmOpen = true
	}

	if errLimit == nil { // mute duration
		if dmOpen {
			s.ChannelMessageSend(muteDMChan.ID, fmt.Sprintf("you were muted from %s | %s.", guild.Name, limitArgs[1]))
		}
		time.Sleep(limit)
		s.GuildMemberRoleRemove(m.GuildID, user.ID, muteRoleID) // change this with unmute func

	} else { // no mute duration
		if dmOpen {
			s.ChannelMessageSend(muteDMChan.ID, fmt.Sprintf("you were muted from %s | %s.", guild.Name, limitArgs[1]))
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
	roleID, err := db.PQ.MuteRoleID(gID)
	if err != sql.ErrNoRows && err != nil {
		fmt.Println("features.muteRole error")
		return "", err
	}
	valid := validRoleID(s, m, roleID)

	if err == sql.ErrNoRows || !valid { // err = new guild, !valid = something wrong with role
		role, err := createMuteRole(s, m)
		if err != nil {
			fmt.Println("create role error")
			return "", err
		}
		err = revokeChannelPerms(s, m, role.ID) // maybe use goroutine
		if err != nil {
			fmt.Println("revoke Channel perms error")
			return "", err
		}
		// db.UpdateMuteRole() || db.InsertMuteRole() use {{upsert insert on conflict }}
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
