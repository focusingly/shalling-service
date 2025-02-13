package db

import (
	"fmt"
	"log"
	"os"
	"space-api/conf"
	"space-domain/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _bizDB *gorm.DB
var _extraDB *gorm.DB

func init() {
	cf := conf.ProjectConf

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second * 0, // Slow SQL threshold
			LogLevel:                  logger.Info,     // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,           // Don't include params in the SQL log
			Colorful:                  true,            // Disable color
		},
	)

	var dialect1 gorm.Dialector
	bizConf := cf.GetBizDBConf()
	switch bizConf.DBType {
	case "postgres":
		dialect1 = postgres.Open(bizConf.Dsn)
	case "sqlite":
		dialect1 = sqlite.Open(bizConf.Dsn)
	default:
		panic(fmt.Errorf("un-support database type: %s", bizConf.DBType))
	}

	if db, err := gorm.Open(dialect1, &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		panic(err)
	} else {
		_bizDB = db
		db.AutoMigrate(model.GetBizMigrateTables()...)
	}

	var dialect gorm.Dialector
	extraConf := cf.GetExtraDBConf()
	switch extraConf.DBType {
	case "postgres":
		dialect = postgres.Open(extraConf.Dsn)
	case "sqlite", "sqlite3":
		dialect = sqlite.Open(extraConf.Dsn)
	default:
		panic(fmt.Errorf("un-support database type: %s", extraConf.DBType))
	}
	if db, err := gorm.Open(dialect, &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		panic(err)
	} else {
		_extraDB = db
		_extraDB.AutoMigrate(model.GetExtraHelperMigrateTables()...)
	}

}

func GetBizDB() *gorm.DB {
	return _bizDB
}

func GetExtraDB() *gorm.DB {
	return _extraDB
}
