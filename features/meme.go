package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly/v2"
)

var (
	leftArrow  = "⬅"
	rightArrow = "➡"
	// emojiC         = "🇨"
	// emojiM         = "🇲"
	emojiNew      = "🆕"
	emojiOne      = "1️⃣"
	emojiTwo      = "2️⃣"
	emojiReturn   = "↩"
	reactLed_meme = "a:led_meme:854054433986052136"
	emojiLed_meme = "<a:led_meme:854054433986052136>"
)

type meme struct {
	ID, name, img string
}

type memeCache struct {
	id, box int
}

func (r *Mux) meme(s *dg.Session, m *dg.MessageCreate) {
	args := fieldsN(m.Content, 2)
	if len(args) == 0 {
		// helpMeme()
		fmt.Println("not enough args") // *TOOD delete this
		return
	}
	memes := r.saveMemeTem(args[1]) //  *TODO todo don't use user input directly use a validater
	// memes := []meme{
	// 	{
	// 		name: "idk",
	// 		img:  "http://i.imgflip.com/2/3bovka.jpg",
	// 	},
	// }
	if len(memes) == 0 {
		msg, _ := s.ChannelMessageSend(m.ChannelID, "not a single soul knows what meme is that try again kiddo")
		time.AfterFunc(30*time.Second, func() {
			s.ChannelMessageDelete(m.ChannelID, msg.ID)
		})
		fmt.Println("no meme")
		return
	}
	meme := memes[0]
	// memeEmb := &dg.MessageEmbed{
	// 	Title: meme.name,
	// 	Image: &dg.MessageEmbedImage{
	// 		URL: meme.img,
	// 	},
	// 	Footer: &dg.MessageEmbedFooter{
	// 		Text: meme.ID,
	// 	},
	// }
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, meme.msgMemeEmbed())
	if err != nil {
		fmt.Println("meme search embed error", err)
		return
	}
	r.memeMap[msg.ID] = memes // *TODO have to delete this after some time ...
	// fmt.Println(len(memes))
	s.MessageReactionAdd(m.ChannelID, msg.ID, leftArrow)
	s.MessageReactionAdd(m.ChannelID, msg.ID, emojiNew)
	s.MessageReactionAdd(m.ChannelID, msg.ID, reactLed_meme)
	s.MessageReactionAdd(m.ChannelID, msg.ID, rightArrow)

	// s.ChannelMessageSend(m.ChannelID /*"https://imgflip.com/meme/"+*/, id)
}

