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
		"desc":    existing.Desc,
	}
	if existing.ID != "" {
		params["id"] = existing.ID
	}
	files := []APIFile{
		{
			Field: "key",
			Data:  key,
		},
		{
			Field: "cert",
			Data:  cert,
		},
	}
	if len(intermediate) != 0 {
		files = append(files, APIFile{
			Field: "inter_cert",
			Data:  intermediate,
		})
	}

	req, err := dsm.makeFileRequest(http.MethodPost, "webapi/entry.cgi", params, files...)
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
