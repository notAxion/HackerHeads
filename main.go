package main

import (
	"fmt"
	"./config"
	"./bot"
)


func main() {
	err := config.ReadConfig()
	
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	bot.Start()
	
	<- make(chan struct{})
	return 
}
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	