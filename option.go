package nxconfig

import "github.com/spf13/pflag"

type options struct {
	args      []string
	env       []string
	envprefix string
	flagset   *pflag.FlagSet
}

// Option .
type Option interface {
	apply(*options)
}

type optionFn func(o *options)

func (fn optionFn) apply(o *options) {
	fn(o)
}

// WithArgs overrides the default argument list (os.Args[1:]) used for parsing the flags.
func WithArgs(args []string) Option {
	return optionFn(func(o *options) {
		o.args = args
	})
}

// WithEnv overrides the default environment (os.Environ).
func WithEnv(env []string) Option {
	return optionFn(func(o *options) {
		o.env = env
	})
}

// WithPrefix sets the env prefix.
func WithPrefix(p string) Option {
	return optionFn(func(o *options) {
		o.envprefix = p
	})
}

// WithFlagset overrides the default flagset the flags will be added to.
func WithFlagset(flagset *pflag.FlagSet) Option {
	return optionFn(func(o *options) {
		o.flagset = flagset
	})
}
