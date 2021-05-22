package bot

import (
	"fmt"

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

	r := features.NewRouter(u.ID)

	goBot.AddHandler(test)
	goBot.AddHandler(r.Ready)
	goBot.AddHandler(r.MessageCreate)
	goBot.AddHandler(r.ManageChannels)
	goBot.AddHandler(r.EventRoleAdd)
	goBot.AddHandler(r.EventRoleRemove)

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
