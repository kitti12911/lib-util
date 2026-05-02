# lib-util

shared Go utility library used across homelab services. provides config loading, structured logging, validation, pagination, formatting, and pointer utilities.

## install

```bash
go get github.com/kitti12911/lib-util/v2
```

## observability

profiling and tracing now live in `github.com/kitti12911/lib-monitor`. use `github.com/kitti12911/lib-monitor/profiling` and `github.com/kitti12911/lib-monitor/tracing` for new code.

## packages

### config

type-safe config loading from file with environment variable overrides and validation.

```go
import libconfig "github.com/kitti12911/lib-util/v2/config"

type Config struct {
    Port int    `mapstructure:"port" env:"PORT" validate:"required"`
    Name string `mapstructure:"name" env:"NAME" validate:"required"`
}

cfg, err := libconfig.Load[Config]("config.yml")
```

- supports yaml, json, toml via viper
- binds env variables using `env` struct tag
- validates struct fields using `validate` struct tag
- supports nested structs

### logger

structured JSON logging built on Go's `slog`. supports opentelemetry trace context injection.

```go
import "github.com/kitti12911/lib-util/v2/logger"

logger.NewFromConfig(logger.Config{
    Level:          logger.LevelInfo,
    AddSource:      true,
    IncludeTraceID: true,
}, "my-service")

slog.Info("server started", "port", 8080)
```

options:

| function             | description                                 |
| -------------------- | ------------------------------------------- |
| `WithLevel(level)`   | set log level (debug, info, warn, error)    |
| `WithServiceName(n)` | add service name to all log entries         |
| `WithSource()`       | include source file and line in logs        |
| `WithTrace()`        | add trace_id and span_id from opentelemetry |

### validator

struct validation with structured error reporting. wraps `go-playground/validator/v10`.

```go
import libvalidator "github.com/kitti12911/lib-util/v2/validator"

v := libvalidator.New()

err := v.Validate(myStruct)

// or get detailed field violations
violations, err := v.ValidateWithErrors(myStruct)
for _, v := range violations {
    fmt.Printf("field: %s, tag: %s, condition: %s\n", v.Field, v.Tag, v.Condition)
}

// register custom validation
v.RegisterCustom("my_tag", func(fl validator.FieldLevel) bool {
    return fl.Field().String() != ""
})
```

### formatter

JSON formatting helpers for debug output, logs, and safe string conversion.

```go
import "github.com/kitti12911/lib-util/v2/formatter"

text, err := formatter.ToJSONStr(payload, true) // pretty-printed JSON
```

| function               | description                        |
| ---------------------- | ---------------------------------- |
| `ToJSONStr(v, indent)` | marshal a value into a JSON string |

### pagination

small helpers for converting page/pageSize input into limit/offset and response metadata.

```go
import "github.com/kitti12911/lib-util/v2/pagination"

input := pagination.ParseInput(page, pageSize)
items, total, err := repo.Find(ctx, input.Limit, input.Offset)
output := pagination.CalcOutput(page, pageSize, total)
```

| function                            | description                                |
| ----------------------------------- | ------------------------------------------ |
| `ParseInput(page, pageSize)`        | normalize page input into limit and offset |
| `CalcOutput(page, pageSize, total)` | calculate response pagination metadata     |

### ptr

pointer helper functions for safely dereferencing pointers with fallback values.

```go
import "github.com/kitti12911/lib-util/v2/ptr"

name := ptr.ValueOr(input.Name, "unnamed") // returns "unnamed" if input.Name is nil
limit := ptr.From(input.Limit)             // returns 0 if input.Limit is nil
```

| function                        | description                                         |
| ------------------------------- | --------------------------------------------------- |
| `From(p *T) T`                  | dereference pointer, returns zero value of T if nil |
| `ValueOr(p *T, defaultVal T) T` | dereference pointer, returns defaultVal if nil      |

## requirements

- go 1.26 or higher

## available commands

```bash
make tidy       # go mod tidy
make fmt        # format code
make test       # run tests with race detector
make cov        # run tests with coverage report
```
