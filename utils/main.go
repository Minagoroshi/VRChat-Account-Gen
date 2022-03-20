package utils

import (
	"encoding/json"
	"github.com/gookit/color"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

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
func Clr() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func Startup() {

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
			Login:      "",
			CaptchaKey: "",
			RememberMe: false,
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

//TimeStamp takes the current time, and converts it to the format of [hh:mm:ss]
func TimeStamp() string {
	now := time.Now()
	h, m, s := now.Hour(), now.Minute(), now.Second()
	hour, min, sec := strconv.Itoa(h), strconv.Itoa(m), strconv.Itoa(s)
	timeStamp := "[" + hour + ":" + min + ":" + sec + "]"
	return timeStamp
}
