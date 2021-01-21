package acme

import (
	"crypto"
	"path"

	"github.com/BeryJu/acme-for-appliances/internal/keys"
	"github.com/BeryJu/acme-for-appliances/internal/storage"
	"github.com/go-acme/lego/v4/registration"
	log "github.com/sirupsen/logrus"
)

const UserKeyName = "user_key.pem"

type User struct {
	Email        string
	Registration *registration.Resource
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u User) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	l := log.WithField("component", "user")
	// Check if we have a key at all
	fullPath := path.Join(storage.PathPrefix(), UserKeyName)
	exists, err := storage.FileExists(fullPath)
	if !exists {
		l.Info("Key does not exist, creating a new key")
		k, err := keys.GenerateKeyAndSaveECDSA(fullPath)
		if err != nil {
			l.WithError(err).Warning("failed to save key")
		}
		return k
	}
	if err != nil {
		l.WithError(err).Warning("failed to stat keys")
	}
	key, err := keys.LoadECDSA(fullPath)
	if err != nil {
		l.WithError(err).Warning("failed to load key")
	}
	l.Info("successfully loaded user's private key")
	return key
}
