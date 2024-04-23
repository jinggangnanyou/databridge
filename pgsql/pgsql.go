package pgsql

import (
	"github.com/jinggangnanyou/dataorm/database"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const Module = "postgresql"

func Connect(logger logger.Interface, dbConfig *database.DatabaseConfig, gormConfig *database.GormConfig) (*gorm.DB, error) {
	dbConfig.Adapter = Module
	return database.Connect(logger, dbConfig, gormConfig)
}
