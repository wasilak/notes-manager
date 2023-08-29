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
)

var (
	err error

	rootCmd = &cobra.Command{
		Use:   "notes-manager",
		Short: "Notes Mana",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetContext(common.CTX)
		},
		Run: func(cmd *cobra.Command, args []string) {

			loggerConfig := loggergo.LoggerGoConfig{
				Level:  viper.GetString("loglevel"),
				Format: viper.GetString("logformat"),
			}

			_, err := loggergo.LoggerInit(loggerConfig)
			if err != nil {
				slog.ErrorContext(common.CTX, err.Error())
				os.Exit(1)
			}

			slog.DebugContext(common.CTX, fmt.Sprintf("%+v", viper.AllSettings()))

			db.DB, err = db.NewMongoDB()
			if err != nil {
				slog.ErrorContext(common.CTX, "Error initializing database:", err)
				panic(err)
			}

			storage.Storage, err = storage.NewS3MinioStorage()
			if err != nil {
				slog.ErrorContext(common.CTX, "Error initializing storage:", err)
				panic(err)
			}

			web.Init()
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

	rootCmd.PersistentFlags().StringVar(&libs.CfgFile, "config", "", "config file (default is $HOME/."+libs.AppName+"/config.yml)")
	rootCmd.PersistentFlags().StringVar(&libs.Listen, "listen", "127.0.0.1:3000", "listen address")

	viper.BindPFlag("listen", rootCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("cacheEnabled", rootCmd.PersistentFlags().Lookup("cacheEnabled"))

	rootCmd.AddCommand(versionCmd)
}
