package synology_dsm

import (
	"time"
)

func (dsm *SynologyDSM) CheckExpiry() (int, error) {
	certs, err := dsm.client.ListCertificates()
	if err != nil {
		return 0, err
	}
	for _, cert := range certs.Data.Certificates {
		dsm.Logger.WithField("cert", cert.Desc).Trace(cert.ValidTill.Human())
		if cert.Desc != dsm.Extension["cert_desc"] {
			continue
		}
		t := time.Until(time.Time(cert.ValidTill))
		dsm.existingCert = &cert
		return int(t.Hours()) / 24, nil
	}
	dsm.Logger.Info("Cert not found")
	return 0, nil
}
