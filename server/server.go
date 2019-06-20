package server

import (
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/errorcheck"
	"log"
	"net"
)

var clients []net.Conn
var count = 0

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
		count++
		if err != nil {
			log.Println(err)
		}

		// Add the client to the connection array
		clients = append(clients, conn)

		go HandleConnection(conn)
	}
}
