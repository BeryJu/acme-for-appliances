package config

import (
	"beryju.org/acme-for-appliances/internal/appliances"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Storage string `mapstructure:"storage"`
	ACME    struct {
		DirectoryURL          string   `mapstructure:"directory_url"`
		Resolvers             []string `mapstructure:"resolvers"`
		TermsAgreed           bool     `mapstructure:"terms_agreed"`
		ChallengeProviderName string   `mapstructure:"challenge_provider_name"`
		UserEmail             string   `mapstructure:"user_email"`
		RefreshThreshold      int      `mapstructure:"refresh_threshold"`
	} `mapstructure:"acme"`
	Appliances map[string]appliances.Appliance `mapstructure:"appliances"`
}

var C Config

func Load(cfgFile string) {
	viper.SetConfigType("toml")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("config.toml")
	}

	viper.SetDefault("storage", "storage")
	viper.SetDefault("acme.directory_url", "https://acme-staging-v02.api.letsencrypt.org/directory")
	viper.SetDefault("acme.refresh_threshold", 15)
	viper.SetDefault("acme.resolvers", []string{})

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.WithError(err).Warning("failed to load config file")
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
}
