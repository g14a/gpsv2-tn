package server

import (
	"gitlab.com/gps2.0/config"
	"gitlab.com/gps2.0/errcheck"
	"log"
	"net"
)

func StartServer() {
	appConfigInstance := config.GetAppConfig()
	tcpAddress := appConfigInstance.TCPSocketConfig.ServerURL + ":" + appConfigInstance.TCPSocketConfig.Port

	ln, err := net.Listen("tcp", tcpAddress)

	errcheck.CheckError(err)

	defer ln.Close()
	go signalHandler()

	log.Println("[SERVER] listening...")

	for {
		conn, err := ln.Accept()
		count++
		if err != nil {
			log.Println(err)
		}

		log.Printf("[SERVER] Client connected %s -> %s -- Number of clients connected (%d)\n", conn.RemoteAddr(), conn.LocalAddr(), count)
		// Add the client to the connection array
		clients = append(clients, conn)

		go HandleConnection(conn)
	}
}

// GTPL $1,867322035130152,A,170319,183004,9.960135,N,76.285625,E,0,20968,140,10,0,21,1,1,00.0250
// GTPL $1,867322035130152,A,170318,183004,9.960135,N,76.285625,E,0,20968,140,10,0,21,1,1,00.0250
