package auth

import (
	"context"
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var ErrMissingOAuthConfig = errors.New("GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET and GOOGLE_REDIRECT_URL must all be set")

type GoogleIdentity struct {
	Subject string
	Email   string
	Name    string
}

type GoogleOAuth struct {
	config   *oauth2.Config
	validate func(ctx context.Context, idToken, audience string) (*idtoken.Payload, error)
}

func NewGoogleOAuth() (*GoogleOAuth, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, ErrMissingOAuthConfig
	}

	return &GoogleOAuth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		validate: func(ctx context.Context, idToken, audience string) (*idtoken.Payload, error) {
			return idtoken.Validate(ctx, idToken, audience)
		},
	}, nil
}

func (g *GoogleOAuth) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)
}

func (g *GoogleOAuth) Exchange(ctx context.Context, code string) (*GoogleIdentity, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, errors.New("missing id_token in oauth response")
	}

	payload, err := g.validate(ctx, rawIDToken, g.config.ClientID)
	if err != nil {
		return nil, err
	}

	identity := &GoogleIdentity{Subject: payload.Subject}
	if email, ok := payload.Claims["email"].(string); ok {
		identity.Email = email
	}
	if name, ok := payload.Claims["name"].(string); ok {
		identity.Name = name
	}
	return identity, nil
}
