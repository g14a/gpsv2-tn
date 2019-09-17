package dbutils

import (
		"bytes"
		"encoding/json"
		"fmt"
		"gitlab.com/gpsv2-tn/errorcheck"
		"gitlab.com/gpsv2-tn/models"
		"net/http"
	)

func InsertGTPLIntoSQL(sqlDevice models.GTPLSQLModel) {

	fixTime := sqlDevice.Fixtime.Format("2006-01-02 15:04:05")
	recvTime := sqlDevice.CreateDate.Format("2006-01-02 15:04:05")

	iq := fmt.Sprintf("INSERT IGNORE INTO location_history_current" +
		" (device_id,gps_validity,created_date,lat_message,north_direction," +
		"lon_message,east_direction,fixTime,speed,odometerreding,cardinal_head," +
		"NumberofSat,box,gsm,battery,ignition_status,anlog,internal_battery," +
		"distance_travelled,pre_lat_message," +
		"pre_lon_message) VALUES ('%d','%s','%v','%f','%s','%f','%s','%v','%s','%s',%f,'%s','%s','%s','%s','%s','%s',%d,%f,%f,%f)",
		sqlDevice.DeviceID, sqlDevice.GPSValidity, recvTime, sqlDevice.LatMessage,
		sqlDevice.NorthDirection, sqlDevice.LonMessage, sqlDevice.EastDirection, fixTime,
		sqlDevice.Speed, sqlDevice.OdometerReading, sqlDevice.CardinalHead, sqlDevice.NumberOfSatellites,
		sqlDevice.Box, sqlDevice.GSM, sqlDevice.MainBatteryStatus, sqlDevice.IgnitionStatus,
		sqlDevice.Analog, sqlDevice.InternalBatteryVoltage, sqlDevice.Distance, sqlDevice.PreLat, sqlDevice.PreLong)

	uq := fmt.Sprintf("UPDATE device_status set gps_validity='%s',time_current='%v',lat_message='%f'," +
							 "north_direction='%s',lon_message='%f',east_direction='%s',fixTime='%v',speed='%s',odometer_reading='%s'," +
							 "cardinal_head=%f,NumberofSat='%s',box='%s',gsm='%s',VehicleBatteryStatus='%s',ignition_status='%s',anlog='%s'," +
							 "internal_battery=%d WHERE GPS_Device_Id='%d'", sqlDevice.GPSValidity, fixTime, sqlDevice.LatMessage,
							 sqlDevice.NorthDirection, sqlDevice.LonMessage, sqlDevice.EastDirection, fixTime, sqlDevice.Speed, sqlDevice.OdometerReading,
							 sqlDevice.CardinalHead, sqlDevice.NumberOfSatellites, sqlDevice.Box, sqlDevice.GSM, sqlDevice.MainBatteryStatus,
							 sqlDevice.IgnitionStatus, sqlDevice.Analog, sqlDevice.InternalBatteryVoltage, sqlDevice.DeviceID)

	insertQuery := &Query{
		Query: iq,
	}

	updateQuery := &Query{
		Query: uq,
	}

	insertJson, _ := json.Marshal(insertQuery)
	updateJson, _ := json.Marshal(updateQuery)

	inreq, _ :=  http.NewRequest("POST", "http://165.22.208.101:3001/avaMysql", bytes.NewReader([]byte(insertJson)))
	upreq, _ :=  http.NewRequest("POST", "http://165.22.208.101:3001/avaMysql", bytes.NewReader([]byte(updateJson)))

	inreq.Header.Add("Content-Type", "application/json")
	upreq.Header.Add("Content-Type", "application/json")

	_, _ = client.Do(inreq)
	_, _ = client.Do(upreq)

	inreq.Close = true
	upreq.Close = true
}

