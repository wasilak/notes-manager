package libs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wasilak/notes-manager/libs/common"
)

var (
	CfgFile string
	Listen  string
)

func InitConfig() {
	godotenv.Load()

	viper.SetEnvPrefix(common.AppName)

	viper.SetDefault("loglevel", "info")
	viper.SetDefault("logformat", "plain")
	viper.SetDefault("profilingEnabled", false)
	viper.SetDefault("profilerApplicationName", common.AppName)
	viper.SetDefault("profilerServerAddress", "")
	viper.SetDefault("openAIEnabled", false)

	if len(os.Getenv("OPENAI_API_KEY")) > 0 {
		viper.SetDefault("openAIEnabled", true)
		viper.SetDefault("openAIAPIKey", os.Getenv("OPENAI_API_KEY"))
	}

	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigType("yaml")
		viper.SetConfigName(common.AppName)
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
