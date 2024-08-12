package consensus

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/mr-tron/base58"
)

// Signer is a type that represents a signer that can sign messages
type Signer struct {
	privateKey *rsa.PrivateKey
}

// NewSigner creates a new signer with the given private key
func NewSigner(privateKey *rsa.PrivateKey) *Signer {
	return &Signer{privateKey: privateKey}
}

// Sign takes a message as input and returns a signature for the message
func (s *Signer) Sign(message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Verifier is a type that represents a verifier that can verify signatures
type Verifier struct {
	publicKey *rsa.PublicKey
}

// NewVerifier creates a new verifier with the given public key
func NewVerifier(publicKey *rsa.PublicKey) *Verifier {
	return &Verifier{publicKey: publicKey}
}

// Verify takes a message and signature as input and returns true if the signature is valid for the message
func (v *Verifier) Verify(message []byte, signature []byte) bool {
	hash := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(v.publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return false
	}
	return true
}

// KeyPair is a type that represents a public-private key pair
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// GenerateKeyPair generates a new public-private key pair
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey
	return &KeyPair{PrivateKey: privateKey, PublicKey: publicKey}, nil
}

// PublicKeyToString converts a public key to a string
func PublicKeyToString(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return base58.Encode(publicKeyPem), nil
}

// StringToPublicKey converts a string to a public key
func StringToPublicKey(publicKeyString string) (*rsa.PublicKey, error) {
	publicKeyPem, err := base58.Decode(publicKeyString)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKeyPem)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to cast public key to RSA public key")
	}
	return rsaPublicKey, nil
}
