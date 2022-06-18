package features

import (
	"fmt"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

// don't commit until webhook ratelimit figures out

type relayConfig struct {
	toChanID  string
	hookID    string
	hookToken string

	// Restrictions
	links bool
	bots  bool
}

func (r *Mux) ConnectServer(s *dg.Session, m *dg.MessageCreate) {
	if oldChan, ok := r.ConnCfg.connectMap[m.Author.ID]; ok {
		// ask for all the configs
		close(r.ConnCfg.StartDone[m.Author.ID])

		// caching the config vars
		hook, err := s.WebhookCreate(oldChan, "idk", "")
		if err != nil {
			fmt.Println("create webhook 1 error")
			fmt.Println(err)
			return
		}
		r.ConnCfg.guildConfig[m.ChannelID] = &relayConfig{
			toChanID:  oldChan,
			hookID:    hook.ID,
			hookToken: hook.Token,
			links:     true,
			bots:      false,
		}
		hook, err = s.WebhookCreate(m.ChannelID, "idk", "")
		if err != nil {
			fmt.Println("create webhook 2 error")
			fmt.Println(err)
			return
		}
		r.ConnCfg.guildConfig[oldChan] = &relayConfig{
			toChanID:  m.ChannelID,
			hookID:    hook.ID,
			hookToken: hook.Token,
			links:     true,
			bots:      false,
		}
		s.ChannelMessageSend(m.ChannelID, ":white_check_mark: connection started")
	} else {
		done := make(chan struct{})
		r.ConnCfg.StartDone[m.Author.ID] = done

		go r.afterFunc(30*time.Second, done, func() {
			delete(r.ConnCfg.connectMap, m.Author.ID)
			delete(r.ConnCfg.StartDone, m.Author.ID)
		})
		r.ConnCfg.connectMap[m.Author.ID] = m.ChannelID
		s.ChannelMessageSend(m.ChannelID, "Connection in progress *You have 30s to start the connection")
	}
}

func (r *Mux) ConnectServerRelay(s *dg.Session, m *dg.MessageCreate) {
	if conf, ok := r.ConnCfg.guildConfig[m.ChannelID]; ok {
		// s.ChannelMessageSend(conf.toChanID, m.Content)
		if m.Author.Bot && !conf.bots {
			// only enters when 1st is true and 2nd is false
			return
		}
		// check links
		if m.Content == "" {
			return
		}
		r.anonUser(s, m, conf.toChanID)
	} else {
		// get it from database
	}
}

func (r *Mux) anonUser(s *dg.Session, m *dg.MessageCreate, toChan string) {
	relayConf := r.ConnCfg.guildConfig[m.ChannelID]
	hookID, hookToken := relayConf.hookID, relayConf.hookToken
	_, err := s.WebhookExecute(hookID, hookToken, true, &dg.WebhookParams{
		Username:  m.Author.Username,
		Content:   m.Content,
		AvatarURL: m.Author.AvatarURL(""),
	})
	if err != nil {
		fmt.Println("send webhook error")
		fmt.Println(err)
		return
	}
	// fmt.Println(msg.Content)

}

func (r *Mux) CleanWebhook(s *dg.Session, m *dg.MessageCreate) { // don't commit
	hooks, _ := s.GuildWebhooks(m.GuildID)
	for _, hook := range hooks {
		fmt.Println("cleaning guild")
		s.WebhookDelete(hook.ID)
	}
	hooks, _ = s.ChannelWebhooks(m.ChannelID)
	for _, hook := range hooks {
		fmt.Println("cleaning channel")
		s.WebhookDelete(hook.ID)
	}
}

// func (r *Mux) ConnectServerStart(cancel context.CancelFunc, s *dg.Session, m *dg.MessageCreate) {
// }

func (r *Mux) afterFunc(dur time.Duration, done chan struct{}, f func()) {
	timer := time.NewTimer(dur)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-done:
	}
	f()
}
