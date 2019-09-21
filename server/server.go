package server

import (
	"gitlab.com/gpsv2-tn/config"
	"gitlab.com/gpsv2-tn/errorcheck"
	"log"
	"net"
)

// StartServer starts listening to the server via TCP protocol
func StartServer() {

	appConfigInstance := config.GetAppConfig()
	tcpAddress := appConfigInstance.TCPSocketConfig.ServerURL + ":" + appConfigInstance.TCPSocketConfig.Port

	ln, err := net.Listen("tcp", tcpAddress)

	errorcheck.CheckError(err)

	defer ln.Close()

	log.Println("[SERVER] listening...")

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println(err)
		}

		go HandleConnection(conn)
	}
}