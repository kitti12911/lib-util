package logger

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func (l Level) toSlog() slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type Options struct {
	Level       Level
	AddSource   bool
	ServiceName string
	EnableTrace bool
}

type Config struct {
	Level          Level `mapstructure:"level"            env:"LOG_LEVEL"            validate:"oneof=debug info warn error"`
	AddSource      bool  `mapstructure:"add_source"       env:"LOG_ADD_SOURCE"`
	IncludeTraceID bool  `mapstructure:"include_trace_id" env:"LOG_INCLUDE_TRACE_ID"`
}

type Option func(*Options)

func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithSource() Option {
	return func(o *Options) {
		o.AddSource = true
	}
}

func WithServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

func WithTrace() Option {
	return func(o *Options) {
		o.EnableTrace = true
	}
}

// NewFromConfig builds a logger from Config and installs it as the process-wide
// default slog logger. Intended for one-time application startup; use New for
// loggers that should not mutate global state.
func NewFromConfig(cfg Config, serviceName string) *slog.Logger {
	opts := []Option{
		WithLevel(cfg.Level),
	}

	if serviceName != "" {
		opts = append(opts, WithServiceName(serviceName))
	}

	if cfg.AddSource {
		opts = append(opts, WithSource())
	}

	if cfg.IncludeTraceID {
		opts = append(opts, WithTrace())
	}

	l := New(opts...)
	SetDefault(l)
	return l
}

// New creates a new JSON-formatted slog.Logger with the given options.
// It does not mutate the process-wide default logger; call SetDefault for that.
func New(opts ...Option) *slog.Logger {
	options := &Options{
		Level: LevelInfo,
	}

	for _, opt := range opts {
		opt(options)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     options.Level.toSlog(),
		AddSource: options.AddSource,
	})

	var baseHandler slog.Handler = handler

	if options.ServiceName != "" {
		baseHandler = handler.WithAttrs([]slog.Attr{
			slog.String("service", options.ServiceName),
		})
	}

	if options.EnableTrace {
		baseHandler = &traceHandler{handler: baseHandler}
	}

	return slog.New(baseHandler)
}

// SetDefault installs l as the process-wide default slog logger. Call this
// once during application startup.
func SetDefault(l *slog.Logger) {
	slog.SetDefault(l)
}

type traceHandler struct {
	handler slog.Handler
}

func (h *traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *traceHandler) Handle(ctx context.Context, record slog.Record) error {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()

	if spanCtx.HasTraceID() {
		record.AddAttrs(slog.String("trace_id", spanCtx.TraceID().String()))
	}

	if spanCtx.HasSpanID() {
		record.AddAttrs(slog.String("span_id", spanCtx.SpanID().String()))
	}

	return h.handler.Handle(ctx, record)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{handler: h.handler.WithGroup(name)}
}
