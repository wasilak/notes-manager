package common

import (
	"context"
	"os"
	"runtime/debug"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	Version       string
	CTX           context.Context
	AppName       = "notesmanager"
	TracerCmd     = otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
	TracerWeb     = otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
	MeterProvider = metric.NewMeterProvider()
)

func getGitRevision() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}

func GetVersion() string {
	if Version != "" {
		return Version
	}

	return getGitRevision()
}

func HandleError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}
