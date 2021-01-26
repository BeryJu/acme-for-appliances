package netapp

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// SwitchClusterCert Switch cluster's certificate to certificate with `uuid`
func (na *NetappAppliance) SwitchClusterCert(uuid string) error {
	r := &ontapCertificateUpdate{
		Certificate: ontapRecord{
			UUID: uuid,
		},
	}

	resp, err := na.req("PATCH", "/api/cluster", r)
	if err != nil {
		return errors.Wrap(err, "failed to send request to rest API")
	}
	if resp.StatusCode != 202 {
		return fmt.Errorf("failed to update cluster certificate")
	}
	return nil
}

func (na *NetappAppliance) SwitchSVMS3Cert(uuid string) error {
	// Because we need the SVM UUID to update the certificate, check first
	if na.SVMUUID == nil {
		return errors.New("failed to update s3 certificate because we don't have a SVM UUID")
	}

	err := na.patchProtocolsS3(ontapSVMServiceUpdate{
		Enabled: false,
	})
	if err != nil {
		na.Logger.WithError(err).Warning("failed to disable S3")
		return err
	}
	na.Logger.Info("successfully disabled S3, waiting")

	time.Sleep(time.Second * 5)

	err = na.patchProtocolsS3(ontapCertificateUpdate{
		Certificate: ontapRecord{
			UUID: uuid,
		},
	})
	if err != nil {
		na.Logger.WithError(err).Warning("failed to update cert")
		// Don't return here, we still need to enable S3
	} else {
		na.Logger.Info("successfully replaced certificate")
	}

	time.Sleep(time.Second * 5)

	err = na.patchProtocolsS3(ontapSVMServiceUpdate{
		Enabled: true,
	})
	if err == nil {
		na.Logger.Info("successfully re-enabled s3")
	}
	return err
}
