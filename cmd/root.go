package cmd

import (
	"fmt"
	"os"

	"github.com/dovydasdo/clerk/config"
	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "clerk",
	Short: "PSEC data ingestion manager",
	Long:  `clerk is a cli tool used to process and save data gathered by PSEC scrapers`,
	Run: func(cmd *cobra.Command, args []string) {
		col := config.Collection{}
		err := viper.UnmarshalKey("collection", &col)
		if err != nil {
			panic("failed to unmarshall config")
		}

		if len(col.Managers) < 1 {
			panic("no souces provided in config")
		}

		var level slog.Level
		configLevel := viper.Get("log_level")
		switch configLevel {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelDebug
		}

		opts := &slog.HandlerOptions{
			Level: level,
		}

		handler := slog.NewJSONHandler(os.Stdout, opts)
		logger := slog.New(handler)

		for _, src := range col.Managers {
			switch src.Type {
			case "rent":
				logger.Info("init", "message", fmt.Sprintf("starting manager of type %+v", src.Type))
			default:
				logger.Warn("init", "message", fmt.Sprintf("source of type %v not implemented", src.Type))
			}
		}
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
	rootCmd.PersistentFlags().StringP("log_level", "l", "info", "Log level (debug, info, error)")
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log_level"))
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
