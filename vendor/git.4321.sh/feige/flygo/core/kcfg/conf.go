package kcfg

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
	"time"

	"git.4321.sh/feige/flygo/core/utils/xcast"
	"git.4321.sh/feige/flygo/core/utils/xmap"
	"github.com/mitchellh/mapstructure"
)

const PackageName = "core.kcfg"

type Configuration struct {
	mu        sync.RWMutex
	override  map[string]interface{}
	keyDelim  string
	rawConfig []byte
	keyMap    *sync.Map
	onChanges []func(*Configuration)

	watchers map[string][]func(*Configuration)
}

const (
	defaultKeyDelim = "."
)

func New() *Configuration {
	return &Configuration{
		override:  make(map[string]interface{}),
		keyDelim:  defaultKeyDelim,
		keyMap:    &sync.Map{},
		onChanges: make([]func(*Configuration), 0),
		watchers:  make(map[string][]func(*Configuration)),
	}
}

func (c *Configuration) SetKeyDelim(delim string) {
	c.keyDelim = delim
}

func (c *Configuration) Sub(key string) *Configuration {
	return &Configuration{
		keyDelim: c.keyDelim,
		override: c.GetStringMap(key),
	}
}

func (c *Configuration) Writkcfgig() error {
	return nil
}

func (c *Configuration) OnChange(fn func(*Configuration)) {
	c.onChanges = append(c.onChanges, fn)
}

func (c *Configuration) LoadFromDataSource(ds DataSource, unmarshaller Unmarshaller, opts ...Option) error {
	for _, opt := range opts {
		opt(&defaultContainer)
	}

	content, err := ds.ReadConfig()
	if err != nil {
		return fmt.Errorf("LoadFromDataSource ReadConfig, err: %w", err)
	}

	if err := c.Load(content, unmarshaller); err != nil {
		return fmt.Errorf("LoadFromDataSource Load, err: %w", err)
	}

	go func() {
		for range ds.IsConfigChanged() {
			if content, err := ds.ReadConfig(); err == nil {
				_ = c.Load(content, unmarshaller)
				for _, change := range c.onChanges {
					change(c)
				}
			}
		}
	}()

	return nil
}

func (c *Configuration) Load(content []byte, unmarshal Unmarshaller) error {
	c.rawConfig = content
	configuration := make(map[string]interface{})
	if err := unmarshal(content, &configuration); err != nil {
		return err
	}
	return c.apply(configuration)
}

func (c *Configuration) LoadFromReader(reader io.Reader, unmarshaller Unmarshaller) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return c.Load(content, unmarshaller)
}

func (c *Configuration) apply(conf map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var changes = make(map[string]interface{})

	xmap.MergeStringMap(c.override, conf)
	for k, v := range c.traverse(c.keyDelim) {
		orig, ok := c.keyMap.Load(k)
		if ok && !reflect.DeepEqual(orig, v) {
			changes[k] = v
		}
		c.keyMap.Store(k, v)
	}

	if len(changes) > 0 {
		c.notifyChanges(changes)
	}

	return nil
}

func (c *Configuration) notifyChanges(changes map[string]interface{}) {
	var changedWatchPrefixMap = map[string]struct{}{}

	for watchPrefix := range c.watchers {
		for key := range changes {
			if strings.HasPrefix(key, watchPrefix) {
				changedWatchPrefixMap[watchPrefix] = struct{}{}
			}
		}
	}

	for changedWatchPrefix := range changedWatchPrefixMap {
		for _, handle := range c.watchers[changedWatchPrefix] {
			go handle(c)
		}
	}
}

func (c *Configuration) Set(key string, val interface{}) error {
	paths := strings.Split(key, c.keyDelim)
	lastKey := paths[len(paths)-1]
	m := deepSearch(c.override, paths[:len(paths)-1])
	m[lastKey] = val
	return c.apply(m)
}

func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		m = m3
	}
	return m
}

func (c *Configuration) Get(key string) interface{} {
	return c.find(key)
}

func GetString(key string) string {
	return defaultConfiguration.GetString(key)
}

func (c *Configuration) GetString(key string) string {
	return xcast.ToString(c.Get(key))
}

func GetBool(key string) bool {
	return defaultConfiguration.GetBool(key)
}

func (c *Configuration) GetBool(key string) bool {
	return xcast.ToBool(c.Get(key))
}

func GetInt(key string) int {
	return defaultConfiguration.GetInt(key)
}

