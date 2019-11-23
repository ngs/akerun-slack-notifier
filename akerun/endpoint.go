package akerun

import (
	"os"

	"golang.org/x/oauth2"
)

// Endpoint Akerun endpoint
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://api.akerun.com/oauth/authorize/",
	TokenURL:  "https://api.akerun.com/oauth/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

// Config Akerun Config
var Config = oauth2.Config{
	Endpoint:     Endpoint,
	ClientID:     os.Getenv("AKERUN_CLIENT_ID"),
	ClientSecret: os.Getenv("AKERUN_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("AKERUN_REDIRECT_URL"),
}
