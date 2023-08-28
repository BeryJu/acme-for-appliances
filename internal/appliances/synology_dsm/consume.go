package synology_dsm

import (
	"errors"

	"beryju.org/acme-for-appliances/internal/appliances/synology_dsm/api"
	"github.com/go-acme/lego/v4/certificate"
)

func (dsm *SynologyDSM) Consume(c *certificate.Resource) error {
	if dsm.existingCert == nil {
		dsm.existingCert = &api.SynologyAPICert{
			Desc: dsm.Extension["cert_desc"].(string),
		}
	}
	res, err := dsm.client.UploadCertificate(dsm.existingCert, c.Certificate, c.IssuerCertificate, c.PrivateKey)
	if !res.Success {
		return errors.New("unsuccessful")
	}
	return err
}
