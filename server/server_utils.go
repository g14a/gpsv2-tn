package server

import (
	"fmt"
	"gitlab.com/gpsv2-withoutrmtesting/models"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"
)


// HandleConnection handles a connection by firing
// up a seperate go routine for a TCP connection net.Conn

// readTCPClient reads data sent by the device(a TCP client)
// and pushes it to the DB in an overview. Read more documentation below
func HandleConnection(conn net.Conn) {

	fmt.Println("Number of goroutines : ", runtime.NumGoroutine())

	fmt.Printf("\n[SERVER] Client connected %s -> %s -- Number of clients connected (%d)\n", conn.RemoteAddr(), conn.LocalAddr(), count)

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
			}
		} else {
			buffer := string(buf)

			if strings.Contains(buffer, "BSTPL") {
				dataslice := strings.Split(string(buf), "#")

				start := time.Now()
				for _, record := range dataslice {

					processBSTPLDevice(record)
				}
				fmt.Println(time.Since(start))

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

	var dbWg sync.WaitGroup

	var (
		gtplDevice  models.GTPLDevice
		mysqlDevice models.GTPLSQLModel
		mssqlDevice models.MSSQLDevice
	)

	dbWg.Add(2)

	gtplDevice = ParseGTPLData(record)

	// ignores if an empty data occurs
	if gtplDevice.Latitude != 0 && gtplDevice.Longitude != 0 && gtplDevice.DeviceID != "" {

		mssqlDevice = ParseMSSQLDeviceFromGTPL(gtplDevice)

		gtplDevice.Distance = mysqlDevice.DistanceTravelled

		go InsertIntoMSSQL(mssqlDevice, &dbWg)
		go insertGTPLDataMongo(&gtplDevice, &dbWg)

		mysqlDevice = ParseGTPLDataSQL(gtplDevice)

		go insertGTPLIntoSQL(mysqlDevice, &dbWg)
		go insertRawDataMongo(record, &dbWg)

		dbWg.Wait()
	}
}

func processBSTPLDevice(record string) {
	var dbWg sync.WaitGroup

	var (
		bstplDevice models.BSTPLDevice
		mssqlDevice models.MSSQLDevice
		mysqlDevice models.BSTPLSQLModel
	)

	dbWg.Add(2)

	bstplDevice = ParseBSTPLData(record)

	if bstplDevice.Latitude != 0 && bstplDevice.Longitude != 0 && bstplDevice.VehicleID != "" {

		recvTime := time.Now()

		bstplDevice.CreatedTime = recvTime

		go insertBSTPLDataMongo(&bstplDevice, &dbWg)

		mssqlDevice = ParseMSSQLDeviceFromBSTPL(bstplDevice)
		mssqlDevice.RecvTime = recvTime

		go InsertIntoMSSQL(mssqlDevice, &dbWg)
		mysqlDevice = ParseBSTPLDataSQL(bstplDevice)

		go insertBSTPLIntoSQL(mysqlDevice, &dbWg)
		go insertRawDataMongo(record, &dbWg)

		dbWg.Wait()
	}
}

func processAIS140Device(record string) {

	var dbWg sync.WaitGroup

	var (
		ais140Device  models.AIS140Device
		mysqlDevice models.AIS140SQLModel
		mssqlDevice models.MSSQLDevice
	)

	dbWg.Add(2)

	ais140Device = ParseAIS140Data(record)

	// ignores if an empty data occurs
	if ais140Device.Latitude != 0 && ais140Device.Longitude != 0 && ais140Device.IMEINumber != "" {

		mssqlDevice = ParseMSSQLDeviceFromAIS140(ais140Device)

		go InsertIntoMSSQL(mssqlDevice, &dbWg)
		go insertAIS140DataIntoMongo(&ais140Device, &dbWg)

		mysqlDevice = ParseAIS140DataSQL(ais140Device)
		go insertAIS140IntoSQL(mysqlDevice, &dbWg)
		go insertRawDataMongo(record, &dbWg)

		dbWg.Wait()
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
