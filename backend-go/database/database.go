package database

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"backend-go/config"
	"backend-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg config.Config) (*gorm.DB, error) {
	driver := strings.ToLower(cfg.DBDriver)

	var dialector gorm.Dialector
	switch driver {
	case "postgres", "postgresql", "pg":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSL, cfg.TimeZone,
		)
		dialector = postgres.Open(dsn)
	case "mysql":
		tz := cfg.TimeZone
		if tz == "" {
			tz = "UTC"
		}

		var userInfo string
		if cfg.DBPass != "" {
			userInfo = url.UserPassword(cfg.DBUser, cfg.DBPass).String()
		} else {
			userInfo = url.User(cfg.DBUser).String()
		}

		dsn := fmt.Sprintf(
			"%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=%s",
			userInfo, cfg.DBHost, cfg.DBPort, cfg.DBName, url.QueryEscape(tz),
		)
		dialector = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("driver de banco não suportado: %s", cfg.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no banco: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("erro ao recuperar sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Content{},
		&models.Category{},
		&models.UserInteraction{},
		&models.Recommendation{},
		&models.ContentCategory{},
	)
}
