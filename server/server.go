package server

import (
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/errcheck"
	"log"
	"net"
)

var clients []net.Conn
var count = 0

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

		// Add the client to the connection array
		clients = append(clients, conn)

		go HandleConnection(conn)
	}
}

// GTPL $1,867322035130152,A,170319,183004,9.960135,N,76.285625,E,0,20968,140,10,0,21,1,1,00.0250
// GTPL $1,867322035130152,A,170318,183004,9.960135,N,76.285625,E,0,20968,140,10,0,21,1,1,00.0250

// $,1,AVA,AVA1.0,NR,01,L,869867031217468,1,1,07062019,103144,11.678080,N,78.163207,E,0,346,08,0327,0.92,1.24,CellOne ,0,1,12.11,4.1,0,C,22,404,004,998,db12,998db1388999df1788999dedb90998db2590,0000,00,000848,179*
