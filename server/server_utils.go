package server

import (
	"fmt"
	"gitlab.com/gpsv2/config"
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
// and pushes it to the DB in an overview. Read more documentation below.
func readTCPClient(conn net.Conn, wg *sync.WaitGroup) {

	fmt.Printf("\n[SERVER] Client connected %s -> %s -- Number of clients connected (%d)\n", conn.RemoteAddr(), conn.LocalAddr(), count)

	defer wg.Done()

	for {
		// Initialize a buffer of 5KB to be read from the client and read using conn.Read
		buf := make([]byte, 5*1024)
		_, err := conn.Read(buf)

		// if an error occurs deal with it
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed EOF")
				_ = conn.Close()
			}
		} else {
			dataMutex.Lock()

			if strings.Contains(string(buf), "GTPL") {

				dataSlice := strings.Split(string(buf), "#")

				var gtplDevice models.GTPLDevice

				for _, individualRecord := range dataSlice {
					fmt.Println(individualRecord)

					insertRawDataMongo(individualRecord)
					insertRawDataSQL(individualRecord)

					gtplDevice = ParseGTPLData(individualRecord)

					// ignores if an empty data occurs
					if (models.GTPLDevice{}) != gtplDevice {

						if gtplDevice.DeviceTimeNow.Day() == time.Now().Day() {
							insertGTPLDataMongo(&gtplDevice)
							insertGTPLIntoSQL(gtplDevice)
						} else {
							insertGTPLHistoryDataMongo(&gtplDevice)
							insertGTPLIntoSQL(gtplDevice)
						}
					}
				}

			} else if strings.Contains(string(buf), "AVA") {
				dataSlice := strings.Split(string(buf), "*")

				var ais140Device models.AIS140Device

				for _, individualRecord := range dataSlice {

					insertRawDataMongo(individualRecord)
					insertRawDataSQL(individualRecord)

					ais140Device = ParseAIS140Data(individualRecord)

					// ignores if an empty data occurs
					if (models.AIS140Device{}) != ais140Device {
						if ais140Device.LiveOrHistoryPacket == "L" || (ais140Device.LiveOrHistoryPacket == "H" && ais140Device.DeviceTime.Day() == time.Now().Day()) {
							insertAIS140DataIntoMongo(&ais140Device)
							insertAIS140IntoSQL(ais140Device)
						} else {
							insertAIS140HistoryDataMongo(&ais140Device)
							insertAIS140IntoSQL(ais140Device)
						}
					}
				}
			}
			dataMutex.Unlock()
		}
	}
}

// SignalHandler notices termination signals or
// interrupts from the command line. Eg: ctrl-c and exits cleanly.
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
