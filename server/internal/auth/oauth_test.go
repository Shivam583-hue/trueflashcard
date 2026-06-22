package auth

import (
	"net/url"
	"testing"
)

func setOAuthEnv(t *testing.T) {
	t.Helper()
	t.Setenv("GOOGLE_CLIENT_ID", "client-123.apps.googleusercontent.com")
	t.Setenv("GOOGLE_CLIENT_SECRET", "secret-xyz")
	t.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback")
}

func TestNewGoogleOAuthRequiresAllEnv(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "only-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "")
	t.Setenv("GOOGLE_REDIRECT_URL", "")
	if _, err := NewGoogleOAuth(); err != ErrMissingOAuthConfig {
		t.Fatalf("expected ErrMissingOAuthConfig, got %v", err)
	}
}

func TestAuthCodeURLContainsConfig(t *testing.T) {
	setOAuthEnv(t)
	g, err := NewGoogleOAuth()
	if err != nil {
		t.Fatalf("NewGoogleOAuth: %v", err)
	}

	raw := g.AuthCodeURL("state-token")
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse auth url: %v", err)
	}
	if u.Host != "accounts.google.com" {
		t.Fatalf("unexpected auth host: %s", u.Host)
	}

	q := u.Query()
	checks := map[string]string{
		"client_id":     "client-123.apps.googleusercontent.com",
		"redirect_uri":  "http://localhost:8080/auth/google/callback",
		"state":         "state-token",
		"response_type": "code",
	}
	for key, want := range checks {
		if got := q.Get(key); got != want {
			t.Fatalf("auth url param %q = %q, want %q", key, got, want)
		}
	}
	if scope := q.Get("scope"); scope == "" {
		t.Fatal("expected scopes in auth url")
	}
}
