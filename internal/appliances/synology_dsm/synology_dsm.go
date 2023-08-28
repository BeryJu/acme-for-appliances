package synology_dsm

import (
	"net/http"

	"beryju.org/acme-for-appliances/internal/appliances"
)

type SynologyDSM struct {
	appliances.Appliance

	client *http.Client
	sid    string

	existingCert *SynologyAPICert
}

func (dsm *SynologyDSM) Init() error {
	dsm.client = dsm.HTTPClient()
	err := dsm.login()
	if err != nil {
		return err
	}
	return nil
}
