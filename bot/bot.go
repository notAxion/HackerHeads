package bot

import (
	"fmt"
	str "strings"

	"../config"
	"../features"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

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

	err = goBot.Open()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Bot is Running")
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	//set the playing status 
	s.UpdateStatus(0, "Bot in Development")
	//since := 1 
	/*
	status := discordgo.UpdateStatusData{
		//IdleSince: &since,
		Game: &discordgo.Game{
						Name: "Pydroid & Termux",
						Type: 3,
						Details: "Making this bot that you are watching in go",
						TimeStamps: discordgo.TimeStamps{
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if str.HasPrefix(m.Content, config.BotPrefix) {

		if m.Author.ID == BotID {
			return
		}
		cmd := command(m.Content)
		fmt.Println(config.BotPrefix)
		switch cmd {
		case "remind":
			features.Remind(s, m)
			break
		case "warn":
			features.Warn(s, m)
			break
		default:
			s.ChannelMessageSend(m.ChannelID, "**WHAT !?** :thinking::thinking:" )
		}
	}

	if m.Content == "hi bot" || m.Content == "Hi bot" {
		s.ChannelMessageSend(m.ChannelID, "```Hello Hackers```")
	}
	if m.ChannelID == "765909517879476224" {
		fmt.Println(m.Content)
	}
}

func test(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	if m.Author.ID == BotID {
		return
	}
	//s.MessageReactionAdd("765909517879476224","768750535464714271","232720527448342530")
	
}

/*
func validID(s *discordgo.Session, m *discordgo.MessageCreate, id string ) bool  {
	
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
