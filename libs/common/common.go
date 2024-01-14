package common

import (
	"context"
	"runtime/debug"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	Version       string
	CTX           context.Context
	AppName       = "notesmanager"
	Tracer        trace.Tracer
	MeterProvider *metric.MeterProvider
)

func GetVersion() string {
	buildInfo, _ := debug.ReadBuildInfo()
	return buildInfo.GoVersion
}
