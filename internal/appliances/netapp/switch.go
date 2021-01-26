package netapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// SwitchClusterCert Switch cluster's certificate to certificate with `uuid`
func (na *NetappAppliance) SwitchClusterCert(uuid string) error {
	r := &ontapCertificateUpdate{
		Certificate: ontapRecord{
			UUID: uuid,
		},
	}
	jsonValue, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "failed to parse request to json")
	}

	resp, err := na.req("PATCH", "/api/cluster", bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.Wrap(err, "failed to send request to rest API")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode != 202 {
		return fmt.Errorf("failed to update cluster certificate: %s", responseData)
	}
	return nil
}

func (na *NetappAppliance) SwitchSVMS3Cert(uuid string) error {
	// Because we need the SVM UUID to update the certificate, check first
	if na.SVMUUID == nil {
		return errors.New("failed to update s3 certificate because we don't have a SVM UUID")
	}

	r := &ontapCertificateUpdate{
		Certificate: ontapRecord{
			UUID: uuid,
		},
	}
	jsonValue, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "failed to parse request to json")
	}

	resp, err := na.req("PATCH", fmt.Sprintf("/api/protocols/s3/services/%s", *na.SVMUUID), bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.Wrap(err, "failed to send request to rest API")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to update SVM S3 certificate: %s", responseData)
	}
	return nil
}
