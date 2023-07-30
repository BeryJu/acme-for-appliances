package internal

import (
	"log"
	"strings"

	"beryju.org/acme-for-appliances/internal/appliances"
	"beryju.org/acme-for-appliances/internal/appliances/citrix_adc"
	"beryju.org/acme-for-appliances/internal/appliances/netapp_ontap"
	"beryju.org/acme-for-appliances/internal/appliances/vmware_vsphere"
)

func GetActual(a *appliances.Appliance) appliances.CertificateConsumer {
	switch strings.ToLower(a.Type) {
	case "netapp":
		return &netapp_ontap.NetappAppliance{
			Appliance: *a,
		}
	case "citrix_adc":
		return &citrix_adc.CitrixADC{
			Appliance: *a,
		}
	case "vmware_vcenter":
		return &vmware_vsphere.VMwareVsphere{
			Appliance: *a,
		}
	default:
		log.Fatalf("Invalid appliance type %s", strings.ToLower(a.Type))
	}
	return nil
}
