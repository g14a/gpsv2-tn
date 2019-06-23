package server

import (
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/db"
	"gitlab.com/gpsv2/models"
	"sync"
)

func insertGTPLIntoSQL(gtplDevice models.GTPLDevice, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	appConfigInstance := config.GetAppConfig()

	sqlDB := db.GetSQLDB()

	gtplTable := appConfigInstance.MysqlConfig.SQLTables.GTPLDevices

	sqlDB.Table(gtplTable).Create(&gtplDevice)
}

func insertAIS140IntoSQL(ais140Device models.AIS140Device, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	appConfigInstance := config.GetAppConfig()

	sqlDB := db.GetSQLDB()

	ais140Table := appConfigInstance.MysqlConfig.SQLTables.AIS140Devices

	sqlDB.Table(ais140Table).Create(ais140Device)
}

func insertRawDataSQL(rawData string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	appConfigInstance := config.GetAppConfig()

	sqlDB := db.GetSQLDB()

	rawDataTable := appConfigInstance.MysqlConfig.SQLTables.RawData

	sqlDB.Table(rawDataTable).Create(&rawData)
}
