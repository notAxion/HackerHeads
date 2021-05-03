package main

import (
	"fmt"

	"github.com/notAxion/HackerHeads/bot"
	"github.com/notAxion/HackerHeads/config"
)

/*
func main() {
	token := ""

	f, _ := os.ReadFile("secrets.json")
	var int interface{}
	if err := json.Unmarshal(f, &int); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", int)
	token = int.(map[string]interface{})["token"].(string)
	fmt.Println(token)
} */

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
