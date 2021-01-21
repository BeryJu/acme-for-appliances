package appliances

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
)

const NetappConfigCertName = "cert_name"
const NetappConfigSVMName = "svm_name"

// NetappDateLayout Parse time from response
// example: "2021-04-20T15:59:37+00:00"
const NetappDateLayout = "2006-01-02T15:04:05Z07:00"

type NetappAppliance struct {
	Appliance
	CertUUID *string
	client   *http.Client
}

type ontapSVMSelector struct {
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type ontapClusterInfo struct {
	Version struct {
		Full string `json:"full"`
	} `json:"version"`
}

type ontapCertificatePOST struct {
	IntermediateCertificates []string         `json:"intermediate_certificates,omitempty"`
	Name                     string           `json:"name"`
	PublicCertificate        string           `json:"public_certificate"`
	PrivateKey               string           `json:"private_key"`
	Type                     string           `json:"type"`
	SVM                      ontapSVMSelector `json:"svm"`
}

type ontapCertificateResponse struct {
	Records []struct {
		UUID       string `json:"uuid"`
		ExpiryTime string `json:"expiry_time"`
	} `json:"records"`
	NumRecords int `json:"num_records"`
}

func (na *NetappAppliance) Init() error {
	// Get a client with required TLS Settings (skip cert check)
	na.client = na.httpClient()
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
	return na.ensureKeys(NetappConfigCertName, NetappConfigSVMName)
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

func (na *NetappAppliance) DeleteOldCert() error {
	url := fmt.Sprintf("%s/api/security/certificates/%s", na.URL, *na.CertUUID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create deletion request")
	}
	req.SetBasicAuth(na.Username, na.Password)
	resp, err := na.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to delete certificate")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode != 200 {
		return errors.New(string(responseData))
	}
	return nil
}

func (na *NetappAppliance) Consume(c *certificate.Resource) error {
	// Create request body
	r := &ontapCertificatePOST{
		Name: na.Extension[NetappConfigCertName].(string),
		SVM: ontapSVMSelector{
			Name: na.Extension[NetappConfigSVMName].(string),
		},
		PublicCertificate: MainCertOnly(c),
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
