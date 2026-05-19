package config

import (
	"log"
	"progas-wms-be/constant"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func ConnectDatabase(maxOpenConn int) *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(GetEnv(constant.DbUrl)),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDb.SetMaxOpenConns(maxOpenConn)

	Migrate(db)
	SeedRBAC(db)

	return db
}
