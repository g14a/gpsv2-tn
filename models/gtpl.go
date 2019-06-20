package models

import "time"

// GTPLDevice is a model of a GTPL device
type GTPLDevice struct {
	Header                      string
	DeviceID                    string
	GPSValidity                 string
	DeviceDate                  string
	DeviceTime                  string
	CurrentDateAndTimeFormatted string
	CreatedAtTimestamp          int32
	Latitude                    float64
	LatitudeDirection           string
	Longitude                   float64
	LongitudeDirection          string
	DeviceTimeNow               time.Time
	TimeNowISO                  string
	FormatTime                  string
	Speed                       string
	GPSOdometer                 string
	Direction                   string
	NumberOfSatellites          string
	BoxStatus                   string
	GSMSignal                   string
	MainBatteryStatus           string
	IgnitionStatus              string
	AnalogVoltage               string
	InternalBTVoltage           string
	Location                    string
	PreLat                      string
	PreLong                     string

	//
	Distance int
}
