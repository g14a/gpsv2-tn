package main

import (
	"github.com/pkg/profile"
	"gitlab.com/gpsv2-tn/server"
)

// main starts the server
func main() {
	defer profile.Start(profile.MemProfile).Stop()

	server.StartServer()
}