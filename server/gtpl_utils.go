package server

import (
	"encoding/csv"
	"gitlab.com/gpsv2/errcheck"
	"gitlab.com/gpsv2/models"
	"strings"
)

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
