package database

import (
	"github.com/abdallahelassal/UserAuth.git/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	Cfg *config.Config
	DB  *gorm.DB
	Log *zap.Logger
}
func NewConnection(cfg *config.Config, log *zap.Logger)*Connection{
	
	return &Connection{
		Cfg: cfg,
		Log: log,
	}
}

func (c *Connection) Connect(){
	dsn := "host=" +c.Cfg.DatabaseConfig.Host+
			" user=" +c.Cfg.DatabaseConfig.User+
			" password=" +c.Cfg.DatabaseConfig.Password+
			" dbname=" +c.Cfg.DatabaseConfig.Name+
			" port=" +c.Cfg.DatabaseConfig.Port+
			" sslmode=" +c.Cfg.DatabaseConfig.SSLMode
	c.Log.Info("connection Database",
		zap.String("host",c.Cfg.DatabaseConfig.Host),
		zap.String("user",c.Cfg.DatabaseConfig.User),
		zap.String("port",c.Cfg.DatabaseConfig.Port),
		zap.String("sslmode",c.Cfg.DatabaseConfig.SSLMode),
	)
	
	db ,err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		c.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	c.DB = db
}
package db

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dsn string) (*sql.DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("open db: %w", err)
    }

    // Connection pool tuning — critical for production
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(2 * time.Minute)

    // Verify connection works before returning
    ctx, cancel := context.WithTimeout(
        context.Background(), 5*time.Second,
    )
    defer cancel()

    if err = db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("ping db: %w", err)
    }

    return db, nil
}