func InsertAIS140IntoSQL(sqlDevice models.AIS140SQLModel) {
	fixTime := sqlDevice.FixTime.Format("2006-01-02 15:04:05")
	recvTime := sqlDevice.CreateDate.Format("2006-01-02 15:04:05")

	iq := fmt.Sprintf("INSERT IGNORE INTO location_history_current (header,vendor_identification,software_version,packate_type," +
							 "packate_identifier,live_history_packet,device_id,vehicle_reg_number," +
							 "gps_fix,fixTime,created_date,lat_message,north_direction,lon_message," +
							 "east_direction,speed,heading,number_of_sat,altitude,pdop,hdop,network," +
		    				 "ignition_status,Vehicle_battery_voltage,vehicle_battery," +
							 "internal_battery_voltage,emergence_status,tamper_alert,gsm,mcc,mnc," +
							 "lac,cel_id,nmr,digital_inputs_status,digital_outputs_status,sequence,checksum,distance_travelled) VALUES(" +
							 "'%s','%s','%s','%s','%s','%s','%s','%s','%s','%v','%v','%f','%s','%f','%s','%d','%d','%s','%s','%s'," +
						     "'%s','%s','%d','%d','%f','%f','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%d')",
						     sqlDevice.Header, sqlDevice.VendorIdentification,sqlDevice.SoftwareVersion,sqlDevice.PacketType,
						     sqlDevice.PacketIdentification,sqlDevice.LiveOrHistoryPacket,sqlDevice.DeviceID,sqlDevice.VehicleRegNumber,
						     sqlDevice.GPSFix,fixTime,recvTime,sqlDevice.LatMessage,sqlDevice.LatitudeDirection,
						     sqlDevice.LonMessage, sqlDevice.LongitudeDirection,sqlDevice.Speed,sqlDevice.CardinalHead,sqlDevice.NumberOfSatellites,
						     sqlDevice.Altitude,sqlDevice.PDOP,sqlDevice.HDOP,sqlDevice.Network,sqlDevice.IgnitionStatus,sqlDevice.VehicleBatteryStatus,
						     sqlDevice.VehicleBatteryVoltage,sqlDevice.InternalBatteryVoltage,sqlDevice.EmergencyStatus,sqlDevice.TamperAlert,
						     sqlDevice.GSM, sqlDevice.MCC,sqlDevice.MNC,sqlDevice.LAC, sqlDevice.CellID,sqlDevice.NMR,sqlDevice.DigitalInputStatus,
						     sqlDevice.DigitalOutputStatus,sqlDevice.SequenceNumber,sqlDevice.CheckSumValue,sqlDevice.DistanceTravelled)

	insertQuery := &Query{
		Query: iq,
	}

	uq := fmt.Sprintf("UPDATE device_status SET packate_type='%s',packate_identifier='%s',live_history_packet='%s',gps_fix='%s'," +
							 "fixTime='%s',lat_message='%f',north_direction='%s',lon_message='%f',east_direction='%s',speed='%d',NumberofSat='%s'," +
							 "pdop='%s',hdop='%s',network='%s',ignition_status='%d',VehicleBatteryStatus='%d',vehicle_battery='%f',internal_battery='%.2f'," +
							 "emergence_status='%s',tamper_alert='%s',gsm='%s',mcc='%s',mnc='%s',lac='%s',cel_id='%s',nmr='%s',digital_inputs_status='%s',digital_outputs_status='%s'," +
		 					 "sequence='%s',checksum='%s', WHERE GPS_Device_Id='%s'",sqlDevice.PacketType, sqlDevice.PacketIdentification,
							 sqlDevice.LiveOrHistoryPacket,sqlDevice.GPSFix,fixTime,sqlDevice.LatMessage,sqlDevice.LatitudeDirection,
							 sqlDevice.LonMessage, sqlDevice.LongitudeDirection, sqlDevice.Speed, sqlDevice.NumberOfSatellites, sqlDevice.PDOP,
							 sqlDevice.HDOP,sqlDevice.Network, sqlDevice.IgnitionStatus,sqlDevice.VehicleBatteryStatus,sqlDevice.VehicleBatteryVoltage,
							 sqlDevice.InternalBatteryVoltage,sqlDevice.EmergencyStatus,sqlDevice.TamperAlert,sqlDevice.GSM,sqlDevice.MCC,sqlDevice.MNC,
							 sqlDevice.LAC,sqlDevice.CellID,sqlDevice.NMR,sqlDevice.DigitalInputStatus,sqlDevice.DigitalOutputStatus,
							 sqlDevice.SequenceNumber,sqlDevice.CheckSumValue, sqlDevice.DeviceID)

	updateQuery := &Query{
		Query: uq,
	}

	insertJson, _ := json.Marshal(insertQuery)
	updateJson, _ := json.Marshal(updateQuery)

	inreq, _ :=  http.NewRequest("POST", "http://165.22.208.101:3001/avaMysql", bytes.NewReader([]byte(insertJson)))
	upreq, _ :=  http.NewRequest("POST", "http://165.22.208.101:3001/avaMysql", bytes.NewReader([]byte(updateJson)))

	inreq.Header.Add("Content-Type", "application/json")
	upreq.Header.Add("Content-Type", "application/json")

	_, _ = client.Do(inreq)
	_, _ = client.Do(upreq)

	inreq.Close = true
	upreq.Close = true
}

