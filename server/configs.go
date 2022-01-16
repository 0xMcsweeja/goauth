package main

// Core properties for oauth flow
type Config struct {
	ClientId     string
	ClientSecret string
	Endpoint     Endpoint
	RedirectUrl  string
	Scopes       []string
}

// the two urls for auth and token request and a flag for secret
type Endpoint struct {
	AuthUrl   string
	TokenUrl  string
	AuthStyle Authstyle
}

type Authstyle int

