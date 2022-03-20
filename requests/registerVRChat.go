package requests

import (
	"VRChat_Account_Generator/Shared"
	"VRChat_Account_Generator/utils"
	"fmt"
	"github.com/gookit/color"
	"github.com/thanhpk/randstr"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//Register takes a 2Captcha Key and a Worker id, it registers an account on VRChat with randomly generated data
//and uses 2Captcha API for captcha solving.
//
//
//Returns Proxy Address, Status Code, Response Body, Username, Password, Email
func Register(key string, worker int) (string, int, string, string, string, string) {

	var year, month, day = "1980", "12", "6"
	var username, password, email string
	var statusCode int
	pink := color.Hex("ff69b4")

Start:
	transport, proxy, err := Shared.PManager.GetRandomProxyTransport()
	if err != nil {
		if proxy != nil {
			proxy.InUse = false
		}
		goto Start
	}
	client := &http.Client{Timeout: 10 * time.Second, Transport: transport}

	//Submit Captcha

	submitURL := "https://2captcha.com/in.php?key=" + key + "&method=userrecaptcha&googlekey=6LfxcQ4UAAAAAGNAOUtX3pADEAu-sCsQL6En2E9S&pageurl=https://vrchat.com/home/register"
	methodGet := "GET"

	req, err := http.NewRequest(methodGet, submitURL, nil)
	if err != nil {
		if proxy != nil {
			proxy.InUse = false
		}
		goto Start
	}
	color.Yellow.Println(utils.TimeStamp(), "Worker", worker, "| Submitting Captcha")
	res, err := client.Do(req)
	if err != nil {
		if proxy != nil {
			proxy.InUse = false
		}
		goto Start
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}
	id := strings.Replace(string(body), "OK|", "", -1)
	time.Sleep(45 * time.Second)

	//Get Solution
	solutionURL := "https://2captcha.com/res.php?key=969a80dc07ede0daf4696ec73d61b3da&action=get&id=" + id

	getSolution, err := http.NewRequest(methodGet, solutionURL, nil)
	if err != nil {
		if proxy != nil {
			proxy.InUse = false
		}
		goto Start
	}
	color.HiBlue.Println(utils.TimeStamp(), "Worker", worker, "| Getting Solution")
	res, err = client.Do(getSolution)
	if err != nil {
		if proxy != nil {
			proxy.InUse = false
		}
		goto Start
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}
	var Success bool
	for Success == false {
		time.Sleep(5 * time.Second)
		getSolution, err = http.NewRequest(methodGet, solutionURL, nil)

		if err != nil {
			fmt.Println(err)
		}
		res, err = client.Do(getSolution)
		if err != nil {
			fmt.Println(err)
		}

		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		//Register Account
		if !strings.Contains(string(body), "NOT_READY") {
			color.HiCyan.Println(utils.TimeStamp(), "Worker", worker, "| Got Solution")
			code := strings.Replace(string(body), "OK|", "", -1)
			username = randstr.String(15)
			password = randstr.String(15)
			email = randstr.String(15) + "@outlook.com"
			url := "https://vrchat.com/api/1/auth/register?apiKey=JlE5Jldo5Jibnk5O5hTx6XVqsJu4WJ26"
			methodPost := "POST"
			payload := strings.NewReader("{\"username\":" + "\"" + username + "\"" +
				",\"password\":" + "\"" + password + "\"" +
				",\"email\":" + "\"" + email + "\"" +
				",\"year\":" + "\"" + year + "\"" +
				",\"month\":" + "\"" + month + "\"" +
				",\"day\":" + "\"" + day + "\"" +
				",\"recaptchaCode\":" + "\"" + code + "\"" +
				",\"day\":" + "\"" + day + "\"" +
				",\"subscribe\":true,\"acceptedTOSVersion\":7}")

			req, err = http.NewRequest(methodPost, url, payload)
			if err != nil {
				if proxy != nil {
					proxy.InUse = false
				}
				goto Start
			}

			if err != nil {
				fmt.Println(err)

			}
			req.Header.Add("Host", " vrchat.com")
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
			req.Header.Add("Origin", "https://vrchat.com")
			req.Header.Add("Referer", "https://vrchat.com/home/register")
			pink.Println(utils.TimeStamp(), "Worker", worker, "| Registering Account")
			res, err = client.Do(req)
			if err != nil {
				if proxy != nil {
					proxy.InUse = false
				}
				goto Start
			}
			body, _ = ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)

			}

			statusCode = res.StatusCode
			Success = true
		}
	}

	proxyAddress := proxy.Address
	return proxyAddress, statusCode, string(body), username, password, email
}
