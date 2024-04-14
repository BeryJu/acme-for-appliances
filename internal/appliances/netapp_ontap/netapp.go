package netapp_ontap

import (
	"encoding/json"
	"net/http"

	"beryju.io/acme-for-appliances/internal/appliances"
	"github.com/pkg/errors"
)

const (
	NetappConfigCertNameA = "cert_name_a"
	NetappConfigCertNameB = "cert_name_b"
	NetappConfigSVMName   = "svm_name"
)

// NetappDateLayout Parse time from response
// example: "2021-04-20T15:59:37+00:00"
const NetappDateLayout = "2006-01-02T15:04:05Z07:00"

type NetappAppliance struct {
	appliances.Appliance

	ActiveCertName  string
	PassiveCertName string

	SVMUUID *string

	client           *http.Client
	clusterName      string
	certIsForCluster bool
}

func (na *NetappAppliance) Init() error {
	// Get a client with required TLS Settings (skip cert check)
	na.client = na.HTTPClient()
	resp, err := na.req("GET", "/api/cluster", nil)
	if err != nil {
		return err
	}
	r := &ontapClusterInfo{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return errors.Wrap(err, "failed parse response")
	}
	na.clusterName = r.Name
	na.Logger.WithField("version", r.Version.Full).WithField("name", r.Name).Info("Successfully connected to cluster")
	if na.clusterName == na.Extension[NetappConfigSVMName].(string) {
		na.certIsForCluster = true
		na.Logger.Debug("Cert is for Cluster SVM")
	} else {
		na.certIsForCluster = false
	}
	return na.EnsureKeys(NetappConfigCertNameA, NetappConfigCertNameB, NetappConfigSVMName)
}
