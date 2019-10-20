package service

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type PublicKey struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	UserID      string `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Value       string `json:"value" yaml:"value"`
}

func ParseKey(key []byte) (PublicKey, error) {
	block, err := armor.Decode(bytes.NewReader(key))
	if err != nil {
		return PublicKey{}, err
	}

	if block.Type != openpgp.PublicKeyType {
		return PublicKey{}, errors.New("not a public key")
	}

	pktReader := packet.NewReader(block.Body)
	pkt, err := pktReader.Next()
	if err != nil {
		return PublicKey{}, err
	}

	pk, ok := pkt.(*packet.PublicKey)
	if !ok {
		return PublicKey{}, errors.New("invalid public key")
	}

	pkt, err = pktReader.Next()
	if err != nil {
		return PublicKey{}, err
	}

	userID, ok := pkt.(*packet.UserId)
	if !ok {
		return PublicKey{}, errors.New("invalid user id data")
	}

	publicKey := PublicKey{
		Value:       string(key),
		ID:          pk.KeyIdString(),
		Fingerprint: fmt.Sprintf("%X", pk.Fingerprint),
		UserID:      userID.Id,
	}
	return publicKey, nil
}
