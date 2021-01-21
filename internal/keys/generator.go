package keys

import "crypto"

type KeyGenerator interface {
	GetPrivateKey(name string) crypto.PrivateKey
}
