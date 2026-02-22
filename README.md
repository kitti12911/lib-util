# lib-util

shared Go utility library used across homelab services. provides config loading, structured logging, tracing, and validation.

## install

```bash
go get github.com/kitti12911/lib-util
```

## packages

### config

type-safe config loading from file with environment variable overrides and validation.

```go
import libconfig "github.com/kitti12911/lib-util/config"

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
import "github.com/kitti12911/lib-util/logger"

logger.New(
    logger.WithLevel(logger.LevelInfo),
    logger.WithServiceName("my-service"),
    logger.WithSource(),
    logger.WithTrace(),
)

slog.Info("server started", "port", 8080)
```

options:

| function            | description                                |
|---------------------|--------------------------------------------|
| `WithLevel(level)`  | set log level (debug, info, warn, error)   |
| `WithServiceName(n)`| add service name to all log entries        |
| `WithSource()`      | include source file and line in logs       |
| `WithTrace()`       | add trace_id and span_id from opentelemetry|

### tracing

opentelemetry tracing setup with OTLP gRPC exporter.

```go
import "github.com/kitti12911/lib-util/tracing"

tp, err := tracing.New(ctx, "my-service", "localhost:4317")
if err != nil {
    log.Fatal(err)
}
defer tracing.Shutdown(ctx, tp)
```

- exports traces via OTLP gRPC (e.g. to alloy, otel collector)
- sets global tracer provider
- supports TraceContext and Baggage propagation

### validator

struct validation with structured error reporting. wraps `go-playground/validator/v10`.

```go
import libvalidator "github.com/kitti12911/lib-util/validator"

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

## requirements

- go 1.26.0 or higher

## available commands

```bash
make tidy       # go mod tidy
make fmt        # format code
make test       # run tests with race detector
make cov        # run tests with coverage report
```
