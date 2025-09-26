package signature_helper

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"os"

	TicketDtos "github.com/ClearingHouse/internal/tickets/dtos"
)

func loadPrivateKey() (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile("private.pem")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParseECPrivateKey(block.Bytes)
}

func loadPublicKey() (*ecdsa.PublicKey, error) {
	data, err := os.ReadFile("public.pem")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*ecdsa.PublicKey), nil
}

func SignTicket(ticket TicketDtos.GliderTicket) (string, error) {
	priv, err := loadPrivateKey()
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(ticket)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)

	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return "", err
	}

	sig := append(r.Bytes(), s.Bytes()...)
	return base64.StdEncoding.EncodeToString(sig), nil
}

func VerifyTicket(ticket TicketDtos.GliderTicket, signature string) (bool, error) {
	pub, err := loadPublicKey()
	if err != nil {
		return false, err
	}

	data, err := json.Marshal(ticket)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	sigLen := len(sigBytes) / 2
	r := new(big.Int).SetBytes(sigBytes[:sigLen])
	s := new(big.Int).SetBytes(sigBytes[sigLen:])

	ok := ecdsa.Verify(pub, hash[:], r, s)
	return ok, nil
}
