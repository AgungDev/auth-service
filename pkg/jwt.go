package pkg

import (
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub   string        `json:"sub"`
	Iss   string        `json:"iss"`
	Exp   int64         `json:"exp"`
	Roles []interface{} `json:"roles"`
	jwt.RegisteredClaims
}

func LoadPrivateKey(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func GenerateJWTRS256(sub, iss string, roles []interface{}, privateKeyPath string, expireMinutes int) (string, error) {
	keyData, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return "", err
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return "", err
	}
	claims := Claims{
		Sub:   sub,
		Iss:   iss,
		Exp:   time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireMinutes) * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privKey)
}
