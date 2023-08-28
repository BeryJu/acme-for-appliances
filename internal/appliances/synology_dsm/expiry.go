package synology_dsm

import (
	"time"
)

// func (dsm *SynologyDSM) debug() {
// 	req, _ := dsm.makeRequest(http.MethodGet, "webapi/query.cgi", map[string]string{
// 		"api":     "SYNO.API.Info",
// 		"version": "1",
// 		"method":  "query",
// 		"query":   "all",
// 	}, nil)
// 	res, _ := dsm.HTTPClient().Do(req)
// 	b, _ := io.ReadAll(res.Body)
// 	fmt.Println(string(b))
// }

func (dsm *SynologyDSM) CheckExpiry() (int, error) {
	certs, err := dsm.listCertificates()
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
