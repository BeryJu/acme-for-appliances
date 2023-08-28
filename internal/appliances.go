package internal

import (
	"log"
	"strings"

	"beryju.io/acme-for-appliances/internal/appliances"
	"beryju.io/acme-for-appliances/internal/appliances/citrix_adc"
	"beryju.io/acme-for-appliances/internal/appliances/netapp_ontap"
	"beryju.io/acme-for-appliances/internal/appliances/synology_dsm"
	"beryju.io/acme-for-appliances/internal/appliances/vmware_vsphere"
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
	case "synology_dsm":
		return &synology_dsm.SynologyDSM{
			Appliance: *a,
		}
	default:
		log.Fatalf("Invalid appliance type %s", strings.ToLower(a.Type))
	}
	return nil
}
