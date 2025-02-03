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

var bizDB *gorm.DB
var extraHelperDB *gorm.DB

type DB struct {
	DbType string `yaml:"dbType"`
	Dsn    string `yaml:"dsn"`
	Mark   string `yaml:"mark"`
}

func init() {
	v := conf.GetProjectViper()
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

	var t1 DB
	if err := v.UnmarshalKey("dataSource.db.bizDB", &t1); err != nil {
		panic(err)
	} else {
		var dialect gorm.Dialector
		switch t1.DbType {
		case "postgres":
			dialect = postgres.Open(t1.Dsn)
		case "sqlite":
			dialect = sqlite.Open(t1.Dsn)
		default:
			panic(fmt.Errorf("un-support database type: %s", t1.DbType))
		}
		if bizDB, err = gorm.Open(dialect, &gorm.Config{
			Logger: newLogger,
		}); err != nil {
			panic(err)
		}
	}
	bizDB.AutoMigrate(model.GetBizMigrateTables()...)

	var t2 DB
	if err := v.UnmarshalKey("dataSource.db.logDB", &t2); err != nil {
		panic(err)
	} else {
		var dialect gorm.Dialector
		switch t1.DbType {
		case "postgres":
			dialect = postgres.Open(t2.Dsn)
		case "sqlite":
			dialect = sqlite.Open(t2.Dsn)
		default:
			panic(fmt.Errorf("un-support database type: %s", t2.DbType))
		}
		if extraHelperDB, err = gorm.Open(dialect, &gorm.Config{
			Logger: newLogger,
		}); err != nil {
			panic(err)
		}
	}

	extraHelperDB.AutoMigrate(model.GetExtraHelperMigrateTables()...)
}

func GetBizDB() *gorm.DB {
	return bizDB
}

func GetExtraHelperDB() *gorm.DB {
	return extraHelperDB
}
