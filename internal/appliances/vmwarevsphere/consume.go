package vmwarevsphere

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BeryJu/acme-for-appliances/internal/appliances"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2/json"
)

func (v *VMwareVsphere) Consume(c *certificate.Resource) error {
	r := &vmwareVsphereTLS{
		Spec: vmwareVsphereTLSSpec{
			Cert:     appliances.MainCertOnly(c),
			Key:      string(c.PrivateKey),
			RootCert: string(c.IssuerCertificate) + v.Extension[VSphereConfigRootCAName].(string),
		},
	}
	jsonValue, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "failed to parse request to json")
	}
	fullURL := fmt.Sprintf("%s/rest/vcenter/certificate-management/vcenter/tls", v.URL)
	req, err := http.NewRequest("PUT", fullURL, bytes.NewBuffer(jsonValue))
	req.Header.Add(VMwareVsphereSessionHeader, v.sessionID)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to create certificate: %s", responseData)
	}
	return nil
}
