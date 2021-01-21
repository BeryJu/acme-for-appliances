package netapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BeryJu/acme-for-appliances/internal/appliances"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
)

func (na *NetappAppliance) Consume(c *certificate.Resource) error {
	// Create request body
	r := &ontapCertificatePOST{
		Name: na.Extension[NetappConfigCertName].(string),
		SVM: ontapSVMSelector{
			Name: na.Extension[NetappConfigSVMName].(string),
		},
		PublicCertificate: appliances.MainCertOnly(c),
		PrivateKey:        string(c.PrivateKey),
		Type:              "server",
	}
	jsonValue, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "failed to parse request to json")
	}

	var url string
	// If we don't have a certificate UUID from above, assume we have to create a new one
	if na.CertUUID != nil {
		na.Logger.Info("Deleting old certificate before installing new certificate")
		err := na.DeleteOldCert()
		if err != nil {
			return errors.Wrap(err, "failed to delete old certificate")
		}
	}
	url = fmt.Sprintf("%s/api/security/certificates", na.URL)
	na.Logger.Info("Creating new certificate object")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return errors.Wrap(err, "failed to create post request")
	}
	req.SetBasicAuth(na.Username, na.Password)

	resp, err := na.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request to rest API")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create certificate: %s", responseData)
	}
	return nil
}
