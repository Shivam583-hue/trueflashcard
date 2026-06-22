package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/auth"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
)

const (
	maxNameLen        = 200
	maxDescriptionLen = 2000
	maxCardTextLen    = 10000
	maxImportCards    = 1000

	pgUniqueViolation = "23505"
)

type Transactor interface {
	WithTx(ctx context.Context, fn func(q dbgen.Querier) error) error
}

type cardInput struct {
	front string
	back  string
}

func validateCardInputs(cards []*flashcardv1.CardInput) ([]cardInput, error) {
	if len(cards) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one card is required")
	}
	if len(cards) > maxImportCards {
		return nil, status.Errorf(codes.InvalidArgument, "too many cards: %d (max %d)", len(cards), maxImportCards)
	}

	out := make([]cardInput, 0, len(cards))
	var problems []string
	for i, c := range cards {
		front := strings.TrimSpace(c.GetFront())
		back := strings.TrimSpace(c.GetBack())
		switch {
		case front == "":
			problems = append(problems, fmt.Sprintf("card %d: front is required", i+1))
		case back == "":
			problems = append(problems, fmt.Sprintf("card %d: back is required", i+1))
		case len(front) > maxCardTextLen || len(back) > maxCardTextLen:
			problems = append(problems, fmt.Sprintf("card %d: front and back must be at most %d characters", i+1, maxCardTextLen))
		default:
			out = append(out, cardInput{front: front, back: back})
		}
	}
	if len(problems) > 0 {
		return nil, status.Error(codes.InvalidArgument, strings.Join(problems, "; "))
	}
	return out, nil
}

func callerID(ctx context.Context) (uuid.UUID, error) {
	id, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return uuid.Nil, status.Error(codes.Unauthenticated, "missing authenticated user")
	}
	return id, nil
}

func parseID(field, value string) (uuid.UUID, error) {
	if strings.TrimSpace(value) == "" {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "%s is required", field)
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "%s is not a valid id", field)
	}
	return id, nil
}

func requireText(field, value string, max int) (string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", status.Errorf(codes.InvalidArgument, "%s is required", field)
	}
	if len(v) > max {
		return "", status.Errorf(codes.InvalidArgument, "%s must be at most %d characters", field, max)
	}
	return v, nil
}

func optionalText(field, value string, max int) (string, error) {
	v := strings.TrimSpace(value)
	if len(v) > max {
		return "", status.Errorf(codes.InvalidArgument, "%s must be at most %d characters", field, max)
	}
	return v, nil
}

func translateError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return status.Error(codes.NotFound, "resource not found")
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return status.Error(codes.Aborted, "a conflicting change occurred, please retry")
	}
	return status.Error(codes.Internal, "internal error")
}
