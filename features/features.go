package features

import (
	"fmt"
	str "strings"
	t "time"

	"github.com/bwmarrin/discordgo"
)

func Remind(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == ".remind" || m.Content == ".remind " {
		helpRemind(s,m)
		return
	}
	rem := make([]string, 3)
	rem = str.SplitN(m.Content, " ", 3)
	timer, err := t.ParseDuration(rem[1])
	if err != nil {
		timeError := &discordgo.MessageEmbed {
			Type: "rich",
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
	remindEmbed := &discordgo.MessageEmbed {
		Type: "rich",
		Title: "Reminder is set",
		Color: 0x00ff00,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, remindEmbed)
	
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Sleep(timer)
	s.ChannelMessageSend("766025683570524190", fmt.Sprintf("<@%s> here is the reminder you asked \n%s", m.Author.ID, rem[2]))

}

func helpRemind(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	//helpStr := fmt.Sprintf("\n**Command**: %s \nDescription: Will tag you at a specified channel so that it can remind you.\nUsage: %cremind [time] [channel tag] [Reminder] \nExample:\n\t%cremind 10m #general turn off microwave\n\t%cremind 1h #poles start a pole", m.Content, m.Content[0], m.Content[0], m.Content[0])
	desc := fmt.Sprintf("\n**Description**: Will tag you at a specified channel so that it can remind you.\n**Usage**: %cremind [time] [channel tag] [Reminder] \n**Example**:\n\t%cremind 10m #general turn off microwave\n\t%cremind 1h #poles start a pole", m.Content[0], m.Content[0], m.Content[0])
	helpEmbed := &discordgo.MessageEmbed {
		Type: "rich",
		Title: fmt.Sprintf("\n**Command**: %s", m.Content),
		Description: desc,
		Color: 0x00ff00,
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, helpEmbed)
	
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func Warn(s *discordgo.Session, m *discordgo.MessageCreate ) {
	
	args := make([]string, 3)
	args = str.SplitN(m.Content, " ", 3)
	warnID, valid := validID(s, m, args[1])
	if !valid {
		idError := &discordgo.MessageEmbed {
			Type: "rich",
			Title: fmt.Sprintf("**I can't find that user, %s**", args[1]),
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
	warnEmbed := &discordgo.MessageEmbed {
		Type: "rich",
		Title: fmt.Sprintf("*%s#%s have been warned : %s ‎*", user.Username, user.Discriminator, args[2]),
		Color: 0xBC4F07,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, warnEmbed)
	
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}


func validID( s *discordgo.Session, m *discordgo.MessageCreate, id string) (string, bool)  {
	
	if m.Author.ID == s.State.User.ID {
		return "", false
	}
	//temp := *id
	//fmt.Printf("%T", *id)
	id = str.Trim(id, "<>&!@#")
	_, err:= s.GuildMember(m.GuildID, id)
	
	if err == nil {
		return id, true
	}
	return id, false
}
