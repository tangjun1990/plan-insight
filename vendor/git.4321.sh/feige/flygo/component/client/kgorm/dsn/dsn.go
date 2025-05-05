package dsn

import (
	"gorm.io/gorm"
)

type DSN struct {
	User     string
	Password string
	Net      string
	Addr     string
	DBName   string
	Params   map[string]string
}

type DSNParser interface {
	GetDialector(dsn string) gorm.Dialector
	ParseDSN(dsn string) (cfg *DSN, err error)
}
