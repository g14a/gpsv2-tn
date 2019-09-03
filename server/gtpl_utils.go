package server

import (
	"encoding/csv"
	"gitlab.com/gpsv2-withoutrm/errorcheck"
	"gitlab.com/gpsv2-withoutrm/models"
	"gitlab.com/gpsv2-withoutrm/utils"
	"strconv"
	"strings"
	"time"
)

// ParseGTPLData parses the raw data sent
// by the GTPL device and marshals into a GTPL device model
func ParseGTPLData(rawData string) models.GTPLDevice {

	r := csv.NewReader(strings.NewReader(rawData))
	csvData, err := r.ReadAll()

	errorcheck.CheckError(err)

	var gtplDevice models.GTPLDevice

	for _, csvArray := range csvData {
		if len(csvArray) >= 18 {
			gtplDevice.Header = csvArray[0]
			gtplDevice.DeviceID = csvArray[1]
			gtplDevice.GPSValidity = csvArray[2]
			gtplDevice.RawDeviceDate = csvArray[3]
			gtplDevice.RawDeviceTime = csvArray[4]
			gtplDevice.Latitude, _ = strconv.ParseFloat(csvArray[5], 10)
			gtplDevice.LatitudeDirection = csvArray[6]
			gtplDevice.Longitude, _ = strconv.ParseFloat(csvArray[7], 10)
			gtplDevice.LongitudeDirection = csvArray[8]
			gtplDevice.Speed, _ = strconv.Atoi(csvArray[9])
			gtplDevice.GPSOdometer,_  = strconv.Atoi(csvArray[10])
			gtplDevice.Direction, _ = strconv.Atoi(csvArray[11])
			gtplDevice.NumberOfSatellites, _ = strconv.Atoi(csvArray[12])
			gtplDevice.BoxStatus = true

			if csvArray[13] == "0" {
				gtplDevice.BoxStatus = false
			}

			gtplDevice.GSMSignal, _ = strconv.Atoi(csvArray[14])
			gtplDevice.MainBatteryStatus = csvArray[15]

			gtplDevice.IgnitionStatus = true
			if csvArray[16] == "0" {
				gtplDevice.IgnitionStatus = false
			}

			gtplDevice.AnalogVoltage, _ = strconv.ParseFloat(csvArray[17], 10)

			// Time Fields
			gtplDevice.DeviceTime = utils.ConvertTimeGTPL(gtplDevice.RawDeviceDate, gtplDevice.RawDeviceTime)
			gtplDevice.CreatedTime = time.Now()

			gtplDevice.ButtonCode = 98

			if utils.GTPLCheckLiveHistory(gtplDevice.DeviceTime) {
				gtplDevice.ButtonCode = 99
			}

			gtplDevice.Port = 7788
		}
	}

	return gtplDevice
}
