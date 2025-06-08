package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/LikhithMar14/management/models"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("ysdfousadfdfr-2sdfsdfdsfsdf")

func validateCredentials(creds models.Credentials) error {
	if creds.Username == "" || creds.Password == "" {
		return errors.New("username and password are required")
	}
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I am in the login handler")
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	valid := validateCredentials(creds)
	if valid != nil {
		http.Error(w, valid.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := generateToken(creds)

	if err != nil {
		http.Error(w, "Failed to generate token",http.StatusInternalServerError	)
		return
	}

	response := models.LoginResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func generateToken(creds models.Credentials) (string, error) {
	expiration := time.Now().Add(24*time.Hour)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiration),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject: creds.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}