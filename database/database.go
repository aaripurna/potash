package database

import (
	"context"
	"fmt"
	"time"

	"github.com/aaripurna/potash/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open connects to dsn and verifies it with a ping, so a bad DSN fails at boot
// rather than on the first request that needs the database. It takes the dsn
// rather than reading config directly so a second database is just a second
// call.
func Open(dsn string) (*gorm.DB, error) {
	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logLevel())}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to the database: %w", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("unable to reach the underlying connection pool: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.DBConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}

	return db, nil
}

func logLevel() logger.LogLevel {
	if config.AppEnv == string(config.AppEnvProduction) {
		return logger.Error
	}

	return logger.Info
}
