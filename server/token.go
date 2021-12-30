package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

func Token(w http.ResponseWriter, r *http.Request) {

	isValid := validateCredentials(r)

	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Credentials"))
	}

	//user authenticated so grant them a token for 5 minutes
	token := TokenStruct{AccessToken: "", Expiry: time.Now().Add(time.Minute * 5)}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"Client":       "clientName",
		"nbf":          time.Date(1995, 04, 22, 12, 0, 0, 0, time.UTC).Unix(),
		"access_token": token,
	})

	jsonBytes, err := json.Marshal(jwtToken)
	if err != nil {
		panic(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func validateCredentials(r *http.Request) bool {
	header := r.Header.Get("authorization")
	split := strings.Split(header, " ")[1]
	decoded, err := base64.StdEncoding.DecodeString(split)
	if err != nil {
		panic(err)
	}

	creds := strings.Split(string(decoded), ":")
	user, pass := creds[0], creds[1]

	if user == "admin" && pass == "password" {
		return true
	} else {
		return false
	}

}

// the definition for the tokens used in flow
type TokenStruct struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

func (t *TokenStruct) isValid() bool {
	if t != nil && t.AccessToken != "" && !t.Expiry.IsZero() {
		return true
	} else {
		return false
	}
}
