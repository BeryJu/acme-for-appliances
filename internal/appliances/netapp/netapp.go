package netapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/BeryJu/acme-for-appliances/internal/appliances"
	"github.com/pkg/errors"
)

const NetappConfigCertName = "cert_name"
const NetappConfigSVMName = "svm_name"

// NetappDateLayout Parse time from response
// example: "2021-04-20T15:59:37+00:00"
const NetappDateLayout = "2006-01-02T15:04:05Z07:00"

type NetappAppliance struct {
	appliances.Appliance
	CertUUID *string
	client   *http.Client
}

func (na *NetappAppliance) Init() error {
	// Get a client with required TLS Settings (skip cert check)
	na.client = na.HTTPClient()
	// To ensure the credential work, we check /cluster
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/cluster", na.URL), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create get request")
	}
	req.SetBasicAuth(na.Username, na.Password)
	resp, err := na.client.Do(req)
	if err != nil {
		return err
	}
	r := &ontapClusterInfo{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return errors.Wrap(err, "failed parse response")
	}
	na.Logger.WithField("version", r.Version.Full).Info("Successfully connected to cluster")
	return na.EnsureKeys(NetappConfigCertName, NetappConfigSVMName)
}

func (na *NetappAppliance) CheckExpiry() (int, error) {
	values := url.Values{}
	values.Add("name", na.Extension[NetappConfigCertName].(string))
	values.Add("svm.name", na.Extension[NetappConfigSVMName].(string))
	values.Add("fields", "name,uuid,expiry_time")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/security/certificates?%s", na.URL, values.Encode()), nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create get request")
	}
	req.SetBasicAuth(na.Username, na.Password)

	resp, err := na.client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "failed to send request to rest API")
	}
	r := &ontapCertificateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return 0, errors.Wrap(err, "failed parse response")
	}
	if r.NumRecords > 1 {
		return 0, fmt.Errorf("expected to get 1 certificate, but got %d", r.NumRecords)
	}
	na.CertUUID = &r.Records[0].UUID

	t, err := time.Parse(NetappDateLayout, r.Records[0].ExpiryTime)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to parse expiry time")
	}
	d := t.Sub(time.Now())
	return int(d.Hours() / 24), nil
}
