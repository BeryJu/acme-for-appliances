package keys

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path"

	"beryju.io/acme-for-appliances/internal/storage"
	log "github.com/sirupsen/logrus"
)

type ECDSAKeyGenerator struct {
	log         *log.Entry
	storageBase string
}

func NewECDSAKeyGenerator(storageBase string) *ECDSAKeyGenerator {
	return &ECDSAKeyGenerator{
		storageBase: storageBase,
		log:         log.WithField("component", "ecdsa-generator"),
	}
}

func (e *ECDSAKeyGenerator) GetPrivateKey(name string) crypto.PrivateKey {
	keyPath := path.Join(storage.PathPrefix(e.storageBase), fmt.Sprintf("%s.pem", name))
	exists, err := storage.FileExists(keyPath)
	if err != nil {
		e.log.WithError(err).Warning("failed to read key")
		return nil
	}
	if !exists {
		k, err := GenerateKeyAndSaveECDSA(keyPath)
		if err != nil {
			e.log.WithError(err).Warning("failed to save key")
		}
		e.log.Info("successfully saved new appliance private key")
		return k
	}
	key, err := LoadECDSA(keyPath)
	if err != nil {
		e.log.WithError(err).Warning("failed to load key")
	}
	e.log.Debug("successfully loaded appliance private key")
	return key
}

func GenerateECDSAKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}

func GenerateKeyAndSaveECDSA(path string) (crypto.PrivateKey, error) {
	key := GenerateECDSAKey()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.WithField("path", path).Warning("failed to open file")
		return key, err
	}
	privkeyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.WithField("path", path).Warning("failed to marshal key")
		return key, err
	}
	err = pem.Encode(f, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privkeyBytes,
	})
	if err != nil {
		log.WithField("path", path).Warning("failed to encode to PEM")
		return key, err
	}
	return key, nil
}

func LoadECDSA(path string) (crypto.PrivateKey, error) {
	f, err := os.Open(path)
	if err != nil {
		log.WithField("path", path).Debug("Failed to open file")
		return GenerateECDSAKey(), err
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		log.WithField("path", path).Debug("Failed to read")
		return GenerateECDSAKey(), err
	}
	block, _ := pem.Decode(contents)
	if block == nil {
		log.WithField("path", path).Debug("Failed to pem decode")
		return GenerateECDSAKey(), err
	}
	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		log.WithField("path", path).Debug("Failed to parse pkcs8 key")
		return GenerateECDSAKey(), err
	}
	return priv, nil
}
