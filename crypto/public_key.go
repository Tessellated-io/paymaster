package crypto

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func Secp256k1PublicKeyFromBytes(bytes []byte) cryptotypes.PubKey {
	return &secp256k1.PubKey{
		Key: bytes,
	}
}
