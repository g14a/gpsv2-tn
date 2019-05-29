package server

import (
	"fmt"
	"gitlab.com/gps2.0/config"
	"gitlab.com/gps2.0/errcheck"
	"net"
)

func StartServer() {
	appConfigInstance := config.GetAppConfig()
	tcpAddress := appConfigInstance.TCPSocketConfig.ServerURL + ":" + appConfigInstance.TCPSocketConfig.Port

	ln, err := net.Listen("tcp", tcpAddress)

	errcheck.CheckError(err)

	fmt.Println("Accept Incoming connection")

	for {
		conn, err := ln.Accept()
		errcheck.CheckError(err)

		fmt.Println("New Client -- ", conn.RemoteAddr(), " --  connected..")
		go HandleConnection(conn)
	}
}

// GTPL $1,867322035130152,A,170319,183004,9.960135,N,76.285625,E,0,20968,140,10,0,21,1,1,00.0250
// GTPL $1,867322035130152,A,170318,183004,9.960135,N,76.285625,E,0,20968,140,10,0,21,1,1,00.0250
