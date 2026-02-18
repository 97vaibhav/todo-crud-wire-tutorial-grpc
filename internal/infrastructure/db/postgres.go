package db

import (
	"fmt"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB is a Wire provider. Wire sees it takes *config.Config and returns
// (*gorm.DB, error), so it knows to call config.Load() first, then pass the result here.
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
