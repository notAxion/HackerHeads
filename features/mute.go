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
func (r *Mux) Mute(s *dg.Session, m *dg.MessageCreate) {
	//  Getting the args
	args := fieldsN(m.Content[1:], -1)
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
	if m.Author.ID == user.ID || user.Bot {
		r.MuteReply(s, m, user, "")
		return
	}
	switch {
	case len(args) == 1:
		r.helpMute(s, m.ChannelID)
	case len(args) == 2:
		r.muteComplex(s, m, user, -1, "")
	case len(args) == 3:
		r.mute3Arg(s, m, user, args)
	case len(args) >= 4:
		r.muteAllArgs(s, m, user, args)
	}

}

func (r *Mux) mute3Arg(s *dg.Session, m *dg.MessageCreate, user *dg.User, args []string) {
	if len(args) != 3 {
		return
	}

	muteDur, err := time.ParseDuration(args[2])
	if err != nil {
		r.muteComplex(s, m, user, -1, args[2])
		return
	}
	r.muteComplex(s, m, user, muteDur, "")
}

func (r *Mux) muteAllArgs(s *dg.Session, m *dg.MessageCreate, user *dg.User, args []string) {
	if len(args) < 4 {
		return
	}

	muteDur, err := time.ParseDuration(args[2])
	if err == nil {
		reason := strings.Join(args[3:], " ")
		r.muteComplex(s, m, user, muteDur, reason)
		return
	}
	lastIndx := len(args) - 1
	muteDur, err = time.ParseDuration(args[lastIndx])
	if err == nil {
		reason := strings.Join(args[2:lastIndx], " ")
		r.muteComplex(s, m, user, muteDur, reason)
		return
	}
	// no mute duration
	reason := strings.Join(args[2:], " ")
	r.muteComplex(s, m, user, -1, reason)

}

func (r *Mux) muteComplex(s *dg.Session, m *dg.MessageCreate, user *dg.User, dur time.Duration, reason string) {
	if reason == "" {
		reason = "for no reason lol"
	}
	if dur > 0 && dur < 30*time.Second {
		s.ChannelMessageSend(m.ChannelID, "Give at least a 30s mute ||u dumb||")
		return
	}

	muteRoleID, err := r.muteRole(s, m.GuildID)
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
	r.MuteReply(s, m, user, reason)
	if dur > 0 {
		tmp := time.Now().Add(dur)
		r.AddTimer(db.FromString(m.GuildID, user.ID), tmp)
		r.PQ.SaveUnmuteTime(m.GuildID, user.ID, tmp)
	}

}

func (r *Mux) MuteReply(s *dg.Session, m *dg.MessageCreate, user *dg.User, reason string) {
	muteEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: fmt.Sprintf(":white_check_mark: *%s#%s has been muted. \u200e*", user.Username, user.Discriminator),
		Footer: &dg.MessageEmbedFooter{
			Text: "bye bye",
		},
		Color: 0x00fa00,
	}
	if m.Author.ID == user.ID || user.Bot {
		uNoob := "https://t4.ftcdn.net/jpg/01/35/86/23/240_F_135862342_Q3LJPMyd8LLhBm4fPPRhemy8CczBzr4G.jpg"
		muteEmbed.Title += "||jk||"
		muteEmbed.Footer.Text = "\u200e"
		muteEmbed.Footer.IconURL = uNoob
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, muteEmbed); err != nil {
		fmt.Println(err)
	}
	if m.Author.ID == user.ID {
		return
	}

	var dmOpen bool
	muteDMChan, err := s.UserChannelCreate(user.ID)
	if err == nil {
		dmOpen = true
	}

	guild, _ := s.Guild(m.GuildID)
	if dmOpen {
		msg := fmt.Sprintf("you were muted from %s | %s.", guild.Name, reason)
		s.ChannelMessageSend(muteDMChan.ID, msg)
	}
}

//												*** 	helpMute	***

func (r *Mux) helpMute(s *dg.Session, chnID string) {
	desc := fmt.Sprintf(`
**Description**: muting a member from a server will revoke them from chatting or talking from a Channel
however then can see the message history and will be able to connect in the channels by default.
**Usage**: %smute [@user] <limit> [reason]  
**Example**:
%smute @raider 3d now cry in a corner
	`, r.botPrefix, r.botPrefix)
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

func (r *Mux) createMuteRole(s *dg.Session, guildID string) (muteRole *dg.Role, err error) {
	muteRole, err = s.GuildRoleCreate(guildID)
	if err != nil {
		return
	}

	//(for  muted role)
	var perm int64 = 0x400 | 0x10000 | 0x100000

	_, err = s.GuildRoleEdit(
		guildID, muteRole.ID, "Muted",
		0x6b6b6b, false, perm, false) // bools are hoist and
	if err != nil {
		return
	}
	return muteRole, nil
}

//												***		muteRole		***

// muteRole
func (r *Mux) muteRole(s *dg.Session, guildID string) (string, error) {
	roleID, err := r.PQ.MuteRoleID(guildID)
	if err != sql.ErrNoRows && err != nil {
		fmt.Println("features.muteRole error")
		return "", err
	}
	valid := r.validRoleID(s, guildID, roleID)

	if err == sql.ErrNoRows || !valid { // err = new guild, !valid = something wrong with role
		newRole, err := r.createMuteRole(s, guildID)
		if err != nil {
			fmt.Println("create role error")
			return "", err
		}

		if err = r.PQ.UpsertRole(guildID, newRole.ID); err != nil {
			fmt.Println("Upsert in DB error")
			return "", err
		}

		if err = r.revokeChannelPerms(s, guildID, newRole.ID); err != nil {
			fmt.Println("revoke Channel perms error")
			return "", err
		}
	}
	return roleID, nil

}

//												***		revokeChannelPerms		***

// revokeChannelPerms will go on each channels of the guild
// when a mute command is called called for the first time in a guild
func (r *Mux) revokeChannelPerms(s *dg.Session, guildID, muteRoleID string) error {
	chans, err := s.GuildChannels(guildID)
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
func (r *Mux) AddMuteRole(s *dg.Session, chans *dg.ChannelCreate) { // ** todo move this to another file
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
