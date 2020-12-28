# NX Config (Go)

[![Go Reference](https://pkg.go.dev/badge/github.com/expo21xx/nxconfig.svg)](https://pkg.go.dev/github.com/expo21xx/nxconfig)

NX Config can be used to map environment variables and flags to struct values.
Using reflection (exported) struct field names. Nested structs will be prefixed with the field name.

## Usage

```go
import (
    "log"
    "time"

    "github.com/expo21xx/nxconfig"
)

type Config struct {
    Host string
    Port uint16
    PGConfig PGConfig `name:"PG"`
}

type PGConfig struct {
    Host string `default:"localhost"`
    Port uint16 `default:"5432"`
    Username string `usage:"postgres username"`
    Password string
    Timeout time.Duration `name:"connection-timeout"`
}

func main() {
    var cfg Config

    err := nxconfig.Load(&cfg)
    if err != nil {
        log.Fatal(err)
    }

    // use cfg
}
```

Starting the binary will now parse all flags, and check all environment variables.

NOTE: Flags override environment variables!

## Options

The `Load` function can be customized with a number of helper functions, passed in after the target.

```go
nxconfig.Load(&target, nxconfig.WithArgs([]string{"--override", "foo"}), nxconfig.WithEnv([]string{}))
```

If you already use [`spf13/pflag`](https://github.com/spf13/pflag) you can pass the `*pflag.FlagSet` to `Load`:

```go
nxconfig.Load(&target, nxconfig.WithFlagSet(myflagset))
```
NOTE: `Load` will call `flagset.Parse()` with the arguments it was passed, either with `WithArgs` or `os.Args[1:]` if not.


## Struct Tags:

- `name`: name in PascalCase (overrides automatic inferring). NOTE: Use an empty `""` to not export the field.
- `default`: default value
- `usage`: usage string to print when `--help` flag is passed

## Supported Types

The currently supported types are:

- `string`
- `int`, `int32`, `int64`
- `uint`, `uint32`, `uint64`
- `flaot32`, `flaot64`
- `bool`
- `time.Duration`

More types will be added in future on a *need to implement* basis.

