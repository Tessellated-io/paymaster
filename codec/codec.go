package codec

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/codec"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core/types"
)

// Provides a singleton codec that can be used across the application

// Mutex to initialize the codec exactly once
var initRelayerOnce sync.Once

// The codec
var cdc *codec.ProtoCodec = nil

func GetCodec() *codec.ProtoCodec {
	initRelayerOnce.Do(func() {
		interfaceRegistry := codectypes.NewInterfaceRegistry()
		ibc.RegisterInterfaces(interfaceRegistry)
		authtypes.RegisterInterfaces(interfaceRegistry)
		banktypes.RegisterInterfaces(interfaceRegistry)
		cryptotypes.RegisterInterfaces(interfaceRegistry)
		cdc = codec.NewProtoCodec(interfaceRegistry)
	})

	return cdc
}
