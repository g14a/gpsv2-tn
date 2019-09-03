package models

import "time"

type MSSQLDevice struct {
	SimID string  //
	SentTime time.Time //
	RecvTime time.Time //
	Latitude float64 //
	Longitude float64 //
	Speed float64 //
	Orientation int //
	OdometerReading int
	OdoDecimal float64
	Fuel int
	Ignition bool //
	BatteryStatus string //
	ButtonCode string //
	VehicleNumber string //
	InternalBattery int //
	Geofence int64
	LiveHistory string
}