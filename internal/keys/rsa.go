package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path"

	"beryju.io/acme-for-appliances/internal/storage"
	log "github.com/sirupsen/logrus"
)

type RSAKeyGenerator struct {
	log         *log.Entry
	storageBase string
}

func NewRSAKeyGenerator(storageBase string) *RSAKeyGenerator {
	return &RSAKeyGenerator{
		storageBase: storageBase,
		log:         log.WithField("component", "rsa-generator"),
	}
}

func (e *RSAKeyGenerator) GetPrivateKey(name string) crypto.PrivateKey {
	keyPath := path.Join(storage.PathPrefix(e.storageBase), fmt.Sprintf("%s.pem", name))
	exists, err := storage.FileExists(keyPath)
	if err != nil {
		e.log.WithError(err).Warning("failed to read key")
		return nil
	}
	if !exists {
		k, err := GenerateKeyAndSaveRSA(keyPath)
		if err != nil {
			e.log.WithError(err).Warning("failed to save key")
		}
		e.log.Info("successfully saved new appliance private key")
		return k
	}
	key, err := LoadRSA(keyPath)
	if err != nil {
		e.log.WithError(err).Warning("failed to load key")
	}
	e.log.Debug("successfully loaded appliance private key")
	return key
}

func GenerateRSAKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}

func GenerateKeyAndSaveRSA(path string) (crypto.PrivateKey, error) {
	key := GenerateRSAKey()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o755)
	if err != nil {
		log.WithField("path", path).Warning("failed to open file")
		return key, err
	}
	privkeyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		log.WithField("path", path).Warning("failed to marshal key")
		return key, err
	}
	err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privkeyBytes,
	})
	if err != nil {
		log.WithField("path", path).Warning("failed to encode to PEM")
		return key, err
	}
	return key, nil
}

func LoadRSA(path string) (crypto.PrivateKey, error) {
	f, err := os.Open(path)
	if err != nil {
		log.WithField("path", path).Debug("Failed to open file")
		return GenerateRSAKey(), err
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		log.WithField("path", path).Debug("Failed to read")
		return GenerateRSAKey(), err
	}
	block, _ := pem.Decode(contents)
	if block == nil {
		log.WithField("path", path).Debug("Failed to pem decode")
		return GenerateRSAKey(), err
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.WithField("path", path).Debug("Failed to parse pkcs8 key")
		return GenerateRSAKey(), err
	}
	return priv, nil
}
