package dsn

import (
	"errors"
	"net/url"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	errInvalidDSNUnescaped           = errors.New("invalid DSN: did you forget to escape a param value")
	errInvalidDSNAddr                = errors.New("invalid DSN: network address not terminated (missing closing brace)")
	errInvalidDSNNoSlash             = errors.New("invalid DSN: missing the slash separating the database name")
	DefaultMysqlDSNParser            = &MysqlDSNParser{}
	_                      DSNParser = (*MysqlDSNParser)(nil)
)

type MysqlDSNParser struct {
}

func (m *MysqlDSNParser) GetDialector(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

func (m *MysqlDSNParser) ParseDSN(dsn string) (cfg *DSN, err error) {
	cfg = new(DSN)
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			foundSlash = true
			var j int

			if i > 0 {
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						parseUsernamePassword(cfg, dsn[:j])
						break
					}
				}
				if err = parseAddrNet(cfg, dsn[j:i]); err != nil {
					return
				}
			}
			for j = i + 1; j < len(dsn); j++ {
				if dsn[j] == '?' {
					if err = parseDSNParams(cfg, dsn[j+1:]); err != nil {
						return
					}
					break
				}
			}
			cfg.DBName = dsn[i+1 : j]

			break
		}
	}
	if !foundSlash && len(dsn) > 0 {
		return nil, errInvalidDSNNoSlash
	}
	return
}

func parseUsernamePassword(cfg *DSN, userPassStr string) {
	for i := 0; i < len(userPassStr); i++ {
		if userPassStr[i] == ':' {
			cfg.Password = userPassStr[i+1:]
			cfg.User = userPassStr[:i]
			break
		}
	}
}

func parseAddrNet(cfg *DSN, addrNetStr string) error {
	for i := 0; i < len(addrNetStr); i++ {
		if addrNetStr[i] == '(' {
			if addrNetStr[len(addrNetStr)-1] != ')' {
				if strings.ContainsRune(addrNetStr[i+1:], ')') {
					return errInvalidDSNUnescaped
				}
				return errInvalidDSNAddr
			}
			cfg.Addr = addrNetStr[i+1 : len(addrNetStr)-1]
			cfg.Net = addrNetStr[1:i]
			break
		}
	}
	return nil
}

func parseDSNParams(cfg *DSN, params string) (err error) {
	for _, v := range strings.Split(params, "&") {
		param := strings.SplitN(v, "=", 2)
		if len(param) != 2 {
			continue
		}
		if cfg.Params == nil {
			cfg.Params = make(map[string]string)
		}
		value := param[1]
		if cfg.Params[param[0]], err = url.QueryUnescape(value); err != nil {
			return
		}
	}
	return
}
