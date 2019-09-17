package server

import (
	"fmt"
	"gitlab.com/gpsv2-tn/dbutils"
	"gitlab.com/gpsv2-tn/models"
	"gitlab.com/gpsv2-tn/modelutils"
	"io"
	"net"
	"strings"
	"time"
)

// HandleConnection handles a connection by firing
// up a seperate go routine for a TCP connection net.Conn
func HandleConnection(conn net.Conn) {

	startTime := time.Now()

	for {
		// Initialize a buffer of 5KB to be read from the client and read using conn.Read
		buf := make([]byte, 2*1024)
		_, err := conn.Read(buf)

		// if an error occurs deal with it
		if err != nil {
			fmt.Println(err.Error())
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
					processBSTPLDevices(record, startTime)
				}

				dataslice = nil

			} else if strings.Contains(buffer, "GTPL") {
				dataslice := strings.Split(string(buf), "#")

				for _, record := range dataslice {
					processGTPLDevices(record, startTime)
				}

				dataslice = nil

			} else if strings.Contains(buffer, "AVA") {
				dataslice := strings.Split(string(buf), "*")

				for _, record := range dataslice {
					processAIS140Device(record, startTime)
				}

				dataslice = nil
			}
		}
	}
}

func processGTPLDevices(record string, startTime time.Time) {

	var (
		gtplDevice  models.GTPLDevice
		mysqlDevice models.GTPLSQLModel
		mssqlDevice models.MSSQLDevice
	)

	gtplDevice = modelutils.ParseGTPLData(record)

	// ignores if an empty data occurs
	if gtplDevice.Latitude != 0 && gtplDevice.Longitude != 0 && gtplDevice.DeviceID != "" {

		fmt.Println(gtplDevice)

		mssqlDevice = dbutils.ParseMSSQLDeviceFromGTPL(gtplDevice)

		gtplDevice.Distance = mysqlDevice.DistanceTravelled

		dbutils.InsertIntoMSSQL(mssqlDevice)
		dbutils.InsertGTPLDataMongo(&gtplDevice, startTime)
		mysqlDevice = modelutils.ParseGTPLDataSQL(gtplDevice)

		dbutils.InsertGTPLIntoSQL(mysqlDevice)
		dbutils.InsertRawDataMongo(record)
	}
}

func processBSTPLDevices(record string, startTime time.Time) {

	var (
		bstplDevice models.BSTPLDevice
		mssqlDevice models.MSSQLDevice
		mysqlDevice models.BSTPLSQLModel
	)

	bstplDevice = modelutils.ParseBSTPLData(record)

	if bstplDevice.Latitude != 0 && bstplDevice.Longitude != 0 && bstplDevice.VehicleID != "" {

		fmt.Println(bstplDevice)

		recvTime := time.Now()

		dbutils.InsertBSTPLDataMongo(&bstplDevice, startTime)
		mssqlDevice = dbutils.ParseMSSQLDeviceFromBSTPL(bstplDevice)
		mssqlDevice.RecvTime = recvTime

		dbutils.InsertIntoMSSQL(mssqlDevice)
		mysqlDevice = modelutils.ParseBSTPLDataSQL(bstplDevice)

		dbutils.InsertBSTPLIntoSQL(mysqlDevice)
		dbutils.InsertRawDataMongo(record)

	}

}

func processAIS140Device(record string, startTime time.Time) {

	var (
		ais140Device  models.AIS140Device
		mysqlDevice models.AIS140SQLModel
		mssqlDevice models.MSSQLDevice
	)

	ais140Device = modelutils.ParseAIS140Data(record)

	// ignores if an empty data occurs
	if ais140Device.Latitude != 0 && ais140Device.Longitude != 0 && ais140Device.IMEINumber != "" {

		fmt.Println(ais140Device)

		mssqlDevice = dbutils.ParseMSSQLDeviceFromAIS140(ais140Device)

		dbutils.InsertIntoMSSQL(mssqlDevice)
		dbutils.InsertAIS140DataIntoMongo(&ais140Device, startTime)

		mysqlDevice = modelutils.ParseAIS140DataSQL(ais140Device)
		dbutils.InsertAIS140IntoSQL(mysqlDevice)
		dbutils.InsertRawDataMongo(record)
	}
}
