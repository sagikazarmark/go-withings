/*
This package demonstrates how OAuth2 authentication works in the Withings API.

Prerequisites

	- Withings account
	- Withings application
	- HTTP tunnel (ngrok, inlets, tunnelto, etc)

Withings account

The first step is creating a Withings account here: https://account.withings.com/connectionuser/account_create

You can use your personal account, but Withings recommends creating new accounts for production applications.

Read more about registration here: https://developer.withings.com/developer-guide/getting-started/register-to-withings-api#get-ready

HTTP tunnel

Application registration requires a redirect URL for the authentication to work,
but unfortunately localhost is not allowed, so you need a "real" URL.
For development purposes, the easiest solution is setting up some sort of tunnel that
proxies HTTP(S) requests to your application running on your machine.

Some well known solutions:
	- Ngrok: https://ngrok.com/
	- Tunnelto: https://tunnelto.dev/

Both solutions offer free tier, but they usually restrict you to non-reserved subdomains,
meaning every time you launch the tunnel you will get a randomly generated URL.
That's probably fine for quick tests, but having to always change the application redirect URL
can quickly become annoying, so you might want to subscribe to one of these solutions
(or host your own solution: https://github.com/anderspitman/awesome-tunneling).
The cheapest plan for Ngrok is $5/month, while Tunnelto is just $2/month.

Personally, I prefer Tunnelto, because it's cheap and ridiculuosly easy to use:
	tunnelto --subdomain withings-$USER-local --port 8080

But Ngrok isn't difficult either:
	ngrok http --subdomain withings-$USER-local 8080

Again, if you don't have a subscription, you might have to lose the `--subdomain`
part and get a randomly generated domain.

Either way, take note of your URL, because you are going to need it for the next step.

Read more about the callback URL here:
	- https://developer.withings.com/developer-guide/faq/faq#can-i-use-an-ip-adress-or-localhost-as-the-callback-url-or-redirect-uri
	- https://developer.withings.com/developer-guide/glossary/glossary/#callback-url

Withings application

Once you have an account and a callback URL ready,
you can register an application here: https://account.withings.com/partner/add_oauth2

Fill out the form with your details
(yes, even the logo is required, I just used my standard profile picture for dev purposes).

Make sure to fill out the "Callback URL" field with your tunnel URL + "/oauth2/callback".
For example `https://withings-YOURNAME-local.tunnelto.dev/oauth2/callback`.

Once you registered your application take note of the client ID and secret.

Running the application

Set the following environment variables:
	export WITHINGS_CLIENT_ID="<YOUR CLIENT ID>"
	export WITHINGS_CLIENT_SECRET="<YOUR CLIENT SECRET>"
	export WITHINGS_REDIRECT_URL="<YOUR CALLBACK URL>"

Then run the application:
	go run main.go

Open the URL that shows up on the console and authorize the application.
(Don't worry, it'll only use a demo account, not your real accont)

If you've done everything right, the access token (and some other info) should show up on in the console:
	Access token: ba3438ebe9bae6b642e429cc7c662bda
	Expiry: 2021-10-25T22:41:26+02:00
	Refresh token: 93e82d64d8279589d506fffb067ed59d
	Token type: Bearer
	User ID: 1234
	Scope: user.activity,user.metrics,user.sleepevents
*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sagikazarmark/go-withings/oauth2"
)

func main() {
	// Initialize oauth2 config
	config := &oauth2.WithingsConfig{
		Config: &oauth2.Config{
			ClientID:     os.Getenv("WITHINGS_CLIENT_ID"),
			ClientSecret: os.Getenv("WITHINGS_CLIENT_SECRET"),
			Scopes:       []string{"user.activity", "user.metrics", "user.sleepevents"},
			Endpoint:     oauth2.Endpoint,
			RedirectURL:  os.Getenv("WITHINGS_REDIRECT_URL"),
		},
	}

	// Initialize auth flow
	url := config.AuthCodeURL("state", oauth2.ModeDemo)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	// Register callback URL
	http.HandleFunc("/oauth2/callback", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("could not parse query: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		code := r.FormValue("code")

		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			fmt.Printf("token exchange failed: %v\n", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)

			return
		}

		fmt.Printf("Access token: %s\n", token.AccessToken)
		fmt.Printf("Expiry: %s\n", token.Expiry.Format(time.RFC3339))
		fmt.Printf("Refresh token: %s\n", token.RefreshToken)
		fmt.Printf("Token type: %s\n", token.TokenType)
		fmt.Printf("User ID: %.0f\n", token.Extra("userid"))
		fmt.Printf("Scope: %s\n", token.Extra("scope"))
	})

	http.ListenAndServe("127.0.0.1:8080", nil)
}
