package features

import (
	"fmt"
	str "strings"
	t "time"

	"../config"

	dg "github.com/bwmarrin/discordgo"
)

//											***		P I N G 		***

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

func Remind(s *dg.Session, m *dg.MessageCreate) {
	if m.Content[1:] == "remind" || m.Content[1:] == "remind " {
		helpRemind(s, m.ChannelID)
		return
	}
	rem := fieldsN(m.Content, 3)
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

	desc := fmt.Sprintf("\n**Description**: Will tag you at the channel you have sent the message so that it can remind you.\n**Usage**: %sremind [time] [Reminder] \n**Example**:\n\t%sremind 10m turn off microwave\n\t%sremind 1h start a pole", config.BotPrefix, config.BotPrefix, config.BotPrefix)
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

//												***		W A R N		***

func Warn(s *dg.Session, m *dg.MessageCreate) {

	if m.Content[1:] == "warn" || m.Content[1:] == "warn " {
		helpWarn(s, m.ChannelID)
	}
	args := fieldsN(m.Content, 3)
	if len(args) == 0 {
		helpWarn(s, m.ChannelID)
	}
	warnID, valid := validUserID(s, m, args[1])
	if !valid {
		idError := &dg.MessageEmbed{
			Type:  "rich",
			Title: fmt.Sprintf(":x: **I can't find that user, %s**", args[1]),
			Color: 0xff0000,
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, idError)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
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
	warnChn, err := s.UserChannelCreate(warnID)
	guild, _ := s.Guild(m.GuildID)
	_, err = s.ChannelMessageSend(warnChn.ID, fmt.Sprintf("You were warned in %s for %s", guild.Name, args[2]))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, warnEmbed)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//												*** 	helpWarn	***

func helpWarn(s *dg.Session, chnID string) {
	desc := fmt.Sprintf("\n**Description**: warn those noobs who don't follow the rules.\n**Usage**: %swarn [@user] [reason]  \n**Example**:\n\t%swarn @noobSpammer stop the spam please", config.BotPrefix, config.BotPrefix)
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

//												***		fieldsN		***

func fieldsN(s string, n int) []string {
	if len(s) == 0 {
		return []string{}
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
			break // i - 1 is the endpoint that i will pass
		}
	}
	if end != n {
		return []string{}
	}
	args := str.Fields(s[:i-1])
	args = append(args, s[i-1:]) // appending rest of the string

	return args
}
