package api

import (
	"encoding/json"
	"net/http"
)

type SynologyAPICert struct {
	ID        string       `json:"id"`
	Desc      string       `json:"desc"`
	ValidTill SynologyDate `json:"valid_till"`
}

type SynologyAPICertListResponse struct {
	Data struct {
		Certificates []SynologyAPICert `json:"certificates"`
	}
	SynologyAPIResponse
}

func (dsm *SynologyAPI) ListCertificates() (*SynologyAPICertListResponse, error) {
	req, err := dsm.makeRequest(http.MethodGet, "webapi/entry.cgi", map[string]string{
		"api":     string(SynoAPICoreCertCRT),
		"version": "1",
		"method":  "list",
	}, nil)
	if err != nil {
		return nil, err
	}
	res, err := dsm.Client.Do(req)
	if err != nil {
		return nil, err
	}
	var r SynologyAPICertListResponse
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (dsm *SynologyAPI) UploadCertificate(existing *SynologyAPICert, cert, intermediate, key []byte) (*SynologyAPIResponse, error) {
	params := map[string]string{
		"api":     string(SynoAPICoreCert),
		"version": "1",
		"method":  "import",
	}
	if existing != nil {
		params["id"] = existing.ID
		params["desc"] = existing.Desc
	}
	req, err := dsm.makeFileRequest(http.MethodPost, "webapi/entry.cgi", params,
		APIFile{
			Name: "key",
			Type: "application/x-iwork-keynote-sffkey",
			Data: key,
		},
		APIFile{
			Name: "inter_cert",
			Type: "application/pkix-cert",
			Data: intermediate,
		},
		APIFile{
			Name: "cert",
			Type: "application/pkix-cert",
			Data: cert,
		},
	)
	if err != nil {
		return nil, err
	}
	res, err := dsm.Client.Do(req)
	if err != nil {
		return nil, err
	}
	var r SynologyAPIResponse
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
