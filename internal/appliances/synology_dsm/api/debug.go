package api

import (
	"fmt"
	"io"
	"net/http"
)

func (dsm *SynologyAPI) Debug() {
	req, _ := dsm.makeRequest(http.MethodGet, "webapi/query.cgi", map[string]string{
		"api":     "SYNO.API.Info",
		"version": "1",
		"method":  "query",
		"query":   "all",
	}, nil)
	res, _ := dsm.Client.Do(req)
	b, _ := io.ReadAll(res.Body)
	fmt.Println(string(b))
}
