package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

func createJWT(w http.ResponseWriter, req *http.Request) {
	keys := "jwt-key"
	issuer := "joey.teng.dev"

	username := req.URL.Path
	if strings.HasPrefix(username, "/") {
		username = strings.Split(username[1:], "/")[1]
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
		Issuer:    issuer,
		Subject:   username,
	}

	// privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	privateKeyBytes, _ := os.ReadFile(keys)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		panic(err)
	}
	publicKeyBytes, _ := os.ReadFile(keys + ".public.pem")

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	ss, _ := token.SignedString(privateKey)

	cookie := http.Cookie{
		Name:     "token",
		Value:    ss,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   24 * 3600,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	w.Header().Set("Content-Type", "text/plain")
	http.SetCookie(w, &cookie)

	w.Write(publicKeyBytes)
}

func verifyJWT(w http.ResponseWriter, req *http.Request) {
	cookies := req.Cookies()
	jwtCookie := ""
	for i := range cookies {
		if (cookies[i].Name == "token"){
			jwtCookie = cookies[i].Value
		}
	}
	if (len(jwtCookie) == 0) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("JWT token not found."))

		return
	}
	token, err := jwt.ParseWithClaims(jwtCookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(claims.Subject))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Invalid JWT claims: "))
		w.Write([]byte(err.Error()))
	}
}

func main() {
	http.HandleFunc("/auth/", createJWT)
	http.HandleFunc("/verify", verifyJWT)
	// http.HandleFunc("/README.txt", readme)
	// http.HandleFunc("/stats", stats)
	http.ListenAndServe(":8080", nil)
}
