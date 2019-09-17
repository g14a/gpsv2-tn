package main

import (
	"gitlab.com/gpsv2-tn/server"
	"runtime/debug"
	"time"
)

// main starts the server
func main() {
	server.StartServer()

	go periodicFreeMemory(1 * time.Minute)
}

func periodicFreeMemory(d time.Duration) {
	tick := time.Tick(d)
	for _ = range tick {
		debug.FreeOSMemory()
	}
}