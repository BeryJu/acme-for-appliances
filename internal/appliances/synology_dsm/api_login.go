package synology_dsm

import (
	"encoding/json"
	"net/http"
)

type SynologyAPILoginResponse struct {
	Data struct {
		Sid string `json:"sid"`
	}
	SynologyAPIResponse
}

func (dsm *SynologyDSM) login() error {
	if dsm.sid != "" {
		return nil
	}
	req, err := dsm.makeRequest(http.MethodGet, "webapi/auth.cgi", map[string]string{
		"api":     string(SynoAPIAuth),
		"version": "3",
		"method":  "login",
		"account": dsm.Username,
		"passwd":  dsm.Password,
		"format":  "sid",
	}, nil)
	if err != nil {
		return err
	}
	res, err := dsm.client.Do(req)
	if err != nil {
		return err
	}
	var b SynologyAPILoginResponse
	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return err
	}
	dsm.Logger.Info("successfully logged into DSM")
	dsm.sid = b.Data.Sid
	return nil
}
