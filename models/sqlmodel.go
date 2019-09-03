package models

import "time"

type GTPLSQLModel struct {
	DeviceID int	`sql:"primary_key"`
	GPSValidity string
	LatMessage float64
	LonMessage float64
	Speed string
	CardinalHead float64
	DistanceTravelled float64
	OdometerReading string
	NumberOfSatellites string
	Box string
	NorthDirection string
	EastDirection string
	Analog string
	GSM string
	IgnitionStatus string
	InternalBatteryVoltage int
	PreLat float64
	PreLong float64
	CreateDate time.Time
	MainBatteryStatus string
	Fixtime time.Time
	Distance float64
}

type AIS140SQLModel struct {
	Header                 string
	VendorIdentification   string
	SoftwareVersion        string
	PacketType             string
	PacketIdentification   string
	LiveOrHistoryPacket    string
	DeviceID               string
	VehicleRegNumber       string
	GPSFix                 string
	DateInDDMMYYYY         string
	TimeInHHMMSS           string
	LatMessage             float64
	LatitudeDirection	   string
	LonMessage             float64
	LongitudeDirection 	   string
	Speed                  int
	CardinalHead           int
	NumberOfSatellites     string
	Altitude		       string
	PDOP                   string
	HDOP                   string
	Network			       string
	IgnitionStatus         int
	VehicleBatteryStatus   int
	VehicleBatteryVoltage  float64
	InternalBatteryVoltage float64
	EmergencyStatus        string
	TamperAlert            string
	GSM            		   string
	MCC                    string
	MNC                    string
	LAC                    string
	CellID                 string
	NMR                    string
	DigitalInputStatus     string
	DigitalOutputStatus    string
	SequenceNumber         string
	CheckSumValue          string
	FixTime 			   time.Time
	CreateDate 			   time.Time
	DistanceTravelled 	   int
}

type BSTPLSQLModel struct {
	DeviceID int	`sql:"primary_key"`
	GPSValidity string
	LatMessage float64
	LatDirection string
	LonMessage float64
	LongDirection float64
	Speed string
	CardinalHead float64
	DistanceTravelled float64
	OdometerReading string
	NumberOfSatellites string
	Box string
	NorthDirection string
	EastDirection string
	Analog string
	GSM string
	IgnitionStatus string
	InternalBatteryVoltage int
	PreLat float64
	PreLong float64
	CreateDate time.Time
	MainBatteryStatus string
	Fixtime time.Time
	Distance float64
}