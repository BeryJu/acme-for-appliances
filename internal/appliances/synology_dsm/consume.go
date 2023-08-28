package synology_dsm

import (
	"errors"

	"github.com/go-acme/lego/v4/certificate"
)

func (dsm *SynologyDSM) Consume(c *certificate.Resource) error {
	res, err := dsm.client.UploadCertificate(dsm.existingCert, c.Certificate, c.IssuerCertificate, c.PrivateKey)
	if !res.Success {
		return errors.New("unsuccessful")
	}
	return err
}
