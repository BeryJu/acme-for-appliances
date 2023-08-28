package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type SynologyAPI struct {
	URL      string
	Username string
	Password string
	Logger   *log.Entry
	Client   *http.Client

	sid string
}

type SynologyAPIResponse struct {
	Success bool `json:"success"`
}

func (dsm *SynologyAPI) makeRequest(method string, path string, params map[string]string, body io.Reader) (*http.Request, error) {
	b, err := url.Parse(dsm.URL)
	if err != nil {
		return nil, err
	}
	b.Path = path
	q := b.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	if dsm.sid != "" {
		q.Set("_sid", dsm.sid)
	}
	b.RawQuery = q.Encode()
	req, err := http.NewRequest(method, b.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	return req, err
}

// APIFile is a struct which contains a file's name, its type and its data.
type APIFile struct {
	Name string
	Type string
	Data []byte
}

func (dsm *SynologyAPI) makeFileRequest(method string, path string, params map[string]string, files ...APIFile) (*http.Request, error) {
	var (
		buf = new(bytes.Buffer)
		w   = multipart.NewWriter(buf)
	)

	for _, f := range files {
		part, err := w.CreateFormFile(f.Type, filepath.Base(f.Name))
		if err != nil {
			return nil, err
		}

		_, err = part.Write(f.Data)
		if err != nil {
			return nil, err
		}
	}

	err := w.Close()
	if err != nil {
		return nil, err
	}

	req, err := dsm.makeRequest(method, path, params, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	return req, nil
}
