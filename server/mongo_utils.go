package server

import (
	"context"
	"gitlab.com/gpsv2/db"
	"gitlab.com/gpsv2/errcheck"
	"gitlab.com/gpsv2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
)

// insertAIS140DataIntoMongo inserts a AIS140 device document
// into the live Mongo DB. It essentially updates the documents in a
// seperate collection which contains the latest state of the device
func insertAIS140DataIntoMongo(ais140Device *models.AIS140Device) error {
	// the live mongo db collection.
	locationHistoriesCollection, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

	// the updating mongo db collection
	vehicleDetailsCollection, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	// config options for the Find API call.
	options := options2.FindOptions{}
	limit := int64(1)
	options.Limit = &limit

	collectionMutex.Lock()

	// Find returns a cursor when you try to update in
	// the collection by first finding if there are any documents
	// already of the device by filtering using the IMEI number.
	cursor, err := vehicleDetailsCollection.Find(vctx, bson.M{"imeinumber": ais140Device.IMEINumber}, &options)

	errcheck.CheckError(err)

	// If the cursor has data, it means there are
	// already documents of the device. So we only need to update.
	if cursor.Next(vctx) {
		_, err := vehicleDetailsCollection.ReplaceOne(vctx, bson.M{"imeinumber": ais140Device.IMEINumber}, ais140Device)
		errcheck.CheckError(err)

	} else {
		// if the cursor doesn't any documents of the devices
		// it means this will be the first document
		_, err = vehicleDetailsCollection.InsertOne(vctx, ais140Device)
		errcheck.CheckError(err)
	}

	// Now insert in the live database. This doesn't have any conditions.
	_, err = locationHistoriesCollection.InsertOne(locCtx, ais140Device)
	collectionMutex.Unlock()

	return err
}

// insertAIS140HistoryDataMongo inserts history data into the history database
func insertAIS140HistoryDataMongo(ais140device *models.AIS140Device) error {

	historyLHcollection, hctx := db.GetHistoryCollectionsWithContext(historyLHcollection)

	collectionMutex.Lock()
	_, err := historyLHcollection.InsertOne(hctx, ais140device)
	errcheck.CheckError(err)

	collectionMutex.Unlock()

	return err
}

// insertGTPLDataMongo inserts a GTPL device document
// into the live Mongo DB. It essentially updates the documents in a
// seperate collection which contains the latest state of the device
func insertGTPLDataMongo(gtplDevice *models.GTPLDevice) error {
	// the live mongo db collection.
	locationHistoriesCollection, locCtx := db.GetMongoCollectionWithContext(locationHistoriesCollection)

	// the updating mongo db collection
	vehicleDetailsCollection, vctx := db.GetMongoCollectionWithContext(vehicleDetailsCollection)

	// config options for the Find API call.
	options := options2.FindOptions{}
	limit := int64(1)
	options.Limit = &limit

	collectionMutex.Lock()

	// Find returns a cursor when you try to update in
	// the collection by first finding if there are any documents
	// already of the device by filtering using the DeviceID.
	cursor, err := vehicleDetailsCollection.Find(vctx, bson.M{"deviceid": gtplDevice.DeviceID}, &options)

	errcheck.CheckError(err)

	// If the cursor has data, it means there are
	// already documents of the device. So we only need to update.
	if cursor.Next(vctx) {
		_, err := vehicleDetailsCollection.ReplaceOne(vctx, bson.M{"deviceid": gtplDevice.DeviceID}, gtplDevice)
		errcheck.CheckError(err)

	} else {
		// if the cursor doesn't any documents of the devices
		// it means this will be the first document
		_, err = vehicleDetailsCollection.InsertOne(vctx, gtplDevice)
		errcheck.CheckError(err)
	}

	// Now insert in the live database. This doesn't have any conditions.
	_, err = locationHistoriesCollection.InsertOne(locCtx, gtplDevice)
	collectionMutex.Unlock()

	return err
}

// insertGTPLHistoryDataMongo inserts history data into the history database
func insertGTPLHistoryDataMongo(gtplDevice *models.GTPLDevice) error {
	historyLHcollection, hctx := db.GetHistoryCollectionsWithContext(historyLHcollection)

	collectionMutex.Lock()
	_, err := historyLHcollection.InsertOne(hctx, gtplDevice)
	errcheck.CheckError(err)

	collectionMutex.Unlock()

	return err
}

// insertRawDataMongo inserts any raw data given by any device
// into an extra collection in the history database
func insertRawDataMongo(rawData string) error {

	rawDataCollection, rctx := db.GetHistoryCollectionsWithContext(rawDataCollection)

	rd := &models.RawData{
		RawData: rawData,
	}

	collectionMutex.Lock()
	_, err := rawDataCollection.InsertOne(rctx, rd)
	errcheck.CheckError(err)

	collectionMutex.Unlock()

	return err
}

// BulkWrite writes a bulk items into the database
// It is now not functional. To deal with errors a lot
func BulkWrite(devices []models.GTPLDevice) {
	var ctx context.Context

	session, err := db.GetSessionFromClient()

	errcheck.CheckError(err)

	err = mongo.WithSession(ctx, session, func(sctx mongo.SessionContext) error {
		_ = sctx.StartTransaction()

		var operations []mongo.WriteModel

		for _, device := range devices {
			operations = append(operations,
				mongo.NewInsertOneModel().SetDocument(device))
		}

		locationHistoriesCollection, _ := db.GetMongoCollectionWithContext(locationHistoriesCollection)

		_, err := locationHistoriesCollection.BulkWrite(sctx, operations)

		errcheck.CheckError(err)

		_ = session.CommitTransaction(sctx)

		return nil
	})

	session.EndSession(ctx)
}
