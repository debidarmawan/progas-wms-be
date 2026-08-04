package config

import (
	"context"
	"fmt"
	"log"
	"net"
	"progas-wms-be/constant"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

func ConnectDatabase(maxOpenConn int) (*gorm.DB, *ssh.Client, error) {
	dsn := GetEnv(constant.DbUrl)

	var (
		sshConn *ssh.Client
		err     error
	)

	if strings.EqualFold(GetEnv(constant.GoEnv), "LOCALHOST") {
		sshUser := GetEnv("SSH_USER")
		sshHost := GetEnv("SSH_HOST")
		sshPort := GetEnv("SSH_PORT")

		sshConfig := &ssh.ClientConfig{
			User: sshUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(GetEnv("SSH_PASSWORD")),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         15 * time.Second,
		}

		sshConn, err = ssh.Dial("tcp", net.JoinHostPort(sshHost, sshPort), sshConfig)
		if err != nil {
			log.Fatalf("Gagal SSH dial ke VPS: %v", err)
		}

		mysqlDriver.RegisterDialContext("gorm+ssh", func(ctx context.Context, addr string) (net.Conn, error) {
			return sshConn.Dial("tcp", addr)
		})

		dsn = fmt.Sprintf("%s:%s@gorm+ssh(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", GetEnv("MYSQL_USER"), GetEnv("MYSQL_PASSWORD"), GetEnv("MYSQL_BIND"), GetEnv("MYSQL_PORT"), GetEnv("MYSQL_DATABASE"))
	}

	db, err := gorm.Open(
		mysql.Open(dsn),
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

	if GetEnv("RUN_MIGRATION_AND_SEED") == "true" {
		log.Println("Running migration and seed...")
		Migrate(db)
		SeedRoles(db)
		SeedRBAC(db)
		SeedBootstrapAdmin(db)
		log.Println("Migration and seed completed.")
	} else {
		log.Println("Skipping migration and seed.")
	}

	return db, sshConn, err
}
