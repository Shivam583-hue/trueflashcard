package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Shivam583-hue/trueflashcard/server/internal/auth"
)

func authedContext() (context.Context, uuid.UUID) {
	userID := uuid.New()
	return auth.WithUserID(context.Background(), userID), userID
}

func requireCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	if got := status.Code(err); got != want {
		t.Fatalf("expected status code %v, got %v (err=%v)", want, got, err)
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireNoCalls(t *testing.T, q *stubQuerier) {
	t.Helper()
	if len(q.calls) != 0 {
		t.Fatalf("expected no database calls, got %v", q.calls)
	}
}

func validTimestamp() pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
}
