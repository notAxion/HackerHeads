package features

import (
	"fmt"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

//												***		validUserID		***

func (r *Mux) validUserID(s *dg.Session, m *dg.MessageCreate, id string) (*dg.User, bool) {

	id = strings.Trim(id, "<>&!@#")
	mem, err := s.GuildMember(m.GuildID, id)

	if err != nil {
		return nil, false
	}
	return mem.User, true
}

// 												***		validChannelID		***

func (r *Mux) validChannelID(s *dg.Session, m *dg.MessageCreate, id string) (string, bool) {

	id = strings.Trim(id, "<>&!@#")
	chn, err := s.Channel(id)
	if chn.GuildID != m.GuildID {
		return "", false
	}
	if err == nil {
		return id, true
	}
	return "", false
}

func (r *Mux) ValidRoleID(s *dg.Session, m *dg.MessageCreate, id string) bool {
	return r.validRoleID(s, m.GuildID, id)
}

// 												***		validRoleID		***

// not good use simple loop roles
func (r *Mux) validRoleID(s *dg.Session, guildID string, id string) bool {

	id = strings.Trim(id, "<>&!@#")
	gRoles, err := s.GuildRoles(guildID)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, gRole := range gRoles {
		if gRole.ID == id {
			return true
		}
	}
	return false
	/*
		err := s.GuildMemberRoleAdd(m.GuildID, s.State.User.ID, id)
		s.GuildMemberRoleRemove(m.GuildID, s.State.User.ID, id)
		if err == nil {
			return true
		}
		return false
	*/
}

func InCodeBlock(str string) string {
	return "```" + str + "```"
}

//												***		fieldsN		***

// fieldN is just a upgraded version of strings.Fields the extra thing that it does is,
// it returns slice of first 'N' substrings of the passed input
// it is similer like Split and SplitN but without any white spaces in any slice like Fields
// #For error checking: if the len of the returned slice is 0 then an error has occured
func fieldsN(s string, n int) []string {
	if len(s) == 0 || n == 0 {
		return []string{}
	}
	if n < 0 {
		return strings.Fields(s)
	}
	if n == 1 {
		return []string{s}
	}
	end, i := 1, 0
	for i = range s {
		if i != 0 && s[i-1] == ' ' {
			continue
		}
		if s[i] == ' ' {
			end++
		}
		if s[i] != ' ' && end == n {
			break
		}
	}
	if end != n {
		return []string{}
	}
	args := strings.Fields(s[:i-1])
	args = append(args, s[i-1:]) // appending rest of the string

	return args
}
