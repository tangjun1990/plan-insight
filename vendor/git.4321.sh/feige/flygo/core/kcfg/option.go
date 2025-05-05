package kcfg

type Option func(o *Container)

func WithTagName(tag ConfigType) Option {
	return func(o *Container) {
		o.TagName = string(tag)
	}
}

func WithWeaklyTypedInput(weaklyTypedInput bool) Option {
	return func(o *Container) {
		o.WeaklyTypedInput = weaklyTypedInput
	}
}
