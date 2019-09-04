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

// ParseAIS140Data parses the raw data sent
// by the AIS140 device and marshals into a AIS140 device model
func ParseAIS140Data(rawData string) models.AIS140Device {

	r := csv.NewReader(strings.NewReader(rawData))

	csvData, err := r.ReadAll()

	errorcheck.CheckError(err)

	var ais140Device models.AIS140Device

	for _, csvarray := range csvData {

		if len(csvarray) >= 37 {
			ais140Device.Header = csvarray[1]
			ais140Device.VendorIdentification = csvarray[2]
			ais140Device.SoftwareVersion = csvarray[3]
			ais140Device.PacketType = csvarray[4]
			ais140Device.PacketIdentification = csvarray[5]
			ais140Device.LiveOrHistoryPacket = csvarray[6]
			ais140Device.IMEINumber = csvarray[7]
			ais140Device.VehicleRegNumber = csvarray[8]
			ais140Device.GPSFix = csvarray[9]
			ais140Device.DateInDDMMYYYY = csvarray[10]
			ais140Device.TimeInHHMMSS = csvarray[11]
			ais140Device.Latitude, _ = strconv.ParseFloat(csvarray[12], 10)
			ais140Device.LatitudeDirection = csvarray[13]
			ais140Device.Longitude, _ = strconv.ParseFloat(csvarray[14], 10)
			ais140Device.LongitudeDirection = csvarray[15]
			ais140Device.Speed, _ = strconv.Atoi(csvarray[16])
			ais140Device.Heading, _ = strconv.Atoi(csvarray[17])
			ais140Device.NumberOfSatellites, _ = strconv.Atoi(csvarray[18])
			ais140Device.AltitudeInMetres, _ = strconv.Atoi(csvarray[19])
			ais140Device.NetworkOperator = csvarray[22]

			ais140Device.IgnitionStatus = true

			if csvarray[23] == "0" {
				ais140Device.IgnitionStatus = false
			}

			ais140Device.VehicleBatteryStatus, _ = strconv.Atoi(csvarray[24])
			ais140Device.VehicleBatteryVoltage, _ = strconv.ParseFloat(csvarray[25], 10)
			ais140Device.InternalBatteryVoltage, _ = strconv.ParseFloat(csvarray[26], 10)

			ais140Device.EmergencyStatus = csvarray[27]
			ais140Device.TamperAlert = csvarray[28]
			ais140Device.GSMStrength, _ = strconv.Atoi(csvarray[29])
			ais140Device.DigitalInputStatus = csvarray[35]
			ais140Device.DigitalOutputStatus = csvarray[36]
			ais140Device.SequenceNumber = csvarray[37]
			ais140Device.CheckSum = csvarray[38]

			//Custom fields
			ais140Device.DeviceTime = utils.ConvertTimeAIS140(ais140Device.DateInDDMMYYYY, ais140Device.TimeInHHMMSS)
			ais140Device.CreatedTime = time.Now()

			ais140Device.ButtonCode = 98

			if ais140Device.LiveOrHistoryPacket == "L" {
				ais140Device.ButtonCode = 99
			}

			ais140Device.Port = 7788
		}
	}

	return ais140Device
}
