package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var (
	jwtMeanEncodingTime time.Duration = time.Duration(0)
	jwtMeanDecodingTime time.Duration = time.Duration(0)
	encodingTimes       int           = 0
	decodingTimes       int           = 0
	authTimes       int           = 0
	verifyTimes       int           = 0
)

func createJWT(w http.ResponseWriter, req *http.Request) {
	authTimes += 1
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

	_before := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, _ := token.SignedString(privateKey)
	_after := time.Now()

	jwtMeanEncodingTime = jwtMeanEncodingTime*time.Duration(encodingTimes) + (_after.Sub(_before))
	encodingTimes += 1
	jwtMeanEncodingTime /= time.Duration(encodingTimes)

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
	verifyTimes += 1
	cookies := req.Cookies()
	jwtCookie := ""
	for i := range cookies {
		if cookies[i].Name == "token" {
			jwtCookie = cookies[i].Value
		}
	}
	if len(jwtCookie) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("JWT token not found."))

		return
	}

	_before := time.Now()
	token, err := jwt.ParseWithClaims(jwtCookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	_after := time.Now()
	jwtMeanDecodingTime = jwtMeanDecodingTime*time.Duration(decodingTimes) + (_after.Sub(_before))
	decodingTimes += 1
	jwtMeanDecodingTime /= time.Duration(decodingTimes)

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

func readme(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "README.txt")
}

func stats(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]int64)
	data["Mean JWT Encode Time (us)"] = jwtMeanEncodingTime.Microseconds()
	data["Mean JWT Decode Time (us)"] = jwtMeanDecodingTime.Microseconds()
	data["Auth Times"] = int64(authTimes)
	data["Verify Times"] = int64(verifyTimes)

	json, _ := json.Marshal(data)
	w.Write(json)
}

func main() {
	http.HandleFunc("/auth/", createJWT)
	http.HandleFunc("/verify", verifyJWT)
	http.HandleFunc("/README.txt", readme)
	http.HandleFunc("/stats", stats)
	http.ListenAndServe(":8080", nil)
}
