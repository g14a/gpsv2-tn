package server

import (
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/errorcheck"
	"log"
	"net"
	"sync"
)

var clients []net.Conn
var count = 0
var countMutex sync.Mutex

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

		// Add the client to the connection array
		clients = append(clients, conn)

		go HandleConnection(conn)
	}
}

// GTPL $1,867322035130152,A,170319,183004,9.960135,N,76.285625,E,0,20968,140,10,0,21,1,1,00.0250