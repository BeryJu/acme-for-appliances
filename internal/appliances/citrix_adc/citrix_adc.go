package citrix_adc

import (
	"beryju.io/acme-for-appliances/internal/appliances"
	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/pkg/errors"
)

const (
	ADCConfigFilenameCert = "filename_cert"
	ADCConfigFilenameKey  = "filename_key"
	ADCConfigPathSSL      = "path_ssl"
	ADCConfigCertName     = "cert_name"
)

type CitrixADC struct {
	appliances.Appliance

	client *netscaler.NitroClient
}

func (adc *CitrixADC) Init() error {
	// Validate Connection Details
	client, err := netscaler.NewNitroClientFromParams(netscaler.NitroParams{
		Url:       adc.URL,
		Username:  adc.GetUsername(),
		Password:  adc.GetPassword(),
		SslVerify: adc.ValidateCerts,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create nitro client")
	}
	adc.client = client
	// Validate that all settings are set
	return adc.EnsureKeys(ADCConfigFilenameCert, ADCConfigFilenameKey, ADCConfigPathSSL, ADCConfigCertName)
}
