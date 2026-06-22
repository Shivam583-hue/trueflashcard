package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestManager(t *testing.T, secret string) *SessionManager {
	t.Helper()
	t.Setenv("JWT_SECRET", secret)
	m, err := NewSessionManager()
	if err != nil {
		t.Fatalf("NewSessionManager: %v", err)
	}
	return m
}

func TestSessionRoundTrip(t *testing.T) {
	m := newTestManager(t, "test-secret")
	userID := uuid.New()

	token, expiresAt, err := m.Issue(userID)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if !expiresAt.After(time.Now()) {
		t.Fatal("expiry should be in the future")
	}

	got, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if got != userID {
		t.Fatalf("round trip mismatch: got %v want %v", got, userID)
	}
}

func TestVerifyRejectsWrongSecret(t *testing.T) {
	issuer := newTestManager(t, "secret-a")
	token, _, err := issuer.Issue(uuid.New())
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	verifier := newTestManager(t, "secret-b")
	if _, err := verifier.Verify(token); err == nil {
		t.Fatal("expected verification to fail with a different secret")
	}
}

func TestVerifyRejectsGarbage(t *testing.T) {
	m := newTestManager(t, "test-secret")
	if _, err := m.Verify("not.a.jwt"); err == nil {
		t.Fatal("expected verification to fail for malformed token")
	}
}

func TestNewSessionManagerRequiresSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	if _, err := NewSessionManager(); err != ErrMissingJWTSecret {
		t.Fatalf("expected ErrMissingJWTSecret, got %v", err)
	}
}
