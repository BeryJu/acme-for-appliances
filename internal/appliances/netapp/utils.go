package netapp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// getCertificateByName Get a single Certificate's info by Name (includes name and expiry)
func (na *NetappAppliance) getCertificateByName(name string) (*ontapRecord, error) {
	values := url.Values{}
	values.Add("name", name)
	if na.certIsForCluster {
		values.Add("scope", "cluster")
	} else {
		values.Add("svm.name", na.Extension[NetappConfigSVMName].(string))
	}
	values.Add("fields", "name,uuid,expiry_time")
	resp, err := na.req("GET", fmt.Sprintf("/api/security/certificates?%s", values.Encode()), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to rest API")
	}
	r := &ontapCertificateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse response")
	}
	if r.NumRecords < 1 {
		return nil, nil
	}
	if r.NumRecords > 1 {
		return nil, fmt.Errorf("expected to get 1 certificate, but got %d", r.NumRecords)
	}
	return &r.Records[0], nil
}

// DeleteCert Delete certificate by uuid
func (na *NetappAppliance) DeleteCert(uuid string) error {
	resp, err := na.req("DELETE", fmt.Sprintf("/api/security/certificates/%s", uuid), nil)
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

func (na *NetappAppliance) req(method string, path string, body io.Reader) (*http.Response, error) {
	// To ensure the credential work, we check /cluster
	url := fmt.Sprintf("%s%s", na.URL, path)
	na.Logger.Tracef("%s '%s'", method, url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create get request")
	}
	req.SetBasicAuth(na.Username, na.Password)
	if method == "POST" || method == "PATCH" {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := na.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
