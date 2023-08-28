package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"beryju.io/acme-for-appliances/internal"
	"beryju.io/acme-for-appliances/internal/config"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var cfgFile string
var force bool
var infinite bool
var checkInterval int
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "acme-for-appliances",
	Version: Version,
	Short:   "ACME Certificates for appliances",
	Long:    `Use ACME Certificates for appliances which don't natively support them.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !infinite {
			internal.Main(force)
			os.Exit(0)
		}
		d := time.Duration(checkInterval) * time.Hour
		ticker := time.NewTicker(d)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		log.Infof("Running, will check in %s", d)
		defer sentry.Flush(2 * time.Second)
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
	config.Load(cfgFile)
	log.WithField("version", Version).Info("acme-for-appliances")
	if os.Getenv("DISABLE_SENTRY") != "true" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              "https://52929b6b754742a3bb6a39ef5839b8e9@sentry.beryju.org/13",
			Release:          fmt.Sprintf("acme-for-appliances@%s", Version),
			AttachStacktrace: true,
			TracesSampleRate: 1,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
	}
}
