package models

import "time"

type GTPLDevice struct {
	Header                      string
	DeviceID                    string
	GPSValidity                 string
	DeviceDate                  string
	DeviceTime                  string
	CurrentDateAndTimeFormatted string
	CreatedAtTimestamp          int32
	Latitude                    string
	LatitudeDirection           string
	Longitude                   string
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
}
