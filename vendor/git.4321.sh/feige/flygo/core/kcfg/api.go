package kcfg

import (
	"io"

	"github.com/davecgh/go-spew/spew"
)

var (
	ConfigTypeJSON ConfigType = "json"
	ConfigTypeToml ConfigType = "toml"
)

type ConfigType string

type DataSource interface {
	Parse(addr string, watch bool) ConfigType
	ReadConfig() ([]byte, error)
	IsConfigChanged() <-chan struct{}
	io.Closer
}

type Unmarshaller = func([]byte, interface{}) error

var defaultConfiguration = New()

func OnChange(fn func(*Configuration)) {
	defaultConfiguration.OnChange(fn)
}

func LoadFromDataSource(ds DataSource, unmarshaller Unmarshaller, opts ...Option) error {
	return defaultConfiguration.LoadFromDataSource(ds, unmarshaller, opts...)
}

func LoadFromReader(r io.Reader, unmarshaller Unmarshaller) error {
	return defaultConfiguration.LoadFromReader(r, unmarshaller)
}

func Apply(conf map[string]interface{}) error {
	return defaultConfiguration.apply(conf)
}

func Reset() {
	defaultConfiguration = New()
}

func Traverse(sep string) map[string]interface{} {
	return defaultConfiguration.traverse(sep)
}

func RawConfig() []byte {
	return defaultConfiguration.raw()
}

func Debug(sep string) {
	spew.Dump("Debug", Traverse(sep))
}

func Get(key string) interface{} {
	return defaultConfiguration.Get(key)
}

func Set(key string, val interface{}) {
	_ = defaultConfiguration.Set(key, val)
}
