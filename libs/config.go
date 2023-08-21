package libs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	AppName = "notesmanager"
	CfgFile string
	Listen  string
)

func InitConfig() {
	godotenv.Load()

	viper.SetEnvPrefix(AppName)

	viper.SetDefault("loglevel", "info")
	viper.SetDefault("logformat", "plain")

	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigType("yaml")
		viper.SetConfigName(AppName)
		viper.AddConfigPath(home)
		viper.AddConfigPath("./")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Printf("%+v\n", err)
	}
}
