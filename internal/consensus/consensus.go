package consensus

import (
	"crypto/sha256"
	"errors"

	"github.com/cloudflare/circl/sign"
	"github.com/cloudflare/circl/sign/schemes"
)

// Consensus is the implementation of Hierarchical Authorization Consensus (HAC).
type Consensus struct {
	ElectionCommissionPubKey *sign.PublicKey
	ReturningOfficerPubKey   *sign.PublicKey
	PollingOfficerPubKey     *sign.PublicKey

	ElectionCommissionPrivateKey *sign.PrivateKey
	ReturningOfficerPrivateKey   *sign.PrivateKey
	PollingOfficerPrivateKey     *sign.PrivateKey
}

// NewConsensus returns a new instance of Consensus with the given public keys.
func NewConsensus(ecPubKey, roPubKey, poPubKey *sign.PublicKey) *Consensus {
	return &Consensus{
		ElectionCommissionPubKey: ecPubKey,
		ReturningOfficerPubKey:   roPubKey,
		PollingOfficerPubKey:     poPubKey,
	}
}

// AuthorizeWithElectionCommission signs the data with the Election Commission's private key.
func (c *Consensus) AuthorizeWithElectionCommission(data []byte) ([]byte, error) {
	if c.ElectionCommissionPrivateKey == nil {
		return nil, errors.New("Election Commission's public key is not set")
	}

	digest := sha256.New()
	digest.Write(data)

	mode := schemes.ByName("Ed25519")

	return mode.Sign(*c.ElectionCommissionPrivateKey, digest.Sum(nil), nil), nil
}

// AuthorizeWithReturningOfficer signs the data with the Returning Officer's private key.
func (c *Consensus) AuthorizeWithReturningOfficer(data []byte) ([]byte, error) {
	if c.ReturningOfficerPrivateKey == nil {
		return nil, errors.New("Election Commission's public key is not set")
	}

	digest := sha256.New()
	digest.Write(data)

	mode := schemes.ByName("Ed25519")

	return mode.Sign(*c.ReturningOfficerPrivateKey, digest.Sum(nil), nil), nil
}

// AuthorizeWithPollingOfficer signs the data with the Polling Officer's private key.
func (c *Consensus) AuthorizeWithPollingOfficer(data []byte) ([]byte, error) {
	if c.ReturningOfficerPrivateKey == nil {
		return nil, errors.New("Election Commission's public key is not set")
	}

	digest := sha256.New()
	digest.Write(data)

	mode := schemes.ByName("Ed25519")

	return mode.Sign(*c.ReturningOfficerPrivateKey, digest.Sum(nil), nil), nil
}

// VerifySignature verifies the signature of the data using the provided public key.
func (c *Consensus) VerifySignature(data, signature []byte, pubKey *sign.PublicKey) bool {
	digest := sha256.New()
	digest.Write(data)
	mode := schemes.ByName("Ed25519")
	return mode.Verify(*pubKey, digest.Sum(nil), signature, nil)

}
