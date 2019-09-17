// Packge dbutils has the required functions to update, insert and manipulate anything
// related to the respected DB.
package dbutils

import (
	"fmt"
	"gitlab.com/gpsv2-tn/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"time"

	"gitlab.com/gpsv2-tn/config"
	"gitlab.com/gpsv2-tn/db"
	"gitlab.com/gpsv2-tn/errorcheck"
	"gitlab.com/gpsv2-tn/models"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	// live database collections
	locationHistoriesCollection= config.GetAppConfig().Mongoconfig.Collections.LocationHistoriesCollection
	vehicleDetailsCollection   = config.GetAppConfig().Mongoconfig.Collections.VehicleDetailsCollection
	FenceDetailsCollection     = config.GetAppConfig().Mongoconfig.Collections.FenceDetailsCollection
)

// InsertGTPLDataMongo inserts a GTPL device document
// into the live Mongo DB. It essentially updates the documents in a
// seperate collection which contains the latest state of the device.
func InsertGTPLDataMongo(gtplDevice *models.GTPLDevice, startTime time.Time) {
	if gtplDevice.DeviceID != "" {

		// the live mongo db collection.
		locationHistoriesColl, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

		// the updating mongo db collection
		vehicleDetailsColl, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

		// If the cursor has data, it means there are
		// already documents of the device. So we only need to update.
		if utils.GTPLCheckLiveHistory(gtplDevice.DeviceTime) {

			prevlat, prevlong, odometerValue := GetPrevLatLong(gtplDevice.DeviceID)
			distance := utils.Distance(prevlat, prevlong, gtplDevice.Latitude, gtplDevice.Longitude)

			gtplDevice.Distance = distance
			odometerValue += distance

			o := fmt.Sprintf("%0.0f", odometerValue)

			gtplDevice.GPSOdometer, _ = strconv.Atoi(o)

			_, err := vehicleDetailsColl.
				UpdateOne(vctx, bson.M{"gps_device_id": gtplDevice.DeviceID},

					bson.M{"$set": bson.M{"ignition_status": gtplDevice.IgnitionStatus,
						"latitude":            gtplDevice.Latitude,
						"longitude":           gtplDevice.Longitude,
						"cardinal_head":       gtplDevice.Direction,
						"gsm_signal_strength": gtplDevice.GSMSignal,
						"internal_bt_voltage": gtplDevice.InternalBTVoltage,
						"lat_direction":       gtplDevice.LatitudeDirection,
						"longitude_direction": gtplDevice.LongitudeDirection,
						"no_of_satellites":    gtplDevice.NumberOfSatellites,
						"speed":               gtplDevice.Speed,
						"gps_odometer": 	   odometerValue,
						"device_time":         gtplDevice.DeviceTime,
						"garage"	 :         false,
						"port"		 :         gtplDevice.Port,
						"isrunning"  :         true,
					}})

			errorcheck.CheckError(err)

			ignitionStatus := getRunningStatus(gtplDevice.DeviceID)

			if ignitionStatus == false {
				if gtplDevice.IgnitionStatus == true {

					_, err := vehicleDetailsColl.UpdateOne(vctx, bson.M{"gps_device_id": gtplDevice.DeviceID},
						bson.M{"$set": bson.M{
							"isrunning": true,
							"start_time": startTime,
						}})

					errorcheck.CheckError(err)
				}
			} else if getRunningStatus(gtplDevice.DeviceID) == true {
				if gtplDevice.IgnitionStatus == false {

					_, err := vehicleDetailsColl.UpdateOne(vctx, bson.M{"gps_device_id": gtplDevice.DeviceID},
						bson.M{"$set": bson.M{
							"isrunning": false,
						}})

					errorcheck.CheckError(err)

					gtplDevice.StartTime = startTime
					gtplDevice.EndTime = time.Now()
					runTime := gtplDevice.EndTime.Unix() - gtplDevice.StartTime.Unix()
					gtplDevice.RunTime = runTime
				}
			}
		}

		insertFencingGTPL(gtplDevice)

		// Now insert in the live database. This doesn't have any conditions.
		_, err := locationHistoriesColl.InsertOne(locCtx, gtplDevice)

		errorcheck.CheckError(err)
	}
}

