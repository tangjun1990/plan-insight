package manager

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"

	"git.4321.sh/feige/flygo/core/kcfg"
	"github.com/BurntSushi/toml"
)

var (
	ErrInvalidDataSource     = errors.New("invalid data source, please make sure the scheme has been registered")
	ErrInvalidUnmarshaller   = errors.New("invalid unmarshaller, please make sure the config type is right")
	ErrDefaultConfigNotExist = errors.New("default config not exist")
	registry                 map[string]kcfg.DataSource
	DefaultScheme            = "file"
	unmarshallerMap          = map[kcfg.ConfigType]kcfg.Unmarshaller{
		kcfg.ConfigTypeJSON: json.Unmarshal,
		kcfg.ConfigTypeToml: toml.Unmarshal,
	}
)

type DataSourceCreatorFunc func() kcfg.DataSource

func init() {
	registry = make(map[string]kcfg.DataSource)
}

func Register(scheme string, creator kcfg.DataSource) {
	registry[scheme] = creator
}

func NewDataSource(configAddr string, watch bool) (kcfg.DataSource, kcfg.Unmarshaller, kcfg.ConfigType, error) {
	scheme := DefaultScheme
	urlObj, err := url.Parse(configAddr)
	if err == nil && len(urlObj.Scheme) > 1 {
		scheme = urlObj.Scheme
	}

	if scheme == DefaultScheme {
		_, err := os.Stat(configAddr)
		if err != nil {
			return nil, nil, "", ErrDefaultConfigNotExist
		}
	}

	creatorFunc, exist := registry[scheme]
	if !exist {
		return nil, nil, "", ErrInvalidDataSource
	}
	tag := creatorFunc.Parse(configAddr, watch)

	parser, flag := unmarshallerMap[tag]
	if !flag {
		return nil, nil, "", ErrInvalidUnmarshaller
	}
	return creatorFunc, parser, tag, nil
}
