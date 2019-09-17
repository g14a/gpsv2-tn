package modelutils

import (
	"encoding/csv"
	"gitlab.com/gpsv2-tn/config"
	"gitlab.com/gpsv2-tn/errorcheck"
	"gitlab.com/gpsv2-tn/models"
	"gitlab.com/gpsv2-tn/utils"
	"strconv"
	"strings"
	"time"
)

// ParseGTPLData parses the raw data sent
// by the GTPL device and marshals into a GTPL device model and returns the model
func ParseBSTPLData(rawData string) models.BSTPLDevice {
	var bstplDevice models.BSTPLDevice

	bstplDevice.CreatedTime = time.Now()

	r := csv.NewReader(strings.NewReader(rawData))
	csvData, err := r.ReadAll()

	errorcheck.CheckError(err)

	for _, csvArray := range csvData {

		if len(csvArray) > 20 {
			bstplDevice.Header = csvArray[0]
			bstplDevice.VehicleID = csvArray[1]
			bstplDevice.GPSValidity = csvArray[2]
			bstplDevice.RawDeviceDate = csvArray[3]
			bstplDevice.RawDeviceTime = csvArray[4]
			bstplDevice.Latitude, _ = strconv.ParseFloat(csvArray[5], 10)
			bstplDevice.LatitudeDirection = csvArray[6]
			bstplDevice.Longitude, _ = strconv.ParseFloat(csvArray[7], 10)
			bstplDevice.LongitudeDirection = csvArray[8]
			bstplDevice.Speed, _ = strconv.ParseFloat(csvArray[9], 10)
			bstplDevice.GPSOdometer,_  = strconv.ParseFloat(csvArray[10], 10)
			bstplDevice.Direction = csvArray[11]
			bstplDevice.NumberOfSatellites, _ = strconv.Atoi(csvArray[12])
			bstplDevice.BoxOpenCloseStatus = "1"

			bstplDevice.DeviceTime = utils.ConvertTimeBSTPL(bstplDevice.RawDeviceDate, bstplDevice.RawDeviceTime)

			if csvArray[13] == "0" {
				bstplDevice.BoxOpenCloseStatus = "0"
			}

			bstplDevice.GSMSignal, _ = strconv.Atoi(csvArray[14])
			bstplDevice.MainBatteryStatus, _ = strconv.Atoi(csvArray[15])

			bstplDevice.DigitalInputStatus = true

			if csvArray[16] == "0" {
				bstplDevice.DigitalInputStatus = false
			}

			bstplDevice.AnalogInput, _ = strconv.ParseFloat(csvArray[17], 10)
			bstplDevice.LiveOrHistoryPacket = csvArray[20]

			bstplDevice.ButtonCode = 98

			if bstplDevice.LiveOrHistoryPacket == "L" {
				bstplDevice.ButtonCode = 99
			}

			bstplDevice.Port, _ = strconv.Atoi(config.GetAppConfig().TCPSocketConfig.Port)
		}
	}

	return bstplDevice
}
