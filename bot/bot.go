package bot

import (
	"fmt"
	str "strings"

	"../config"
	"../features"
	dg "github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *dg.Session

func Start() {
	goBot, err := dg.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err)
	}

	BotID = u.ID

	goBot.AddHandler(messageCreate)

	goBot.AddHandler(ready)

	goBot.AddHandler(test)

	goBot.AddHandler(manageChannels)

	goBot.AddHandler(features.EventRoleAdd)

	//goBot.AddHandler(maxmsgCount)
	//st = *dg.State
	//st.MaxMessageCount = 500
	err = goBot.Open()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Bot is Running")
}

func ready(s *dg.Session, event *dg.Ready) {
	//set the playing status
	s.UpdateStatus(0, "Bot in Development")
	s.StateEnabled = true
	st := s.State
	st.MaxMessageCount = 500
	//since := 1

	/*
		status := dg.UpdateStatusData{
			//IdleSince: &since,
			Game: &dg.Game{
							Name: "Pydroid & Termux",
							Type: 3,
							Details: "Making this bot that you are watching in go",
							TimeStamps: dg.TimeStamps{
													StartTimestamp: 1602918742000,
												},
						},
			AFK: false,
			Status: "test",
		}
		err := s.UpdateStatusComplex(status)
		if err != nil {
			fmt.Println(err)
			return
		}
	*/
	//s.UpdateListeningStatus("Mr.Makra")

}

func messageCreate(s *dg.Session, m *dg.MessageCreate) {

	if str.HasPrefix(m.Content, config.BotPrefix) {

		if m.Author.ID == BotID || m.Author.Bot {
			return
		}
		cmd := command(m.Content)
		switch cmd {
		case "event":
			features.EventStart(s, m)
			break
		case "mute":
			features.Mute(s, m)
			break
		case "ping":
			features.Ping(s, m)
			break
		case "remind":
			features.Remind(s, m)
			break
		case "warn":
			features.Warn(s, m)
			break
		default:
			s.ChannelMessageSend(m.ChannelID, "**WHAT !?** :thinking::thinking:")
		}
	}

	if m.Content == "hi bot" || m.Content == "Hi bot" {
		s.ChannelMessageSend(m.ChannelID, "```Hello Hackers```")
	} /*
		if m.ChannelID == "765909517879476224" {
			fmt.Println(m.Content)
		}*/
}

func test(s *dg.Session, m *dg.MessageCreate) {

	if m.Content != "test" {
		return
	}
	// id , valid := features.ValidRoleID(s, m, m.Content[5:])
	// if valid {
	// 	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@&%s> exists", id))
	// } else {
	// 	s.ChannelMessageSend(m.ChannelID, "Doesnt exist")
	// }
}

func manageChannels(s *dg.Session, chans *dg.ChannelCreate) {

	features.AddMuteRole(s, chans)

}

func maxmsgCount(s *dg.Session, st *dg.State) {
	st.MaxMessageCount = 500
}

/*
func validID(s *dg.Session, m *dg.MessageCreate, id string ) bool  {

	if m.Author.ID == s.State.User.ID {
		return false
	}
	id = str.Trim(id, "<>&!@#")
	_, err:= s.GuildMember(m.GuildID, m.Content)

	if err == nil {
		return true
	}
	return false
}*/

func command(s string) string {
	cmd := make([]rune, 0, 10)
	for i, val := range s {
		if i == 0 {
			continue
		}
		if val != ' ' {
			cmd = append(cmd, val)
		} else {
			break
		}
	}
	return string(cmd)
}
