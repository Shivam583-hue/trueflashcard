package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
)

func cards(pairs ...[2]string) []*flashcardv1.CardInput {
	out := make([]*flashcardv1.CardInput, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, &flashcardv1.CardInput{Front: p[0], Back: p[1]})
	}
	return out
}

func TestImportDeckSuccess(t *testing.T) {
	ctx, userID := authedContext()
	folderID := uuid.New()
	deckID := uuid.New()
	var deckParams dbgen.CreateDeckParams
	var insertedDeckIDs []uuid.UUID
	q := &stubQuerier{
		createDeckFn: func(_ context.Context, p dbgen.CreateDeckParams) (dbgen.Deck, error) {
			deckParams = p
			return dbgen.Deck{ID: deckID, FolderID: p.FolderID, Name: p.Name}, nil
		},
		createFlashcardFn: func(_ context.Context, p dbgen.CreateFlashcardParams) (dbgen.Flashcard, error) {
			insertedDeckIDs = append(insertedDeckIDs, p.DeckID)
			return dbgen.Flashcard{ID: uuid.New(), DeckID: p.DeckID, Front: p.Front, Back: p.Back}, nil
		},
	}

	resp, err := NewDeckService(q, directTx{q}).ImportDeck(ctx, &flashcardv1.ImportDeckRequest{
		FolderId: folderID.String(),
		Name:     "  Biology  ",
		Cards:    cards([2]string{" Q1 ", " A1 "}, [2]string{"Q2", "A2"}),
	})
	requireNoError(t, err)

	if deckParams.FolderID != folderID || deckParams.OwnerID != userID {
		t.Fatalf("deck not scoped to caller: %+v", deckParams)
	}
	if deckParams.Name != "Biology" {
		t.Fatalf("name not trimmed: %q", deckParams.Name)
	}
	if resp.GetCreatedCount() != 2 {
		t.Fatalf("expected created_count 2, got %d", resp.GetCreatedCount())
	}
	if len(insertedDeckIDs) != 2 {
		t.Fatalf("expected 2 card inserts, got %d", len(insertedDeckIDs))
	}
	for _, id := range insertedDeckIDs {
		if id != deckID {
			t.Fatalf("card inserted into wrong deck: %v", id)
		}
	}
}

func TestImportDeckTrimsCards(t *testing.T) {
	ctx, _ := authedContext()
	var fronts []string
	q := &stubQuerier{
		createDeckFn: func(_ context.Context, p dbgen.CreateDeckParams) (dbgen.Deck, error) {
			return dbgen.Deck{ID: uuid.New(), FolderID: p.FolderID, Name: p.Name}, nil
		},
		createFlashcardFn: func(_ context.Context, p dbgen.CreateFlashcardParams) (dbgen.Flashcard, error) {
			fronts = append(fronts, p.Front)
			return dbgen.Flashcard{ID: uuid.New(), DeckID: p.DeckID, Front: p.Front, Back: p.Back}, nil
		},
	}
	_, err := NewDeckService(q, directTx{q}).ImportDeck(ctx, &flashcardv1.ImportDeckRequest{
		FolderId: uuid.NewString(),
		Name:     "Deck",
		Cards:    cards([2]string{"  spaced  ", "back"}),
	})
	requireNoError(t, err)
	if len(fronts) != 1 || fronts[0] != "spaced" {
		t.Fatalf("front not trimmed: %v", fronts)
	}
}

