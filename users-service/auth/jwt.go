package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type (
	JwtWrapper struct {
		SecretKey       string
		Issuer          string
		ExpirationHours int64
	}
	JwtClaim struct {
		UserType int
		Username string
		jwt.StandardClaims
	}
)

func (j *JwtWrapper) CreateJwToken(username string, userType int) (string, error) {
	claims := &JwtClaim{
		UserType: userType,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Issuer:    j.Issuer,
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JwtWrapper) CheckToken(providedToken string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		providedToken,
		&JwtClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, Ok := token.Claims.(*JwtClaim)
	if !Ok {
		return nil, errors.New("Failed to parse claims")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("JWT is expired")
	}
	return claims, nil
}
