package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/models"
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

	dbWg = sync.WaitGroup{}
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

					go insertRawDataMongo(individualRecord, &dbWg)
					go insertRawDataSQL(individualRecord, &dbWg)

					gtplDevice = ParseGTPLData(individualRecord)

					// ignores if an empty data occurs
					if (models.GTPLDevice{}) != gtplDevice {

						if gtplDevice.DeviceTimeNow.Day() == time.Now().Day() {
							go insertGTPLDataMongo(&gtplDevice, &dbWg)
							go insertGTPLIntoSQL(gtplDevice, &dbWg)
						} else {
							go insertGTPLHistoryDataMongo(&gtplDevice, &dbWg)
							go insertGTPLIntoSQL(gtplDevice, &dbWg)
						}
					}
				}

			} else if strings.Contains(string(buf), "AVA") {
				dataSlice := strings.Split(string(buf), "*")

				var ais140Device models.AIS140Device

				for _, individualRecord := range dataSlice {

					go insertRawDataMongo(individualRecord, &dbWg)
					go insertRawDataSQL(individualRecord, &dbWg)

					ais140Device = ParseAIS140Data(individualRecord)

					// ignores if an empty data occurs
					if (models.AIS140Device{}) != ais140Device {
						if ais140Device.LiveOrHistoryPacket == "L" || (ais140Device.LiveOrHistoryPacket == "H" && ais140Device.DeviceTime.Day() == time.Now().Day()) {
							go insertAIS140DataIntoMongo(&ais140Device, &dbWg)
							go insertAIS140IntoSQL(ais140Device, &dbWg)
						} else {
							go insertAIS140HistoryDataMongo(&ais140Device, &dbWg)
							go insertAIS140IntoSQL(ais140Device, &dbWg)
						}
					}
				}
			}
			dbWg.Wait()
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
