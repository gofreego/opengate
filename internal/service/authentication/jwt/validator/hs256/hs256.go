package hs256

import (
	"github.com/gofreego/goutils/customerrors"
	"github.com/golang-jwt/jwt/v5"
)

type HS256Validator struct {
	secretKey string
}

func NewHS256Validator(secretKey string) *HS256Validator {
	return &HS256Validator{secretKey: secretKey}
}

func (v *HS256Validator) UpdateSecretKey(secretKey string) {
	v.secretKey = secretKey
}

func (v *HS256Validator) Validate(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, customerrors.BAD_REQUEST_ERROR("invalid signing method for jwt token, need HMAC")
		}
		return []byte(v.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, customerrors.BAD_REQUEST_ERROR("invalid jwt token")

}
