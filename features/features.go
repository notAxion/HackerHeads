package features

import (
	"fmt"
	str "strings"
	t "time"

	"github.com/notAxion/HackerHeads/config"

	dg "github.com/bwmarrin/discordgo"
)

var botPrefix string = config.BotPrefix

// #todo add a variable aka var = botPrefix and replace every botPrefix with that var
// #and also change every help funcs and use different variables for the desc.

// stores the eventRoleID for each channel (map[ChannelID]eventRoleID)
// go to EventStop for more info # thanks to chanbakjsd from gophers
var eventMap map[string]string = make(map[string]string)

// 													***		E V E N T 	***

// EventStart will start an instance of an event for that channel
// So afterwards if any member of that event types any message it will give that member a role which should be specified when event start command was sent
// and removes the role when a message is deleted within the event period
func EventStart(s *dg.Session, m *dg.MessageCreate) { //# todo check the role hierarchy to check if person that is sending the command does himself has the perms to add that role someone else
	if m.Content[1:] == "event" || m.Content[1:] == "event " {
		helpEvent(s, m.ChannelID)
		return
	}
	args := fieldsN(m.Content[1:], 3)
	if len(args) == 0 {
		//helpEvent(s, m.ChannelID)
		return
	}
	if args[1] != "start" {
		helpEvent(s, m.ChannelID)
		return
	}
	eventRoleID, valid := validRoleID(s, m, args[2])
	if !valid {
		roleinvalidEmbed := &dg.MessageEmbed{
			Type:        "rich",
			Title:       "Invvalid Role ID",
			Description: "Role doesn't exist \n or I don't have perms to add that role to anyone tip: check the role hierarchy and put it below than my highest role",
			Color:       0xff0000,
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, roleinvalidEmbed)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
	eventMap[m.ChannelID] = eventRoleID
	s.ChannelMessageSend(m.ChannelID, "event is started")

}

// EventRoleAdd will be adding roles to the users after the event is started also this will handle the event stop command
func EventRoleAdd(s *dg.Session, m *dg.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if eventMap[m.ChannelID] == "" {
		return
	}

	//This if part will check for the event stop command and will stop the event
	if str.HasPrefix(m.Content, botPrefix) {
		args := fieldsN(m.Content[1:], 2)
		if len(args) == 0 {
			return
		}
		if args[0] == "event" && args[1] == "stop" {
			eventMap[m.ChannelID] = ""
			s.ChannelMessageSend(m.ChannelID, "event has stopped")
			return
		}
	}
	eventRoleID := eventMap[m.ChannelID]
	err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, eventRoleID)
	if err != nil {
		fmt.Println(err.Error())
	}
}
func EventRoleRemove(s *dg.Session, m *dg.MessageDelete) {
	if eventMap[m.BeforeDelete.ChannelID] == "" {
		return
	}
	err := s.GuildMemberRoleRemove(m.BeforeDelete.GuildID, m.BeforeDelete.Author.ID, eventMap[m.BeforeDelete.ChannelID])
	fmt.Println(m.BeforeDelete.GuildID, m.BeforeDelete.Author.ID, eventMap[m.BeforeDelete.ChannelID])
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 											***		helpEvent		***

func helpEvent(s *dg.Session, chnID string) {
	desc := fmt.Sprintf("\n**Description**:  EventStart will start an instance of an event for that channel \nSo afterwards if any member of that event types any message it will give that member a role which should be specified when event start command was sent \nand removes the role when a message is deleted within the event period.\n*It is advisable to create a new channel and then start the event and dont reuse it for other events* \n**Usage**: %sevent < [start] || [stop] > {@role}  \n**Example**:\n\t%sevent start @participant \n\t%sevent stop", botPrefix, botPrefix, botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       fmt.Sprintf("\n**Command**: event"),
		Description: desc,
		Color:       0x00fa00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
	}
}

//											***		M U T E		***

// Mute command that will mute the user so that he can't talk or chat in any channel however they an join the VC and will be able to see the message history by default
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
	if len(limitArgs) == 0 { // if the length is 0 then there was just the reason there was reason entered no time limit
		//goto noLimit
	}
	limit, errLimit := t.ParseDuration(limitArgs[0])
	err = s.GuildMemberRoleAdd(m.GuildID, user.ID, muteRoleID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("sorry i cant update the role  of %s#%s", user.Username, user.Discriminator))
		return
	}
	muteEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: fmt.Sprintf(":white_check_mark: *%s#%s has been muted. ‎*", user.Username, user.Discriminator),
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

	if errLimit == nil { // this mean limit exists
		if dmOpen {
			s.ChannelMessageSend(muteDMChan.ID, fmt.Sprintf("you were muted from %s | %s.", guild.Name, limitArgs[1]))
		}
		t.Sleep(limit)
		s.GuildMemberRoleRemove(m.GuildID, user.ID, muteRoleID)

	} else {
		if dmOpen {
			s.ChannelMessageSend(muteDMChan.ID, fmt.Sprintf("you were muted from %s | %s.", guild.Name, limitArgs[1]))
		}
	}
}

