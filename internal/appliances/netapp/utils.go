package netapp

import (
	"bytes"
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

func (na *NetappAppliance) patchProtocolsS3(u interface{}) error {
	// Because we need the SVM UUID to update the certificate, check first
	if na.SVMUUID == nil {
		return errors.New("failed to update s3 certificate because we don't have a SVM UUID")
	}

	resp, err := na.req("PATCH", fmt.Sprintf("/api/protocols/s3/services/%s", *na.SVMUUID), u)
	if err != nil {
		return errors.Wrap(err, "failed to send request to rest API")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to update SVM S3 certificate: %s", responseData)
	}
	return nil
}

func (na *NetappAppliance) req(method string, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", na.URL, path)
	na.Logger.Tracef("%s '%s'", method, url)

	var bodyReader io.Reader
	if body != nil {
		jsonValue, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse request to json")
		}
		bodyReader = bytes.NewBuffer(jsonValue)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create get request")
	}
	req.SetBasicAuth(na.Username, na.Password)
	req.Header.Add("Accept", "application/json")
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := na.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		responseData, err := ioutil.ReadAll(resp.Body)
		na.Logger.Trace(string(responseData))
		if err != nil {
			return resp, errors.Wrap(err, "failed to read response body")
		}
	}
	return resp, nil
}
