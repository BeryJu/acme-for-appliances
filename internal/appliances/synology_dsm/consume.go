package synology_dsm

import (
	"github.com/go-acme/lego/v4/certificate"
)

func (dsm *SynologyDSM) Consume(c *certificate.Resource) error {
	dsm.uploadCertificate(dsm.existingCert, c.Certificate, c.IssuerCertificate, c.PrivateKey)
	return nil
}