//												*** 	helpMute	***

func helpMute(s *dg.Session, chnID string) {
	desc := fmt.Sprintf("\n**Description**: muting a member from a server will revoke them from chatting or talking from a Channel, however then can see the message history and will be able to connect in the channels by default.\n**Usage**: %smute [@user] <limit> [reason]  \n**Example**:\n\t%smute @raider 3d be happy with muted", botPrefix, botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       fmt.Sprintf("\n**Command**: mute"),
		Description: desc,
		Color:       0x00ff00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
	}
}

//											***		P I N G 	***

//Ping : checks the ping of the bot in millisecond
func Ping(s *dg.Session, m *dg.MessageCreate) {
	// get current time, send message, subtract new current time with old, update said message to show that time
	u := t.Now()
	msg, _ := s.ChannelMessageSend(m.ChannelID, "pong -")
	//v, _ := msg.Timestamp.Parse()
	v := t.Now()
	ping := v.Sub(u)
	//fmt.Println(u, msg.Timestamp)
	//s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("pong - `%v`", ))
	s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("pong - `%dms`", ping.Milliseconds()))
}

// 											***		R E M I N D		***

//Remind command will remind after a specified time and you should have to add some reminder description to get the reminder
func Remind(s *dg.Session, m *dg.MessageCreate) {
	if m.Content[1:] == "remind" || m.Content[1:] == "remind " {
		helpRemind(s, m.ChannelID)
		return
	}
	rem := fieldsN(m.Content[1:], 3)
	if len(rem) == 0 {
		helpRemind(s, m.ChannelID)
		return
	}
	timer, err := t.ParseDuration(rem[1])
	if err != nil {
		timeError := &dg.MessageEmbed{
			Type:  "rich",
			Title: "**Wrong time format for remind ex-> 1m , 12h50m , 10h50m55s**",
			Color: 0xff0000,
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, timeError)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//s.ChannelMessageSend(m.ChannelID, "**Wrong time format for remind ex-> 1m , 12h50m , 10h50m55s**" )
		return
	}

	remindEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: "Reminder is set",
		Color: 0x00ff00,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, remindEmbed)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Sleep(timer)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> here is the reminder you asked \n%s", m.Author.ID, rem[2]))

}

//												***		helpRemind		***

