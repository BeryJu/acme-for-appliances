package citrix_adc

import (
	"github.com/chiradeep/go-nitro/config/ssl"
	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/mitchellh/mapstructure"
)

func (adc *CitrixADC) CheckExpiry() (int, error) {
	certs, err := adc.client.FindResourceArray(netscaler.Sslcertkey.Type(), adc.Extension[ADCConfigCertName].(string))
	if err != nil {
		return -1, nil
	}
	if len(certs) < 1 {
		// No cert found, return -1 without error
		return -1, nil
	}
	var cert ssl.Sslcertkey
	mapstructure.Decode(certs[0], &cert)
	return cert.Daystoexpiration, nil
}
