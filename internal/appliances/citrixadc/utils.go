package citrixadc

import (
	"net/url"

	"github.com/chiradeep/go-nitro/config/ssl"
	"github.com/chiradeep/go-nitro/config/system"
	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/pkg/errors"
)

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
