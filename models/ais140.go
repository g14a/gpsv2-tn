// Package models contains the models of specified devices and raw data
package models

import "time"

// AIS140Device is a model of an AIS140 device
// and contains the following fields
type AIS140Device struct {
	Header                 string
	VendorIdentification   string
	SoftwareVersion        string
	PacketType             string
	PacketIdentification   string
	LiveOrHistoryPacket    string
	IMEINumber             string
	VehicleRegNumber       string
	GPSFix                 string
	DateInDDMMYYYY         string
	TimeInHHMMSS           string
	Latitude               float64
	Longitude              float64
	Speed                  int
	Heading                int
	NumberOfSatellites     int
	AltitudeInMetres       int
	PDOP                   float64
	HDOP                   float64
	NetworkOperator        string
	IgnitionStatus         int
	VehicleBatteryStatus   int
	VehicleBatteryVoltage  float64
	InternalBatteryVoltage float64
	EmergencyStatus        string
	TamperAlert            string
	GSMStrength            int
	MCC                    int
	MNC                    int
	LAC                    int
	CellID                 string
	NMR                    string
	DigitalInputStatus     string
	DigitalOutputStatus    string
	SequenceNumber         string
	CheckSum               string

	// Custom fields
	DeviceTime        time.Time
	InsertedTimeStamp time.Time
	Distance          int
}
