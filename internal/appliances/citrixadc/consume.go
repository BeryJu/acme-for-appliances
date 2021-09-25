package citrixadc

import (
	"encoding/base64"
	"fmt"

	"beryju.org/acme-for-appliances/internal/appliances"
	"github.com/chiradeep/go-nitro/config/ssl"
	"github.com/chiradeep/go-nitro/config/system"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
)

func (adc *CitrixADC) Consume(c *certificate.Resource) error {
	certFile := system.Systemfile{
		Filename:     adc.Extension[ADCConfigFilenameCert].(string),
		Filelocation: adc.Extension[ADCConfigPathSSL].(string),
		Filecontent:  base64.StdEncoding.EncodeToString([]byte(appliances.MainCertOnly(c))),
		Fileencoding: "BASE64",
	}
	err := adc.CreateOrUpdateFile(certFile)
	if err != nil {
		return errors.Wrap(err, "failed to upload certificate file")
	}

	keyFile := system.Systemfile{
		Filename:     adc.Extension[ADCConfigFilenameKey].(string),
		Filelocation: adc.Extension[ADCConfigPathSSL].(string),
		Filecontent:  base64.StdEncoding.EncodeToString([]byte(c.PrivateKey)),
		Fileencoding: "BASE64",
	}

	err = adc.CreateOrUpdateFile(keyFile)
	if err != nil {
		return errors.Wrap(err, "failed to upload private key file")
	}

	certKey := ssl.Sslcertkey{
		Certkey:       adc.Extension[ADCConfigCertName].(string),
		Cert:          fmt.Sprintf("%s%s", adc.Extension[ADCConfigPathSSL].(string), adc.Extension[ADCConfigFilenameCert].(string)),
		Key:           fmt.Sprintf("%s%s", adc.Extension[ADCConfigPathSSL].(string), adc.Extension[ADCConfigFilenameKey].(string)),
		Nodomaincheck: true,
	}
	err = adc.CreateOrUpdateCert(certKey)
	if err != nil {
		return errors.Wrap(err, "failed to create cert")
	}

	err = adc.client.SaveConfig()
	if err != nil {
		return errors.Wrap(err, "failed to save config")
	}

	return nil
}
