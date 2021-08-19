package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// StartDaemon daemon process will deal with SIGTERM, Daemon process will be exit
// when received that signal .
func StartDaemon() {
	c := make(chan os.Signal, 1)
	// SIGKILL will force kill the Progress,Unable to handle and ignore the signal.
	// SIGTERM elegant, allowing to catch signals and do processing after exit
	signal.Notify(c, syscall.SIGTERM)
	<-c
	close(c)
	// there maybe send the status to the server in the future
	os.Exit(0)
}
