package Util

import (
	"fmt"
	"os"
	"os/signal"

	"antegr.al/chatanium-bot/v1/src/Util/Log"
)

// this function is blocking (main) thread until a signal (Interrupt, Kill) is received.
// if end this function called in main, excute defer function in main.
func WaitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c // Wait for a signal

	fmt.Println() // Print a new line (for better readability)
	Log.Info.Println("Starting shutdown process. Please wait...")
}
