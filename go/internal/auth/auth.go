package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("12345")

// hash a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// check if a password is correct
func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// generate a JWT token
func GenerateJWT(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"exp":      JWTExpiration().Unix(),
		"iat":      time.Now().Unix(),
		"username": username,
		"user_id":  userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// parse a JWT token
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// genearate a refresh token
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// helper functions to set expiration times
func RefreshTokenExpiration() time.Time {
	return time.Now().Add(30 * 24 * time.Hour)
}
func JWTExpiration() time.Time {
	return time.Now().Add(time.Minute * 15)
}

// parse the user id from an expired JWT token
func ParseID(tokenStr string) (uint, error) {
	claims := jwt.MapClaims{}

	// parse the token without verifying
	_, _, err := new(jwt.Parser).ParseUnverified(tokenStr, claims)
	if err != nil {
		return 0, err
	}

	// Safely extract the user ID
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found or invalid type")
	}

	return uint(userIDFloat), nil
}

// extract the JWT from the request header
func ExtractJWT(r *http.Request) (string, error) {
	// extract the token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("missing token")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return "", fmt.Errorf("missing token")
	}
	return tokenString, nil
}
