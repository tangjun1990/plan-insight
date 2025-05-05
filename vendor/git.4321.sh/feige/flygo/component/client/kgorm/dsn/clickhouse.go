package dsn

import (
    "errors"
    "gorm.io/driver/clickhouse"
    "gorm.io/gorm"
    "strings"
)

var (
    errInvalidDSNError         = errors.New("invalid DSN: missing parameter")
    errInvalidDSNNetAddr       = errors.New("invalid DSN: net or addr error")
    DefaultClickhouseDSNParser = &ClickhouseDSNParser{}
)

// clickhouse

type ClickhouseDSNParser struct{}

func (c *ClickhouseDSNParser) GetDialector(dsn string) gorm.Dialector {
    return clickhouse.Open(dsn)
}

// ParseDSN 解析 DSN
func (c *ClickhouseDSNParser) ParseDSN(dsn string) (cfg *DSN, err error) {
    cfg = new(DSN)
    dsnArr := strings.Split(dsn, "?")
    if len(dsnArr) < 2 {
        err = errInvalidDSNError
        return
    }
    // dsnArr[0]
    netAndAddr := strings.Split(dsnArr[0], "://")
    if len(netAndAddr) < 2 {
        err = errInvalidDSNNetAddr
        return
    }
    cfg.Net = netAndAddr[0]
    cfg.Addr = netAndAddr[1]
    cfg.Params = make(map[string]string)
    
    // dsnArr[1] database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20
    params := strings.Split(dsnArr[1], "&")
    for _, param := range params {
        kvArr := strings.Split(param, "=")
        if len(kvArr) == 2 {
            switch kvArr[0] {
            case "database":
                cfg.DBName = kvArr[1]
            case "username":
                cfg.User = kvArr[1]
            case "password":
                cfg.Password = kvArr[1]
            default:
                cfg.Params[kvArr[0]] = kvArr[1]
            }
        }
    }
    
    return
}
