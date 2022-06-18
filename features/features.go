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

type connConfig struct {
	// this two strings should be m.Author.ID
	connectMap  map[string]string
	StartDone   map[string]chan struct{}
	MsgCh       chan *dg.Message
	guildConfig map[string]*relayConfig
}

// Mux have all the handlers for every features
// use NewRouter to create a new Mux
type Mux struct {
	botPrefix string
	botID     string
	PQ        *db.DB
	// stores the eventRoleID for each channel (map[ChannelID]eventRoleID)
	// go to EventStop for more info # thanks to chanbakjsd from gophers
	eventMap    map[string]string
	memeMap     map[string][]meme
	topMemes    map[string]*memeCache
	ConnCfg     *connConfig
	muteDone    map[db.User]chan struct{}
	Cooldown    map[string]bool
	GuildPrefix map[string]string
}

func NewRouter(bID string) *Mux {
	r := &Mux{
		botPrefix: config.BotPrefix,
		botID:     bID,
		PQ:        db.NewDB(),
		eventMap:  make(map[string]string),
		memeMap:   make(map[string][]meme),
		topMemes:  make(map[string]*memeCache),
		ConnCfg: &connConfig{
			connectMap:  make(map[string]string),
			StartDone:   make(map[string]chan struct{}),
			guildConfig: make(map[string]*relayConfig),
		},
		muteDone:    make(map[db.User]chan struct{}),
		Cooldown:    make(map[string]bool),
		GuildPrefix: make(map[string]string),
	}
	r.SetAllMutedTimer()
	r.SetAllPrefix()
	go r.cacheTopMemes()
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

func (r *Mux) Whois(s *dg.Session, m *dg.MessageCreate) {
	args := fieldsN(m.Content[1:], -1)
	if len(args) < 2 {
		return
	}
	user, valid := r.validUserID(s, m, args[1])
	if !valid {
		s.ChannelMessageSend(m.ChannelID, "nope")
		return
	}
	s.ChannelMessageSend(m.ChannelID, user.Username)
}
