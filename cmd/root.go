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
	"github.com/wasilak/profilego"
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

			if viper.GetBool("profilingEnabled") {
				ProfileGoConfig := profilego.ProfileGoConfig{
					ApplicationName: viper.GetString("profilerApplicationName"),
					ServerAddress:   viper.GetString("profilerServerAddress"),
					Tags:            map[string]string{"version": common.GetVersion()},
				}

				err := profilego.Init(ProfileGoConfig)
				if err != nil {
					common.HandleError(ctx, err)
					panic(err)
				}
			}

			if viper.GetBool("otelEnabled") {
				otelGoTracingConfig := otelgotracer.OtelGoTracingConfig{
					HostMetricsEnabled:    true,
					RuntimeMetricsEnabled: true,
				}
				ctx, _, err := otelgotracer.Init(ctx, otelGoTracingConfig)
				if err != nil {
					common.HandleError(ctx, err)
					panic(err)
				}

				otelGoMetricsConfig := otelgometrics.OtelGoMetricsConfig{}

				ctx, common.MeterProvider, err = otelgometrics.Init(ctx, otelGoMetricsConfig)
				if err != nil {
					common.HandleError(ctx, err)
					panic(err)
				}
			}

			ctx, span := common.TracerCmd.Start(ctx, "rootCmd")
			defer span.End()

			loggerConfig := loggergo.Config{
				Level:  loggergo.LogLevelFromString(viper.GetString("loglevel")),
				Format: loggergo.LogFormatFromString(viper.GetString("logformat")),
			}

			if viper.GetBool("otelEnabled") {
				loggerConfig.OtelServiceName = common.GetAppName()
				loggerConfig.OtelLoggerName = "github.com/wasilak/go-hello-world"
				loggerConfig.OtelTracingEnabled = false
			}

			ctx, spanLoggerGo := common.TracerCmd.Start(ctx, "loggergo.LoggerInit")
			_, err := loggergo.LoggerInit(ctx, loggerConfig)
			if err != nil {
				common.HandleError(ctx, err)
				panic(err)
			}
			spanLoggerGo.End()

			slog.DebugContext(ctx, fmt.Sprintf("%+v", viper.AllSettings()))

			db.DB, err = db.NewMongoDB(ctx)
			if err != nil {
				common.HandleError(ctx, err)
				panic(err)
			}

			ctx, spanNewS3MinioStorage := common.TracerCmd.Start(ctx, "NewS3MinioStorage")
			storage.Storage, err = storage.NewS3MinioStorage(ctx)
			if err != nil {
				common.HandleError(ctx, err)
				panic(err)
			}
			spanNewS3MinioStorage.End()

			web.Init(ctx)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	common.CTX = context.Background()

	cobra.OnInitialize(libs.InitConfig)

	rootCmd.PersistentFlags().StringVar(&libs.CfgFile, "config", "", "config file (default is $HOME/."+common.GetAppName()+"/config.yml)")
	rootCmd.PersistentFlags().StringVar(&libs.Listen, "listen", "127.0.0.1:3000", "listen address")

	viper.BindPFlag("listen", rootCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("cacheEnabled", rootCmd.PersistentFlags().Lookup("cacheEnabled"))

	rootCmd.AddCommand(versionCmd)
}
