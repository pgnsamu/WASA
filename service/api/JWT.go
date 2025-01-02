package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Chiave segreta per firmare il token
var jwtKey = []byte("secret-key")

// Funzione per generare un JWT
func GenerateJWT(username string, id int) (string, error) {

	// Definizione dei claims
	claims := jwt.MapClaims{
		"username": username,
		"id":       id,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Creazione del token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Creazione del token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Parsing del token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Controlla il metodo di firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// Estrai i claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Conversione di claims["id"] da float64 a int
		if idFloat, ok := claims["id"].(float64); ok {
			claims["id"] = int(idFloat)
		}

		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
