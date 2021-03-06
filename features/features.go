package features

import (
	"errors"
	"fmt"
	"time"

	"github.com/notAxion/HackerHeads/config"
	"github.com/notAxion/HackerHeads/db"

	dg "github.com/bwmarrin/discordgo"
)

func init() {
	// SetAllMutedTimer()
}

type User struct {
	GID int64
	UID int64

}

// Mux have all the handlers for every features
// use NewRouter to create a new Mux
type Mux struct {
	botPrefix string
	botID     string
	PQ        *db.DB
	// stores the eventRoleID for each channel (map[ChannelID]eventRoleID)
	// go to EventStop for more info # thanks to chanbakjsd from gophers
	eventMap map[string]string
	muteDone map[db.User]chan struct{} // change it to a type of two 2 int64(s)
}

func NewRouter(bID string) *Mux {
	r := &Mux{
		botPrefix: config.BotPrefix,
		botID:     bID,
		PQ:        db.NewDB(),
		eventMap:  make(map[string]string),
		muteDone:  make(map[db.User]chan struct{}),
	}
	r.SetAllMutedTimer()
	return r
}

var ErrArgs = errors.New("features: wrong argument passed")
var ErrPerms = errors.New("features: don't have the damn perms")

//											***		P I N G 	***

//Ping : checks the ping of the bot in millisecond
func (r *Mux) Ping(s *dg.Session, m *dg.MessageCreate) {
	u := time.Now()
	msg, _ := s.ChannelMessageSend(m.ChannelID, "pong -")
	//v, _ := msg.Timestamp.Parse()
	v := time.Now()
	ping := v.Sub(u)
	s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("pong - `%dms`", ping.Milliseconds()))
}