func (r *Mux) ReactionListen(s *dg.Session, e *dg.MessageReactionAdd) {
	if e.UserID == r.botID {
		return
	}
	msg, err := s.State.Message(e.ChannelID, e.MessageID)
	if err != nil {
		msg, err = s.ChannelMessage(e.ChannelID, e.MessageID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if memes, ok := r.memeMap[msg.ID]; ok {
		embeds := msg.Embeds
		if len(embeds) < 1 {
			fmt.Println("no embed to begin with")
			return
		}
		embed := embeds[0]
		in, err := strconv.Atoi(embed.Title[0:1])
		if err != nil {
			fmt.Println("add index lol")
			return
		}

		in -= 1 // index correction
		meme := memes[in]
		switch e.Emoji.Name {
		case leftArrow: // *TODO have delete it from the cache after certain time if they noone is interacting
			s.MessageReactionRemove(e.ChannelID, e.MessageID, leftArrow, e.UserID)
			if in == 0 {
				return
			}
			// in, err := strconv.Atoi(embed.Title[0:1])
			// if err != nil {
			// 	fmt.Println("add index lol")
			// 	return
			// }
			// // correct index
			// in -= 1
			meme := memes[in-1]
			s.ChannelMessageEditEmbed(e.ChannelID, e.MessageID, meme.msgMemeEmbed())
		case rightArrow:
			s.MessageReactionRemove(e.ChannelID, e.MessageID, rightArrow, e.UserID)
			// in, err := strconv.Atoi(embed.Title[0:1])
			// if err != nil {
			// fmt.Println("add index lol")
			// return
			// }
			if in == len(memes)-1 {
				return
			}
			// correct index
			// in -= 1
			meme := memes[in+1]
			s.ChannelMessageEditEmbed(e.ChannelID, e.MessageID, meme.msgMemeEmbed())
		case emojiNew: // *TODO delete the cache and show info how to caption the meme
			s.MessageReactionRemove(e.ChannelID, e.MessageID, emojiNew, e.UserID)
			// in, err := strconv.Atoi(embed.Title[0:1])
			// if err != nil {
			// 	fmt.Println("add index lol")
			// 	return
			// }

			// in -= 1 // index correction
			// s.ChannelMessageSend(e.ChannelID, memes[in].ID)
			// meme := memes[in]

			// note to have a undo button as well
			// edit the embed for info to makememe
			captionHelpEmbed := &dg.MessageEmbed{
				Title: embed.Title,
				Thumbnail: &dg.MessageEmbedThumbnail{
					URL: meme.img,
				},
				Description: `So you wonna make meme alright ...
				But where would you write the texts ? glad u asked
				at the moment you have 2 options 
				:one: use ` + r.botPrefix + "makememe" + ` command after getting the id
				:two: type every texts in each message
				react on the option for futher info `,
			}
			s.MessageReactionsRemoveAll(e.ChannelID, e.MessageID)
			/*msgHelp, _ :=*/ s.ChannelMessageEditEmbed(e.ChannelID, e.MessageID, captionHelpEmbed)
			s.MessageReactionAdd(e.ChannelID, e.MessageID, emojiOne)
			s.MessageReactionAdd(e.ChannelID, e.MessageID, emojiTwo)
			// make a undo button that would restore everything
			s.MessageReactionAdd(e.ChannelID, e.MessageID, emojiReturn)

		case emojiOne: // this for commands - show id
			s.MessageReactionRemove(e.ChannelID, e.MessageID, emojiOne, e.UserID)
			s.ChannelMessageSend(e.ChannelID, meme.ID)
		case emojiTwo: // msg version - show help - have to reply within certain period of time
			s.MessageReactionRemove(e.ChannelID, e.MessageID, emojiTwo, e.UserID)
			captionHelp2 := &dg.MessageEmbed{
				Title: embed.Title,
				Thumbnail: &dg.MessageEmbedThumbnail{
					URL: meme.img,
				},
				Description: `Ok now listen up you gotta idk what to write 
				so just send msg and return meme`,
			}
			s.ChannelMessageEditEmbed(e.ChannelID, e.MessageID, captionHelp2)
			msgCh := make(chan string, 5)
			cancelStr := "return meme"
			remove := r.userHandler(s, e.UserID+e.ChannelID, cancelStr, msgCh)
			t := time.AfterFunc(2*time.Minute, func() { close(msgCh) })

			texts := make([]string, 0, 6)
			for text := range msgCh {
				t.Reset(2 * time.Minute)
				texts = append(texts, text)
			}

			if !t.Stop() { // *TODO delete everything user no makememe
				content := e.Emoji.User.Mention() + "forgot about your meme nice. ||u suck||"
				s.ChannelMessageSend(e.ChannelID, content)
			}
			remove()

			// sending the meme
			tempMsg, _ := s.ChannelMessageSend(e.ChannelID, "your meme is on the way hol up")
			memeUrl, err := r.memeFactory(meme.ID, texts...)
			if err != nil {
				fmt.Println(err)
				return
			}
			go s.ChannelMessageDelete(e.ChannelID, tempMsg.ID)
			s.ChannelMessageSend(e.ChannelID, memeUrl) // *TODO maybe do a embed with mention
			s.MessageReactionsRemoveAll(e.ChannelID, e.MessageID)
		} // *TODO add a reaction to delete things immediately
	}
}

// would save the meme templates
func (r *Mux) saveMemeTem(args string) []meme {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("exit scraper")
		}
	}()
	v := url.Values{}
	api := "https://imgflip.com/search?"
	if args == "" {
		fmt.Println("nothing to search +_+")
		return nil
	}
	v.Set("q", args)
	api += v.Encode()

	memes := make([]meme, 0, 7)
	c := colly.NewCollector(
		colly.AllowedDomains("imgflip.com"),
	)

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL)
	})
	c.OnHTML("a.s-result.clearfix[href]", func(e *colly.HTMLElement) {
		if id := e.Attr("href"); strings.HasPrefix(id, `/meme`) {
			if len(memes) >= 7 {
				c.OnHTMLDetach("a.s-result.clearfix[href]")
				return
			}
			num := fmt.Sprintf("%d. ", len(memes)+1)
			name := e.ChildText("div.s-result-title")
			var tid string
			id = strings.TrimPrefix(id, "/meme/")
			in := strings.Index(id, "/")

			if mID, ok := r.topMemes[name]; ok {
				tid = strconv.Itoa(mID.id)
				// fmt.Println("from cache", tid)
			} else if in >= 0 { //  use regex here
				tid = id[:in]
				// fmt.Println("from url", tid)
			} else {
				tid = r.getMemeID(id)
				fmt.Println("from scraping", tid) // delete this println(s)
			}

			memes = append(memes, meme{
				ID:   tid,
				name: num + name,
				img:  "http:" + e.ChildAttr("img", "src"),
			})
			// c.Visit("https://imgflip.com/memetemplate/"+ id)
		}
	})
	c.Visit(api)
	// fmt.Println(memes)
	return memes

}

