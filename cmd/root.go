package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/BeryJu/acme-for-appliances/internal"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var force bool
var infinite bool
var checkInterval int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "acme-for-appliances",
	Short: "ACME Certificates for appliances",
	Long:  `Use ACME Certificates for appliances which don't natively support them.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !infinite {
			internal.Main(force)
			os.Exit(0)
		}
		d := time.Duration(checkInterval) * time.Hour
		ticker := time.NewTicker(d)
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		log.Infof("Running, will check in %s", d)
		for {
			select {
			case <-ticker.C:
				internal.Main(force)
			case <-quit:
				ticker.Stop()
				os.Exit(0)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force renewal")
	rootCmd.PersistentFlags().BoolVarP(&infinite, "infinite", "i", false, "Infinite mode, keep running the program infinitley and check every interval.")
	rootCmd.PersistentFlags().IntVarP(&checkInterval, "check-interval", "n", 24, "Interval for infinite mode, in hours")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.SetLevel(log.DebugLevel)
	viper.SetConfigType("toml")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("config.toml")
	}
	viper.AddConfigPath("/config")

	viper.SetDefault("storage", "storage")
	viper.SetDefault("acme.directory_url", "https://acme-staging-v02.api.letsencrypt.org/directory")
	viper.SetDefault("acme.refresh_threshold", 15)
	viper.SetDefault("acme.resolvers", []string{})

	viper.SetEnvPrefix("a4a")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.WithError(err).Warning("failed to load config file")
	}
}
