package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.com/gpsv2/server"
)

// main starts the server
func main() {
	server.StartServer()

	//data := "GTPL $1,869867030663092,A,220619,055539,18.720370,N,80.081505,E,0,52,046,08,0,11,0,0,00.0000,03.88"
	//
	//gtpl := server.ParseGTPLData(data)
	//
	//server.InsertGTPLIntoSQL(gtpl)
}
