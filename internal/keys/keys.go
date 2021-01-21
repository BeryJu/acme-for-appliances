package keys

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

func GenerateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}

func GenerateKeyAndSave(path string) (crypto.PrivateKey, error) {
	key := GenerateKey()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return key, err
	}
	privkeyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return key, err
	}
	err = pem.Encode(f, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privkeyBytes,
	})
	if err != nil {
		return key, err
	}
	return key, nil
}

func Load(path string) (crypto.PrivateKey, error) {
	f, err := os.Open(path)
	if err != nil {
		return GenerateKey(), err
	}
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return GenerateKey(), err
	}
	block, _ := pem.Decode([]byte(contents))
	if block == nil {
		return GenerateKey(), err
	}
	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return GenerateKey(), err
	}
	return priv, nil
}
