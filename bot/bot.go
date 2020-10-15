package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"../config"
	"../features"
	str "strings"
)

var BotID string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New ("Bot " + config.Token )
	
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
	
	err = goBot.Open()
	
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Bot is Running")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	if str.HasPrefix(m.Content, config.BotPrefix) {
		
		if m.Author.ID == BotID {
			return
		}
		cmd := command(m.Content)
		switch cmd {
			case "remind" :
				features.Remind(s, m )
			default :
				s.ChannelMessageSend(m.ChannelID,"**BAD INPUT**")
		}
	}
	
	if m.Content == "hi bot" || m.Content == "Hi bot" {
		s.ChannelMessageSend(m.ChannelID, "```Hello Hackers```")
	}
	if m.ChannelID == "765909517879476224" {
		fmt.Println(m.Content)
	}
}

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