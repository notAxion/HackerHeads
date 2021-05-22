package features

import (
	"fmt"
	str "strings"

	dg "github.com/bwmarrin/discordgo"
)

// 													***		E V E N T 	***

// EventStart will start an instance of an event for that channel
// So afterwards if any member of that event types any message it will give that member a role which should be specified when event start command was sent
// and removes the role when a message is deleted within the event period
func (r *Mux) EventStart(s *dg.Session, m *dg.MessageCreate) { //# todo check the role hierarchy to check if person that is sending the command does himself has the perms to add that role someone else
	if m.Content[1:] == "event" || m.Content[1:] == "event " {
		r.helpEvent(s, m.ChannelID)
		return
	}
	args := fieldsN(m.Content[1:], 3)
	if len(args) == 0 {
		//r.helpEvent(s, m.ChannelID)
		return
	}
	if args[1] != "start" {
		r.helpEvent(s, m.ChannelID)
		return
	}
	eventRoleID := args[2]
	valid := r.validRoleID(s, m, eventRoleID)
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

// EventRoleAdd will be adding roles to the users after the event is started
// also this will handle the event stop command
func (r *Mux) EventRoleAdd(s *dg.Session, m *dg.MessageCreate) {
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
func (r *Mux) EventRoleRemove(s *dg.Session, m *dg.MessageDelete) {
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

func (r *Mux) helpEvent(s *dg.Session, chnID string) {
	desc := fmt.Sprintf(`
**Description**:  EventStart will start an instance of an event for that channel
So afterwards if any member of that event types any message
`+` it will give that member a role which should be specified when event start command was sent
and removes the role when a message is deleted within the event period.
*It is advisable to create a new channel and then start the event and dont reuse it for other events* 
**Usage**: %sevent < [start] || [stop] > {@role}
**Example**:
%sevent start @participant 
%sevent stop
	`, botPrefix, botPrefix, botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       "\n**Command**: event",
		Description: desc,
		Color:       0x00fa00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
	}
}
