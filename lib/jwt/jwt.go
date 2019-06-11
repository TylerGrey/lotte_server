package jwt

import (
	"time"

	"github.com/TylerGrey/lotte_server/lib/model"
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken ...
func GenerateToken(jwtParameter *model.Jwt, expireTime int) (string, error) {
	mySigningKey := []byte("TylerGrey")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JwtClaims{
		UserID:    jwtParameter.UserID,
		Email:     jwtParameter.Email,
		FirstName: jwtParameter.FirstName,
		LastName:  jwtParameter.LastName,
		Role:      jwtParameter.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expireTime)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
