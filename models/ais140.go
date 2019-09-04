// Package models contains the models of specified devices and raw data
package models

import "time"

// The AIS140Device is a model of an AIS140 device
// and contains the following fields.

type AIS140Device struct {
	Header                 string
	VendorIdentification   string	`bson:""`
	SoftwareVersion        string
	PacketType             string	`bson:"packate_type"`
	PacketIdentification   string	`bson:"packate_identifier"`
	LiveOrHistoryPacket    string	`bson:"live_history_packet"`
	IMEINumber             string	`bson:"gps_device_id"`
	VehicleRegNumber       string	`bson:"vehicle_reg_number"`
	GPSFix                 string	`bson:"gps_fix"`
	DateInDDMMYYYY         string	`bson:""`
	TimeInHHMMSS           string
	Latitude               float64	`bson:"latitude"`
	LatitudeDirection 	   string	`bson:"lat_direction"`
	Longitude              float64	`bson:"longitude"`
	LongitudeDirection	   string	`bson:"lon_direction"`
	Speed                  int		`bson:"speed"`
	Heading                int		`bson:"heading"`
	NumberOfSatellites     int		`bson:"no_of_satellites"`
	AltitudeInMetres       int		`bson:"altitude"`
	NetworkOperator        string	`bson:"network_operator_name"`
	IgnitionStatus         bool		`bson:"ignition_status"`
	VehicleBatteryStatus   int		`bson:"vehicle_btry_status"`
	VehicleBatteryVoltage  float64	`bson:""`
	InternalBatteryVoltage float64	`bson:"internal_btry_voltage"`
	EmergencyStatus        string	`bson:"emergency_status"`
	TamperAlert            string	`bson:"tamper_alert"`
	GSMStrength            int		`bson:"gsm_signal_strength"`

	DigitalInputStatus     string	`bson:"digital_input_status"`
	DigitalOutputStatus    string	`bson:"digital_output_status"`
	SequenceNumber         string	`bson:"sequence_number_message"`
	CheckSum               string	`bson:"checksum"`

	// Custom fields
	DeviceTime        time.Time		`bson:"device_time"`
	CreatedTime		  time.Time		`bson:"l_created_date"`
	Distance          float64 		`bson:"distance"`
	ButtonCode 		  int64 		`bson:"button_code"`
	GeoFenceID 		  int64			`bson:"geo_code"`
	GPSOdometer 	  int			`bson:"gps_odometer"`
	Port 			  int 			`bson:"port"`
	StartTime 		  time.Time		`bson:"start_time"`
	EndTime			  time.Time		`bson:"end_time"`
	RunTime 		  int 			`bson:"run_time"`
}
