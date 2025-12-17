package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecretKey = []byte("rahasia_super_aman_uas_backend")

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID, username, role string, permissions []string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":     userID,
		"username":    username,
		"role":        role,
		"permissions": permissions,
		"iat":         now.Unix(),                     // Issued At
		"exp":         now.Add(time.Hour * 24).Unix(), // Expires At
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecretKey)
}

func ExtractToken(authHeader string) string {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, errors.New("token expired")
		}
	}

	return claims, nil
}

// GetTokenExpiration mengambil expiration time dari JWT token
func GetTokenExpiration(tokenString string) (time.Time, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return time.Time{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}, errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		return time.Unix(int64(exp), 0), nil
	}

	return time.Time{}, errors.New("expiration not found in token")
}

// GetTokenIssuedAt mengambil issued at time dari JWT token
func GetTokenIssuedAt(tokenString string) (time.Time, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return time.Time{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}, errors.New("invalid token claims")
	}

	if iat, ok := claims["iat"].(float64); ok {
		return time.Unix(int64(iat), 0), nil
	}

	return time.Time{}, errors.New("issued at not found in token")
}
