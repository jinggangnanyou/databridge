package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormConfig struct {
	LogLevel                 int  `yaml:"log_level"`
	InitEnable               bool `yaml:"init_enable"`
	Debug                    bool `yaml:"debug"`
	IgnoreMigrateErrorEnable bool `yaml:"ignore_migrate_error_enable"`
	MaxIdleConns             int  `yaml:"max_idle_conns"`
	MaxOpenConns             int  `yaml:"max_open_conns"`
	ConnMaxLifetime          int  `yaml:"conn_max_lifetime"`
}

type DatabaseConfig struct {
	Adapter        string `yaml:"adapter"`
	Address        string `yaml:"address"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	Port           int    `yaml:"port"`
	Database       string `yaml:"database"`
	FeatureProbeDB string `yaml:"feature_probe_db" json:"feature_probe_db"`
}

func Connect(logger logger.Interface, dbConfig *DatabaseConfig, gorm *GormConfig) (*gorm.DB, error) {
	return initDB(logger,
		SetDbAddr(dbConfig.Address),
		SetDbUser(dbConfig.Username),
		SetDbPwd(dbConfig.Password),
		SetPort(dbConfig.Port),
		SetDbType(dbConfig.Adapter),
		SetDbName(dbConfig.Database),
		SetConnMaxLifetime(gorm.ConnMaxLifetime),
		SetMaxIdleConns(gorm.MaxIdleConns),
		SetMaxOpenConns(gorm.MaxOpenConns),
		SetGormLogLevel(gorm.LogLevel),
	)
}

// initDB 初始化数据库
func initDB(logger logger.Interface, opt ...DbOption) (*gorm.DB, error) {
	options := makeOptions(opt)
	dialector, gormConfig, err := connectDB(options)
	if err != nil {
		return nil, err
	}
	if logger != nil {
		gormConfig.Logger = logger
	}
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, err
	}
	err = setDBPool(db, options)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 连接数据库 如果要连接其他类型的数据库 这里要适配
// 支持其他模块自己进行实例化，如feature_probe
func connectDB(options DbOptions) (gorm.Dialector, *gorm.Config, error) {
	var dsn string
	var dialector gorm.Dialector
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: options.DisableForeignKeyConstraintWhenMigrating,
		CreateBatchSize:                          options.CreateBatchSize,
	}
	switch options.DbType {
	case "mysql":
		dsn = getMysqlDsn(options)
		dialector = mysql.Open(dsn)
		return dialector, gormConfig, nil
	case "postgres", "postgresql":
		dsn = getPostgresDsn(options)
		dialector = postgres.Open(dsn)
		return dialector, gormConfig, nil
	default:
		return nil, nil, fmt.Errorf(`unsupported database type: %v`, options.DbType)
	}

}

// 设置数据库连接池
func setDBPool(db *gorm.DB, options DbOptions) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(options.MaxIdleConns)                                    // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxOpenConns(options.MaxOpenConns)                                    // SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(options.ConnMaxLifetime)) // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	return nil
}

// 获取mysql数据库dsn
func getMysqlDsn(options DbOptions) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		options.DbUser, options.DbPwd, options.DbAddr, options.Port, options.DbName)
}

// 获取postgres数据库dsn
func getPostgresDsn(options DbOptions) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s  port=%d sslmode=disable TimeZone=Asia/Shanghai",
		options.DbAddr, options.DbUser, options.DbPwd, options.DbName, options.Port)
}
