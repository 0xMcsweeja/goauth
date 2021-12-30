package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/rand"
	"net/http"
	"strings"
)

var dummyUser = sha256.Sum256([]byte("admin"))
var dummyPassword = sha256.Sum256([]byte("password"))

func Authorize(w http.ResponseWriter, r *http.Request) {
	if !extractLogin(r) {
		errors.New("couldn't validate caller auth credentials")
	}
	dummyCode := sha256.Sum256([]byte(string(rand.Intn(999999))))
	slice := dummyCode[1 : len(dummyCode)-1]
	x := string(slice)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("code:" + x))
}

func extractLogin(r *http.Request) bool {
	header := r.Header.Get("authorization")

	split := strings.Split(header, " ")[1]
	if strings.Split(header, " ")[0] != "Basic" {
		errors.New("couldn't parse auth header")
	}

	decoded, err := base64.StdEncoding.DecodeString(split)
	if err != nil {
		panic(err)
	}

	creds := strings.Split(string(decoded), ":")
	user, pass := creds[0], creds[1]

	hashedUser := sha256.Sum256([]byte(user))
	hashedPass := sha256.Sum256([]byte(pass))

	return compareHash(hashedUser, dummyUser) && compareHash(hashedPass, dummyPassword)

}

func compareHash(x, y [32]byte) bool {
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
