package kcfg

type Container struct {
	TagName          string
	WeaklyTypedInput bool
}

var defaultContainer = Container{
	TagName:          "mapstructure",
	WeaklyTypedInput: false,
}
