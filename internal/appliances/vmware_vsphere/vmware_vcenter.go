package vmware_vsphere

import (
	"fmt"
	"net/http"
	"time"

	"beryju.io/acme-for-appliances/internal/appliances"
	"beryju.io/acme-for-appliances/internal/config"
	"beryju.io/acme-for-appliances/internal/keys"
	"gopkg.in/square/go-jose.v2/json"
)

const VMwareVsphereSessionHeader = "vmware-api-session-id"

// VMwareVsphereDateLayout Parse time from response
// example: "2022-10-07T21:50:42.000Z"
const VMwareVsphereDateLayout = "2006-01-02T15:04:05Z07:00"

const VSphereConfigRootCAName = "root_ca"

type VMwareVsphere struct {
	appliances.Appliance

	client    *http.Client
	sessionID string
}

func (v *VMwareVsphere) Init() error {
	v.client = v.HTTPClient()
	// Set an extra long timeout since all the services have to restart
	v.client.Timeout = time.Minute * 5

	fullURL := fmt.Sprintf("%s/rest/com/vmware/cis/session", v.URL)
	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(v.GetUsername(), v.GetPassword())
	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	var respBody vmwareVsphereSessionResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}
	v.Logger.Info("Successfully authenticated to vCenter")
	v.sessionID = respBody.Value
	return v.EnsureKeys(VSphereConfigRootCAName)
}

func (v *VMwareVsphere) GetKeyGenerator(storageBase string) keys.KeyGenerator {
	return keys.NewRSAKeyGenerator(config.C.Storage)
}
