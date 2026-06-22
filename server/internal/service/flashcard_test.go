package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
)

func TestCreateFlashcardSuccess(t *testing.T) {
	ctx, userID := authedContext()
	deckID := uuid.New()
	cardID := uuid.New()
	var captured dbgen.CreateFlashcardParams
	q := &stubQuerier{
		createFlashcardFn: func(_ context.Context, p dbgen.CreateFlashcardParams) (dbgen.Flashcard, error) {
			captured = p
			return dbgen.Flashcard{ID: cardID, DeckID: p.DeckID, Front: p.Front, Back: p.Back, Position: 0}, nil
		},
	}
	resp, err := NewFlashcardService(q).CreateFlashcard(ctx, &flashcardv1.CreateFlashcardRequest{
		DeckId: deckID.String(),
		Front:  "  Q  ",
		Back:   "  A  ",
	})
	requireNoError(t, err)
	if captured.DeckID != deckID || captured.OwnerID != userID {
		t.Fatalf("flashcard not scoped correctly: %+v", captured)
	}
	if captured.Front != "Q" || captured.Back != "A" {
		t.Fatalf("text not trimmed: front=%q back=%q", captured.Front, captured.Back)
	}
	if resp.GetFlashcard().GetId() != cardID.String() {
		t.Fatalf("unexpected card id")
	}
}

func TestCreateFlashcardMissingBack(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	_, err := NewFlashcardService(q).CreateFlashcard(ctx, &flashcardv1.CreateFlashcardRequest{
		DeckId: uuid.NewString(),
		Front:  "Q",
		Back:   "   ",
	})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestCreateFlashcardDeckNotOwned(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		createFlashcardFn: func(_ context.Context, _ dbgen.CreateFlashcardParams) (dbgen.Flashcard, error) {
			return dbgen.Flashcard{}, pgx.ErrNoRows
		},
	}
	_, err := NewFlashcardService(q).CreateFlashcard(ctx, &flashcardv1.CreateFlashcardRequest{
		DeckId: uuid.NewString(),
		Front:  "Q",
		Back:   "A",
	})
	requireCode(t, err, codes.NotFound)
}

func TestListFlashcardsSuccess(t *testing.T) {
	ctx, userID := authedContext()
	deckID := uuid.New()
	var captured dbgen.ListFlashcardsParams
	q := &stubQuerier{
		listFlashcardsFn: func(_ context.Context, p dbgen.ListFlashcardsParams) ([]dbgen.Flashcard, error) {
			captured = p
			return []dbgen.Flashcard{
				{ID: uuid.New(), DeckID: p.DeckID, Front: "Q1", Back: "A1", Position: 0},
				{ID: uuid.New(), DeckID: p.DeckID, Front: "Q2", Back: "A2", Position: 1},
			}, nil
		},
	}
	resp, err := NewFlashcardService(q).ListFlashcards(ctx, &flashcardv1.ListFlashcardsRequest{DeckId: deckID.String()})
	requireNoError(t, err)
	if captured.DeckID != deckID || captured.OwnerID != userID {
		t.Fatalf("list not scoped correctly: %+v", captured)
	}
	if len(resp.GetFlashcards()) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(resp.GetFlashcards()))
	}
}

func TestUpdateFlashcardNegativePosition(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	_, err := NewFlashcardService(q).UpdateFlashcard(ctx, &flashcardv1.UpdateFlashcardRequest{
		Id:       uuid.NewString(),
		Front:    "Q",
		Back:     "A",
		Position: -1,
	})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestDeleteFlashcardSuccess(t *testing.T) {
	ctx, userID := authedContext()
	var captured dbgen.DeleteFlashcardParams
	q := &stubQuerier{
		deleteFlashcardFn: func(_ context.Context, p dbgen.DeleteFlashcardParams) (int64, error) {
			captured = p
			return 1, nil
		},
	}
	_, err := NewFlashcardService(q).DeleteFlashcard(ctx, &flashcardv1.DeleteFlashcardRequest{Id: uuid.NewString()})
	requireNoError(t, err)
	if captured.OwnerID != userID {
		t.Fatalf("delete not scoped to caller")
	}
}

func TestDeleteFlashcardNotFound(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		deleteFlashcardFn: func(_ context.Context, _ dbgen.DeleteFlashcardParams) (int64, error) {
			return 0, nil
		},
	}
	_, err := NewFlashcardService(q).DeleteFlashcard(ctx, &flashcardv1.DeleteFlashcardRequest{Id: uuid.NewString()})
	requireCode(t, err, codes.NotFound)
}
