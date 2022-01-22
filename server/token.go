package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

func Token(w http.ResponseWriter, r *http.Request) {

	// take a jwt as input and verify the signature. If it's valid then
	// generate a new access token holding the claims provided plus new ones ?

	//user authenticated so grant them a token for 5 minutes
	token := TokenStruct{AccessToken: "", Expiry: time.Now().Add(time.Minute * 5)}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"Client":       "clientName",
		"nbf":          time.Date(1995, 04, 22, 12, 0, 0, 0, time.UTC).Unix(),
		"access_token": token,
	})
	jsonBytes, err := json.Marshal(jwtToken)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(jsonBytes))

	cookie := http.Cookie{
		Name:    "ACCESS_TOKEN",
		Value:   encoded,
		Domain:  "http://localhostcad",
		Expires: time.Now().Add(time.Minute * 10),
	}
	http.SetCookie(w, &cookie)
	if err != nil {
		panic(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
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
