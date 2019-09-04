package server

import (
	"strconv"
	"time"

	"gitlab.com/gpsv2-withoutrm/models"
)

// ParseGTPLData parses the raw data sent
// by the GTPL device and marshals into a GTPL device model
func ParseGTPLDataSQL(gtplDevice models.GTPLDevice) models.GTPLSQLModel {

	var sqlDevice models.GTPLSQLModel

	sqlDevice.CreateDate = time.Now()
	sqlDevice.Fixtime = gtplDevice.DeviceTime
	sqlDevice.DeviceID, _ = strconv.Atoi(gtplDevice.DeviceID)
	sqlDevice.GPSValidity = gtplDevice.GPSValidity
	sqlDevice.LatMessage = gtplDevice.Latitude
	sqlDevice.NorthDirection = gtplDevice.LatitudeDirection
	sqlDevice.LonMessage = gtplDevice.Longitude
	sqlDevice.EastDirection = gtplDevice.LongitudeDirection
	sqlDevice.Speed = strconv.Itoa(gtplDevice.Speed)
	sqlDevice.OdometerReading = strconv.Itoa(gtplDevice.GPSOdometer)
	sqlDevice.CardinalHead = float64(gtplDevice.Direction)
	sqlDevice.NumberOfSatellites = strconv.Itoa(gtplDevice.NumberOfSatellites)
	sqlDevice.Box = strconv.FormatBool(gtplDevice.BoxStatus)
	sqlDevice.GSM = strconv.Itoa(gtplDevice.GSMSignal)
	sqlDevice.MainBatteryStatus = gtplDevice.MainBatteryStatus
	sqlDevice.IgnitionStatus = strconv.FormatBool(gtplDevice.IgnitionStatus)
	sqlDevice.Analog = strconv.FormatFloat(gtplDevice.AnalogVoltage, 'E', -1, 32)
	sqlDevice.InternalBatteryVoltage = 0

	return sqlDevice
}

func ParseAIS140DataSQL(ais140Device models.AIS140Device) models.AIS140SQLModel {

	var sqlDevice models.AIS140SQLModel

	sqlDevice.Header = ais140Device.Header
	sqlDevice.DeviceID = ais140Device.IMEINumber
	sqlDevice.VendorIdentification = ais140Device.VendorIdentification
	sqlDevice.SoftwareVersion = ais140Device.SoftwareVersion
	sqlDevice.PacketType = ais140Device.PacketType
	sqlDevice.PacketIdentification = ais140Device.PacketIdentification
	sqlDevice.VehicleRegNumber = ais140Device.VehicleRegNumber
	sqlDevice.LatMessage = ais140Device.Latitude
	sqlDevice.LatitudeDirection = ais140Device.LatitudeDirection
	sqlDevice.LonMessage = ais140Device.Longitude
	sqlDevice.LongitudeDirection = ais140Device.LongitudeDirection
	sqlDevice.Speed = ais140Device.Speed
	sqlDevice.LiveOrHistoryPacket = ais140Device.LiveOrHistoryPacket
	sqlDevice.GPSFix = ais140Device.GPSFix
	sqlDevice.Altitude = strconv.Itoa(ais140Device.AltitudeInMetres)
	sqlDevice.CardinalHead = ais140Device.Heading
	sqlDevice.NumberOfSatellites = strconv.Itoa(ais140Device.NumberOfSatellites)
	sqlDevice.GSM = strconv.Itoa(ais140Device.GSMStrength)
	sqlDevice.DigitalOutputStatus = ais140Device.DigitalOutputStatus
	sqlDevice.DigitalInputStatus = ais140Device.DigitalInputStatus
	sqlDevice.SequenceNumber = ais140Device.SequenceNumber
	sqlDevice.CheckSumValue = ais140Device.CheckSum
	sqlDevice.TamperAlert = ais140Device.TamperAlert
	sqlDevice.InternalBatteryVoltage = ais140Device.InternalBatteryVoltage
	sqlDevice.VehicleBatteryVoltage = ais140Device.VehicleBatteryVoltage

	sqlDevice.FixTime = ais140Device.DeviceTime
	sqlDevice.CreateDate = time.Now()

	sqlDevice.Network = ais140Device.NetworkOperator

	sqlDevice.IgnitionStatus = 0

	if ais140Device.IgnitionStatus {
		sqlDevice.IgnitionStatus = 1
	}

	sqlDevice.NumberOfSatellites = strconv.Itoa(ais140Device.NumberOfSatellites)
	sqlDevice.EmergencyStatus = ais140Device.EmergencyStatus

	return sqlDevice
}

func ParseBSTPLDataSQL(bstplDevice models.BSTPLDevice) models.BSTPLSQLModel {

	var sqlDevice models.BSTPLSQLModel

	sqlDevice.CreateDate = time.Now()
	sqlDevice.DeviceID, _ = strconv.Atoi(bstplDevice.VehicleID)
	sqlDevice.GPSValidity = bstplDevice.GPSValidity
	sqlDevice.LatMessage = bstplDevice.Latitude
	sqlDevice.NorthDirection = bstplDevice.LatitudeDirection
	sqlDevice.LonMessage = bstplDevice.Longitude
	sqlDevice.EastDirection = bstplDevice.LongitudeDirection
	sqlDevice.Speed = strconv.FormatFloat(bstplDevice.Speed, 'E', -1, 32)
	sqlDevice.OdometerReading = strconv.FormatFloat(bstplDevice.GPSOdometer, 'E', -1, 32)
	sqlDevice.CardinalHead, _ = strconv.ParseFloat(bstplDevice.Direction, 10)
	sqlDevice.NumberOfSatellites = strconv.Itoa(bstplDevice.NumberOfSatellites)
	sqlDevice.Box = bstplDevice.BoxOpenCloseStatus
	sqlDevice.GSM = strconv.Itoa(bstplDevice.GSMSignal)
	sqlDevice.MainBatteryStatus = string(bstplDevice.MainBatteryStatus)
	sqlDevice.IgnitionStatus = strconv.Itoa(bstplDevice.DigitalInputStatus)
	sqlDevice.Analog = strconv.FormatFloat(bstplDevice.AnalogInput, 'E', -1, 32)
	sqlDevice.InternalBatteryVoltage = 0
	sqlDevice.CreateDate = time.Now()
	sqlDevice.Fixtime = bstplDevice.DeviceTime

	return sqlDevice
}
