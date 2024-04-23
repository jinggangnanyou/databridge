package database

type DbOption func(dbOptions *DbOptions)

type DbOptions struct {
	DbType                                   string
	DbAddr                                   string
	DbUser                                   string
	DbPwd                                    string
	DbName                                   string
	Port                                     int
	MaxIdleConns                             int
	MaxOpenConns                             int
	ConnMaxLifetime                          int
	GormLogLevel                             int
	CreateBatchSize                          int
	DisableForeignKeyConstraintWhenMigrating bool
}

func SetDbType(dbType string) DbOption {
	return func(dbOptions *DbOptions) {
		dbOptions.DbType = dbType
	}
}

func SetDbAddr(addr string) DbOption {
	return func(dbOptions *DbOptions) {
		dbOptions.DbAddr = addr
	}
}

// setDbUser
func SetDbUser(user string) DbOption {
	return func(dbOptions *DbOptions) {
		dbOptions.DbUser = user
	}
}

// setDbPwd
func SetDbPwd(pwd string) DbOption {
	return func(dbOptions *DbOptions) {
		dbOptions.DbPwd = pwd

	}
}

// setDbName
func SetDbName(name string) DbOption {
	return func(dbOptions *DbOptions) {
		dbOptions.DbName = name
	}
}

// setPort
func SetPort(port int) DbOption {
	return func(dbOptions *DbOptions) {
		dbOptions.Port = port
	}
}

// setMaxIdleConns
func SetMaxIdleConns(maxIdleConns int) DbOption {
	return func(dbOptions *DbOptions) {
		if maxIdleConns > 0 {
			dbOptions.MaxIdleConns = maxIdleConns
		}
	}
}

// setMaxOpenConns
func SetMaxOpenConns(maxOpenConns int) DbOption {
	return func(dbOptions *DbOptions) {
		if maxOpenConns > 0 {
			dbOptions.MaxOpenConns = maxOpenConns
		}

	}
}

// setConnMaxLifetime
func SetConnMaxLifetime(connMaxLifetime int) DbOption {
	return func(dbOptions *DbOptions) {
		if connMaxLifetime > 0 {
			dbOptions.ConnMaxLifetime = connMaxLifetime
		}
	}
}

// setGormLogLevel
func SetGormLogLevel(gormLogLevel int) DbOption {
	return func(dbOptions *DbOptions) {
		if gormLogLevel > 0 {
			dbOptions.GormLogLevel = gormLogLevel
		}
	}
}

// setCreateBatchSize
func SetCreateBatchSize(createBatchSize int) DbOption {
	return func(dbOptions *DbOptions) {
		if createBatchSize > 0 {
			dbOptions.CreateBatchSize = createBatchSize
		}

	}
}

// makeOptions 生成options
func makeOptions(opt []DbOption) DbOptions {
	options := DbOptions{
		MaxIdleConns:    5,
		MaxOpenConns:    400,
		ConnMaxLifetime: 30,
		CreateBatchSize: 0,
	}
	for _, o := range opt {
		if o != nil {
			o(&options)
		}
	}
	return options
}
