package server

import (
	"gitlab.com/gpsv2-tn/dbutils"
	"gitlab.com/gpsv2-tn/errorcheck"
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

	// Initialize a buffer of 1KB to be read from the client and read using conn.Read
	buf := make([]byte, 1024)

	for {
		_, err := conn.Read(buf)

		// if an error occurs deal with it
		if err != nil {
			if err == io.EOF {
				errorcheck.CheckError(err)
				_ = conn.Close()
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

	gtplDevice := modelutils.ParseGTPLData(record)

	// ignores if an empty data occurs
	if gtplDevice.Latitude != 0 && gtplDevice.Longitude != 0 && gtplDevice.DeviceID != "" {

		mssqlDevice := dbutils.ParseMSSQLDeviceFromGTPL(gtplDevice)

		go dbutils.InsertIntoMSSQL(mssqlDevice)
		go dbutils.InsertGTPLDataMongo(&gtplDevice, startTime)

		mysqlDevice := modelutils.ParseGTPLDataSQL(gtplDevice)

		dbutils.InsertGTPLIntoSQL(mysqlDevice)
	}
}

func processBSTPLDevices(record string, startTime time.Time) {

	bstplDevice := modelutils.ParseBSTPLData(record)

	if bstplDevice.Latitude != 0 && bstplDevice.Longitude != 0 && bstplDevice.VehicleID != "" {

		recvTime := time.Now()

		go dbutils.InsertBSTPLDataMongo(&bstplDevice, startTime)

		mssqlDevice := dbutils.ParseMSSQLDeviceFromBSTPL(bstplDevice)
		mssqlDevice.RecvTime = recvTime

		go dbutils.InsertIntoMSSQL(mssqlDevice)
		mysqlDevice := modelutils.ParseBSTPLDataSQL(bstplDevice)

		dbutils.InsertBSTPLIntoSQL(mysqlDevice)
	}
}

func processAIS140Device(record string, startTime time.Time) {

	ais140Device := modelutils.ParseAIS140Data(record)

	// ignores if an empty data occurs
	if ais140Device.Latitude != 0 && ais140Device.Longitude != 0 && ais140Device.IMEINumber != "" {

		mssqlDevice := dbutils.ParseMSSQLDeviceFromAIS140(ais140Device)

		go dbutils.InsertIntoMSSQL(mssqlDevice)
		go dbutils.InsertAIS140DataIntoMongo(&ais140Device, startTime)

		mysqlDevice := modelutils.ParseAIS140DataSQL(ais140Device)
		dbutils.InsertAIS140IntoSQL(mysqlDevice)
	}
}