func TestImportDeckRejectsInvalidCards(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	_, err := NewDeckService(q, directTx{q}).ImportDeck(ctx, &flashcardv1.ImportDeckRequest{
		FolderId: uuid.NewString(),
		Name:     "Deck",
		Cards:    cards([2]string{"Q1", "A1"}, [2]string{"Q2", "   "}),
	})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestImportDeckRejectsEmptyCards(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	_, err := NewDeckService(q, directTx{q}).ImportDeck(ctx, &flashcardv1.ImportDeckRequest{
		FolderId: uuid.NewString(),
		Name:     "Deck",
		Cards:    nil,
	})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestImportDeckTooManyCards(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	big := make([]*flashcardv1.CardInput, maxImportCards+1)
	for i := range big {
		big[i] = &flashcardv1.CardInput{Front: "f", Back: "b"}
	}
	_, err := NewDeckService(q, directTx{q}).ImportDeck(ctx, &flashcardv1.ImportDeckRequest{
		FolderId: uuid.NewString(),
		Name:     "Deck",
		Cards:    big,
	})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestImportDeckFolderNotOwned(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		createDeckFn: func(_ context.Context, _ dbgen.CreateDeckParams) (dbgen.Deck, error) {
			return dbgen.Deck{}, pgx.ErrNoRows
		},
	}
	_, err := NewDeckService(q, directTx{q}).ImportDeck(ctx, &flashcardv1.ImportDeckRequest{
		FolderId: uuid.NewString(),
		Name:     "Deck",
		Cards:    cards([2]string{"Q", "A"}),
	})
	requireCode(t, err, codes.NotFound)
}

func TestImportFlashcardsSuccess(t *testing.T) {
	ctx, userID := authedContext()
	deckID := uuid.New()
	var captured []dbgen.CreateFlashcardParams
	q := &stubQuerier{
		lockDeckFn: func(_ context.Context, p dbgen.LockDeckForUpdateParams) (uuid.UUID, error) {
			return p.ID, nil
		},
		createFlashcardFn: func(_ context.Context, p dbgen.CreateFlashcardParams) (dbgen.Flashcard, error) {
			captured = append(captured, p)
			return dbgen.Flashcard{ID: uuid.New(), DeckID: p.DeckID, Front: p.Front, Back: p.Back}, nil
		},
	}
	resp, err := NewFlashcardService(q, directTx{q}).ImportFlashcards(ctx, &flashcardv1.ImportFlashcardsRequest{
		DeckId: deckID.String(),
		Cards:  cards([2]string{"Q1", "A1"}, [2]string{"Q2", "A2"}),
	})
	requireNoError(t, err)
	if len(resp.GetFlashcards()) != 2 {
		t.Fatalf("expected 2 returned cards, got %d", len(resp.GetFlashcards()))
	}
	for _, p := range captured {
		if p.DeckID != deckID || p.OwnerID != userID {
			t.Fatalf("card not scoped correctly: %+v", p)
		}
	}
}

func TestImportFlashcardsDeckNotOwned(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		lockDeckFn: func(_ context.Context, _ dbgen.LockDeckForUpdateParams) (uuid.UUID, error) {
			return uuid.Nil, pgx.ErrNoRows
		},
	}
	_, err := NewFlashcardService(q, directTx{q}).ImportFlashcards(ctx, &flashcardv1.ImportFlashcardsRequest{
		DeckId: uuid.NewString(),
		Cards:  cards([2]string{"Q", "A"}),
	})
	requireCode(t, err, codes.NotFound)
}

func TestImportFlashcardsPositionConflictAborted(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		lockDeckFn: func(_ context.Context, p dbgen.LockDeckForUpdateParams) (uuid.UUID, error) {
			return p.ID, nil
		},
		createFlashcardFn: func(_ context.Context, _ dbgen.CreateFlashcardParams) (dbgen.Flashcard, error) {
			return dbgen.Flashcard{}, &pgconn.PgError{Code: "23505"}
		},
	}
	_, err := NewFlashcardService(q, directTx{q}).ImportFlashcards(ctx, &flashcardv1.ImportFlashcardsRequest{
		DeckId: uuid.NewString(),
		Cards:  cards([2]string{"Q", "A"}),
	})
	requireCode(t, err, codes.Aborted)
}

func TestImportDeckUnauthenticated(t *testing.T) {
	q := &stubQuerier{}
	_, err := NewDeckService(q, directTx{q}).ImportDeck(context.Background(), &flashcardv1.ImportDeckRequest{
		FolderId: uuid.NewString(),
		Name:     "Deck",
		Cards:    cards([2]string{"Q", "A"}),
	})
	requireCode(t, err, codes.Unauthenticated)
	requireNoCalls(t, q)
}