func InsertAIS140DataIntoMongo(ais140Device *models.AIS140Device, startTime time.Time) {
	// the live mongo db collection.
	var distance float64

	locationHistoriesColl, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

	// the updating mongo db collection
	vehicleDetailsColl, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	if ais140Device.LiveOrHistoryPacket == "L" {

		prevlat, prevlong, odometerValue := GetPrevLatLong(ais140Device.IMEINumber)
		distance = utils.Distance(prevlat, prevlong, ais140Device.Latitude, ais140Device.Longitude)

		ais140Device.Distance = distance
		odometerValue += distance

		o := fmt.Sprintf("%0.0f", odometerValue)

		ais140Device.GPSOdometer, _ = strconv.Atoi(o)

		_, err := vehicleDetailsColl.
			UpdateOne(vctx, bson.M{"gps_device_id": ais140Device.IMEINumber},

				bson.M{"$set": bson.M{"ignition_status": ais140Device.IgnitionStatus,
					"latitude":            ais140Device.Latitude,
					"longitude":           ais140Device.Longitude,
					"cardinal_head":       ais140Device.Heading,
					"gsm_signal_strength": ais140Device.GSMStrength,
					"internal_bt_voltage": ais140Device.InternalBatteryVoltage,
					"no_of_satellites":    ais140Device.NumberOfSatellites,
					"speed":               ais140Device.Speed,
					"gps_odometer": 	   odometerValue,
					"device_time":         ais140Device.DeviceTime,
					"garage": 			   false,
					"port":				   ais140Device.Port,
				}})

		errorcheck.CheckError(err)

		ignitionStatus := getRunningStatus(ais140Device.IMEINumber)

		if ignitionStatus == false {
			if ais140Device.IgnitionStatus == true {

				_, err := vehicleDetailsColl.UpdateOne(vctx, bson.M{"gps_device_id": ais140Device.IMEINumber},
					bson.M{"$set": bson.M{
						"start_time": startTime,
					}})

				errorcheck.CheckError(err)
			}
		} else if ignitionStatus == true {
			if ais140Device.IgnitionStatus == false {

				errorcheck.CheckError(err)

				ais140Device.StartTime = startTime
				ais140Device.EndTime = time.Now()
				runTime := ais140Device.EndTime.Unix() - startTime.Unix()
				ais140Device.RunTime = runTime
			}
		}
	}

	insertFencingAIS140(ais140Device)

	_, err := locationHistoriesColl.InsertOne(locCtx, ais140Device)

	errorcheck.CheckError(err)
}

