package netapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/BeryJu/acme-for-appliances/internal/appliances"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
)

func (na *NetappAppliance) Consume(c *certificate.Resource) error {
	// First off, we need to check if the passive cert exists
	// (and attempt to delete it)
	passiveCert, err := na.getCertificateByName(na.PassiveCertName)
	if err != nil {
		// Failed to check certificate, fail early to ensure we don't error later
		return err
	}
	if passiveCert != nil {
		// Passive cert exists, lets try to delete it
		err := na.DeleteCert(passiveCert.UUID)
		if err != nil {
			// Don't fail if we're not successful
			na.Logger.WithError(err).Warning("Failed to deleted passive cert")
		}
	}

	cert, err := na.CreateCert(c)
	if err != nil {
		return err
	}
	na.Logger.WithField("certUUID", cert.UUID).Debug("Got UUID for new cert")

	// We've now successfully created the passive cert.
	// Now we need to switch either cluster or S3 over to the new cert.

	if na.certIsForCluster {
		err = na.SwitchClusterCert(cert.UUID)
	} else {
		err = na.SwitchSVMS3Cert(cert.UUID)
	}

	if err != nil {
		na.Logger.WithError(err).Warning("failed to switch cluster/svm certificate")
		return err
	}

	a := na.ActiveCertName
	b := na.PassiveCertName
	na.ActiveCertName = b
	na.PassiveCertName = a

	// Sleep a second, if we send too many requests without pause
	// we get a connection reset
	time.Sleep(time.Second * 1)

	// Now we've successfully swapped the certificates.
	// Let's the delete the (now passive) certificate

	passiveCert, err = na.getCertificateByName(na.PassiveCertName)
	if err != nil {
		na.Logger.WithError(err).Warning("failed to get passive cert")
		return nil
	}
	if passiveCert == nil {
		na.Logger.Info("couldn't find passive cert")
		return nil
	}
	err = na.DeleteCert(passiveCert.UUID)
	if err == nil {
		na.Logger.WithField("passive", passiveCert.Name).Info("Successfully deleted passive cert")
	}
	return err
}

func (na *NetappAppliance) CreateCert(c *certificate.Resource) (*ontapRecord, error) {
	// Create request body
	r := &ontapCertificatePOST{
		Name: na.PassiveCertName,
		SVM: ontapSVMSelector{
			Name: na.Extension[NetappConfigSVMName].(string),
		},
		PublicCertificate: appliances.MainCertOnly(c),
		PrivateKey:        string(c.PrivateKey),
		Type:              "server",
	}

	na.Logger.Info("Creating new certificate object")
	resp, err := na.req("POST", "/api/security/certificates?return_records=true", r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to rest API")
	}
	rec := &ontapCertificateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse response")
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("failed to create certificate")
	}
	return &rec.Records[0], nil
}
