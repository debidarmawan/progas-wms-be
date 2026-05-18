package main

import (
	"progas-wms-be/config"
	"progas-wms-be/constant"
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
	db := config.ConnectDatabase(maxPool)

	waitGroup := sync.WaitGroup{}

	waitGroup.Go(func() {
		server.ServeHTTP(db)
	})

	waitGroup.Wait()
}
