package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type State struct {
	Domains []string `json:"domains"`
}

func (s *State) CompareDomains(domains []string) bool {
	if len(s.Domains) != len(domains) {
		return false
	}
	for i, v := range s.Domains {
		if v != domains[i] {
			return false
		}
	}
	return true
}

func (s *State) Write(base string, name string) {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log.WithError(err).Warning("Failed to write state")
		return
	}
	err = os.WriteFile(path.Join(PathPrefix(base), fmt.Sprintf("%s.json", name)), data, 0o644)
	if err != nil {
		log.WithError(err).Warning("Failed to write state")
	}
}

func GetState(base string, name string) *State {
	jsonFile, err := os.Open(path.Join(PathPrefix(base), fmt.Sprintf("%s.json", name)))
	if err != nil {
		log.WithError(err).Warning("Failed to read state")
		return &State{}
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.WithError(err).Warning("failed to close state")
		}
	}()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.WithError(err).Warning("Failed to read state")
		return &State{}
	}
	var state State
	err = json.Unmarshal(byteValue, &state)
	if err != nil {
		log.WithError(err).Warning("failed to parse state")
		return &State{}
	}
	return &state
}