func InsertBSTPLDataMongo(bstplDevice *models.BSTPLDevice, startTime time.Time) {

	var distance float64

	if bstplDevice.VehicleID != "" {

		// the live mongo db collection.
		locationHistoriesColl, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

		if bstplDevice.LiveOrHistoryPacket == "L" {

			// the updating mongo db collection
			vehicleDetailsColl, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

			prevlat, prevlong, odometerValue := GetPrevLatLong(bstplDevice.VehicleID)
			distance = utils.Distance(prevlat, prevlong, bstplDevice.Latitude, bstplDevice.Longitude)

			bstplDevice.Distance = distance
			odometerValue += distance

			bstplDevice.GPSOdometer = odometerValue

			_, err := vehicleDetailsColl.
				UpdateOne(vctx, bson.M{"gps_device_id": bstplDevice.VehicleID},

					bson.M{"$set": bson.M{
						"device_time":         bstplDevice.DeviceTime,
						"ignition_status": 	   bstplDevice.DigitalInputStatus,
						"latitude":            bstplDevice.Latitude,
						"longitude":           bstplDevice.Longitude,
						"cardinal_head":       bstplDevice.Direction,
						"gsm_signal_strength": bstplDevice.GSMSignal,
						"internal_bt_voltage": bstplDevice.InternalBatteryVoltage,
						"lat_direction":       bstplDevice.LatitudeDirection,
						"longitude_direction": bstplDevice.LongitudeDirection,
						"no_of_satellites":    bstplDevice.NumberOfSatellites,
						"speed":               bstplDevice.Speed,
						"gps_odometer": 	   odometerValue,
						"garage" : 			   false,
						"port"	:			   bstplDevice.Port,
					}})

			errorcheck.CheckError(err)

			ignitionStatus := getRunningStatus(bstplDevice.VehicleID)

			if ignitionStatus == false {
				if bstplDevice.DigitalInputStatus == true {

					_, err := vehicleDetailsColl.UpdateOne(vctx, bson.M{"gps_device_id":bstplDevice.VehicleID},
						bson.M{"$set": bson.M{
							"start_time": startTime,
						}})
					errorcheck.CheckError(err)
				}
			} else if ignitionStatus == true {
				if bstplDevice.DigitalInputStatus == false {

					errorcheck.CheckError(err)

					bstplDevice.StartTime = startTime
					bstplDevice.EndTime = time.Now()
					runTime := bstplDevice.EndTime.Unix() - startTime.Unix()
					bstplDevice.RunTime = runTime
				}
			}
		}

		insertFencingBSTPL(bstplDevice)

		_, err := locationHistoriesColl.InsertOne(locCtx, bstplDevice)

		errorcheck.CheckError(err)
	}
}

func getVehicleRegNo(deviceID string) string {

	vehicleDetailsColl, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	cursor, err := vehicleDetailsColl.Find(ctx, bson.M{"gps_device_id": deviceID}, &options.FindOptions{
		Projection: bson.D {
			{"vehicle_reg_no", 1},
		},
	})

	var device Device

	if cursor != nil {
		for cursor.Next(ctx) {

			err = cursor.Decode(&device)

			if err != nil  {
				log.Println(err)
			}
		}
	}

	return device.VehicleRegNo

}

func GetPrevLatLong(deviceID string) (float64, float64, float64) {

	vehicleDetailsColl, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	cursor, err := vehicleDetailsColl.Find(ctx, bson.M{"gps_device_id": deviceID}, &options.FindOptions{
		Projection: bson.D {
			{"latitude", 1},
			{"longitude", 1},
			{"gps_odometer", 1},
		},
	})

	var device latlong

	if cursor != nil {
		for cursor.Next(ctx) {

			err = cursor.Decode(&device)

			if err != nil  {
				log.Println(err)
			}
		}
	}

	return device.Latitude, device.Longitude, device.GPSOdometer
}

func InsertRawDataMongo(rawData string) {

	rawDataCollection, rctx := db.GetRawCollectionWithContext()

	rd := &models.RawData{
		RawData: rawData,
	}

	_, err := rawDataCollection.InsertOne(rctx, rd)
	errorcheck.CheckError(err)
}

func insertFencingBSTPL(bstplDevice *models.BSTPLDevice) {

	if bstplDevice != nil {

		FenceDetailsColl, ctx := db.GetMongoCollectionWithContext(FenceDetailsCollection)

		// Get the lat, long and radius of the fence with fence id
		cursor, err := FenceDetailsColl.Find(ctx, bson.M{"gps_device_id": bstplDevice.VehicleID}, &options.FindOptions{
			Projection: bson.D{
				{"latitude", 1},
				{"longitude", 1},
				{"radius", 1},
				{"gps_device_id", 1},
			},
		})

		var deviceFence deviceFence

		if cursor != nil {
			for cursor.Next(ctx) {
				err = cursor.Decode(&deviceFence)

				errorcheck.CheckError(err)
			}

			deviceLat, deviceLong, _ := GetPrevLatLong(bstplDevice.VehicleID)

			bstplDevice.GeoFenceID = 22
			// if it is true, meaning the device reached the fence
			if utils.GeoFence(deviceFence.Latitude, deviceFence.Longitude, deviceLat, deviceLong, float64(deviceFence.Radius)) {
				bstplDevice.GeoFenceID = 33 // to Base
			}
		}
	}
}

