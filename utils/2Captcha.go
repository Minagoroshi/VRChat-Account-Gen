package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func KeyPrompt() {
	var key string
	fmt.Println("Input 2Captcha Key...")
	_, err := fmt.Scanln(&key)
	if err != nil {
		return
	}
	userConfig := &UserConfig{
		CaptchaKey: key,
	}

	data, _ := json.Marshal(userConfig)
	err = os.WriteFile("config.json", data, 0755)
	if err != nil {
		return
	}

}