func InsertBSTPLIntoSQL(sqlDevice models.BSTPLSQLModel) {

	fixTime := sqlDevice.Fixtime.Format("2006-01-02 15:04:05")
	recvTime := sqlDevice.CreateDate.Format("2006-01-02 15:04:05")

	iq := fmt.Sprintf("INSERT IGNORE INTO location_history_current" +
		" (device_id,gps_validity,created_date,lat_message,north_direction," +
		"lon_message,east_direction,fixTime,speed,odometerreding,cardinal_head," +
		"NumberofSat,box,gsm,battery,ignition_status,anlog,internal_battery," +
		"distance_travelled,pre_lat_message," +
		"pre_lon_message) VALUES ('%d','%s','%v','%f','%s','%f','%s','%v','%s','%s',%f,'%s','%s','%s','%s','%s','%s',%d,%f,%f,%f)",
		sqlDevice.DeviceID, sqlDevice.GPSValidity, recvTime, sqlDevice.LatMessage,
		sqlDevice.NorthDirection, sqlDevice.LonMessage, sqlDevice.EastDirection, fixTime,
		sqlDevice.Speed, sqlDevice.OdometerReading, sqlDevice.CardinalHead, sqlDevice.NumberOfSatellites,
		sqlDevice.Box, sqlDevice.GSM, sqlDevice.MainBatteryStatus, sqlDevice.IgnitionStatus,
		sqlDevice.Analog, sqlDevice.InternalBatteryVoltage, sqlDevice.Distance, sqlDevice.PreLat, sqlDevice.PreLong)

	uq := fmt.Sprintf("UPDATE device_status SET gps_validity='%s',time_current='%v',lat_message='%f'," +
		"north_direction='%s',lon_message='%f',east_direction='%s',fixTime='%v',speed='%s',odometer_reading='%s'," +
		"cardinal_head=%f,NumberofSat='%s',box='%s',gsm='%s',VehicleBatteryStatus='%s',ignition_status='%s',anlog='%s'," +
		"internal_battery=%d WHERE GPS_Device_Id='%d'", sqlDevice.GPSValidity, fixTime, sqlDevice.LatMessage,
		sqlDevice.NorthDirection, sqlDevice.LonMessage, sqlDevice.EastDirection, fixTime, sqlDevice.Speed, sqlDevice.OdometerReading,
		sqlDevice.CardinalHead, sqlDevice.NumberOfSatellites, sqlDevice.Box, sqlDevice.GSM, sqlDevice.MainBatteryStatus,
		sqlDevice.IgnitionStatus, sqlDevice.Analog, sqlDevice.InternalBatteryVoltage, sqlDevice.DeviceID)

	insertQuery := &Query{
		Query: iq,
	}

	updateQuery := &Query{
		Query: uq,
	}

	insertJson, _ := json.Marshal(insertQuery)
	updateJson, _ := json.Marshal(updateQuery)

	inreq, _ :=  http.NewRequest("POST", "http://165.22.208.101:3001/avaMysql", bytes.NewReader([]byte(insertJson)))
	upreq, _ :=  http.NewRequest("POST", "http://165.22.208.101:3001/avaMysql", bytes.NewReader([]byte(updateJson)))

	inreq.Header.Add("Content-Type", "application/json")
	upreq.Header.Add("Content-Type", "application/json")

	_, err := client.Do(inreq)
	_, err = client.Do(upreq)

	inreq.Close = true
	upreq.Close = true

	errorcheck.CheckError(err)
}
