package models

import "time"

type BSTPLDevice struct {
	Header string
	VehicleID string				`bson:"gps_device_id"`
	GPSValidity string				`bson:"gps_validity"`
	RawDeviceDate 	    string
	RawDeviceTime 		string
	Latitude float64				`bson:"latitude"`
	LatitudeDirection string		`bson:"lat_direction"`
	Longitude float64				`bson:"longitude"`
	LongitudeDirection string		`bson:"longitude_direction"`
	Speed float64					`bson:"speed"`
	GPSOdometer float64				`bson:"gps_odometer"`
	Direction string				`bson:"cardinal_head"`
	NumberOfSatellites int			`bson:"no_of_satellites"`
	BoxOpenCloseStatus string		`bson:"box_status"`
	GSMSignal int					`bson:"gsm_signal_strength"`
	MainBatteryStatus int			`bson:"main_bt_status"`
	DigitalInputStatus int			`bson:"ignition_status"`	//Ignition status
	AnalogInput float64
	GPSAccuracy int
	InternalBatteryVoltage float64	`bson:"internal_bt_voltage"`
	LiveOrHistoryPacket    string	`bson:"live_history_packet"`

	//Custom Fields
	DeviceTime time.Time			`bson:"device_time"`
	CreatedTime time.Time			`bson:"l_created_date"`
	Distance float64 				`bson:"distance"`

	// based on live history packet
	ButtonCode 	int64				`bson:"button_code"`
	GeoFenceID  int64				`bson:"geo_code"`
	Port 		int 				`bson:"port"`
}