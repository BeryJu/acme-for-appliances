package vmwarevsphere

type vmwareVsphereSessionResponse struct {
	Value string `json:"value"`
}

// Body used to set the certificate
type vmwareVsphereTLS struct {
	Spec vmwareVsphereTLSSpec `json:"spec"`
}
type vmwareVsphereTLSSpec struct {
	Cert     string `json:"cert"`
	Key      string `json:"key"`
	RootCert string `json:"root_cert"`
}

// Response when checking the current certificate
type vmwareVsphereTLSResponse struct {
	Value vmwareVsphereTLSResponseValue `json:"value"`
}
type vmwareVsphereTLSResponseValue struct {
	ValidTo string `json:"valid_to"`
}