func (r *Mux) CaptionMeme(s *dg.Session, m *dg.MessageCreate) { // TEMP
	args := fieldsN(m.Content, 3)
	// args[2] = strings.Join(args[2:], " ")
	if len(args) < 3 {
		// helpCapMeme()
		s.ChannelMessageSend(m.ChannelID, "at least write 3 args") // delete this
		fmt.Printf("%#v", args)
		return
	}
	texts := strings.Split(args[2], ",")
	memeUrl, err := r.memeFactory(args[1], texts...)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, memeUrl, m.Reference())

}

func (r *Mux) memeFactory(tID string, texts ...string) (string, error) {
	uname, pass := os.Getenv("MEME_USERNAME"), os.Getenv("MEME_PASS")
	v := url.Values{}

	v.Set("template_id", tID)
	for i, caption := range texts {
		k := fmt.Sprintf("boxes[%d][text]", i)
		v.Add(k, caption)
	}
	api := fmt.Sprintf("https://api.imgflip.com/caption_image?username=%s&password=%s", uname, pass)
	api += "&" + v.Encode()

	res, err := http.Post(api, "", nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	obj := struct {
		Data struct {
			Url string
		}
	}{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&obj)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return obj.Data.Url, nil
}

func (r *Mux) getMemeID(meme string) string {
	url := "https://imgflip.com/memetemplate/" + meme

	c := colly.NewCollector(
		colly.AllowedDomains("imgflip.com"),
	)
	var id string
	c.OnHTML("div.ibox", func(e *colly.HTMLElement) {
		// fmt.Println(e.ChildTexts("p"))
		ps := e.ChildTexts("p")
		for _, p := range ps {
			if strings.Contains(p, "ID") {
				id = strings.TrimPrefix(p, "Template ID: ")
				c.OnHTMLDetach("div.ibox")
			}
		}
	})
	c.Visit(url)
	return id
}

func (r *Mux) cacheTopMemes() {
	url := "https://api.imgflip.com/get_memes"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	obj := struct {
		Data struct {
			Memes []struct {
				Id, Name  string
				Box_count int
			}
		}
	}{}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	ignore := func(i int, _ error) int {
		return i
	}
	for _, m := range obj.Data.Memes {
		r.topMemes[m.Name] = &memeCache{
			id:  ignore(strconv.Atoi(m.Id)),
			box: m.Box_count,
		}
	}
}

func (meme *meme) msgMemeEmbed() *dg.MessageEmbed {
	return &dg.MessageEmbed{
		Title: meme.name,
		Image: &dg.MessageEmbedImage{
			URL: meme.img,
		},
		// Footer: &dg.MessageEmbedFooter{
		// },
		Description: fmt.Sprintf(`
You can make %s meme
or browse more %s of this type
			`, emojiNew, emojiLed_meme),
	}
}
