package appliances

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BeryJu/acme-for-appliances/internal/keys"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2/json"
)

const VMwareVsphereSessionHeader = "vmware-api-session-id"

// VMwareVsphereDateLayout Parse time from response
// example: "2022-10-07T21:50:42.000Z"
const VMwareVsphereDateLayout = "2006-01-02T15:04:05Z07:00"

const VSphereConfigRootCAName = "root_ca"

type vmwareVsphereSessionResponse struct {
	Value string `json:"value"`
}

// Body used to set the certificate
type vmwareVsphereTLS struct {
	Spec vmwareVsphereTLSSpec `json:"spec"`
}
type vmwareVsphereTLSSpec struct {
	Cert     string `json:"cert"`
	Key      string `json:"key"`
	RootCert string `json:"root_cert"`
}

// Response when checking the current certificate
type vmwareVsphereTLSResponse struct {
	Value vmwareVsphereTLSResponseValue `json:"value"`
}
type vmwareVsphereTLSResponseValue struct {
	ValidTo string `json:"valid_to"`
}

type VMwareVsphere struct {
	Appliance

	client    *http.Client
	sessionID string
}

func (v *VMwareVsphere) Init() error {
	v.client = v.httpClient()
	// Set an extra long timeout since all the services have to restart
	v.client.Timeout = time.Minute * 5

	fullURL := fmt.Sprintf("%s/rest/com/vmware/cis/session", v.URL)
	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(v.Username, v.Password)
	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	var respBody vmwareVsphereSessionResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}
	v.Logger.Info("Successfully authenticated to vCenter")
	v.sessionID = respBody.Value
	return v.ensureKeys(VSphereConfigRootCAName)
}

func (v *VMwareVsphere) CheckExpiry() (int, error) { // Create request body
	fullURL := fmt.Sprintf("%s/rest/vcenter/certificate-management/vcenter/tls", v.URL)
	req, err := http.NewRequest("GET", fullURL, nil)
	req.Header.Add(VMwareVsphereSessionHeader, v.sessionID)
	if err != nil {
		return -1, err
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return -1, err
	}
	var respBody vmwareVsphereTLSResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return -1, err
	}

	t, err := time.Parse(VMwareVsphereDateLayout, respBody.Value.ValidTo)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to parse expiry time")
	}
	d := t.Sub(time.Now())
	return int(d.Hours() / 24), nil
}

func (v *VMwareVsphere) GetKeyGenerator() keys.KeyGenerator {
	return keys.NewRSAKeyGenerator()
}

func (v *VMwareVsphere) Consume(c *certificate.Resource) error {
	r := &vmwareVsphereTLS{
		Spec: vmwareVsphereTLSSpec{
			Cert:     MainCertOnly(c),
			Key:      string(c.PrivateKey),
			RootCert: string(c.IssuerCertificate) + v.Extension[VSphereConfigRootCAName].(string),
		},
	}
	jsonValue, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "failed to parse request to json")
	}
	v.Logger.Debug(string(jsonValue))
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
