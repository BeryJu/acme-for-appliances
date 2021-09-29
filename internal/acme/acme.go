package acme

import (
	"beryju.org/acme-for-appliances/internal/appliances"
	"beryju.org/acme-for-appliances/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	llog "github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
)

type Client struct {
	*lego.Client
}

func NewClient(u *User) *Client {
	llog.Logger = log.WithField("component", "acme")
	lc := lego.NewConfig(u)
	lc.CADirURL = config.C.ACME.DirectoryURL
	lc.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(lc)
	if err != nil {
		log.Fatal(err)
	}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{
		TermsOfServiceAgreed: config.C.ACME.TermsAgreed,
	})
	if err != nil {
		log.Fatal(err)
	}
	u.Registration = reg

	return &Client{
		client,
	}
}

func (c *Client) GetCerts(app appliances.CertificateConsumer) (*certificate.Resource, error) {
	provider, err := dns.NewDNSChallengeProviderByName(config.C.ACME.ChallengeProviderName)
	if err != nil {
		log.Fatal(err)
	}

	opts := []dns01.ChallengeOption{}

	if len(config.C.ACME.Resolvers) > 0 {
		log.WithField("resolvers", config.C.ACME.Resolvers).Debug("Using custom resolvers")
		opts = append(opts, dns01.AddRecursiveNameservers(config.C.ACME.Resolvers))
	}

	err = c.Challenge.SetDNS01Provider(provider, opts...)
	if err != nil {
		log.Fatal(err)
	}

	pk := app.GetKeyGenerator(config.C.Storage).GetPrivateKey(app.GetName())
	request := certificate.ObtainRequest{
		Domains:    app.GetDomains(),
		Bundle:     false,
		PrivateKey: pk,
	}
	return c.Certificate.Obtain(request)
}
