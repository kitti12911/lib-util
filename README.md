# lib-util

shared Go utility library used across homelab services. provides config
loading, structured logging, validation, pagination, field masks, formatting,
Huma helpers, protobuf helpers, query parsing, and pointer utilities.

## install

```bash
go get github.com/kitti12911/lib-util/v3
```

## observability

profiling and tracing now live in
[`github.com/kitti12911/lib-monitor`](https://github.com/kitti12911/lib-monitor).
use `github.com/kitti12911/lib-monitor/profiling` and
`github.com/kitti12911/lib-monitor/tracing` for new code.

## project structure

```bash
lib-util/
├── apperror/    # transport-neutral application error helpers
├── config/       # type-safe config loading
├── fieldmask/    # protobuf field mask validation and extraction
├── formatter/    # JSON/string formatting helpers
├── huma/         # Huma response, PATCH, and gRPC error helpers
├── logger/       # slog setup and trace context logging
├── pagination/   # page, limit, offset, and response metadata helpers
├── protoutil/    # protobuf conversion helpers
├── ptr/          # pointer helpers
├── query/        # string query parsing helpers
├── validator/    # go-playground validator wrapper
├── Makefile
├── go.mod
└── README.md
```

## packages

### apperror

transport-neutral application errors for service/domain layers.

```go
import "github.com/kitti12911/lib-util/v3/apperror"

return apperror.InvalidInput("invalid request", err)
```

| function/type          | description                                  |
| ---------------------- | -------------------------------------------- |
| `Error`                | typed application error with message + cause |
| `Code`                 | stable application error category            |
| `Internal(message, e)` | create internal error                        |
| `NotFound(message, e)` | create not found error                       |
| `AlreadyExist(...)`    | create duplicate resource error              |
| `InvalidInput(...)`    | create validation/input error                |
| `Unauthorized(...)`    | create authentication error                  |
| `Forbidden(...)`       | create authorization error                   |

### config

type-safe config loading from file with environment variable overrides and validation.

```go
import libconfig "github.com/kitti12911/lib-util/v3/config"

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
import "github.com/kitti12911/lib-util/v3/logger"

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
import libvalidator "github.com/kitti12911/lib-util/v3/validator"

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
import "github.com/kitti12911/lib-util/v3/formatter"

text, err := formatter.ToJSONStr(payload, true) // pretty-printed JSON
```

| function               | description                        |
| ---------------------- | ---------------------------------- |
| `ToJSONStr(v, indent)` | marshal a value into a JSON string |

### pagination

small helpers for converting page/pageSize input into limit/offset and response metadata.

```go
import "github.com/kitti12911/lib-util/v3/pagination"

input := pagination.ParseInput(page, pageSize)
items, total, err := repo.Find(ctx, input.Limit, input.Offset)
output := pagination.CalcOutput(page, pageSize, total)
```

| function                            | description                                |
| ----------------------------------- | ------------------------------------------ |
| `ParseInput(page, pageSize)`        | normalize page input into limit and offset |
| `CalcOutput(page, pageSize, total)` | calculate response pagination metadata     |

### fieldmask

protobuf field mask helpers for PATCH-style APIs.

```go
import "github.com/kitti12911/lib-util/v3/fieldmask"

immutable := map[string]bool{
    "id":         true,
    "created_at": true,
}

if err := fieldmask.ValidateMask(req.GetUpdateMask(), req.GetUser(), immutable); err != nil {
    return err
}
```

- `ValidateMask(mask, msg, immutable)`: reject empty masks, unknown fields, and messages
- `ExtractChanges(mask, msg)`: read masked protobuf values into a map
- `ExtractNestedChanges(changes, fields, n)`: select one nested object from extracted changes

### huma

small helpers for Huma-based HTTP APIs.

```go
import humautil "github.com/kitti12911/lib-util/v3/huma"

huma.Patch(api, "/users/{id}", handler, humautil.WithTag("Users"))
```

| function/type        | description                                       |
| -------------------- | ------------------------------------------------- |
| `Patch[T]`           | tri-state PATCH value: omitted, null, or value    |
| `GRPCError(err)`     | map common gRPC status errors to Huma HTTP errors |
| `WithTag(tag)`       | set one OpenAPI operation tag                     |
| `StatusCreated(op)`  | set default response status to 201                |
| `AffectedRows(rows)` | build a common affected rows response body        |

### protoutil

protobuf conversion helpers.

```go
import "github.com/kitti12911/lib-util/v3/protoutil"

createdAt := protoutil.TimeFromProto(resp.GetCreatedAt())
```

| function               | description                                  |
| ---------------------- | -------------------------------------------- |
| `TimeFromProto(value)` | convert protobuf timestamp to Go `time.Time` |

### query

string parsers for API query parameters that map to the shared query enum
numbers used by `proto-sandbox`.

```go
import "github.com/kitti12911/lib-util/v3/query"

op := query.FilterOpFromString[commonv1.FilterOp]("like_ci")
direction := query.OrderDirectionFromString[commonv1.OrderDirection]("desc")
```

| function                      | description                                  |
| ----------------------------- | -------------------------------------------- |
| `FilterOpFromString[O](op)`   | parse exact, like, gt, lt, null, in, between |
| `OrderDirectionFromString[O]` | parse asc or desc                            |

### ptr

pointer helper functions for safely dereferencing pointers with fallback values.

```go
import "github.com/kitti12911/lib-util/v3/ptr"

name := ptr.ValueOr(input.Name, "unnamed") // returns "unnamed" if input.Name is nil
limit := ptr.From(input.Limit)             // returns 0 if input.Limit is nil
```

| function                        | description                                         |
| ------------------------------- | --------------------------------------------------- |
| `From(p *T) T`                  | dereference pointer, returns zero value of T if nil |
| `ValueOr(p *T, defaultVal T) T` | dereference pointer, returns defaultVal if nil      |

## requirements

- go 1.26 or higher

Optional:

- [prettier](https://prettier.io/) for Markdown, YAML, JSON, and JSONC formatting

## available commands

```bash
make tidy       # go mod tidy
make fmt        # format code
make pretty     # format docs and config with Prettier
make format     # run fmt and pretty
make test       # run tests with race detector
make cov        # run tests with coverage report
```
