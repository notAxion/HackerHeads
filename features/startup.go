package features

import (
	dg "github.com/bwmarrin/discordgo"
)

func (r *Mux) Ready(s *dg.Session, event *dg.Ready) {
	//set the playing status
	s.UpdateGameStatus(0, "Bot in Development")
	s.StateEnabled = true
	st := s.State
	st.MaxMessageCount = 500
	//since := 1

	/*
		status := dg.UpdateStatusData{
			//IdleSince: &since,
			Game: &dg.Game{
							Name: "Pydroid & Termux",
							Type: 3,
							Details: "Making this bot that you are watching in go",
							TimeStamps: dg.TimeStamps{
													StartTimestamp: 1602918742000,
												},
						},
			AFK: false,
			Status: "test",
		}
		err := s.UpdateStatusComplex(status)
		if err != nil {
			fmt.Println(err)
			return
		}
	*/
	//s.UpdateListeningStatus("Mr.Makra")

}

// func invite(s *dg.Session, inv *dg.Invite) {
// 	fmt.Println(inv.Inviter.Username)
// 	s.ChannelMessageSend("765909517879476224", "invite created")
// }

func (r *Mux) ManageChannels(s *dg.Session, chans *dg.ChannelCreate) {

	r.AddMuteRole(s, chans)

}

/*
func maxmsgCount(s *dg.Session, st *dg.State) {
	st.MaxMessageCount = 500
} */

/*
func validID(s *dg.Session, m *dg.MessageCreate, id string ) bool  {

	if m.Author.ID == s.State.User.ID {
		return false
	}
	id = str.Trim(id, "<>&!@#")
	_, err:= s.GuildMember(m.GuildID, m.Content)

	if err == nil {
		return true
	}
	return false
}*/

func command(s string) string {
	cmd := make([]byte, 0, 7)
	var i int
	for i = 0; i < len(s); i++ {
		if s[i] != ' ' {
			break
		}
	}
	for ; i < len(s); i++ {
		val := s[i]
		if val != ' ' {
			if 'A' <= val && val <= 'Z' {
				val += 'a' - 'A'
			}
			cmd = append(cmd, val)
		} else {
			break
		}
	}
	return string(cmd)
}
