package main

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"math/rand"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

// SetConsoleTitle sets the console title
func SetConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer func(handle syscall.Handle) {
		_ = syscall.FreeLibrary(handle)
	}(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}
	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}

// Clr clears the console
func Clr() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// LoadConfig loads the config file and the output directory
func LoadConfig() {

	//Check if output directory exists
	if stat, err := os.Stat("./Output"); err == nil && stat.IsDir() {
		color.Success.Println("Output directory loaded.")
	} else {
		color.Red.Println("Output directory not found, creating it.")
		err := os.Mkdir("Output", 0755)
		if err != nil {
			return
		}
	}

	if _, err := os.Stat("config.json"); err == nil {
		color.Success.Println("Config Loaded")
	} else {
		color.Red.Println("Config not found, creating it")
		cfg := &UserConfig{
			CaptchaKey: "",
		}
		data, _ := json.Marshal(cfg)
		err := os.WriteFile("config.json", data, 0755)
		if err != nil {
			return
		}
		if err != nil {
			return
		}
	}

}

// Scanln is similar to Scan, but stops scanning at a newline and
// after the final item there must be a newline or EOF.
func Scanln(a ...interface{}) (n int, err error) {
	return fmt.Fscanln(os.Stdin, a...)
}

func KeyPrompt() {
	var key string
	fmt.Println("Input 2Captcha Key...")
	_, err := Scanln(&key)
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

//randInt generates a random integer between min and max
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
