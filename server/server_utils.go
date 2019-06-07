package server

import (
	"encoding/csv"
	"errors"
	"fmt"
	"gitlab.com/gps2.0/config"
	"gitlab.com/gps2.0/db"
	"gitlab.com/gps2.0/errcheck"
	"gitlab.com/gps2.0/models"
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
)

var clients []net.Conn
var count = 0

func HandleConnection(conn net.Conn) {

	defer removeClient(conn)
	errorChan := make(chan error)
	dataChan := make(chan []byte)

	go readWrapper(conn, dataChan, errorChan)
	go connCheckForShutdown(conn)

	for {
		select {
		case data := <-dataChan:

			log.Printf("[SERVER} Client %s sent: %s", conn.RemoteAddr(), string(data))
			gtplDevice := ParseGTPLData(string(data))

			fmt.Println(gtplDevice)

		case err := <-errorChan:
			log.Println("[SERVER] An error occured:", err.Error())
			return
		}
	}
}

var (
	locationHistoriesCollection = config.GetAppConfig().Mongoconfig.Collections.LocationHistoriesCollection
	vehicleDetailsCollection = config.GetAppConfig().Mongoconfig.Collections.VehicleDetailsCollection
	collectionMutex = &sync.Mutex{}
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

func ParseGTPLData(rawData string) models.GTPLDevice {

	r := csv.NewReader(strings.NewReader(rawData))

	csvData, err := r.ReadAll()

	errcheck.CheckError(err)

	var gtplDevice models.GTPLDevice

	for _, csvArray := range csvData {
		gtplDevice.Header = csvArray[0]
		gtplDevice.DeviceID = csvArray[1]
		gtplDevice.GPSValidity = csvArray[2]
		gtplDevice.DeviceDate = csvArray[3]
		gtplDevice.DeviceTime = csvArray[4]
		gtplDevice.Latitude = csvArray[5]
		gtplDevice.LatitudeDirection = csvArray[6]
		gtplDevice.Longitude = csvArray[7]
		gtplDevice.LongitudeDirection = csvArray[8]
		gtplDevice.Speed = csvArray[9]
		gtplDevice.GPSOdometer = csvArray[10]
		gtplDevice.Direction = csvArray[11]
		gtplDevice.NumberOfSatellites = csvArray[12]
		gtplDevice.BoxStatus = csvArray[13]
		gtplDevice.GSMSignal = csvArray[14]
		gtplDevice.MainBatteryStatus = csvArray[15]
		gtplDevice.IgnitionStatus = csvArray[16]
		gtplDevice.AnalogVoltage = csvArray[17]
		gtplDevice.DeviceTimeNow = ConvertToUnixTS(gtplDevice.DeviceDate, gtplDevice.DeviceTime)
	}

	return gtplDevice
}

func readWrapper(conn net.Conn, dataChan chan []byte, errorChan chan error) {
	for {
		buf := make([]byte, 5*1024)
		reqLen, err := conn.Read(buf)
		if err != nil {
			errorChan <- err
			return
		}
		dataChan <- buf[:reqLen]
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

func removeClient(conn net.Conn) {
	log.Printf("[SERVER] Client %s disconnected", conn.RemoteAddr())
	count--
	conn.Close()
	//remove client from clients here
}

func InsertGTPLDataIntoMongo(gtplDevice *models.GTPLDevice) error {
	locationHistoriesCollection, ctx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

	collectionMutex.Lock()

	_, err := locationHistoriesCollection.InsertOne(ctx, gtplDevice)

	collectionMutex.Unlock()
	errcheck.CheckError(err)

	return err
}

func UpdateGTPLDataIntoMongo(gtpldevice *models.GTPLDevice) error {

	vehicleDetailsCollection, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	options := options2.FindOptions{}
	limit := int64(1)
	options.Limit = &limit

	cursor, err := vehicleDetailsCollection.Find(ctx, bson.M{"deviceid": gtpldevice.DeviceID}, &options)

	if cursor.Next(ctx) {

		collectionMutex.Lock()

		_, err = vehicleDetailsCollection.ReplaceOne(ctx, bson.M{"deviceid": gtpldevice.DeviceID}, &gtpldevice)

		collectionMutex.Unlock()
		errcheck.CheckError(err)

		return err
	}

	collectionMutex.Lock()

	_, err = vehicleDetailsCollection.InsertOne(ctx, gtpldevice)

	collectionMutex.Unlock()
	errcheck.CheckError(err)

	return err
}