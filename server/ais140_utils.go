package server

import (
	"encoding/csv"
	"gitlab.com/gpsv2/errcheck"
	"gitlab.com/gpsv2/models"
	"strconv"
	"strings"
)

func ParseAIS140Data(rawData string) models.AIS140Device {
	r := csv.NewReader(strings.NewReader(rawData))

	csvData, err := r.ReadAll()

	errcheck.CheckError(err)

	var ais140Device models.AIS140Device

	for _, csvarray := range csvData {
		ais140Device.Header = csvarray[0]
		ais140Device.VendorIdentification = csvarray[1]
		ais140Device.SoftwareVersion = csvarray[2]
		ais140Device.PacketType = csvarray[3]
		ais140Device.PacketIdentification = csvarray[4]
		ais140Device.LiveOrHistoryPacket = csvarray[5]
		ais140Device.IMEINumber, _ = strconv.Atoi(csvarray[6])
		ais140Device.VehicleRegNumber = csvarray[7]
		ais140Device.GPSFix = csvarray[8]
		ais140Device.DateInDDMMYYYY = csvarray[9]
		ais140Device.TimeInHHMMSS = csvarray[10]
		ais140Device.Latitude, _ = strconv.ParseFloat(csvarray[11], 10)
		ais140Device.Longitude, _ = strconv.ParseFloat(csvarray[13], 10)
		ais140Device.Speed, _ = strconv.Atoi(csvarray[14])
		ais140Device.Heading, _ = strconv.Atoi(csvarray[15])
		ais140Device.NumberOfSatellites, _ = strconv.Atoi(csvarray[16])
		ais140Device.AltitudeInMetres, _ = strconv.Atoi(csvarray[17])
		ais140Device.PDOP, _ = strconv.ParseFloat(csvarray[18], 10)
		ais140Device.HDOP, _ = strconv.ParseFloat(csvarray[19], 10)
		ais140Device.NetworkOperator = csvarray[20]
		ais140Device.IgnitionStatus, _ = strconv.Atoi(csvarray[21])
		ais140Device.VehicleBatteryStatus, _ = strconv.Atoi(csvarray[22])
		ais140Device.VehicleBatteryVoltage, _ = strconv.ParseFloat(csvarray[23], 10)
		ais140Device.InternalBatteryVoltage, _ = strconv.ParseFloat(csvarray[24], 10)
		ais140Device.EmergencyStatus = csvarray[25]
		ais140Device.TamperAlert = csvarray[26]
		ais140Device.GSMStrength, _ = strconv.Atoi(csvarray[27])
		ais140Device.MCC, _ = strconv.Atoi(csvarray[28])
		ais140Device.MNC, _ = strconv.Atoi(csvarray[29])
		ais140Device.LAC, _ = strconv.Atoi(csvarray[30])
		ais140Device.CellID = csvarray[31]
		ais140Device.NMR = csvarray[32]
		ais140Device.DigitalInputStatus = csvarray[33]
		ais140Device.DigitalOutputStatus = csvarray[34]
		ais140Device.SequenceNumber = csvarray[35]
		ais140Device.CheckSum = csvarray[36]
	}

	return ais140Device
}
