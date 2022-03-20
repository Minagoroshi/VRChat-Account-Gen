package Shared

import (
	"os"
	"sync"
)

var (
	PManager      = &ProxyManager{}
	BotCount  int = 20
	WaitGroup     = sync.WaitGroup{}
	Semaphore chan int
	Worker    chan int
	OutFile   *os.File
)
