package netapp_ontap

type ontapSVMSelector struct {
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type ontapClusterInfo struct {
	Name        string      `json:"name"`
	Certificate ontapRecord `json:"certificate"`
	Version     struct {
		Full string `json:"full"`
	} `json:"version"`
}

type ontapSVMServiceUpdate struct {
	Enabled bool `json:"enabled"`
}

type ontapCertificateUpdate struct {
	Certificate ontapRecord `json:"certificate"`
}

type ontapSVMS3Info struct {
	Records []struct {
		SVM         ontapSVMSelector `json:"svm"`
		Certificate ontapRecord      `json:"certificate"`
	} `json:"records"`
	NumRecords int `json:"num_records"`
}

type ontapCertificatePOST struct {
	IntermediateCertificates []string         `json:"intermediate_certificates,omitempty"`
	Name                     string           `json:"name"`
	PublicCertificate        string           `json:"public_certificate"`
	PrivateKey               string           `json:"private_key"`
	Type                     string           `json:"type"`
	SVM                      ontapSVMSelector `json:"svm"`
}

type ontapRecord struct {
	Name       string `json:"name,omitempty"`
	UUID       string `json:"uuid"`
	ExpiryTime string `json:"expiry_time,omitempty"`
}

type ontapCertificateResponse struct {
	Records    []ontapRecord `json:"records"`
	NumRecords int           `json:"num_records"`
}
