package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type UserConfig struct {
	CaptchaKey string `json:"CaptchaKey"`
}

func ReadConfig() UserConfig {
	var config UserConfig
	f, _ := os.ReadFile("config.json")
	err := json.Unmarshal(f, &config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}
