package server

import (
	"errors"
	"fmt"
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/db"
	"gitlab.com/gpsv2/errcheck"
	"gitlab.com/gpsv2/models"
	"go.mongodb.org/mongo-driver/bson"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func HandleConnection(conn net.Conn) {

	var wg sync.WaitGroup

	wg.Add(1)
	go readWrapper(conn, &wg)
	wg.Wait()

}

var (
	locationHistoriesCollection = config.GetAppConfig().Mongoconfig.Collections.LocationHistoriesCollection
	vehicleDetailsCollection    = config.GetAppConfig().Mongoconfig.Collections.VehicleDetailsCollection

	// backups collections
	historyLHcollection = config.GetAppConfig().HistoryMongoConfig.BackupCollections.BackupLocationHistoriesColl
	rawDataCollection   = config.GetAppConfig().HistoryMongoConfig.BackupCollections.RawDataCollection

	collectionMutex = &sync.Mutex{}
	dataMutex       = &sync.Mutex{}
)

func connCheckForShutdown(conn net.Conn) error {
	var (
		n    int
		err  error
		buff [1]byte
	)

	sconn, ok := conn.(syscall.Conn)

	if !ok {
		return nil
	}

	rc, err := sconn.SyscallConn()

	if err != nil {
		return err
	}

	rerr := rc.Read(func(fd uintptr) bool {
		n, err = syscall.Read(int(fd), buff[:])
		return true
	})

	switch {
	case rerr != nil:
		return rerr
	case n == 0 && err == nil:
		return io.EOF
	case n > 0:
		return errors.New("unexpected read from socket")
	case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
		return nil
	default:
		return err
	}
}

func readWrapper(conn net.Conn, wg *sync.WaitGroup) {

	fmt.Printf("\n[SERVER] Client connected %s -> %s -- Number of clients connected (%d)\n", conn.RemoteAddr(), conn.LocalAddr(), count)

	defer wg.Done()

	for {
		buf := make([]byte, 5*1024)
		_, err := conn.Read(buf)

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

					err = InsertRawDataMongo(individualRecord)
					fmt.Println(individualRecord)

					gtplDevice = ParseGTPLData(individualRecord)

					if gtplDevice.DeviceTimeNow.Day() == time.Now().Day() {
						err = InsertGTPLDataMongo(&gtplDevice)
						errcheck.CheckError(err)
					} else {
						err = InsertGTPLHistoryDataMongo(&gtplDevice)
						errcheck.CheckError(err)
					}
				}
			} else if strings.Contains(string(buf), "AVA") {

				dataSlice := strings.Split(string(buf), "*")

				var ais140Device models.AIS140Device

				for _, individualRecord := range dataSlice {

					err = InsertRawDataMongo(individualRecord)
					fmt.Println(individualRecord)
					ais140Device = ParseAIS140Data(individualRecord)

					if ais140Device.LiveOrHistoryPacket == "L" || (ais140Device.LiveOrHistoryPacket == "H" && ais140Device.DeviceTime.Day() == time.Now().Day()) {
						err = InsertAIS140DataIntoMongo(&ais140Device)
						errcheck.CheckError(err)
					} else {
						err = InsertAIS140HistoryDataMongo(&ais140Device)
						errcheck.CheckError(err)
					}
				}
			}
			dataMutex.Unlock()
		}
	}
}

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

func InsertAIS140DataIntoMongo(ais140Device *models.AIS140Device) error {

	locationHistoriesCollection, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)
	//vehicleDetailsCollection, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	options := options2.FindOptions{}
	limit := int64(1)
	options.Limit = &limit

	collectionMutex.Lock()

	_, err := locationHistoriesCollection.InsertOne(locCtx, ais140Device)
	errcheck.CheckError(err)

	collectionMutex.Unlock()

	return err
}

func InsertAIS140HistoryDataMongo(ais140device *models.AIS140Device) error {

	historyLHcollection, hctx := db.GetHistoryCollectionsWithContext(historyLHcollection)

	collectionMutex.Lock()
	_, err := historyLHcollection.InsertOne(hctx, ais140device)
	errcheck.CheckError(err)

	collectionMutex.Unlock()

	return err
}

func InsertGTPLDataMongo(gtplDevice *models.GTPLDevice) error {

	locationHistoriesCollection, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)
	vehicleDetailsCollection, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	options := options2.FindOptions{}
	limit := int64(1)
	options.Limit = &limit

	collectionMutex.Lock()

	cursor, err := vehicleDetailsCollection.Find(vctx, bson.M{"deviceid": gtplDevice.DeviceID}, &options)

	errcheck.CheckError(err)

	if cursor.Next(vctx) {
		_, err := vehicleDetailsCollection.ReplaceOne(vctx, bson.M{"deviceid": gtplDevice.DeviceID}, gtplDevice)
		errcheck.CheckError(err)

	} else {
		_, err = vehicleDetailsCollection.InsertOne(vctx, gtplDevice)
		errcheck.CheckError(err)
	}

	_, err = locationHistoriesCollection.InsertOne(locCtx, gtplDevice)
	collectionMutex.Unlock()

	return err
}

func InsertGTPLHistoryDataMongo(gtplDevice *models.GTPLDevice) error {
	historyLHcollection, hctx := db.GetHistoryCollectionsWithContext(historyLHcollection)

	collectionMutex.Lock()
	_, err := historyLHcollection.InsertOne(hctx, gtplDevice)
	errcheck.CheckError(err)

	collectionMutex.Unlock()

	return err
}

func InsertRawDataMongo(rawData string) error {

	rawDataCollection, rctx := db.GetHistoryCollectionsWithContext(rawDataCollection)

	rd := &models.RawData{
		RawData: rawData,
	}

	collectionMutex.Lock()
	_, err := rawDataCollection.InsertOne(rctx, rd)
	errcheck.CheckError(err)

	collectionMutex.Unlock()

	return err
}
