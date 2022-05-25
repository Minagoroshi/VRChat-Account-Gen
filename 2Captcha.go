package main

import (
	"github.com/2captcha/2captcha-go"
	"github.com/gookit/color"
	"log"
	"time"
)

// SolveHCaptcha takes a 2Captcha Key and solves a captcha using the 2Captcha API
func SolveHCaptcha(worker int) string {
	if ReadConfig().CaptchaKey != "" {
		client := api2captcha.NewClient(ReadConfig().CaptchaKey)

		captcha := api2captcha.HCaptcha{
			SiteKey: "67ab3710-6883-4d15-8614-db0861382a92",
			Url:     "https://vrchat.com/home/register",
		}

		id, err := client.Send(captcha.ToRequest())
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second * 30)
		// Wait for captcha to be solved
		var ready bool
		for ready == false {
			time.Sleep(10 * time.Second)

			code, err := client.GetResult(id)
			if err != nil {
				log.Fatal(err)
			}

			if code == nil {
				WorkerLog("info", worker, "Waiting for captcha to be solved")
			}

			ready = true

			return *code
		}

	} else {
		color.Error.Println("No 2Captcha Key found in config.json")
	}
	return ""
}
