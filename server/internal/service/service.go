package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Shivam583-hue/trueflashcard/server/internal/auth"
)

const (
	maxNameLen        = 200
	maxDescriptionLen = 2000
	maxCardTextLen    = 10000
)

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
	return status.Error(codes.Internal, "internal error")
}
