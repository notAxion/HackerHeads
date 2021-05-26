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

	var c struct {
		BotPrefix string `json:"bot_prefix"`
	}
	var s struct {
		Token string `json:"token"`
	}
	if err = json.Unmarshal(cbs, &c); err != nil {
		return err
	}

	if err = json.Unmarshal(sbs, &s); err != nil {
		return err
	}

	BotPrefix = c.BotPrefix
	Token = s.Token

	// Config = ConfigsType{Token: s.Token, BotPrefix: c.BotPrefix}
	return nil
}
