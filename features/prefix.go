package features

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

func (r *Mux) Prefix(s *dg.Session, m *dg.MessageCreate) {
	args := fieldsN(m.Content[1:], 2)
	if len(args) < 2 {
		// r.helpPrefix() *TODO
		fmt.Println("not enough args")
		return
	}
	if len(args[1]) > 10 {
		s.ChannelMessageSendReply(m.ChannelID, InCodeBlock(`
Hey hooman,
be a hooman and use something up to 10 characters.
`), m.Reference())
		return
	}
	if len(args[1]) != 1 {
		s.ChannelMessageSendReply(m.ChannelID, InCodeBlock(`
only one character prefix is supported
contact the dev if you are seeing this.
`), m.Reference()) // *TODO make a router
		fmt.Printf("%#v\n", args)
		return
	}
	err := r.PQ.UpsertPrefix(m.GuildID, args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	r.GuildPrefix[m.GuildID] = args[1]
	s.ChannelMessageSendReply(m.ChannelID, InCodeBlock(`
The prefix has been changed to `+args[1]), m.Reference())
}

func (r *Mux) SetAllPrefix() {

}
