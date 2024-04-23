package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ModuleMySQL = "mysql"

func ConnectMySQL(logger logger.Interface, dbConfig *DatabaseConfig, gormConfig *GormConfig) (*gorm.DB, error) {
	dbConfig.Adapter = ModuleMySQL
	return Connect(logger, dbConfig, gormConfig)
}
