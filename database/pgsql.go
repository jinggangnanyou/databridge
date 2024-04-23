package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ModulePgSQL = "postgresql"

func ConnectPgSQL(logger logger.Interface, dbConfig *DatabaseConfig, gormConfig *GormConfig) (*gorm.DB, error) {
	dbConfig.Adapter = ModulePgSQL
	return Connect(logger, dbConfig, gormConfig)
}
