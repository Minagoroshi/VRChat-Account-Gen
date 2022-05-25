package main

import (
	"VRChat_Account_Generator/Shared"
	"fmt"
	"github.com/gookit/color"
	"log"
	"time"
)

// WorkerLog is a function that prints custom logs for the worker to the console, on success logs the data to the output file
// Takes a logType, a message, and a worker id
// LogType can be "info", "error", "warning", "success", "failure", or "debug"
func WorkerLog(logType string, worker int, message string) {

	now := time.Now()

	switch logType {
	case "info":
		color.Info.Printf("[%s] [Worker %d] [%s]\n", now.Format("2006-01-02 15:04:05"), worker, message)
	case "error":
		color.Error.Printf("[%s] [Worker %d] [%s]\n", now.Format("2006-01-02 15:04:05"), worker, message)
	case "warning":
		color.Warn.Printf("[%s] [Worker %d] [%s]\n", now.Format("2006-01-02 15:04:05"), worker, message)
	case "success":
		successLog := fmt.Sprintf("[%s] [Worker %d] [%s]\n", now.Format("2006-01-02 15:04:05"), worker, message)
		color.Success.Printf(successLog)
		_, err := Shared.OutFile.WriteString(successLog)
		if err != nil {
			log.Fatal(err)
		}

	case "failure":
		color.Error.Printf("[%s] [Worker %d] [%s]\n", now.Format("2006-01-02 15:04:05"), worker, message)
	case "debug":
		color.Debug.Printf("[%s] [Worker %d] [%s]\n", now.Format("2006-01-02 15:04:05"), worker, message)
	}

}
