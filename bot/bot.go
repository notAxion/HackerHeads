package bot

import (
	"fmt"
	"strings"
	"time"

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
	// goBot.AddHandler(getEmoji)
	goBot.AddHandler(r.ReactionListen)
	// goBot.AddHandler(testSlash)
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

	// if len(m.Content) < 4 || m.Content[:4] != "test" {
	if m.Content != "test" {
		return
	}
	trackuser := m.ChannelID + m.Author.ID
	// done := make(chan struct{})
	// func (done chan<- struct{})  {
	// 	close(done)
	// }(done)
	msgCh := make(chan string, 10)
	remove := s.AddHandler(func(s *dg.Session, m *dg.MessageCreate) {
		handleruser := m.ChannelID + m.Author.ID
		if trackuser == handleruser {
			if m.Content == "return 0" {
				close(msgCh)
				return
			}
			msgCh <- m.Content
			// test2(s, m)
			// close(done)
		}
	})
	// fmt.Println("it does")
	// <-done
	arr := make([]string, 0, 10)
	t := time.AfterFunc(10*time.Second, func() { close(msgCh) })
	for msg := range msgCh {
		fmt.Println(t.Reset(10 * time.Second))
		arr = append(arr, msg)
	}
	fmt.Println(t.Stop())
	remove()
	s.ChannelMessageSend(m.ChannelID, strings.Join(arr, " + "))

}

func getEmoji(s *dg.Session, e *dg.MessageReactionAdd) {
	fmt.Println(e.Emoji.Name)
	fmt.Printf("%#v\n", e.Emoji)
}
