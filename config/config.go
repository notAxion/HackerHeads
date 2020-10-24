package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Token     string 
	BotPrefix string

	// Config ConfigsType
)

var c struct{ 
	BotPrefix string `json:"bot_prefix"`
}
var s struct{ Token string `json:"token"` }


// type ConfigsType struct {
// 	Token     string `json:"token"`
// 	BotPrefix string `json:"bot_prefix"`
// }

func ReadConfig() error {
	fmt.Println("Reading from config file...")

	cbs, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}

	sbs, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		return err
	}

	fmt.Println("reading working")


	err = json.Unmarshal(cbs, &c)
	if err != nil {
		return err
	}

	err = json.Unmarshal(sbs, &s)
	if err != nil {
		return err
	}

	BotPrefix = c.BotPrefix
	Token = s.Token

	// Config = ConfigsType{Token: s.Token, BotPrefix: c.BotPrefix}
	return nil
}
