package vmware_vsphere

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

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
	d := time.Until(t)
	return int(d.Hours() / 24), nil
}
