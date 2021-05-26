package bot

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
	"github.com/notAxion/HackerHeads/config"
	"github.com/notAxion/HackerHeads/features"
)

// Start works like listenAndServe it will block until exit is called
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
	<-make(chan struct{})
}

func test(s *dg.Session, m *dg.MessageCreate) {
	if m.Author.Bot {
		return
	}

	// if len(m.Content) > 3 && m.Content[:4] != "test" {
	if len(m.Content) < 3 && m.Content != "test" {
		return
	}

}

