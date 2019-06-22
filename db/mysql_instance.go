package db

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/errorcheck"
	"sync"
)

var (
	sqlDB   *gorm.DB
	sqlOnce sync.Once
)

func GetSQLDB() *gorm.DB {
	sqlOnce.Do(func() {
		connectMySQL()
	})

	return sqlDB
}

func connectMySQL() {

	appConfigInstance := config.GetAppConfig()

	db, err := gorm.Open("mysql", appConfigInstance.MysqlConfig.URL)
	errorcheck.CheckError(err)

	sqlDB = db
}
