package common

import (
	"context"
	"os"
	"runtime/debug"

	"go.opentelemetry.io/otel"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

var (
	Version       string
	CTX           context.Context
	AppName       = "notesmanager"
	TracerCmd     = otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
	TracerWeb     = otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
	MeterProvider = sdk.NewMeterProvider()
)

func GetVersion() string {
	buildInfo, _ := debug.ReadBuildInfo()
	return buildInfo.GoVersion
}
