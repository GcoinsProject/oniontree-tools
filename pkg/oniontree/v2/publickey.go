package oniontree

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"unicode"
)

type PublicKey struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	UserID      string `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Value       string `json:"value" yaml:"value"`
}

func ParseKey(key []byte) (PublicKey, error) {
	keyClean := bytes.TrimLeftFunc(key, unicode.IsSpace)
	el, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(keyClean))
	if err != nil {
		return PublicKey{}, err
	}

	publicKey := PublicKey{}
	for _, e := range el {
		userID := ""
		for _, ident := range e.Identities {
			userID = ident.Name
		}
		pk := e.PrimaryKey
		publicKey = PublicKey{
			Value:       string(keyClean),
			ID:          pk.KeyIdString(),
			Fingerprint: fmt.Sprintf("%X", pk.Fingerprint),
			UserID:      userID,
		}
	}
	return publicKey, nil
}