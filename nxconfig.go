package nxconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

// ErrTargetNotAPointer is returned if the target passed to Load is not pointer.
var ErrTargetNotAPointer = errors.New("target must be pointer")

// ErrTargetNotAStruct is returned if the target passed to Load* is not pointer to a struct type.
var ErrTargetNotAStruct = errors.New("target must be struct pointer")

var durationType = reflect.ValueOf(time.Duration(0)).Type()

const tagkey = "nxconfig"

// Load config values automatically from os.Environ() (no prefix) and os.Args[1:].
func Load(target interface{}, opts ...Option) error {
	options := options{
		args: os.Args[1:],
		env:  os.Environ(),
	}

	for _, opt := range opts {
		opt.apply(&options)
	}

	return load(&options, target)
}

func load(opts *options, target interface{}) error {
	envmap := make(map[string]string, len(opts.env))
	envprefix := opts.envprefix + "_"
	for _, e := range opts.env {
		v := strings.Split(e, "=")
		if opts.envprefix != "" {
			v[0] = strings.Replace(v[0], envprefix, "", 1)
		}
		envmap[v[0]] = v[1]
	}

	val := reflect.ValueOf(target)

	var flagset = opts.flagset
	if flagset == nil {
		flagset = pflag.NewFlagSet(val.Type().Name(), pflag.ContinueOnError)
		flagset.ParseErrorsWhitelist.UnknownFlags = true
	}

	err := loadIntoStruct(flagset, envmap, "", val)
	if err != nil {
		return err
	}

	return flagset.Parse(opts.args)
}

func loadIntoStruct(flagset *pflag.FlagSet, envmap map[string]string, fieldPrefix string, target reflect.Value) error {
	elem := target.Elem()

	if target.Kind() != reflect.Ptr {
		return ErrTargetNotAPointer
	}

	if elem.Kind() != reflect.Struct {
		return ErrTargetNotAStruct
	}

	numfields := elem.NumField()
	elemType := elem.Type()
	for i := 0; i < numfields; i++ {
		field := elem.Field(i)

		if !field.CanSet() {
			continue
		}

		name := fieldPrefix + elemType.Field(i).Name

		tag, ok := elemType.Field(i).Tag.Lookup(tagkey)
		if ok {
			name = tag
		}

		if field.Kind() == reflect.Ptr || field.Kind() == reflect.Struct {
			err := loadIntoStruct(flagset, envmap, name, field.Addr())
			if err != nil {
				return err
			}
			continue
		}

		err := flagForValue(flagset, envmap, field, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func flagForValue(flagset *pflag.FlagSet, envmap map[string]string, val reflect.Value, name string) error {
	flagName := toKebabCase(name)

	defValue, err := envValueForValue(envmap, val, name)
	if err != nil {
		return err
	}

	switch val.Type() {
	case durationType:
		flagset.DurationVar(val.Addr().Interface().(*time.Duration), flagName, defValue.(time.Duration), name)
		return nil
	}

	switch val.Type().Kind() {
	case reflect.String:
		flagset.StringVar(val.Addr().Interface().(*string), flagName, defValue.(string), name)
	case reflect.Int:
		flagset.IntVar(val.Addr().Interface().(*int), flagName, defValue.(int), name)
	case reflect.Int32:
		flagset.Int32Var(val.Addr().Interface().(*int32), flagName, defValue.(int32), name)
	case reflect.Int64:
		flagset.Int64Var(val.Addr().Interface().(*int64), flagName, defValue.(int64), name)
	case reflect.Uint:
		flagset.UintVar(val.Addr().Interface().(*uint), flagName, defValue.(uint), name)
	case reflect.Uint16:
		flagset.Uint16Var(val.Addr().Interface().(*uint16), flagName, defValue.(uint16), name)
	case reflect.Uint32:
		flagset.Uint32Var(val.Addr().Interface().(*uint32), flagName, defValue.(uint32), name)
	case reflect.Uint64:
		flagset.Uint64Var(val.Addr().Interface().(*uint64), flagName, defValue.(uint64), name)
	case reflect.Float32:
		flagset.Float32Var(val.Addr().Interface().(*float32), flagName, defValue.(float32), name)
	case reflect.Float64:
		flagset.Float64Var(val.Addr().Interface().(*float64), flagName, defValue.(float64), name)
	case reflect.Bool:
		flagset.BoolVar(val.Addr().Interface().(*bool), flagName, defValue.(bool), name)
	}

	return nil
}

func envValueForValue(env map[string]string, val reflect.Value, name string) (interface{}, error) {
	envname := toUpperSnakeCase(name)
	typ := val.Type()
	return convertFromStr(env[envname], typ)
}

func convertFromStr(in string, typ reflect.Type) (interface{}, error) {
	switch typ {
	case durationType:
		if in == "" {
			return time.Duration(0), nil
		}
		return time.ParseDuration(in)
	}

	switch typ.Kind() {
	case reflect.String:
		return in, nil
	case reflect.Int:
		if in == "" {
			return 0, nil
		}
		return strconv.Atoi(in)
	case reflect.Int32:
		if in == "" {
			return int32(0), nil
		}
		return strconv.ParseInt(in, 10, 32)
	case reflect.Int64:
		if in == "" {
			return int64(0), nil
		}
		return strconv.ParseInt(in, 10, 64)
	case reflect.Uint:
		if in == "" {
			return uint(0), nil
		}
		v, err := strconv.ParseUint(in, 10, 64)
		return uint(v), err
	case reflect.Uint16:
		if in == "" {
			return uint16(0), nil
		}
		v, err := strconv.ParseUint(in, 10, 16)
		return uint16(v), err
	case reflect.Uint32:
		if in == "" {
			return uint32(0), nil
		}
		return strconv.ParseUint(in, 10, 32)
	case reflect.Uint64:
		if in == "" {
			return uint64(0), nil
		}
		return strconv.ParseUint(in, 10, 64)
	case reflect.Float32:
		if in == "" {
			return float32(0), nil
		}
		v, err := strconv.ParseFloat(in, 32)
		return float32(v), err
	case reflect.Float64:
		if in == "" {
			return float64(0), nil
		}
		return strconv.ParseFloat(in, 64)
	case reflect.Bool:
		if in == "" {
			return false, nil
		}
		return strconv.ParseBool(in)
	}

	return nil, fmt.Errorf("don't know how to convert %v", typ.Name())
}

// inspired by https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toKebabCase(str string) string {
	kebab := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	kebab = matchAllCap.ReplaceAllString(kebab, "${1}-${2}")
	return strings.ToLower(kebab)
}

func toUpperSnakeCase(str string) string {
	uppersnake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	uppersnake = matchAllCap.ReplaceAllString(uppersnake, "${1}_${2}")
	return strings.ToUpper(uppersnake)
}
