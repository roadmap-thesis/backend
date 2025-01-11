package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/HotPotatoC/roadmap_gen/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(claims map[string]any) (string, error) {
	claims["exp"] = time.Now().Add(config.JWTSecretExpiresIn()).Unix()
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func VerifyJWT(token string) (*jwt.Token, jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.JWTSecretKey()), nil
	})
	if err != nil {
		return nil, nil, err
	}

	return t, t.Claims.(jwt.MapClaims), nil
}
