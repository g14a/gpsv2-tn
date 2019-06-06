package server

import (
	"bufio"
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
	"net"
	"strings"
	"sync"
	"syscall"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Message received: ", message)

		gtplDevice := ParseGTPLData(message)

		fmt.Println(gtplDevice)

		err := InsertGTPLDataIntoMongo(&gtplDevice)
		errcheck.CheckError(err)

		err = UpdateGTPLDataIntoMongo(&gtplDevice)
		errcheck.CheckError(err)

		go connCheckForShutdown(conn)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error: ", err)
	}
}

var (
	locationHistoriesCollection = config.GetAppConfig().Mongoconfig.Collections.LocationHistoriesCollection
	vehicleDetailsCollection = config.GetAppConfig().Mongoconfig.Collections.VehicleDetailsCollection
	collectionMutex = &sync.Mutex{}
)

func connCheckForShutdown(c net.Conn) error {
	var (
		n    int
		err  error
		buff [1]byte
	)

	sconn, ok := c.(syscall.Conn)

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
		gtplDevice.CurrentDateAndTime = csvArray[3]
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
	}

	return gtplDevice
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