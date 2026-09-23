package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"backend-api-kpmacademy/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	logMode := logger.Silent
	if cfg.App.Debug {
		logMode = logger.Info
	}

	buildDSN := func(network, addr string) string {
		return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DB.User,
			cfg.DB.Password,
			network,
			addr,
			cfg.DB.Name,
		)
	}

	var dsn string
	switch {
	case cfg.DB.Socket != "":
		dsn = buildDSN("unix", cfg.DB.Socket)
	case strings.HasPrefix(cfg.DB.Host, "/"):
		dsn = buildDSN("unix", cfg.DB.Host)
	default:
		dsn = buildDSN("tcp", fmt.Sprintf("%s:%s", cfg.DB.Host, cfg.DB.Port))
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		// If TCP connection failed and host is localhost/127.0.0.1, try common local UNIX sockets
		if cfg.DB.Socket == "" && (cfg.DB.Host == "127.0.0.1" || cfg.DB.Host == "localhost") {
			commonSockets := []string{
				"/var/lib/mysql/mysql.sock",
				"/var/run/mysqld/mysqld.sock",
				"/tmp/mysql.sock",
			}
			for _, sock := range commonSockets {
				if _, statErr := os.Stat(sock); statErr == nil {
					log.Printf("TCP connection failed, falling back to UNIX socket: %s", sock)
					fallbackDSN := buildDSN("unix", sock)
					if fbDB, fbErr := gorm.Open(mysql.Open(fallbackDSN), &gorm.Config{
						Logger: logger.Default.LogMode(logMode),
					}); fbErr == nil {
						db = fbDB
						err = nil
						break
					}
				}
			}
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DB.ConnMaxLifetime) * time.Minute)

	DB = db
	log.Println("Database connected successfully")
}

func GetDB() *gorm.DB {
	return DB
}
