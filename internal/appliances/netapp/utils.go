package netapp

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

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
