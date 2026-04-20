package database

import (
	"context"
	"time"

	"github.com/abdallahelassal/UserAuth/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	Cfg *config.Config
	DB  *gorm.DB
	Log *zap.Logger
}

func NewConnection(cfg *config.Config, log *zap.Logger) *Connection {

	return &Connection{
		Cfg: cfg,
		Log: log,
	}
}

func (c *Connection) Connect() error {
	dsn := "host=" + c.Cfg.DatabaseConfig.Host +
		" user=" + c.Cfg.DatabaseConfig.User +
		" password=" + c.Cfg.DatabaseConfig.Password +
		" dbname=" + c.Cfg.DatabaseConfig.Name +
		" port=" + c.Cfg.DatabaseConfig.Port +
		" sslmode=" + c.Cfg.DatabaseConfig.SSLMode
	c.Log.Info("connection Database",
		zap.String("host", c.Cfg.DatabaseConfig.Host),
		zap.String("user", c.Cfg.DatabaseConfig.User),
		zap.String("port", c.Cfg.DatabaseConfig.Port),
		zap.String("sslmode", c.Cfg.DatabaseConfig.SSLMode),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		c.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	// Verify connection works before returning

	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = sqlDB.PingContext(ctx); err != nil {
		return err
	}
	c.DB = db
	return nil
}
