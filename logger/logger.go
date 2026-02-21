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

// New creates a new JSON-formatted slog.Logger with the given options.
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

	l := slog.New(baseHandler)
	slog.SetDefault(l)

	return l
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
