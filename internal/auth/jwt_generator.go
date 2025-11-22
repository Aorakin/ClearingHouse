package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var RefreshSecret []byte
var PublicKey *rsa.PublicKey

func init() {
	secret := os.Getenv("REFRESH_SECRET")
	RefreshSecret = []byte(secret)

	publicKeyBytes, err := os.ReadFile("token_public.pem")
	if err != nil {
		panic(fmt.Errorf("failed to read public key: %w", err))
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		panic("invalid public key format")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("failed to parse public key: %w", err))
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("public key is not RSA")
	}

	PublicKey = rsaPub

}

// GenerateAccessToken - RS512, 15 min expiry
func GenerateAccessToken(userID, email, firstName, lastName string, privateKey *rsa.PrivateKey) (string, error) {
	claims := AccessClaims{
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "clearing-house-auth",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return token.SignedString(privateKey)
}

// GenerateRefreshToken - HMAC, 7 days expiry
func GenerateRefreshToken(userID string) (string, error) {
	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "clearing-house-auth",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(RefreshSecret)
}

func VerifyRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return RefreshSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*RefreshClaims), nil
}
