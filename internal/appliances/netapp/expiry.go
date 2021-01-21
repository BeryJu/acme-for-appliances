package netapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

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
