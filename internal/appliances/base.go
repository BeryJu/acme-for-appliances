package appliances

import (
	"crypto"
	"crypto/tls"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/BeryJu/acme-for-appliances/internal/keys"
	"github.com/BeryJu/acme-for-appliances/internal/storage"
	"github.com/go-acme/lego/v4/certificate"
	log "github.com/sirupsen/logrus"
)

type Appliance struct {
	Name          string
	Domains       []string
	Type          string
	URL           string
	Username      string
	Password      string
	ValidateCerts bool
	Extension     map[string]interface{}
	Logger        *log.Entry
}

func NewAppliance() *Appliance {
	return &Appliance{
		ValidateCerts: true,
		Extension:     make(map[string]interface{}),
	}
}

func (a *Appliance) GetActual() CertificateConsumer {
	switch strings.ToLower(a.Type) {
	case "netapp":
		return &NetappAppliance{
			Appliance: *a,
		}
	case "citrix_adc":
		return &CitrixADC{
			Appliance: *a,
		}
	default:
		log.Fatalf("Invalid appliance type %s", strings.ToLower(a.Type))
	}
	return nil
}

func (a *Appliance) httpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !a.ValidateCerts},
	}
	return &http.Client{Transport: tr}
}

func (a *Appliance) ensureKeys(keys ...string) error {
	for _, key := range keys {
		if _, ok := a.Extension[key]; !ok {
			return fmt.Errorf("no value for setting %s set", key)
		}
	}
	return nil
}

func (a *Appliance) GetPrivateKey() crypto.PrivateKey {
	keyPath := path.Join(storage.PathPrefix(), fmt.Sprintf("%s.pem", a.Name))
	exists, err := storage.FileExists(keyPath)
	if err != nil {
		a.Logger.WithError(err).Warning("failed to read key")
		return nil
	}
	if !exists {
		k, err := keys.GenerateKeyAndSaveECDSA(keyPath)
		if err != nil {
			a.Logger.WithError(err).Warning("failed to save key")
		}
		a.Logger.Info("successfully saved new appliance private key")
		return k
	}
	key, err := keys.LoadECDSA(keyPath)
	if err != nil {
		a.Logger.WithError(err).Warning("failed to load key")
	}
	a.Logger.Info("successfully loaded appliance private key")
	return key
}

type CertificateConsumer interface {
	Init() error
	CheckExpiry() (int, error)
	Consume(*certificate.Resource) error
}
