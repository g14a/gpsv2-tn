package main

import (
	"encoding/csv"
	"fmt"
	"gitlab.com/gpsv2-withoutrmtesting/models"
	"gitlab.com/gpsv2-withoutrmtesting/utils"
	"strconv"
	"strings"
	"time"
)

var raw = []string{"BSTPL$1,869867039013448,A,030919,110054,12.504572,N,79.603163,E,0,8557,275,09,0,21,1,1,13.91,4.12,AVA17A_V1.1,L",
				   "BSTPL$1,869867039013448,A,030919,110054,12.504572,N,79.603163,E,0,8557,275,09,0,21,1,1,13.91,4.12,AVA17A_V1.1,L",
				   "BSTPL$1,869867039013448,A,030919,110054,12.504572,N,79.603163,E,0,8557,275,09,0,21,1,1,13.91,4.12,AVA17A_V1.1,L",
				   "BSTPL$1,869867039013448,A,030919,110054,12.504572,N,79.603163,E,0,8557,275,09,0,21,1,1,13.91,4.12,AVA17A_V1.1,L",
				   "BSTPL$1,869867039013448,A,030919,110054,12.504572,N,79.603163,E,0,8557,275,09,0,21,1,1,13.91,4.12,AVA17A_V1.1,L"}

func main() {
	for i:=0;i<1000;i++ {
		start := time.Now()
		for k:=0;k<len(raw);k++ {
			_ = ParseBSTPLData(raw[k])
		}
		fmt.Println(time.Since(start))
	}
}

func ParseBSTPLData(rawData string) models.BSTPLDevice {
	var bstplDevice models.BSTPLDevice

	bstplDevice.CreatedTime = time.Now()

	r := csv.NewReader(strings.NewReader(rawData))
	csvData, _ := r.ReadAll()

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

			bstplDevice.DigitalInputStatus = 1

			if csvArray[16] == "0" {
				bstplDevice.DigitalInputStatus = 0
			}

			bstplDevice.AnalogInput, _ = strconv.ParseFloat(csvArray[17], 10)
			bstplDevice.LiveOrHistoryPacket = csvArray[20]

			bstplDevice.ButtonCode = 98

			if bstplDevice.LiveOrHistoryPacket == "L" {
				bstplDevice.ButtonCode = 99
			}

			bstplDevice.Port = 7788
		}
	}

	return bstplDevice
}
