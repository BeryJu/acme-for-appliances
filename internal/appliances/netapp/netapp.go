package netapp

import (
	"encoding/json"
	"fmt"
	"net/http"

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
