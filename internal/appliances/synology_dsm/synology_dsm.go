package synology_dsm

import (
	"beryju.io/acme-for-appliances/internal/appliances"
	"beryju.io/acme-for-appliances/internal/appliances/synology_dsm/api"
)

type SynologyDSM struct {
	appliances.Appliance

	client *api.SynologyAPI

	existingCert *api.SynologyAPICert
}

func (dsm *SynologyDSM) Init() error {
	dsm.client = &api.SynologyAPI{
		URL:      dsm.URL,
		Username: dsm.Username,
		Password: dsm.Password,
		Client:   dsm.HTTPClient(),
		Logger:   dsm.Logger.WithField("component", "api"),
	}
	err := dsm.client.Login()
	if err != nil {
		return err
	}
	return nil
}
