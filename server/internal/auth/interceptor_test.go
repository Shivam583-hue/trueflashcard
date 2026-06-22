package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func runInterceptor(t *testing.T, m *SessionManager, md metadata.MD) (uuid.UUID, bool) {
	t.Helper()
	ctx := context.Background()
	if md != nil {
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	var gotID uuid.UUID
	var gotOK bool
	handler := func(ctx context.Context, _ any) (any, error) {
		gotID, gotOK = UserIDFromContext(ctx)
		return nil, nil
	}
	_, err := UnaryAuthInterceptor(m)(ctx, nil, &grpc.UnaryServerInfo{}, handler)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}
	return gotID, gotOK
}

func TestInterceptorInjectsValidSession(t *testing.T) {
	m := newTestManager(t, "test-secret")
	userID := uuid.New()
	token, _, err := m.Issue(userID)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	gotID, ok := runInterceptor(t, m, metadata.Pairs("authorization", "Bearer "+token))
	if !ok {
		t.Fatal("expected user id to be injected")
	}
	if gotID != userID {
		t.Fatalf("got %v want %v", gotID, userID)
	}
}

func TestInterceptorIgnoresMissingMetadata(t *testing.T) {
	m := newTestManager(t, "test-secret")
	if _, ok := runInterceptor(t, m, nil); ok {
		t.Fatal("expected no user id without metadata")
	}
}

func TestInterceptorIgnoresInvalidToken(t *testing.T) {
	m := newTestManager(t, "test-secret")
	if _, ok := runInterceptor(t, m, metadata.Pairs("authorization", "Bearer garbage")); ok {
		t.Fatal("expected no user id for invalid token")
	}
}
