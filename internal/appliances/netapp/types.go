package netapp

type ontapSVMSelector struct {
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type ontapClusterInfo struct {
	Version struct {
		Full string `json:"full"`
	} `json:"version"`
}

type ontapCertificatePOST struct {
	IntermediateCertificates []string         `json:"intermediate_certificates,omitempty"`
	Name                     string           `json:"name"`
	PublicCertificate        string           `json:"public_certificate"`
	PrivateKey               string           `json:"private_key"`
	Type                     string           `json:"type"`
	SVM                      ontapSVMSelector `json:"svm"`
}

type ontapCertificateResponse struct {
	Records []struct {
		UUID       string `json:"uuid"`
		ExpiryTime string `json:"expiry_time"`
	} `json:"records"`
	NumRecords int `json:"num_records"`
}
