package vmware_esxi

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http/cookiejar"
	"strings"

	"beryju.io/acme-for-appliances/internal/appliances"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type VMwareESXi struct {
	appliances.Appliance

	Client *vim25.Client
}

func (v *VMwareESXi) Init() error {
	j, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return err
	}
	hc := v.HTTPClient()
	hc.Jar = j

	u, err := soap.ParseURL(v.URL)
	if err != nil {
		return fmt.Errorf("failed to parse ESXi URL: %v", err)
	}

	// Create vim25 client
	soapClient := soap.NewClient(u, true)
	soapClient.Client = *hc

	vimClient, err := vim25.NewClient(context.Background(), soapClient)
	if err != nil {
		return fmt.Errorf("failed to create vim25 client: %v", err)
	}

	v.Client = vimClient
	// Create login request
	req := types.Login{
		This:     *v.Client.ServiceContent.SessionManager,
		UserName: v.Username,
		Password: v.Password,
		Locale:   "en-US",
	}

	// Execute login
	resp, err := methods.Login(context.Background(), v.Client, &req)
	if err != nil {
		return fmt.Errorf("login request failed: %v", err)
	}

	if resp.Returnval.Key != "" {
		v.Logger.Infof("Successfully logged in as: %s (%s)\n",
			resp.Returnval.FullName, resp.Returnval.UserName)
		return nil
	}
	return nil
}

func (v *VMwareESXi) Obtain(c *lego.Client, storage string) (*certificate.Resource, error) {
	res, err := methods.GenerateCertificateSigningRequest(context.Background(), v.Client, &types.GenerateCertificateSigningRequest{
		This: types.ManagedObjectReference{
			Type:  "HostCertificateManager",
			Value: "ha-certificate-manager",
		},
		UseIpAddressAsCommonName: false,
	})
	if err != nil {
		return nil, err
	}
	pcsr, _ := pem.Decode([]byte(res.Returnval))
	csr, err := x509.ParseCertificateRequest(pcsr.Bytes)
	if err != nil {
		return nil, err
	}
	// Let's encrypt doesn't want trailing domain dots
	for i, dns := range csr.DNSNames {
		csr.DNSNames[i] = strings.TrimSuffix(dns, ".")
	}
	csr.Subject.CommonName = strings.TrimSuffix(csr.Subject.CommonName, ".")
	request := certificate.ObtainForCSRRequest{
		Bundle: false,
		CSR:    csr,
	}
	return c.Certificate.ObtainForCSR(request)
}
