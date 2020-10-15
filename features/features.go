package features

import (
	"fmt"
	str "strings"
	t "time"
	"github.com/bwmarrin/discordgo"
	
)

func Remind(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	rem := make([]string, 3)
	rem = str.SplitN(m.Content, " ", 3)
	timer, err := t.ParseDuration(rem[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "*Wrong time format for remind ex-> 1m , 12h50m , 10h50m55s*")
		return 
	}
	t.Sleep(timer)
	s.ChannelMessageSend("766025683570524190", fmt.Sprintf("<@%s> here is the reminder you asked \n%s", m.Author.ID, rem[2]))

}