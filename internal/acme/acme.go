package acme

import (
	log "github.com/sirupsen/logrus"

	"github.com/BeryJu/acme-for-appliances/internal/appliances"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	llog "github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
	"github.com/spf13/viper"
)

type Client struct {
	*lego.Client
}

func NewClient(u *User) *Client {
	llog.Logger = log.WithField("component", "acme")
	config := lego.NewConfig(u)
	config.CADirURL = viper.GetString("acme.directory_url")
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{
		TermsOfServiceAgreed: viper.GetBool("acme.terms_agreed"),
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
	provider, err := dns.NewDNSChallengeProviderByName(viper.GetString("acme.challenge_provider_name"))
	if err != nil {
		log.Fatal(err)
	}

	resolvers := viper.GetStringSlice("acme.resolvers")
	opts := make([]dns01.ChallengeOption, 0)

	if len(resolvers) > 1 {
		opts = append(opts, dns01.AddRecursiveNameservers(resolvers))
	}

	err = c.Challenge.SetDNS01Provider(provider, opts...)
	if err != nil {
		log.Fatal(err)
	}

	pk := app.GetKeyGenerator().GetPrivateKey(app.GetName())
	request := certificate.ObtainRequest{
		Domains:    app.GetDomains(),
		Bundle:     false,
		PrivateKey: pk,
	}
	return c.Certificate.Obtain(request)
}
