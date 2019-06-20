package server

import (
	"fmt"
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/errorcheck"
	"gitlab.com/gpsv2/models"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

// HandleConnection handles a connection by firing
// up a seperate go routine for a TCP connection net.Conn
func HandleConnection(conn net.Conn) {

	var wg sync.WaitGroup

	wg.Add(1)
	go readTCPClient(conn, &wg)
	wg.Wait()

}

var (
	// live database collections
	locationHistoriesCollection = config.GetAppConfig().Mongoconfig.Collections.LocationHistoriesCollection
	vehicleDetailsCollection    = config.GetAppConfig().Mongoconfig.Collections.VehicleDetailsCollection

	// backups collections
	historyLHcollection = config.GetAppConfig().HistoryMongoConfig.BackupCollections.BackupLocationHistoriesColl
	rawDataCollection   = config.GetAppConfig().HistoryMongoConfig.BackupCollections.RawDataCollection

	collectionMutex = &sync.Mutex{}
	dataMutex       = &sync.Mutex{}
)

// readTCPClient reads data sent by the device(a TCP client)
// and pushes it to the DB in an overview. Read more documentation below
func readTCPClient(conn net.Conn, wg *sync.WaitGroup) {

	fmt.Printf("\n[SERVER] Client connected %s -> %s -- Number of clients connected (%d)\n", conn.RemoteAddr(), conn.LocalAddr(), count)

	defer wg.Done()

	for {
		// Initialize a buffer of 5KB to be read from the client and read using conn.Read
		buf := make([]byte, 5*1024)
		_, err := conn.Read(buf)

		// if an error occurs deal with it
		if err != nil {
			// if the error is EOF which means the client
			// is not sending any data, close the connection
			if err == io.EOF {
				fmt.Println("Connection closed EOF")
				_ = conn.Close()
			}
		} else {
			// Data is successfully read here.
			dataMutex.Lock()

			// If the Data contains GTPL Data, proceed with GTPL functions.
			if strings.Contains(string(buf), "GTPL") {
				// Split the data if it comes in a bulk using the
				// GTPL delimiter # and put it into an array or a slice in Go.
				dataSlice := strings.Split(string(buf), "#")

				// Initialize a gtpl device for each record.
				var gtplDevice models.GTPLDevice

				// Iterate over the slice and put each record into the db.
				for _, individualRecord := range dataSlice {

					// First put the raw data.
					err = insertRawDataMongo(individualRecord)
					fmt.Println(individualRecord)

					// Parse raw data into a device.
					gtplDevice = ParseGTPLData(individualRecord)

					// Filter live and history data.
					if gtplDevice.DeviceTimeNow.Day() == time.Now().Day() {
						err = insertGTPLDataMongo(&gtplDevice)
						errorcheck.CheckError(err)
					} else {
						err = insertGTPLHistoryDataMongo(&gtplDevice)
						errorcheck.CheckError(err)
					}
				}
			} else if strings.Contains(string(buf), "AVA") {
				// If the raw data contains AIS140 data split it using *
				dataSlice := strings.Split(string(buf), "*")

				var ais140Device models.AIS140Device

				for _, individualRecord := range dataSlice {

					err = insertRawDataMongo(individualRecord)
					fmt.Println(individualRecord)
					ais140Device = ParseAIS140Data(individualRecord)

					// Filter history packet using L and H field
					if ais140Device.LiveOrHistoryPacket == "L" || (ais140Device.LiveOrHistoryPacket == "H" && ais140Device.DeviceTime.Day() == time.Now().Day()) {
						err = insertAIS140DataIntoMongo(&ais140Device)
						errorcheck.CheckError(err)
					} else {
						err = insertAIS140HistoryDataMongo(&ais140Device)
						errorcheck.CheckError(err)
					}
				}
			}
			dataMutex.Unlock()
		}
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
