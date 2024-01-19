package Util

import (
	"os"
	"os/signal"

	"antegr.al/chatanium-bot/v1/src/Util/Log"
)

// WaitSignal waits for a signal to shutdown (Interrupt, Kill)
func WaitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c // Wait for a signal

	Log.Info.Println("Starting shutdown process. Please wait...")
}
