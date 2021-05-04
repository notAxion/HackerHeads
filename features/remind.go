package features

import (
	"fmt"
	t "time"

	dg "github.com/bwmarrin/discordgo"
)

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

	desc := fmt.Sprintf(`
**Description**: Will tag you at the channel you have sent the message so that it can remind you.
**Usage**: %sremind [time] [Reminder] 
**Example**:
%sremind 10m turn off microwave
%sremind 1h start a pole
	`, botPrefix, botPrefix, botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       "\n**Command**: remind",
		Description: desc,
		Color:       0x00ff00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
