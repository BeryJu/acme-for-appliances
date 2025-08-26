package vmware_esxi

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"
)

func (v *VMwareESXi) CheckExpiry() (int, error) {
	v.Logger.Debug("Checking current certificate for ESXi host")

	// Create a custom transport that captures the certificate
	var serverCert *x509.Certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			VerifyConnection: func(cs tls.ConnectionState) error {
				if len(cs.PeerCertificates) > 0 {
					serverCert = cs.PeerCertificates[0]
				}
				return nil
			},
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	// Make a simple request to capture the certificate
	_, err := client.Get(v.URL)
	if err != nil {
		return -1, fmt.Errorf("failed to connect to ESXi host: %v", err)
	}

	if serverCert == nil {
		return -1, fmt.Errorf("no certificate found from ESXi host")
	}

	// Calculate days until expiry
	now := time.Now()
	daysUntilExpiry := int(serverCert.NotAfter.Sub(now).Hours() / 24)
	return daysUntilExpiry, nil
}
