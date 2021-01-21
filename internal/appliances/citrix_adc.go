package appliances

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/chiradeep/go-nitro/config/ssl"
	"github.com/chiradeep/go-nitro/config/system"
	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
)

const ADCConfigFilenameCert = "filename_cert"
const ADCConfigFilenameKey = "filename_key"
const ADCConfigPathSSL = "path_ssl"
const ADCConfigCertName = "cert_name"

type CitrixADC struct {
	Appliance

	client *netscaler.NitroClient
}

func (adc *CitrixADC) CheckExpiry() (int, error) {
	// TODO
	return -1, nil
}

func (adc *CitrixADC) Init() error {
	// Validate Connection Details
	client, err := netscaler.NewNitroClientFromParams(netscaler.NitroParams{
		Url:       adc.URL,
		Username:  adc.Username,
		Password:  adc.Password,
		SslVerify: adc.ValidateCerts,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create nitro client")
	}
	adc.client = client
	// Validate that all settings are set
	return adc.ensureKeys(ADCConfigFilenameCert, ADCConfigFilenameKey, ADCConfigPathSSL, ADCConfigCertName)
}

func (adc *CitrixADC) CreateOrUpdateFile(f system.Systemfile) error {
	adc.Logger.WithField("filename", f.Filename).Debug("checking if file exists")
	mf, err := adc.client.FindResourceArrayWithParams(netscaler.FindParams{
		ResourceType: netscaler.Systemfile.Type(),
		ArgsMap: map[string]string{
			"filename":     url.PathEscape(f.Filename),
			"filelocation": url.PathEscape(adc.Extension[ADCConfigPathSSL].(string)),
		},
	})
	if err != nil {
		return err
	}
	if len(mf) > 0 {
		adc.Logger.WithField("filename", f.Filename).Debug("cert exists, deleting it")
		err := adc.client.DeleteResourceWithArgsMap(netscaler.Systemfile.Type(), f.Filename, map[string]string{
			"filelocation": url.PathEscape(adc.Extension[ADCConfigPathSSL].(string)),
		})
		if err != nil {
			// We don't abort here if an error occurs
			// to ensure that *some* certificate exists
			adc.Logger.WithError(err).Warning("error during file deletion")
		}
	}
	_, err = adc.client.AddResource(netscaler.Systemfile.Type(), "", &f)
	return err
}

func (adc *CitrixADC) CreateOrUpdateCert(s ssl.Sslcertkey) error {
	_, err := adc.client.AddResource(netscaler.Sslcertkey.Type(), adc.Extension[ADCConfigCertName].(string), &s)
	if err != nil {
		_, err := adc.client.ChangeResource(netscaler.Sslcertkey.Type(), adc.Extension[ADCConfigCertName].(string), &s)
		if err != nil {
			return errors.Wrap(err, "failed to update cert")
		}
	}
	return nil
}

func (adc *CitrixADC) Consume(c *certificate.Resource) error {
	certFile := system.Systemfile{
		Filename:     adc.Extension[ADCConfigFilenameCert].(string),
		Filelocation: adc.Extension[ADCConfigPathSSL].(string),
		Filecontent:  base64.StdEncoding.EncodeToString([]byte(MainCertOnly(c))),
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
		Certkey: adc.Extension[ADCConfigCertName].(string),
		Cert:    fmt.Sprintf("%s%s", adc.Extension[ADCConfigPathSSL].(string), adc.Extension[ADCConfigFilenameCert].(string)),
		Key:     fmt.Sprintf("%s%s", adc.Extension[ADCConfigPathSSL].(string), adc.Extension[ADCConfigFilenameKey].(string)),
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
