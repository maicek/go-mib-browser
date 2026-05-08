package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Printf("close")
		os.Exit(0)
	}()

	CreateOsWindow()
}
