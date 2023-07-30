package appliances

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"beryju.org/acme-for-appliances/internal/keys"
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

func (a *Appliance) HTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !a.ValidateCerts},
	}
	return &http.Client{Transport: tr}
}

func (a *Appliance) EnsureKeys(keys ...string) error {
	for _, key := range keys {
		if _, ok := a.Extension[key]; !ok {
			return fmt.Errorf("no value for setting %s set", key)
		}
	}
	return nil
}

func (a *Appliance) GetName() string {
	return a.Name
}

func (a *Appliance) GetDomains() []string {
	return a.Domains
}

func (a *Appliance) GetUsername() string {
	if strings.HasPrefix(a.Username, "env:") {
		envName := strings.Split(a.Username, "env:")
		a.Logger.WithField("env", envName[1]).Trace("Got username from env")
		return os.Getenv(envName[1])
	}
	return a.Username
}

func (a *Appliance) GetPassword() string {
	if strings.HasPrefix(a.Password, "env:") {
		envName := strings.Split(a.Password, "env:")
		a.Logger.WithField("env", envName[1]).Trace("Got password from env")
		return os.Getenv(envName[1])
	}
	return a.Password
}

func (a *Appliance) GetKeyGenerator(storageBase string) keys.KeyGenerator {
	return keys.NewECDSAKeyGenerator(storageBase)
}

type CertificateConsumer interface {
	Init() error
	CheckExpiry() (int, error)
	Consume(*certificate.Resource) error
	GetName() string
	GetDomains() []string
	GetKeyGenerator(storageBase string) keys.KeyGenerator
}
