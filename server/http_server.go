package server

import (
	"fmt"
	"log"
	"os"
	"progas-wms-be/constant"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ServeHTTP(db *gorm.DB) {
	f := fiber.New()
	f.Use(RequestId)

	Routes(f, db)

	var port int
	defaultPort := 3131
	if portEnv, ok := os.LookupEnv(constant.Port); ok {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			port = defaultPort
		} else {
			port = portInt
		}
	} else {
		port = defaultPort
	}

	listenerPort := fmt.Sprintf(":%d", port)
	log.Fatal(f.Listen(listenerPort))
}
