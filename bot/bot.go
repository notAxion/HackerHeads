package bot

import (
	"fmt"
	str "strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/notAxion/HackerHeads/config"
	"github.com/notAxion/HackerHeads/features"
)

var BotID string

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
	// goBot.AddHandler(invite)

	goBot.AddHandler(manageChannels)

	goBot.AddHandler(features.EventRoleAdd)
	goBot.AddHandler(features.EventRoleRemove)

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
	s.UpdateGameStatus(0, "Bot in Development")
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
		case "mute":
			features.Mute(s, m)
		case "ping":
			features.Ping(s, m)
		case "remind":
			features.Remind(s, m)
		case "warn":
			features.Warn(s, m)
		default:
			s.ChannelMessageSend(m.ChannelID, "**WHAT !?** :thinking::thinking:")
		}
	}

	if m.Content == "hi bot" || m.Content == "Hi bot" {
		s.ChannelMessageSendReply(m.ChannelID, "```Hey```", m.MessageReference)
	}

	// if m.ChannelID == "765909517879476224" {
	// fmt.Println(len(m.Attachments))
	// }
}

// func invite(s *dg.Session, inv *dg.Invite) {
// 	fmt.Println(inv.Inviter.Username)
// 	s.ChannelMessageSend("765909517879476224", "invite created")
// }

func test(s *dg.Session, m *dg.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}
	// if len(m.Content) < 4 {
	// return
	// }

	// if len(m.Content) > 3 && m.Content[:4] != "test" {
	if len(m.Content) > 3 && m.Content != "test" {
		return
	}
	inv, err := s.GuildInvites(m.GuildID)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, iv := range inv {
		code := iv.Code
		tempInv, err := s.Invite(code)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%#v", tempInv)
	}

	// if features.ValidRoleID(s, m, m.Content[5:]) {
	// s.ChannelMessageSend(m.ChannelID, "exists")
	// } else {
	// s.ChannelMessageSend(m.ChannelID, "nope")
	// }
	// args := [2]string{m.Content[5:23], m.Content[24:31]}
	// fmt.Println(args[0], args[1])
	// s.ChannelMessageSend(args[0], args[1])

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

/*
func maxmsgCount(s *dg.Session, st *dg.State) {
	st.MaxMessageCount = 500
} */

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
