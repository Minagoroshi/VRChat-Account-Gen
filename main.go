//go:generate goversioninfo -icon=icon.ico

package main

import (
	"VRChat_Account_Generator/Shared"
	"VRChat_Account_Generator/requests"
	"VRChat_Account_Generator/utils"
	"bufio"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gookit/color"
	"github.com/gosuri/uilive"
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
	//Startup
	var key string
	utils.Startup()
	//Menu Prompts
	_, _ = utils.SetConsoleTitle("VRChat Account Generator GO by top#2222")
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
	utils.Clr()
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
	pCount := 0
	if x != 5 {
		pCount, err = Shared.PManager.LoadFromFile("proxies.txt", x-1)
		if err != nil {
			log.Fatal(err.Error())
		}
		color.Success.Println("Loaded", pCount, "proxies")
	}
	//Key Check
	cfg := utils.ReadConfig()
	if cfg.CaptchaKey != "" {
		key = cfg.CaptchaKey
	} else {
		utils.KeyPrompt()
		utils.ReadConfig()
		key = cfg.CaptchaKey
	}
	//Bot Prompt
	color.FgLightCyan.Print(fmt.Sprintf("Number of bots? [Default: %d] [Max: 200]: ", Shared.BotCount))
	e, err := reader.ReadString('\n')
	utils.Clr()
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
	t := time.Now()
	//Make Output File
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
		go WorkerFunc(key)

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
				_, _ = utils.SetConsoleTitle(fmt.Sprintf("VRChat Generator GO | Proxies %d | Accounts Made %d | Fails %d | Bots %d", pCount, successes, fails, len(Shared.Semaphore)))
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

func WorkerFunc(key string) {
	worker := <-Shared.Worker
	//Start Worker Loop
	for {
		//Register
		proxy, statusCode, body, username, password, email := requests.Register(key, worker)
		//Parse Response
		if strings.Contains(body, "currentAvatar") {
			successes++
			color.Success.Println(utils.TimeStamp(), "Worker", worker, "| Success! |  Email : ", email, " | Username : ", username, " | Password : ", password, "| Status Code:", statusCode, "| Proxy:", proxy)
			_, err := Shared.OutFile.WriteString(username + ":" + password + "\n")
			if err != nil {
				return
			}
		} else {
			fails++
			message, _ := jsonparser.GetString([]byte(body), "error", "message")
			if message == "" {
				message = "None"
			}
			color.Error.Println(utils.TimeStamp(), "Worker", worker, "| Failed! | Message:", message, "| Status Code:", statusCode, "| Proxy:", proxy)
		}
	}
}
