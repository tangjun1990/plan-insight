package server

import (
	"context"
	"fmt"
	"git.4321.sh/feige/flygo/core"
	"git.4321.sh/feige/flygo/core/kflag"

	"git.4321.sh/feige/flygo/core/kapp"
)

type Option func(c *ServiceInfo)

type ConfigInfo struct {
	Routes []Route
}

type ServiceInfo struct {
	Name       string              `json:"name"`
	Scheme     string              `json:"scheme"`
	Address    string              `json:"address"`
	Weight     float64             `json:"weight"`
	Enable     bool                `json:"enable"`
	Healthy    bool                `json:"healthy"`
	Metadata   map[string]string   `json:"metadata"`
	Region     string              `json:"region"`
	Zone       string              `json:"zone"`
	Kind       kapp.ServiceKind    `json:"kind"`
	Deployment string              `json:"deployment"`
	Group      string              `json:"group"`
	Services   map[string]*Service `json:"services" toml:"services"`
}

type Service struct {
	Namespace string            `json:"namespace" toml:"namespace"`
	Name      string            `json:"name" toml:"name"`
	Labels    map[string]string `json:"labels" toml:"labels"`
	Methods   []string          `json:"methods" toml:"methods"`
}

func (si ServiceInfo) Label() string {
	return fmt.Sprintf("%s://%s", si.Scheme, si.Address)
}

type Server interface {
	core.Component
	GraceShutdown(ctx context.Context) error
	Info() *ServiceInfo
}

type Route struct {
	WeightGroups []WeightGroup
	Method       string
}

type WeightGroup struct {
	Group  string
	Weight int
}

func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}

func WithMetaData(key, value string) Option {
	return func(c *ServiceInfo) {
		c.Metadata[key] = value
	}
}

func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

func WithName(name string) Option {
	return func(c *ServiceInfo) {
		c.Name = name
	}
}

func WithKind(kind kapp.ServiceKind) Option {
	return func(c *ServiceInfo) {
		c.Kind = kind
	}
}

func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{
		Name:       kapp.Name(),
		Weight:     100,
		Enable:     true,
		Healthy:    true,
		Metadata:   make(map[string]string),
		Region:     kapp.AppRegion(),
		Zone:       kapp.AppZone(),
		Kind:       0,
		Deployment: "",
		Group:      "",
	}
	si.Metadata["appMode"] = kapp.AppMode()
	si.Metadata["appHost"] = kflag.String("host")
	si.Metadata["startTime"] = kapp.StartTime()
	si.Metadata["buildTime"] = kapp.BuildTime()
	si.Metadata["appVersion"] = kapp.AppVersion()
	si.Metadata["flygoVersion"] = kapp.FlygoVersion()
	return si
}
