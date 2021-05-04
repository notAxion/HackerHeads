package features

import (
	"fmt"
	str "strings"
	t "time"

	"github.com/notAxion/HackerHeads/config"

	dg "github.com/bwmarrin/discordgo"
)

// feature : .dc @mention to disconnect someone @me diconnect the user
// how about a website which would control the bot 🤔

var botPrefix string = config.BotPrefix

// #todo add a variable aka var = botPrefix and replace every botPrefix with that var
// #and also change every help funcs and use different variables for the desc.

// stores the eventRoleID for each channel (map[ChannelID]eventRoleID)
// go to EventStop for more info # thanks to chanbakjsd from gophers
var eventMap map[string]string = make(map[string]string)

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
