package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunMigrations(databaseURL, migrationPath string, logger *zap.Logger) error {
	m , err := migrate.New(
		fmt.Sprintf("file://%s",migrationPath),
		databaseURL,
	)
	if err != nil {
        return fmt.Errorf("failed to create migrate instance: %w", err)
    }
    defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err,migrate.ErrNoChange){
		return err
	}

	version , dirty , err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNoChange){
		return err
	}
	
	if dirty {
			logger.Warn("Database is dirty",
			zap.Uint("version", version),
		)
	} else {
			logger.Info("Database migrated successfully",
			zap.Uint("version", version),
		)
	}
	
	return nil
}