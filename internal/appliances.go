package internal

import (
	"log"
	"strings"

	"github.com/BeryJu/acme-for-appliances/internal/appliances"
	"github.com/BeryJu/acme-for-appliances/internal/appliances/citrixadc"
	"github.com/BeryJu/acme-for-appliances/internal/appliances/netapp"
	"github.com/BeryJu/acme-for-appliances/internal/appliances/vmwarevsphere"
)

func GetActual(a *appliances.Appliance) appliances.CertificateConsumer {
	switch strings.ToLower(a.Type) {
	case "netapp":
		return &netapp.NetappAppliance{
			Appliance: *a,
		}
	case "citrix_adc":
		return &citrixadc.CitrixADC{
			Appliance: *a,
		}
	case "vmware_vcenter":
		return &vmwarevsphere.VMwareVsphere{
			Appliance: *a,
		}
	default:
		log.Fatalf("Invalid appliance type %s", strings.ToLower(a.Type))
	}
	return nil
}
