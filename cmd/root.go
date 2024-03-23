package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "clerk",
	Short: "PSEC data queue manager",
	Long:  `clerk is a cli tool used to process and save data gathered by PSEC scrapers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("from Run")
		val := viper.Get("collection")
		fmt.Printf("found in cfg: %v \n", val)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	// For specifying where goes what and when
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASS")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_DEBUG")
	viper.BindEnv("NATS_CREDS")
	viper.BindEnv("NATS_SERVER")
	viper.BindEnv("TURSO_URL")
	viper.BindEnv("TURSO_TOKEN")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
