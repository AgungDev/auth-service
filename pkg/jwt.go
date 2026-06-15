package pkg

import (
	"errors"
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

func LoadPublicKey(path string) ([]byte, error) {
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

func ValidateJWTRS256(tokenString string, publicKeyPath string) (*Claims, error) {
	keyData, err := LoadPublicKey(publicKeyPath)
	if err != nil {
		return nil, err
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}
	parsedToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
