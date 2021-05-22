package features

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

//												***		W A R N		***

// Warn command that will warn that user upon the specified reason and dm them if possible (todo add logging to a DB)
func (r *Mux) Warn(s *dg.Session, m *dg.MessageCreate) {
	// Help warn
	if m.Content[1:] == "warn" || m.Content[1:] == "warn " {
		r.helpWarn(s, m.ChannelID)
		return
	}
	//  Getting the args
	args := fieldsN(m.Content[1:], 3)
	if len(args) == 0 {
		r.helpWarn(s, m.ChannelID)
	}
	user, valid := r.validUserID(s, m, args[1])
	if !valid { // provided id is not valid
		idError := &dg.MessageEmbed{
			Type:  "rich",
			Title: fmt.Sprintf(":x: **I can't find that user, %s**", args[1]),
			Color: 0xff0000,
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, idError)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	warnEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: fmt.Sprintf(":white_check_mark: *%s#%s has been warned. \u200e*", user.Username, user.Discriminator),
		Color: 0x00fa00,
	}

	//  Creates dm channel for the the warned user
	warnDMChn, err := s.UserChannelCreate(user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	guild, _ := s.Guild(m.GuildID) // ** todo try some other way to get the guild name cause i have to get all values just to get the guild name
	_, err = s.ChannelMessageSend(warnDMChn.ID, fmt.Sprintf("You were warned in %s for %s", guild.Name, args[2]))
	if err != nil {
		fmt.Println(err)
		return
	}

	//  Sending the embeded text
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, warnEmbed)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//												*** 	helpWarn	***

func (r *Mux) helpWarn(s *dg.Session, chnID string) {
	desc := fmt.Sprintf("\n**Description**: warn those noobs who don't follow the rules.\n**Usage**: %swarn [@user] [reason]  \n**Example**:\n\t%swarn @noobSpammer stop the spam please", botPrefix, botPrefix)
	helpEmbed := &dg.MessageEmbed{
		Type:        "rich",
		Title:       "\n**Command**: warn",
		Description: desc,
		Color:       0x00ff00,
	}
	_, err := s.ChannelMessageSendEmbed(chnID, helpEmbed)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
