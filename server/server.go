package server

import (
	"gitlab.com/gpsv2-withoutrmtesting/config"
	"gitlab.com/gpsv2-withoutrmtesting/errorcheck"
	"log"
	"net"
	"sync"
)

var (
	count = 0
	countMutex sync.Mutex
)

// StartServer starts listening to the server via TCP protocol
func StartServer() {

	appConfigInstance := config.GetAppConfig()
	tcpAddress := appConfigInstance.TCPSocketConfig.ServerURL + ":" + appConfigInstance.TCPSocketConfig.Port

	ln, err := net.Listen("tcp", tcpAddress)

	errorcheck.CheckError(err)

	defer ln.Close()
	go signalHandler()

	log.Println("[SERVER] listening...")

	for {
		conn, err := ln.Accept()

		countMutex.Lock()
		count++
		countMutex.Unlock()

		if err != nil {
			log.Println(err)
		}

		go HandleConnection(conn)
	}
}