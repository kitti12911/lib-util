package logger

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestLevelToSlog(t *testing.T) {
	tests := map[Level]slog.Level{
		LevelDebug: slog.LevelDebug,
		LevelInfo:  slog.LevelInfo,
		LevelWarn:  slog.LevelWarn,
		LevelError: slog.LevelError,
		"":         slog.LevelInfo,
	}

	for level, want := range tests {
		assert.Equal(t, want, level.toSlog())
	}
}

func TestOptions(t *testing.T) {
	opts := &Options{}
	WithLevel(LevelDebug)(opts)
	WithSource()(opts)
	WithServiceName("svc")(opts)
	WithTrace()(opts)

	assert.Equal(t, LevelDebug, opts.Level)
	assert.True(t, opts.AddSource)
	assert.Equal(t, "svc", opts.ServiceName)
	assert.True(t, opts.EnableTrace)
}

func TestNewAndNewFromConfig(t *testing.T) {
	out := captureStdout(t, func() {
		l := New(
			WithLevel(LevelDebug),
			WithServiceName("svc"),
			WithTrace(),
		)

		ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: trace.TraceID{1, 2, 3},
			SpanID:  trace.SpanID{4, 5, 6},
		}))
		l.InfoContext(ctx, "hello")

		NewFromConfig(Config{
			Level:          LevelInfo,
			AddSource:      true,
			IncludeTraceID: true,
		}, "svc-from-config").Info("from config")
	})

	assert.Contains(t, out, `"service":"svc"`)
	assert.Contains(t, out, `"trace_id"`)
	assert.Contains(t, out, `"span_id"`)
	assert.Contains(t, out, `"service":"svc-from-config"`)
}

func TestNewDoesNotMutateDefault(t *testing.T) {
	original := slog.Default()
	t.Cleanup(func() { slog.SetDefault(original) })

	sentinel := slog.New(slog.NewJSONHandler(io.Discard, nil))
	slog.SetDefault(sentinel)

	_ = New(WithLevel(LevelDebug))
	assert.Same(t, sentinel, slog.Default(), "New must not mutate slog default")
}

func TestNewFromConfigSetsDefault(t *testing.T) {
	original := slog.Default()
	t.Cleanup(func() { slog.SetDefault(original) })

	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, nil)))

	captureStdout(t, func() {
		l := NewFromConfig(Config{Level: LevelInfo}, "init")
		assert.Same(t, l, slog.Default(), "NewFromConfig must install l as the default")
	})
}

func TestSetDefault(t *testing.T) {
	original := slog.Default()
	t.Cleanup(func() { slog.SetDefault(original) })

	target := slog.New(slog.NewJSONHandler(io.Discard, nil))
	SetDefault(target)
	assert.Same(t, target, slog.Default())
}

func TestTraceHandlerPartialSpan(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		wantTrace bool
		wantSpan  bool
	}{
		{
			name:      "trace id only",
			ctx:       trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}})),
			wantTrace: true,
			wantSpan:  false,
		},
		{
			name:      "span id only",
			ctx:       trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{SpanID: trace.SpanID{1}})),
			wantTrace: false,
			wantSpan:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			h := &traceHandler{handler: slog.NewJSONHandler(&out, nil)}
			require.NoError(t, h.Handle(tt.ctx, slog.NewRecord(time.Time{}, slog.LevelInfo, "x", 0)))

			assert.Equal(t, tt.wantTrace, bytes.Contains(out.Bytes(), []byte(`"trace_id"`)))
			assert.Equal(t, tt.wantSpan, bytes.Contains(out.Bytes(), []byte(`"span_id"`)))
		})
	}
}

func TestTraceHandlerWithoutSpan(t *testing.T) {
	var out bytes.Buffer
	h := &traceHandler{handler: slog.NewJSONHandler(&out, nil)}

	assert.True(t, h.Enabled(context.Background(), slog.LevelInfo))
	require.NoError(t, h.Handle(context.Background(), slog.NewRecord(time.Time{}, slog.LevelInfo, "hello", 0)))
	assert.Contains(t, out.String(), `"msg":"hello"`)

	withAttrs := h.WithAttrs([]slog.Attr{slog.String("k", "v")})
	withGroup := h.WithGroup("group")
	assert.NotNil(t, withAttrs)
	assert.NotNil(t, withGroup)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	fn()

	require.NoError(t, w.Close())
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	return string(data)
}
