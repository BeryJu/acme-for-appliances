package netapp

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

func (na *NetappAppliance) CheckExpiry() (int, error) {
	// First off, we need to figure out which certificate is currently active.
	var activeCertUUID *string
	if na.certIsForCluster {
		activeCertUUID = na.GetClusterCertificate()
		na.Logger.Debug("Got certificate UUID for cluster")
	} else {
		// if cert is not for cluster, I assume it's used for SVM S3
		activeCertUUID = na.GetS3Certificate()
		na.Logger.Debug("Got certificate UUID for S3")
	}
	if activeCertUUID == nil {
		na.Logger.Warning("Failed to get UUID for certificate")
		return -1, nil
	}

	// Actually check the certificates status
	resp, err := na.req("GET", fmt.Sprintf("/api/security/certificates/%s", *activeCertUUID), nil)
	r := &ontapRecord{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return -1, errors.Wrap(err, "failed parse response")
	}

	if r.Name == NetappConfigCertNameA {
		na.ActiveCertName = NetappConfigCertNameA
		na.PassiveCertName = NetappConfigCertNameB
	} else {
		na.ActiveCertName = NetappConfigCertNameB
		na.PassiveCertName = NetappConfigCertNameA
	}

	na.Logger.WithField("active", na.ActiveCertName).Debug("Found active cert")

	t, err := time.Parse(NetappDateLayout, r.ExpiryTime)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to parse expiry time")
	}
	d := t.Sub(time.Now())
	return int(d.Hours() / 24), nil
}

// GetClusterCertificate Get certificate UUID of the cluster
func (na *NetappAppliance) GetClusterCertificate() *string {
	resp, err := na.req("GET", "/api/cluster", nil)
	if err != nil {
		return nil
	}
	r := &ontapClusterInfo{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil
	}
	return &r.Certificate.UUID
}

// GetS3Certificate Get UUID of certificate used for S3 Protocol
func (na *NetappAppliance) GetS3Certificate() *string {
	values := url.Values{}
	values.Add("svm.name", na.Extension[NetappConfigSVMName].(string))
	values.Add("fields", "certificate.uuid,svm.uuid")
	resp, err := na.req("GET", fmt.Sprintf("/api/protocols/s3/services?%s", values.Encode()), nil)
	if err != nil {
		return nil
	}
	r := &ontapSVMS3Info{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil
	}
	if len(r.Records) < 1 {
		return nil
	}
	na.SVMUUID = &r.Records[0].SVM.UUID
	return &r.Records[0].Certificate.UUID
}
