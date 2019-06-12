package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TylerGrey/lotte_server/lib/model"
	"github.com/TylerGrey/lotte_server/lib/redis"
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken ...
func GenerateToken(jwtParameter *model.Jwt, expireTime int) (string, error) {
	mySigningKey := []byte("TylerGrey")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JwtClaims{
		UserID: jwtParameter.UserID,
		Email:  jwtParameter.Email,
		Name:   jwtParameter.Name,
		Role:   jwtParameter.Role,
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

// RequireTokenAuthentication ...
func RequireTokenAuthentication(ctx context.Context, tokenstring string) (context.Context, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte("TylerGrey"), nil
	})

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 세션정보에서 있는지 체크한다.
			redisReplConn := redis.MasterConnect()
			defer redisReplConn.Close()

			value, err := redisReplConn.GetValue(claims["email"].(string))
			if err != nil {
				return ctx, err
			}

			if strings.Compare(value, tokenstring) != 0 {
				ctx = context.WithValue(ctx, "userId", claims["userId"])
				ctx = context.WithValue(ctx, "email", claims["email"])
				ctx = context.WithValue(ctx, "name", claims["name"])
				ctx = context.WithValue(ctx, "role", claims["role"])
				ctx = context.WithValue(ctx, "token", token)
				ctx = context.WithValue(ctx, "error", "4005")
				return ctx, err
			}

			ctx = context.WithValue(ctx, "userId", claims["userId"])
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "name", claims["name"])
			ctx = context.WithValue(ctx, "role", claims["role"])
			ctx = context.WithValue(ctx, "token", token)

			return ctx, nil
		} else {
			ctx = context.WithValue(ctx, "error", "4000")
			return ctx, errors.New("jwt parse error")
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("thnat's not token", jwt.ValidationErrorMalformed)
			ctx = context.WithValue(ctx, "error", "4003")
			return ctx, errors.New("인증 양식 에러")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			fmt.Println("Timing is everything", jwt.ValidationErrorExpired)
			ctx = context.WithValue(ctx, "error", "4002")
			return ctx, errors.New("인증 시간 만료")
		} else {
			fmt.Println("Couldn't handle this token:", err)
			ctx = context.WithValue(ctx, "error", "4003")
			return ctx, errors.New("인증 토큰 에러")
		}
	}

	ctx = context.WithValue(ctx, "error", "4000")
	return ctx, errors.New("unkown error")
}
