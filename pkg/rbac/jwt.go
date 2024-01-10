package rbac

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

func JWTAuthValidate[USER database.TableWithID[IDTYPE], IDTYPE int64 | string](authHeader, secret string) (*USER, error) {
	if authHeader == "" {
		return nil, errors.New("Authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("Authorization header is missing Bearer prefix")
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		return nil, errors.New("Authorization header format is invalid")
	}

	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	// validate required claims
	for _, k := range []string{"user", "exp", "iat"} {
		if _, ok := claims[k]; !ok {
			return nil, errors.New("invalid token")
		}
	}

	exp := claims["exp"].(float64)
	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}

	user := new(USER)
	err = json.Unmarshal([]byte(claims["user"].(string)), user)

	return user, err
}

func GenerateJWT(user interface{}, duration time.Duration, secret string) (string, error) {
	// json encode user
	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": string(userJSON),
		"exp":  time.Now().Add(duration).Unix(),
		"iat":  time.Now(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
