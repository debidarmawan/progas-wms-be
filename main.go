package main

import (
	"log"
	"progas-wms-be/config"
	"progas-wms-be/constant"
	_ "progas-wms-be/docs"
	"progas-wms-be/server"
	"strconv"
	"sync"
)

//	@title			PROGAS WMS Backend API
//	@version		1.0
//	@description	This is an API documentation of Progas WMS
//	@contact.name	DeboZero Corp Tech Team
//	@contact.url
//	@contact.email	debidarmawan1998@gmail.com

// @securityDefinitions.apiKey	Bearer
// @in							header
// @name						Authorization
func main() {
	config.Init()

	maxPool, _ := strconv.Atoi(config.GetEnv(constant.DbMaxPool))
	db, sshConn, err := config.ConnectDatabase(maxPool)
	if err != nil {
		log.Fatal(err)
	}

	if sshConn != nil {
		defer func() {
			log.Println("Menutup SSH Tunnel...")
			sshConn.Close()
		}()
	}

	if db != nil {
		defer func() {
			sqlDB, err := db.DB()
			if err == nil {
				log.Println("Menutup koneksi database pool...")
				sqlDB.Close()
			}
		}()
	}

	waitGroup := sync.WaitGroup{}

	waitGroup.Go(func() {
		server.ServeHTTP(db)
	})

	waitGroup.Wait()
}
