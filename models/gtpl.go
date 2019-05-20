package models

type GTPLDevice struct {
	Header string
	DeviceID string
	GPSValidity string
	CurrentDateAndTime string
	CurrentDateAndTimeFormatted string
	CreatedAtTimestamp int32
	Latitude string
	LatitudeDirection string
	Longitude string
	LongitudeDirection string
	TimeNow string
	TimeNowISO string
	FormatTime string
	Speed string
	GPSOdometer string
	Direction string
	NumberOfSatellites string
	BoxStatus string
	GSMSignal string
	MainBatteryStatus string
	IgnitionStatus string
	AnalogVoltage string
	InternalBTVoltage string
	RawData string
	Location string
	PreLat string
	PreLong string

	// geofencing details
	CurrentFenceID string
	CurrentFenceExist bool
	FenceIDsOfDevice []string
}
