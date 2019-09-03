package server

import (
	"fmt"
	"gitlab.com/gpsv2-withoutrmtesting/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
	"time"

	"gitlab.com/gpsv2-withoutrmtesting/config"
	"gitlab.com/gpsv2-withoutrmtesting/db"
	"gitlab.com/gpsv2-withoutrmtesting/errorcheck"
	"gitlab.com/gpsv2-withoutrmtesting/models"
	"go.mongodb.org/mongo-driver/bson"
)


var (
	// live database collections
	locationHistoriesCollection = config.GetAppConfig().Mongoconfig.Collections.LocationHistoriesCollection
	vehicleDetailsCollection    = config.GetAppConfig().Mongoconfig.Collections.VehicleDetailsCollection
	FenceDetailsCollection = config.GetAppConfig().Mongoconfig.Collections.FenceDetailsCollection
)

// insertGTPLDataMongo inserts a GTPL device document
// into the live Mongo DB. It essentially updates the documents in a
// seperate collection which contains the latest state of the device.
func insertGTPLDataMongo(gtplDevice *models.GTPLDevice, wg *sync.WaitGroup) {
	defer wg.Done()

	if gtplDevice.DeviceID != "" {

		// the live mongo db collection.
		locationHistoriesCollection, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

		// the updating mongo db collection
		vehicleDetailsCollection, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

		// If the cursor has data, it means there are
		// already documents of the device. So we only need to update.
		if utils.GTPLCheckLiveHistory(gtplDevice.DeviceTime) {

			prevlat, prevlong, odometerValue := GetPrevLatLong(gtplDevice.DeviceID)
			distance := utils.Distance(prevlat, prevlong, gtplDevice.Latitude, gtplDevice.Longitude)

			gtplDevice.Distance = distance
			odometerValue += distance

			o := fmt.Sprintf("%0.0f", odometerValue)

			gtplDevice.GPSOdometer, _ = strconv.Atoi(o)

			_, err := vehicleDetailsCollection.
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
					}})

			errorcheck.CheckError(err)
		}

		// Now insert in the live database. This doesn't have any conditions.
		_, err := locationHistoriesCollection.InsertOne(locCtx, gtplDevice)

		insertFencingGTPL(gtplDevice)
		calculateRunTime(gtplDevice.DeviceID)

		errorcheck.CheckError(err)
	}
}

func insertAIS140DataIntoMongo(ais140Device *models.AIS140Device, wg *sync.WaitGroup) {
	// the live mongo db collection.
	defer wg.Done()
	// the live mongo db collection.
	var distance float64

	locationHistoriesCollection, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

	// the updating mongo db collection
	vehicleDetailsCollection, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	if ais140Device.LiveOrHistoryPacket == "L" {

		prevlat, prevlong, odometerValue := GetPrevLatLong(ais140Device.IMEINumber)
		distance = utils.Distance(prevlat, prevlong, ais140Device.Latitude, ais140Device.Longitude)

		ais140Device.Distance = distance
		odometerValue += distance

		o := fmt.Sprintf("%0.0f", odometerValue)

		ais140Device.GPSOdometer, _ = strconv.Atoi(o)

		_, err := vehicleDetailsCollection.
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
	}

	_, err := locationHistoriesCollection.InsertOne(locCtx, ais140Device)

	insertFencingAIS140(ais140Device)

	calculateRunTime(ais140Device.IMEINumber)

	errorcheck.CheckError(err)
}

func insertBSTPLDataMongo(bstplDevice *models.BSTPLDevice, wg *sync.WaitGroup) {

	defer wg.Done()
	var distance float64

	if bstplDevice.VehicleID != "" {

		// the live mongo db collection.
		locationHistoriesCollection, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

		if bstplDevice.LiveOrHistoryPacket == "L" {

			// the updating mongo db collection
			vehicleDetailsCollection, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

			prevlat, prevlong, odometerValue := GetPrevLatLong(bstplDevice.VehicleID)
			distance = utils.Distance(prevlat, prevlong, bstplDevice.Latitude, bstplDevice.Longitude)

			bstplDevice.Distance = distance
			odometerValue += distance

			bstplDevice.GPSOdometer = odometerValue

			_, err := vehicleDetailsCollection.
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
		}

		_, err := locationHistoriesCollection.InsertOne(locCtx, bstplDevice)

		insertFencingBSTPL(bstplDevice)
		calculateRunTime(bstplDevice.VehicleID)

		errorcheck.CheckError(err)
	}
}

func getVehicleRegNo(deviceID string) string {

	vehicleDetailsCollection, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	cursor, err := vehicleDetailsCollection.Find(ctx, bson.M{"gps_device_id": deviceID}, &options.FindOptions{
		Projection: bson.D {
			{"vehicle_reg_no", 1},
		},
	})

	var device Device

	if cursor.Next(ctx); cursor.Err() == nil {

		err = cursor.Decode(&device)

		if err != nil  {
			log.Println(err)
		}
	}

	return device.VehicleRegNo
}

