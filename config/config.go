package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Token string
	BotPrefix string
	
	config *newConfig
)

type newConfig struct {
	Token       string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading from config file...")
	
	file, err := ioutil.ReadFile("./config.json")
	
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	
	fmt.Println("reading working")
	
	err = json.Unmarshal(file, &config)
	
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	
	Token = config.Token
	BotPrefix = config.BotPrefix
	
	return nil
}
























