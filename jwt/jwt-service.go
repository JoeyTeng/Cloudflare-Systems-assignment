package main

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

func createJWT(w http.ResponseWriter, req *http.Request) {
	issuer := "joey.teng.dev"

	username := req.URL.Path
	if strings.HasPrefix(username, "/") {
		username = strings.Split(username[1:], "/")[0]
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
		Issuer:    issuer,
		Subject:   username,
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	ss, _ := token.SignedString(privateKey)

	cookie := http.Cookie{
		Name:     "token",
		Value:    ss,
		Expires:  time.Now().UTC().AddDate(0, 0, 2),
		MaxAge:   24 * 3600,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(privateKey.PublicKey.N.Bytes())
}

func main() {
	http.HandleFunc("/auth/", createJWT)
	// http.HandleFunc("/verify", verifyJWT)
	// http.HandleFunc("/README.txt", readme)
	// http.HandleFunc("/stats", stats)
	http.ListenAndServe(":8090", nil)
}
