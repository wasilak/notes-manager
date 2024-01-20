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
