package cmd

import (
	"context"
	"fmt"
	"os"

	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wasilak/loggergo"
	"github.com/wasilak/notes-manager/libs"
	"github.com/wasilak/notes-manager/libs/common"
	"github.com/wasilak/notes-manager/libs/providers/db"
	"github.com/wasilak/notes-manager/libs/providers/storage"
	"github.com/wasilak/notes-manager/libs/web"
	otelgometrics "github.com/wasilak/otelgo/metrics"
	otelgotracer "github.com/wasilak/otelgo/tracing"
)

var (
	rootCmd = &cobra.Command{
		Use:   "notes-manager",
		Short: "Notes Manager",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetContext(common.CTX)
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			if viper.GetBool("otelEnabled") {
				otelGoTracingConfig := otelgotracer.OtelGoTracingConfig{
					HostMetricsEnabled: true,
				}
				ctx, _, err := otelgotracer.Init(ctx, otelGoTracingConfig)
				if err != nil {
					slog.ErrorContext(ctx, err.Error())
					os.Exit(1)
				}

				otelGoMetricsConfig := otelgometrics.OtelGoMetricsConfig{}

				var errMetrics error
				ctx, common.MeterProvider, errMetrics = otelgometrics.Init(ctx, otelGoMetricsConfig)
				if errMetrics != nil {
					slog.ErrorContext(ctx, errMetrics.Error())
					os.Exit(1)
				}
			}

			ctx, span := common.TracerCmd.Start(ctx, "rootCmd")

			loggerConfig := loggergo.LoggerGoConfig{
				Level:  viper.GetString("loglevel"),
				Format: viper.GetString("logformat"),
			}

			ctx, spanLoggerGo := common.TracerCmd.Start(ctx, "loggergo.LoggerInit")
			_, err := loggergo.LoggerInit(loggerConfig)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				os.Exit(1)
			}
			spanLoggerGo.End()

			slog.DebugContext(ctx, fmt.Sprintf("%+v", viper.AllSettings()))

			db.DB, err = db.NewMongoDB(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "Error initializing database:", err)
				panic(err)
			}

			ctx, spanNewS3MinioStorage := common.TracerCmd.Start(ctx, "NewS3MinioStorage")
			storage.Storage, err = storage.NewS3MinioStorage(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "Error initializing storage:", err)
				panic(err)
			}
			spanNewS3MinioStorage.End()

			span.End()

			web.Init(ctx)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	common.CTX = context.Background()

	cobra.OnInitialize(libs.InitConfig)

	rootCmd.PersistentFlags().StringVar(&libs.CfgFile, "config", "", "config file (default is $HOME/."+common.AppName+"/config.yml)")
	rootCmd.PersistentFlags().StringVar(&libs.Listen, "listen", "127.0.0.1:3000", "listen address")

	viper.BindPFlag("listen", rootCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("cacheEnabled", rootCmd.PersistentFlags().Lookup("cacheEnabled"))

	rootCmd.AddCommand(versionCmd)
}
