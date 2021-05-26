package main

import (
	"fmt"

	"github.com/notAxion/HackerHeads/bot"
	"github.com/notAxion/HackerHeads/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println("Config fail", err)
		return
	}

	bot.Start()
}
