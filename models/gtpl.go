// Package models contains the models of specified devices and raw data
package models

import "time"

// The GTPLDevice is a model of a GTPL device.
type GTPLDevice struct {
	Header             string
	DeviceID           string	`bson:"gps_device_id"`
	GPSValidity        string	`bson:"gps_validity"`
	RawDeviceDate      string
	RawDeviceTime      string
	Latitude           float64	`bson:"latitude"`
	LatitudeDirection  string	`bson:"lat_direction"`
	Longitude          float64	`bson:"longitude"`
	LongitudeDirection string	`bson:"longitude_direction"`
	DeviceTime         time.Time	`bson:"device_time"`
	Speed              int		`bson:"speed"`
	GPSOdometer        int		`bson:"gps_odometer"`
	Direction          int		`bson:"cardinal_head"`
	NumberOfSatellites int		`bson:"no_of_satellites"`
	BoxStatus          bool		`bson:"box_status"`
	GSMSignal          int		`bson:"gsm_signal_strength"`
	MainBatteryStatus  string	`bson:"main_bt_status"`
	IgnitionStatus     bool		`bson:"ignition_status"`
	AnalogVoltage      float64	`bson:""`
	InternalBTVoltage  float64	`bson:"internal_bt_voltage"`
	Location           string

	CompanyID 		   string	`bson:"company_id"`
	BuID 			   string	`bson:"bu_id"`
	ProjectID 		   string	`bson:"project_id"`
	DistrictName 	   string 	`bson:"district_name"`
	Branch			   string 	`bson:"branch"`

	// Custom Fields
	Distance    float64 		`bson:"distance,omitempty"`
	CreatedTime time.Time		`bson:"l_created_date"`
	TimeStamp   int64

	//MSSQL fields
	GeoFenceID int64 				`bson:"geo_code"`
	ButtonCode int64				`bson:"button_code"`
	Port 	   int					`bson:"port"`
	StartTime	time.Time 			`bson:"start_time,omitempty"`
	EndTime		time.Time			`bson:"end_time,omitempty"`
	RunTime 	int64				`bson:"run_time,omitempty"`
}
