package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"rest-skeleton/internal/pkg/myctx"
	"runtime"
	"time"

	"github.com/bytedance/sonic"
	"github.com/grafana/loki-client-go/loki"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

type Logger struct {
	Log              *log.Logger
	LokiClient       *loki.Client
	Format           LoggerFormat
	ErrorCountMetric metric.Int64Counter
}
type LoggerFormat struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Context   string `json:"context"`
}

func New(filename string) *Logger {
	if os.Getenv("APP_ENV") == "production" {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return &Logger{Log: log.New(file, "", 0), LokiClient: nil}
	}
	return &Logger{Log: log.New(os.Stdout, "", 0), LokiClient: nil}
}

func (l *Logger) Error(ctx context.Context, err error) error {
	if ok := l.format(ctx, "ERROR", err.Error()); ok {
		l.ErrorCountMetric.Add(ctx, 1)
		message, _ := sonic.Marshal(l.Format)
		l.Log.Println(string(message))

		// span := trace.SpanFromContext(ctx)
		_, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "Error")
		defer span.End()

		span.SetAttributes(attribute.String("error", err.Error()))
		span.SetStatus(codes.Error, "Error handling request")
		span.RecordError(err)
	} else {
		l.Log.Println(err.Error())
	}
	return err
}
func (l *Logger) Info(ctx context.Context, msg string) {
	if ok := l.format(ctx, "INFO", msg); ok {
		message, _ := sonic.Marshal(l.Format)
		l.Log.Println(string(message))
	} else {
		l.Log.Println(msg)
	}
}

func (l *Logger) Fatal(ctx context.Context, err error) {
	l.Error(ctx, err)
	os.Exit(1)
}

func (l *Logger) format(ctx context.Context, level string, msg string) bool {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		l.Format = LoggerFormat{
			Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			Level:     level,
			Message:   msg,
			File:      path.Base(file),
			Line:      line,
			Context:   ctx.Value(myctx.Key("traceID")).(string),
		}
	}
	return ok
}
