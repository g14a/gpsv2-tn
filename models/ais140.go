package models

type AIS140Device struct {
	Header string
	VendorIdentification string
	SoftwareVersion string
	PacketType string
	PacketIdentification string
	LiveOrHistoryPacket string
	IMEINumber int
	VehicleRegNumber string
	GPSFix string
	DateInDDMMYYYY string
	TimeInHHMMSS string
	Latitude float64
	Longitude float64
	Speed int
	Heading int
	NumberOfSatellites int
	AltitudeInMetres int
	PDOP float64
	HDOP float64
	NetworkOperator string
	IgnitionStatus string
	VehicleBatteryStatus int
	VehicleBatteryVoltage float64
	InternalBatteryVoltage float64
	EmergencyStatus string
	TamperAlert string
	GSMStrength int
	MCC int
	MNC int
	LAC int
	CellID int
	NMR string
	DigitalInputStatus string
	DigitalOutputStatus string
	SequenceNumber string
	CheckSum string

	// Done
}