func helpRemind(s *dg.Session, chnID string) {

	desc := fmt.Sprintf("\n**Description**: Will tag you at the channel you have sent the message so that it can remind you.\n**Usage**: %sremind [time] [Reminder] \n**Example**:\n\t%sremind 10m turn off microwave\n\t%sremind 1h start a pole", botPrefix, botPrefix, botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       fmt.Sprintf("\n**Command**: remind"),
		Description: desc,
		Color:       0x00ff00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// 												***		U N M U T E 	***

//
/*
func Unmute(s *dg.Session, m *dg.MessageCreate) {

}
*/
//												***		W A R N		***

// Warn command that will warn that user upon the specified reason and dm them if possible (todo add logging to a DB)
func Warn(s *dg.Session, m *dg.MessageCreate) {
	// Help warn
	if m.Content[1:] == "warn" || m.Content[1:] == "warn " {
		helpWarn(s, m.ChannelID)
		return
	}
	//  Getting the args
	args := fieldsN(m.Content[1:], 3)
	if len(args) == 0 {
		helpWarn(s, m.ChannelID)
	}
	warnID, valid := validUserID(s, m, args[1])
	if !valid { // provided id is not valid
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
	user, err := s.User(warnID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	warnEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: fmt.Sprintf(":white_check_mark: *%s#%s has been warned. ‎*", user.Username, user.Discriminator),
		Color: 0x00fa00,
	}

	//  Creates dm channel for the the warned user
	warnDMChn, err := s.UserChannelCreate(warnID)
	guild, _ := s.Guild(m.GuildID) // ** todo try some other way to get the guild name cause i have to get all values just to get the guild name
	_, err = s.ChannelMessageSend(warnDMChn.ID, fmt.Sprintf("You were warned in %s for %s", guild.Name, args[2]))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//  Sending the embeded text
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, warnEmbed)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//												*** 	helpWarn	***

func helpWarn(s *dg.Session, chnID string) {
	desc := fmt.Sprintf("\n**Description**: warn those noobs who don't follow the rules.\n**Usage**: %swarn [@user] [reason]  \n**Example**:\n\t%swarn @noobSpammer stop the spam please", botPrefix, botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       fmt.Sprintf("\n**Command**: warn"),
		Description: desc,
		Color:       0x00ff00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

//-----------------------------------------------------------||-----------------------------------------------------------

//												***		createMuteRole		***

func createMuteRole(s *dg.Session, m *dg.MessageCreate) error {
	muteRole, err := s.GuildRoleCreate(m.GuildID)
	if err != nil {
		return err
	}

	var perm int64 = 0x400 | 0x10000 | 0x100000                                                  //(for  muted role)
	_, err = s.GuildRoleEdit(m.GuildID, muteRole.ID, "test1Muted", 0x6b6b6b, false, perm, false) // bools are hoist and
	if err != nil {
		return err
	}
	return nil
}

// revokeChannelPerms will go on each channels of the guild when a mute command is called called for the first time in a guild
func revokeChannelPerms(s *dg.Session, m *dg.MessageCreate) error {
	muteRoleID := "772777995025907732"

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

//												***		validUserID		***

func validUserID(s *dg.Session, m *dg.MessageCreate, id string) (string, bool) {

	id = str.Trim(id, "<>&!@#")
	_, err := s.GuildMember(m.GuildID, id)

	if err == nil {
		return id, true
	}
	return "", false
}

// 												***		validChannelID		***

func validChannelID(s *dg.Session, m *dg.MessageCreate, id string) (string, bool) {

	id = str.Trim(id, "<>&!@#")
	chn, err := s.Channel(id)
	if chn.GuildID != m.GuildID {
		return "", false
	}

	if err == nil {
		return id, true
	}
	return "", false
}

// 												***		validRoleID		***

func validRoleID(s *dg.Session, m *dg.MessageCreate, id string) (string, bool) {

	id = str.Trim(id, "<>&!@#")
	err := s.GuildMemberRoleAdd(m.GuildID, s.State.User.ID, id)
	s.GuildMemberRoleRemove(m.GuildID, s.State.User.ID, id)
	if err == nil {
		return id, true
	}
	return "", false
}

//												***		fieldsN		***
// fieldN is just a upgraded version of strings.Fields the extra thing that it does is,
// it returns slice of first 'N' substrings of the passed input
// it is similer like Split and SplitN but without any white spaces in any slice like Fields
// #For error checking: if the len of the returned slice is 0 then an error has occured
func fieldsN(s string, n int) []string {
	if len(s) == 0 || n < 1 {
		return []string{}
	}
	if n == 1 {
		return []string{s}
	}
	var end int = 1
	var i int
	for i = range s {
		if i != 0 && s[i-1] == ' ' {
			continue
		}
		if s[i] == ' ' {
			end++
		}
		if s[i] != ' ' && end == n {
			break
		}
	}
	if end != n {
		return []string{}
	}
	args := str.Fields(s[:i-1])
	args = append(args, s[i-1:]) // appending rest of the string

	return args
}
