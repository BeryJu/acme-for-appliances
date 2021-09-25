package internal

import (
	"beryju.org/acme-for-appliances/internal/acme"
	"beryju.org/acme-for-appliances/internal/appliances"
	"beryju.org/acme-for-appliances/internal/storage"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Main(force bool) {
	l := log.WithField("component", "main")
	u := &acme.User{
		Email: viper.GetString("acme.user_email"),
	}
	c := acme.NewClient(u)
	cAppliances := viper.GetStringMap("appliances")
	threshold := viper.GetInt("acme.refresh_threshold")
	for appName, appMap := range cAppliances {
		// app is map[string]interface{}
		var app appliances.Appliance
		mapstructure.Decode(appMap, &app)
		app.Name = appName
		al := l.WithField("appliance", appName).WithField("type", app.Type)
		app.Logger = al
		appHandler := GetActual(&app)
		// Init handler, check connection, check validity of settings
		err := appHandler.Init()
		if err != nil {
			al.WithError(err).Warning("Appliance failed to init")
			continue
		}
		expiry, err := appHandler.CheckExpiry()
		if err != nil {
			al.WithError(err).Warning("CheckExpiry() failed")
			continue
		}
		if expiry == -1 {
			al.Info("CheckExpiry() returned -1, assuming cert doesn't exist")
		}
		// Check if domains are the same, otherwise renew
		s := storage.GetState(appName)
		if !s.CompareDomains(app.Domains) {
			al.WithField("expiry", expiry).Info("Domains changed, forcing renewal")
		} else if expiry >= threshold && !force {
			al.WithField("threshold", threshold).WithField("expiry", expiry).Info("Cert doesn't need to be renewed")
			continue
		}
		al.Infof("Cert expires in %d days", expiry)
		al.Debug("Starting cert renewal")
		certs, err := c.GetCerts(appHandler)
		if err != nil {
			al.WithError(err).Warning("Failed to get certs for appliance")
			continue
		}
		err = appHandler.Consume(certs)
		if err != nil {
			al.WithError(err).Warning("Appliance failed to consume certificates")
			continue
		}
		al.Info("appliance successfully consumed certificate")
		// Update state and write it back
		s.Domains = app.Domains
		s.Write(appName)
	}
}