func (c *Configuration) GetInt(key string) int {
	return xcast.ToInt(c.Get(key))
}

func GetInt64(key string) int64 {
	return defaultConfiguration.GetInt64(key)
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Configuration) GetInt64(key string) int64 {
	return xcast.ToInt64(c.Get(key))
}

func GetFloat64(key string) float64 {
	return defaultConfiguration.GetFloat64(key)
}

func (c *Configuration) GetFloat64(key string) float64 {
	return xcast.ToFloat64(c.Get(key))
}

func GetTime(key string) time.Time {
	return defaultConfiguration.GetTime(key)
}

func (c *Configuration) GetTime(key string) time.Time {
	return xcast.ToTime(c.Get(key))
}

func GetDuration(key string) time.Duration {
	return defaultConfiguration.GetDuration(key)
}

func (c *Configuration) GetDuration(key string) time.Duration {
	return xcast.ToDuration(c.Get(key))
}

func GetStringSlice(key string) []string {
	return defaultConfiguration.GetStringSlice(key)
}

func (c *Configuration) GetStringSlice(key string) []string {
	return xcast.ToStringSlice(c.Get(key))
}

func GetSlice(key string) []interface{} {
	return defaultConfiguration.GetSlice(key)
}

func (c *Configuration) GetSlice(key string) []interface{} {
	return xcast.ToSlice(c.Get(key))
}

func GetStringMap(key string) map[string]interface{} {
	return defaultConfiguration.GetStringMap(key)
}

func (c *Configuration) GetStringMap(key string) map[string]interface{} {
	return xcast.ToStringMap(c.Get(key))
}

func GetStringMapString(key string) map[string]string {
	return defaultConfiguration.GetStringMapString(key)
}

func (c *Configuration) GetStringMapString(key string) map[string]string {
	return xcast.ToStringMapString(c.Get(key))
}

func (c *Configuration) GetSliceStringMap(key string) []map[string]interface{} {
	return xcast.ToSliceStringMap(c.Get(key))
}

func GetStringMapStringSlice(key string) map[string][]string {
	return defaultConfiguration.GetStringMapStringSlice(key)
}

func (c *Configuration) GetStringMapStringSlice(key string) map[string][]string {
	return xcast.ToStringMapStringSlice(c.Get(key))
}

func UnmarshalWithExpect(key string, expect interface{}) interface{} {
	return defaultConfiguration.UnmarshalWithExpect(key, expect)
}

func (c *Configuration) UnmarshalWithExpect(key string, expect interface{}) interface{} {
	err := c.UnmarshalKey(key, expect)
	if err != nil {
		return expect
	}
	return expect
}

func UnmarshalKey(key string, rawVal interface{}, opts ...Option) error {
	return defaultConfiguration.UnmarshalKey(key, rawVal, opts...)
}

var ErrInvalidKey = errors.New("invalid key, maybe not exist in config")

func (c *Configuration) UnmarshalKey(key string, rawVal interface{}, opts ...Option) error {
	var options = defaultContainer
	for _, opt := range opts {
		opt(&options)
	}

	config := mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		Result:           rawVal,
		TagName:          options.TagName,
		WeaklyTypedInput: options.WeaklyTypedInput,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		c.mu.RLock()
		defer c.mu.RUnlock()
		return decoder.Decode(c.override)
	}

	value := c.Get(key)
	if value == nil {
		return fmt.Errorf(key+",err: %w", ErrInvalidKey)
	}

	return decoder.Decode(value)
}

func (c *Configuration) find(key string) interface{} {
	dd, ok := c.keyMap.Load(key)
	if ok {
		return dd
	}

	paths := strings.Split(key, c.keyDelim)
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := xmap.DeepSearchInMap(c.override, paths[:len(paths)-1]...)
	dd = m[paths[len(paths)-1]]
	c.keyMap.Store(key, dd)
	return dd
}

func lookup(prefix string, target map[string]interface{}, data map[string]interface{}, sep string) {
	for k, v := range target {
		pp := fmt.Sprintf("%s%s%s", prefix, sep, k)
		if prefix == "" {
			pp = k
		}
		if dd, err := xcast.ToStringMapE(v); err == nil {
			lookup(pp, dd, data, sep)
		} else {
			data[pp] = v
		}
	}
}

func (c *Configuration) traverse(sep string) map[string]interface{} {
	data := make(map[string]interface{})
	lookup("", c.override, data, sep)
	return data
}

func (c *Configuration) raw() []byte {
	return c.rawConfig
}