func GetPrevLatLong(deviceID string) (float64, float64, float64) {

	vehicleDetailsCollection, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	cursor, err := vehicleDetailsCollection.Find(ctx, bson.M{"gps_device_id": deviceID}, &options.FindOptions{
		Projection: bson.D {
			{"latitude", 1},
			{"longitude", 1},
			{"gps_odometer", 1},
		},
	})

	var device latlong

	if cursor.Next(ctx); cursor.Err() == nil {

		err = cursor.Decode(&device)

		if err != nil  {
			log.Println(err)
		}
	}

	return device.Latitude, device.Longitude, device.GPSOdometer
}

func insertRawDataMongo(rawData string, wg *sync.WaitGroup) {
	defer wg.Done()

	rawDataCollection, rctx := db.GetRawCollectionWithContext()

	rd := &models.RawData{
		RawData: rawData,
	}

	_, err := rawDataCollection.InsertOne(rctx, rd)
	errorcheck.CheckError(err)
}

func insertFencingBSTPL(bstplDevice *models.BSTPLDevice) {

	FenceDetailsCollection, ctx := db.GetMongoCollectionWithContext(FenceDetailsCollection)

	// Get the lat, long and radius of the fence with fence id
	cursor, err := FenceDetailsCollection.Find(ctx, bson.M{"gps_device_id": bstplDevice.VehicleID}, &options.FindOptions{
		Projection: bson.D {
			{"latitude", 1},
			{"longitude", 1},
			{"radius", 1},
			{"gps_device_id", 1},
		},
	})

	var deviceFence deviceFence

	if cursor.Next(ctx); cursor.Err() == nil {
		err = cursor.Decode(&deviceFence)

		errorcheck.CheckError(err)
	}

	deviceLat, deviceLong, _ := GetPrevLatLong(bstplDevice.VehicleID)

	bstplDevice.GeoFenceID = 22
	// if it is true, meaning the device reached the fence
	if utils.GeoFence(deviceFence.Latitude, deviceFence.Longitude, deviceLat, deviceLong, float64(deviceFence.Radius)) {
		bstplDevice.GeoFenceID = 33			// to Base
	}
}

func insertFencingGTPL(gtplDevice *models.GTPLDevice) {

	FenceDetailsCollection, ctx := db.GetMongoCollectionWithContext(FenceDetailsCollection)

	// Get the lat, long and radius of the fence with fence id
	cursor, err := FenceDetailsCollection.Find(ctx, bson.M{"gps_device_id": gtplDevice.DeviceID}, &options.FindOptions{
		Projection: bson.D {
			{"latitude", 1},
			{"longitude", 1},
			{"radius", 1},
			{"gps_device_id", 1},
		},
	})

	var deviceFence deviceFence

	if cursor.Next(ctx); cursor.Err() == nil  {
		err = cursor.Decode(&deviceFence)

		if err != nil  {
			log.Println(err)
		}
	}

	deviceLat, deviceLong, _ := GetPrevLatLong(gtplDevice.DeviceID)

	gtplDevice.GeoFenceID = 22
	// if it is true, meaning the device reached the fence
	if utils.GeoFence(deviceFence.Latitude, deviceFence.Longitude, deviceLat, deviceLong, float64(deviceFence.Radius)) {
		gtplDevice.GeoFenceID = 33			// to Base
	}
}

func insertFencingAIS140(ais140Device *models.AIS140Device) {

	FenceDetailsCollection, ctx := db.GetMongoCollectionWithContext(FenceDetailsCollection)

	// Get the lat, long and radius of the fence with fence id
	cursor, err := FenceDetailsCollection.Find(ctx, bson.M{"gps_device_id": ais140Device.IMEINumber}, &options.FindOptions{
		Projection: bson.D {
			{"latitude", 1},
			{"longitude", 1},
			{"radius", 1},
			{"gps_device_id", 1},
		},
	})

	var deviceFence deviceFence

	if cursor.Next(ctx); cursor.Err() == nil {
		err = cursor.Decode(&deviceFence)

		if err != nil  {
			log.Println(err)
		}
	}

	deviceLat, deviceLong, _ := GetPrevLatLong(ais140Device.IMEINumber)

	ais140Device.GeoFenceID = 22

	// if it is true, meaning the device reached the fence
	if utils.GeoFence(deviceFence.Latitude, deviceFence.Longitude, deviceLat, deviceLong, float64(deviceFence.Radius)) {
		ais140Device.GeoFenceID = 33			// to Base
	}
}

func calculateRunTime(deviceID string) {
	vehicleDetailsCollection, ctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	cursor, err := vehicleDetailsCollection.Find(ctx, bson.M{"gps_device_id": deviceID}, &options.FindOptions{
		Projection: bson.D {
			{"device_time", 1},
		},
	})

	var device Dt

	if cursor.Next(ctx) {

		err = cursor.Decode(&device)

		if err != nil  {
			log.Println(err)
		}
	}

	fmt.Println(device.DeviceTime)
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
	DeviceTime time.Time `bson:"device_time"`
}

type Device struct {
	VehicleRegNo string `bson:"vehicle_reg_no"`
}
