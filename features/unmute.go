package features

import (
	"fmt"
	"strings"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/notAxion/HackerHeads/config"
	"github.com/notAxion/HackerHeads/db"
)

// 												***		U N M U T E 	***

// *todo make a helpUnmute
func (r *Mux) Unmute(s *dg.Session, m *dg.MessageCreate) {
	args := fieldsN(m.Content, -1)
	if len(args) < 2 {
		r.helpMute(s, m.ChannelID) //!valid or just checking help mute
		return
	}
	user, valid := r.validUserID(s, m, args[1])
	if !valid {
		r.helpMute(s, m.ChannelID)
		msg := fmt.Sprintf("I can't find that user, %s", args[1])
		s.ChannelMessageSend(m.ChannelID, InCodeBlock(msg))
		return
	}
	k := db.FromString(m.GuildID, user.ID)
	if done, ok := r.muteDone[k]; ok {
		close(done)
	}
	err := r.UnmuteUser(s, m.GuildID, user.ID)
	if err == ErrPerms {
		r.helpMute(s, m.ChannelID)
		msg := fmt.Sprintf("can't remove role of %s#%s", user.Username, user.Discriminator)
		s.ChannelMessageSend(m.ChannelID, InCodeBlock(msg))
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	var reason string
	if len(args) > 2 {
		reason = strings.Join(args[3:], " ")
	}

	// sending final message of success to the channel
	err = r.UnmuteReply(s, m.ChannelID, user, reason)
	if err != nil {
		fmt.Println(err)
	}

}

func (r *Mux) UnmuteUser(s *dg.Session, guildID, userID string) error {
	muteRoleID, err := r.muteRole(s, guildID)
	if err != nil {
		fmt.Println("muteRole err")
		return err
	}
	if err = s.GuildMemberRoleRemove(guildID, userID, muteRoleID); err != nil {
		fmt.Println("Remove role error.\n", err)
		return ErrPerms
	}
	return nil
}

func (r *Mux) UnmuteReply(s *dg.Session, channelID string, user *dg.User, res string) error {
	title := fmt.Sprintf(`
		:white_check_mark: *%s#%s has been unmuted %s %s*
		`, user.Username, user.Discriminator, res, "\u200e")
	unmuteEmbed := &dg.MessageEmbed{
		Type:  "rich",
		Title: title,
		Color: 0x00fa00,
	}
	if _, err := s.ChannelMessageSendEmbed(channelID, unmuteEmbed); err != nil {
		return err
	}
	return nil
}

// SetAllMutedTimer is only for init which will set the timer
// of all muted user which were in the db
func (r *Mux) SetAllMutedTimer() {
	users, err := r.PQ.TGetMutedUsers()
	if err != nil {
		// don't know what to do with it
		fmt.Println(err)
	}
	for user, t := range users {
		r.AddTimer(user, t)
	}

}

func (r *Mux) AddTimer(user db.User, t time.Time) {
	done := make(chan struct{})
	r.muteDone[user] = done
	go r.timer(user, t, done)
	go func() {
		<-done
		delete(r.muteDone, user)
	}()
}

func (r *Mux) timer(user db.User, t time.Time, done chan struct{}) {
	s, err := dg.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err)
		return
	}
	dur := time.Until(t)
	timer := time.NewTimer(dur)
	defer timer.Stop() // doubt
	select {
	case <-timer.C:
		gID, uID := user.ToString()
		r.UnmuteUser(s, gID, uID)
		close(done)
	case <-done:
		// fmt.Println("early unmute")
	}
	err = r.PQ.DeleteUnmuteTime(user.GID, user.UID)
	if err != nil {
		fmt.Println(err)
		return
	}
}
