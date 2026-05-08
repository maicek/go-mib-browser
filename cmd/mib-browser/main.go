package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go CreateOsWindow()

	<-signalChan

	fmt.Printf("close")
}
