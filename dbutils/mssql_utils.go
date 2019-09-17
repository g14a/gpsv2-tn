package dbutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/gpsv2-tn/models"
)

var (
	client = &http.Client{
		Transport: &http.Transport {
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout: 0,
				KeepAlive: 0,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
)

// ParseMSSQLDeviceFromGTPL parses the GTPL device model into a model which is
// accepted by MSSQL queries.
func ParseMSSQLDeviceFromGTPL(device models.GTPLDevice) models.MSSQLDevice {
	var mssqlDevice models.MSSQLDevice

	mssqlDevice.RecvTime = device.CreatedTime
	mssqlDevice.SentTime = time.Now()

	vehicleNumber := getVehicleRegNo(device.DeviceID)

	mssqlDevice.ButtonCode = strconv.Itoa(int(device.ButtonCode))
	mssqlDevice.Geofence = device.GeoFenceID

	mssqlDevice.VehicleNumber = vehicleNumber
	mssqlDevice.Longitude = device.Longitude
	mssqlDevice.Latitude = device.Latitude
	mssqlDevice.Speed = float64(device.Speed)
	mssqlDevice.Ignition = device.IgnitionStatus
	mssqlDevice.BatteryStatus = device.MainBatteryStatus

	i := fmt.Sprintf("%.0f", device.InternalBTVoltage)

	mssqlDevice.InternalBattery, _ = strconv.Atoi(i)

	mssqlDevice.Orientation = device.Direction
	mssqlDevice.SimID = device.DeviceID
	mssqlDevice.Fuel = 0

	return mssqlDevice
}

// ParseMSSQLDeviceFromAIS140 parses the AIS140 device model into a model which is
// accepted by MSSQL queries.
func ParseMSSQLDeviceFromAIS140(device models.AIS140Device) models.MSSQLDevice {
	var mssqlDevice models.MSSQLDevice

	vehicleNumber := getVehicleRegNo(device.IMEINumber)

	mssqlDevice.RecvTime = device.CreatedTime
	mssqlDevice.SentTime = time.Now()

	mssqlDevice.ButtonCode = strconv.Itoa(int(device.ButtonCode))
	mssqlDevice.Geofence = device.GeoFenceID
	mssqlDevice.LiveHistory = device.LiveOrHistoryPacket

	mssqlDevice.VehicleNumber = vehicleNumber
	mssqlDevice.Longitude = device.Longitude
	mssqlDevice.Latitude = device.Latitude
	mssqlDevice.Speed = float64(device.Speed)
	mssqlDevice.Ignition = device.IgnitionStatus
	mssqlDevice.BatteryStatus = string(device.VehicleBatteryStatus)

	i := fmt.Sprintf("%.0f", device.InternalBatteryVoltage)

	mssqlDevice.InternalBattery, _ = strconv.Atoi(i)

	mssqlDevice.Orientation = device.Heading
	mssqlDevice.SimID = device.IMEINumber

	return mssqlDevice
}


// ParseMSSQLDeviceFromBSTPL parses the BSTPL device model into a model which is
// accepted by MSSQL queries.
func ParseMSSQLDeviceFromBSTPL(device models.BSTPLDevice) models.MSSQLDevice {

	var mssqlDevice models.MSSQLDevice

	mssqlDevice.RecvTime = device.CreatedTime
	mssqlDevice.SentTime = time.Now()

	mssqlDevice.LiveHistory = device.LiveOrHistoryPacket

	vehicleNumber := getVehicleRegNo(device.VehicleID)

	mssqlDevice.ButtonCode = strconv.Itoa(int(device.ButtonCode))
	mssqlDevice.Geofence = device.GeoFenceID

	mssqlDevice.VehicleNumber = vehicleNumber
	mssqlDevice.Longitude = device.Longitude
	mssqlDevice.Latitude = device.Latitude
	mssqlDevice.Speed = device.Speed
	mssqlDevice.Ignition = true

	if device.DigitalInputStatus {
		mssqlDevice.Ignition = false
	}

	mssqlDevice.BatteryStatus = strconv.Itoa(device.MainBatteryStatus)

	i := fmt.Sprintf("%.0f", device.InternalBatteryVoltage)

	mssqlDevice.InternalBattery, _ = strconv.Atoi(i)

	mssqlDevice.Orientation, _ = strconv.Atoi(device.Direction)
	mssqlDevice.SimID = device.VehicleID

	return mssqlDevice
}


// InsertIntoMSSQL inserts the MSSQL model created from the above functions and performs
// a HTTP call to the endpoint posting the query as the body of the request. There are two insert queries
// and one update query. The query is formed using Sprintf and marshalled into JSON, then passed into the body.
func InsertIntoMSSQL(device models.MSSQLDevice) {

	if device.LiveHistory == "L" {
		ignition := 0

		if device.Ignition {
			ignition = 1
		}

		sentTime := device.SentTime.Format("2006-01-02 15:04:05")
		recvTime := device.RecvTime.Format("2006-01-02 15:04:05")

		iq1 := fmt.Sprintf("INSERT INTO T_Data(SIM_ID, Sent_Time, Rec_Time, Latitude, Longitude, Speed, Orientation," +
			"OdoRead, OdoDecimal, Fuel, Ignition, Battery_Status, Button_Code, VehicleNumber, " +
			"Internal_battery, Geofence) VALUES ('%s','%v','%v','%f','%f','%d','%d','%d','%d',%d,%d,%s,%s,'%s',%d,%v)",
			device.SimID, sentTime, recvTime, device.Latitude, device.Longitude,
			int(device.Speed), device.Orientation, device.OdometerReading, int(device.OdoDecimal),
			device.Fuel, ignition, device.BatteryStatus, device.ButtonCode, device.VehicleNumber,
			device.InternalBattery, device.Geofence)

		iq2 := fmt.Sprintf("INSERT INTO T_LatestData7 (SIM_ID, Sent_Time, Rec_Time, Latitude, Longitude, Speed, Orientation," +
			"OdoRead, OdoDecimal, Fuel, Ignition, Battery_Status, Button_Code, VehicleNumber, " +
			"Internal_battery, Geofence) VALUES ('%s','%v','%v','%f','%f','%d','%d','%d','%d',%d,%d,%s,%s,'%s',%d,%v)",
			device.SimID, sentTime, recvTime, device.Latitude, device.Longitude,
			int(device.Speed), device.Orientation, device.OdometerReading, int(device.OdoDecimal),
			device.Fuel, ignition, device.BatteryStatus, device.ButtonCode, device.VehicleNumber,
			device.InternalBattery, device.Geofence)

		uq := fmt.Sprintf("UPDATE T_LatestData SET SIM_ID='%s',Sent_Time='%v',Rec_Time='%v',Latitude='%f',Longitude='%f',Speed=%f,Orientation='%d',OdoRead='%d',OdoDecimal='%f',Fuel=%d,Ignition=%d,Button_Code=%s,Battery_Status=%s,Geofence=%v,Internal_battery=%d where vehiclenumber='%s'",
			device.SimID, sentTime, recvTime, device.Latitude, device.Longitude, device.Speed,
			device.Orientation, device.OdometerReading, device.OdoDecimal, device.Fuel,
			ignition, device.ButtonCode, device.BatteryStatus,device.Geofence,device.InternalBattery, device.VehicleNumber)

		insertQuery1 := &Query{
			Query: iq1,
		}

		insertQuery2 := &Query{
			Query: iq2,
		}

		updateQuery := &Query{
			Query: uq,
		}

		insertJson1, _ := json.Marshal(insertQuery1)
		insertJson2, _ := json.Marshal(insertQuery2)
		updateJson, _ := json.Marshal(updateQuery)

		inreq1, _ :=  http.NewRequest("POST", "http://206.189.137.125:3009/postGpsGovt", bytes.NewReader([]byte(insertJson1)))

		inreq2, _ :=  http.NewRequest("POST", "http://206.189.137.125:3009/postGpsGovt", bytes.NewReader([]byte(insertJson2)))

		upreq, _ := http.NewRequest("POST", "http://206.189.137.125:3010/postGpsGovt", bytes.NewReader([]byte(updateJson)))

		inreq1.Header.Add("Content-Type", "application/json")
		inreq2.Header.Add("Content-Type", "application/json")
		upreq.Header.Add("Content-Type", "application/json")

		_, _ = client.Do(inreq1)
		_, _ = client.Do(inreq2)

		_, _ = client.Do(upreq)

		inreq1.Close = true
		inreq2.Close = true

		upreq.Close = true

	}
}

type Query struct {
	Query string `json:"query1"`
}
