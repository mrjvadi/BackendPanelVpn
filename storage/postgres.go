package storage

import (
	"fmt"
	"github.com/mrjvadi/BackendPanelVpn/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Database struct {
	Db *gorm.DB
}

func ConnectPostgres(cfg config.PostgresConfig) *Database {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)
	dbi, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return &Database{
		Db: dbi,
	}
}
