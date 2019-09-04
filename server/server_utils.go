package server

import (
	"fmt"
	"gitlab.com/gpsv2-withoutrm/models"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)


// HandleConnection handles a connection by firing
// up a seperate go routine for a TCP connection net.Conn

// readTCPClient reads data sent by the device(a TCP client)
// and pushes it to the DB in an overview. Read more documentation below
func HandleConnection(conn net.Conn) {

	fmt.Println(runtime.NumGoroutine(), " goroutines and ", count, " clients connected")

	for {
		// Initialize a buffer of 5KB to be read from the client and read using conn.Read
		buf := make([]byte, 2*1024)
		_, err := conn.Read(buf)

		// if an error occurs deal with it
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed EOF")
				_ = conn.Close()
				count--
				return
			}
		} else {
			buffer := string(buf)

			if strings.Contains(buffer, "BSTPL") {
				dataslice := strings.Split(string(buf), "#")

				for _, record := range dataslice {
					processBSTPLDevices(record)
				}

			} else if strings.Contains(buffer, "GTPL") {
				dataslice := strings.Split(string(buf), "#")

				for _, record := range dataslice {
					processGTPLDevices(record)
				}

			} else if strings.Contains(buffer, "AVA") {
				dataslice := strings.Split(string(buf), "*")

				for _, record := range dataslice {
					processAIS140Device(record)
				}
			}
		}
	}
}

func processGTPLDevices(record string) {

	var (
		gtplDevice  models.GTPLDevice
		mysqlDevice models.GTPLSQLModel
		mssqlDevice models.MSSQLDevice
	)

	gtplDevice = ParseGTPLData(record)

	// ignores if an empty data occurs
	if gtplDevice.Latitude != 0 && gtplDevice.Longitude != 0 && gtplDevice.DeviceID != "" {

		mssqlDevice = ParseMSSQLDeviceFromGTPL(gtplDevice)

		gtplDevice.Distance = mysqlDevice.DistanceTravelled

		InsertIntoMSSQL(mssqlDevice)
		insertGTPLDataMongo(&gtplDevice)

		mysqlDevice = ParseGTPLDataSQL(gtplDevice)

		insertGTPLIntoSQL(mysqlDevice)
		insertRawDataMongo(record)
	}
}

func processBSTPLDevices(record string) {

	var (
		bstplDevice models.BSTPLDevice
		mssqlDevice models.MSSQLDevice
		mysqlDevice models.BSTPLSQLModel
	)

	bstplDevice = ParseBSTPLData(record)

	if bstplDevice.Latitude != 0 && bstplDevice.Longitude != 0 && bstplDevice.VehicleID != "" {

		recvTime := time.Now()

		insertBSTPLDataMongo(&bstplDevice)
		mssqlDevice = ParseMSSQLDeviceFromBSTPL(bstplDevice)
		mssqlDevice.RecvTime = recvTime

		InsertIntoMSSQL(mssqlDevice)
		mysqlDevice = ParseBSTPLDataSQL(bstplDevice)

		insertBSTPLIntoSQL(mysqlDevice)
		insertRawDataMongo(record)

	}
}

func processAIS140Device(record string) {

	var (
		ais140Device  models.AIS140Device
		mysqlDevice models.AIS140SQLModel
		mssqlDevice models.MSSQLDevice
	)

	ais140Device = ParseAIS140Data(record)

	// ignores if an empty data occurs
	if ais140Device.Latitude != 0 && ais140Device.Longitude != 0 && ais140Device.IMEINumber != "" {

		mssqlDevice = ParseMSSQLDeviceFromAIS140(ais140Device)

		InsertIntoMSSQL(mssqlDevice)
		insertAIS140DataIntoMongo(&ais140Device)

		mysqlDevice = ParseAIS140DataSQL(ais140Device)
		insertAIS140IntoSQL(mysqlDevice)
		insertRawDataMongo(record)
	}
}

// signalHandler notices termination signals or
// interrupts from the command line. Eg: ctrl-c and exits cleanly
func signalHandler() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	go func() {
		for sig := range sigchan {
			log.Printf("[SERVER] Closing due to Signal: %s", sig)
			log.Printf("[SERVER] Graceful shutdown")

			fmt.Println("Done.")

			// Exit cleanly
			os.Exit(0)
		}
	}()
}
