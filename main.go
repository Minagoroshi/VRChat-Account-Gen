//go:generate goversioninfo -icon=icon.ico

package main

import (
	"VRChat_Account_Generator/Shared"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/gosuri/uilive"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	successes = 0
	fails     = 0
)

func main() {

	LoadConfig()
	//Menu Prompts
	_, _ = SetConsoleTitle("VRChat Account Generator GO")
	reader := bufio.NewReader(os.Stdin)

	x := 1
	Blue := color.FgLightBlue.Render
	fmt.Println("Proxy Type: \n",
		" 1.", Blue("Http \n"),
		" 2.", Blue("Socks4 \n"),
		" 3.", Blue("Socks4a \n"),
		" 4.", Blue("Socks5 \n"),
		" 5.", Blue("Proxyless (Don't use unless testing)"))
	pType, err := reader.ReadString('\n')
	Clr()
	pType = strings.ReplaceAll(strings.ReplaceAll(pType, "\n", ""), "\r", "")
	if len(pType) > 0 {
		x, err = strconv.Atoi(pType)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		if x == 0 {
			x = 1
		}
	}

	/// Load Proxies
	pCount := 0
	if x != 5 {
		pCount, err = Shared.PManager.LoadFromFile("proxies.txt", x-1)
		if err != nil {
			log.Fatal(err.Error())
		}
		color.Success.Println("Loaded", pCount, "proxies")
	}

	//Key Check
	cfg := ReadConfig()
	if cfg.CaptchaKey == "" {
		KeyPrompt()
		ReadConfig()
	}

	//Bot Prompt
	color.FgLightCyan.Print(fmt.Sprintf("Number of bots? [Default: %d] [Max: 200]: ", Shared.BotCount))
	e, err := reader.ReadString('\n')
	Clr()
	e = strings.ReplaceAll(strings.ReplaceAll(e, "\n", ""), "\r", "")
	if len(e) > 0 {
		Shared.BotCount, err = strconv.Atoi(e)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		if Shared.BotCount <= 0 {
			Shared.BotCount = 1
		}
		if Shared.BotCount > 200 {
			color.Red.Println("Bots over 200, defaulting to 200")
			Shared.BotCount = 200
		}
	}

	// Setup Output File
	t := time.Now()
	Shared.OutFile, err = os.Create("Output\\" + t.Format("2006-01-02_15-04-05"+".txt"))
	defer func(OutFile *os.File) {
		_ = OutFile.Close()
	}(Shared.OutFile)

	//Start Workers
	Shared.Semaphore = make(chan int, Shared.BotCount)
	Shared.Worker = make(chan int, Shared.BotCount)
	Shared.WaitGroup = sync.WaitGroup{}
	Shared.WaitGroup.Add(Shared.BotCount)
	for i := 0; i < Shared.BotCount; i++ {
		Shared.Worker <- i
		Shared.Semaphore <- 0
		go WorkerFunc()

	}

	//Title Update
	go func() {
		var writer *uilive.Writer
		if runtime.GOOS != "windows" {
			writer = uilive.New()
			writer.Start()
			defer writer.Stop()
		}
		for len(Shared.Semaphore) > 0 {
			if runtime.GOOS == "windows" {
				_, _ = SetConsoleTitle(fmt.Sprintf("VRChat Generator GO | Proxies %d | Accounts Made %d | Fails %d | Bots %d", pCount, successes, fails, len(Shared.Semaphore)))
			} else {
				_, _ = fmt.Fprintf(writer, "\t Accounts Made %d | fails %d  Bots %d\r\n", successes, fails, len(Shared.Semaphore))
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()

	time.Sleep(500 * time.Millisecond)
	Shared.WaitGroup.Wait()
	color.Info.Println("Done!")
	select {}
}

func WorkerFunc() {
	worker := <-Shared.Worker
	//Start Worker Loop
	for {

		//Register
		WorkerLog("debug", worker, "Attempting to register an account")
		res, username, password, email := RegisterVRC(worker)

		//Read res.Body into a string
		body, _ := ioutil.ReadAll(res.Body)

		//Unmarshal the body into a struct
		var responseBody RegistrationResponse
		err := json.Unmarshal(body, &responseBody)
		if err != nil {
			WorkerLog("error", worker, "Failed to unmarshal response body")
		}
		if responseBody.CurrentAvatar != "" {
			successes++
			WorkerLog("success", worker, fmt.Sprintf("Username: %s | Password: %s | Email: %s", username, password, email))
		} else {
			fails++
			WorkerLog("failure", worker, fmt.Sprintf("Failed to register account"))
		}

	}
}
