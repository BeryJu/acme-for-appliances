package appliances

import (
	"strings"

	"github.com/go-acme/lego/v4/certificate"
)

// MainCertOnly return the main certificate only, without any chain
func MainCertOnly(c *certificate.Resource) string {
	chain := string(c.Certificate)
	ca := string(c.IssuerCertificate)
	return strings.ReplaceAll(chain, ca, "")
}