func insertFencingGTPL(gtplDevice *models.GTPLDevice) {

	if gtplDevice != nil {
		FenceDetailsColl, ctx := db.GetMongoCollectionWithContext(FenceDetailsCollection)

		// Get the lat, long and radius of the fence with fence id
		cursor, err := FenceDetailsColl.Find(ctx, bson.M{"gps_device_id": gtplDevice.DeviceID}, &options.FindOptions{
			Projection: bson.D{
				{"latitude", 1},
				{"longitude", 1},
				{"radius", 1},
				{"gps_device_id", 1},
			},
		})

		var deviceFence deviceFence

		if cursor != nil {
			for cursor.Next(ctx) {
				err = cursor.Decode(&deviceFence)

				if err != nil {
					log.Println(err)
				}
			}

			deviceLat, deviceLong, _ := GetPrevLatLong(gtplDevice.DeviceID)

			gtplDevice.GeoFenceID = 22
			// if it is true, meaning the device reached the fence
			if utils.GeoFence(deviceFence.Latitude, deviceFence.Longitude, deviceLat, deviceLong, float64(deviceFence.Radius)) {
				gtplDevice.GeoFenceID = 33 // to Base
			}

		}
	}
}

func insertFencingAIS140(ais140Device *models.AIS140Device) {

	if ais140Device != nil {

		FenceDetailsColl, ctx := db.GetMongoCollectionWithContext(FenceDetailsCollection)

		// Get the lat, long and radius of the fence with fence id
		cursor, err := FenceDetailsColl.Find(ctx, bson.M{"gps_device_id": ais140Device.IMEINumber}, &options.FindOptions{
			Projection: bson.D {
				{"latitude", 1},
				{"longitude", 1},
				{"radius", 1},
				{"gps_device_id", 1},
			},
		})

		var deviceFence deviceFence

		if cursor != nil {
			for cursor.Next(ctx) {
				err = cursor.Decode(&deviceFence)

				if err != nil {
					log.Println(err)
				}
			}

			deviceLat, deviceLong, _ := GetPrevLatLong(ais140Device.IMEINumber)

			ais140Device.GeoFenceID = 22

			// if it is true, meaning the device reached the fence
			if utils.GeoFence(deviceFence.Latitude, deviceFence.Longitude, deviceLat, deviceLong, float64(deviceFence.Radius)) {
				ais140Device.GeoFenceID = 33 // to Base
			}
		}
	}
}

func getRunningStatus(deviceID string) bool {

	vehicleDetailsColl, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	//get previous ignition state "isrunning"
	singleResult := vehicleDetailsColl.FindOne(ctx, bson.M{"gps_device_id": deviceID}, &options.FindOneOptions{
		Projection: bson.D{
			{"ignition_status", 1},
			{"start_time", 1},
		},
	})

	var device Dt

	if singleResult != nil {
		singleResult.Decode(&device)
	}

	return device.IgnitionStatus
}

type deviceFence struct {
	Latitude  float64  `bson:"latitude"`
	Longitude float64  `bson:"longitude"`
	Radius    int      `bson:"radius"`
	DeviceID  string   `bson:"gps_device_id"`
}

type latlong struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	GPSOdometer float64 `bson:"gps_odometer"`
}

type Dt struct {
	IgnitionStatus bool  `bson:"ignition_status"`
}

type Device struct {
	VehicleRegNo string `bson:"vehicle_reg_no"`
}
