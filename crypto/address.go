package crypto

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func PubKeyToAddress(publicKey cryptotypes.PubKey, prefix string) (string, error) {
	rawAddress := sdk.AccAddress(publicKey.Address())
	address, err := bech32.ConvertAndEncode(prefix, rawAddress)
	if err != nil {
		return "", err
	}

	return address, nil
}
