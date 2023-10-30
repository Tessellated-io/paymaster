package entity

import (
	"github.com/tessellated-io/pickaxe/crypto"
	"github.com/tessellated-io/router/router"
)

// Entity defines an entity on a specific chain.
type Entity interface {
	Address() string
	PublicKey() string
}

// Private implementation of enttity
type entity struct {
	chain  router.Chain
	signer crypto.BytesSigner
}

// Ensure entity structs conform to Entity interface
var _ Entity = (*entity)(nil)

// NewEntity creates a new entity
func NewEntity(signer crypto.BytesSigner, chain router.Chain) (Entity, error) {
	return &entity{
		signer: signer,
		chain:  chain,
	}, nil
}

// Entity Interacted

func (e *entity) Address() string {
	prefix := e.chain.Bech32Prefix()
	return e.signer.GetAddress(prefix)
}

func (e *entity) PublicKey() string {
	return e.signer.GetPublicKey().String()
}

func (e *entity) ChainID() string {
	return e.chain.ChainID()
}
