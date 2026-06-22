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

func TestCreateDeckSuccess(t *testing.T) {
	ctx, userID := authedContext()
	folderID := uuid.New()
	deckID := uuid.New()
	var captured dbgen.CreateDeckParams
	q := &stubQuerier{
		createDeckFn: func(_ context.Context, p dbgen.CreateDeckParams) (dbgen.Deck, error) {
			captured = p
			return dbgen.Deck{ID: deckID, FolderID: p.FolderID, Name: p.Name, Description: p.Description}, nil
		},
	}

	resp, err := NewDeckService(q, directTx{q}).CreateDeck(ctx, &flashcardv1.CreateDeckRequest{
		FolderId:    folderID.String(),
		Name:        "Cell Biology",
		Description: "  intro  ",
	})
	requireNoError(t, err)

	if captured.FolderID != folderID || captured.OwnerID != userID {
		t.Fatalf("deck not scoped correctly: %+v", captured)
	}
	if captured.Description != "intro" {
		t.Fatalf("description not trimmed: %q", captured.Description)
	}
	if resp.GetDeck().GetCardCount() != 0 {
		t.Fatalf("new deck should have 0 cards, got %d", resp.GetDeck().GetCardCount())
	}
}

func TestCreateDeckMissingFolderID(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	_, err := NewDeckService(q, directTx{q}).CreateDeck(ctx, &flashcardv1.CreateDeckRequest{Name: "X"})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestCreateDeckFolderNotOwned(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		createDeckFn: func(_ context.Context, _ dbgen.CreateDeckParams) (dbgen.Deck, error) {
			return dbgen.Deck{}, pgx.ErrNoRows
		},
	}
	_, err := NewDeckService(q, directTx{q}).CreateDeck(ctx, &flashcardv1.CreateDeckRequest{
		FolderId: uuid.NewString(),
		Name:     "X",
	})
	requireCode(t, err, codes.NotFound)
}

func TestGetDeckSuccess(t *testing.T) {
	ctx, _ := authedContext()
	deckID := uuid.New()
	q := &stubQuerier{
		getDeckFn: func(_ context.Context, _ dbgen.GetDeckParams) (dbgen.GetDeckRow, error) {
			return dbgen.GetDeckRow{ID: deckID, FolderID: uuid.New(), Name: "Deck", CardCount: 5}, nil
		},
	}
	resp, err := NewDeckService(q, directTx{q}).GetDeck(ctx, &flashcardv1.GetDeckRequest{Id: deckID.String()})
	requireNoError(t, err)
	if resp.GetDeck().GetCardCount() != 5 {
		t.Fatalf("expected card_count 5, got %d", resp.GetDeck().GetCardCount())
	}
}

func TestListDecksSuccess(t *testing.T) {
	ctx, userID := authedContext()
	folderID := uuid.New()
	var captured dbgen.ListDecksParams
	q := &stubQuerier{
		listDecksFn: func(_ context.Context, p dbgen.ListDecksParams) ([]dbgen.ListDecksRow, error) {
			captured = p
			return []dbgen.ListDecksRow{{ID: uuid.New(), FolderID: p.FolderID, Name: "D", CardCount: 2}}, nil
		},
	}
	resp, err := NewDeckService(q, directTx{q}).ListDecks(ctx, &flashcardv1.ListDecksRequest{FolderId: folderID.String()})
	requireNoError(t, err)
	if captured.FolderID != folderID || captured.OwnerID != userID {
		t.Fatalf("list not scoped correctly: %+v", captured)
	}
	if len(resp.GetDecks()) != 1 {
		t.Fatalf("expected 1 deck, got %d", len(resp.GetDecks()))
	}
}

func TestUpdateDeckSuccess(t *testing.T) {
	ctx, _ := authedContext()
	deckID := uuid.New()
	q := &stubQuerier{
		updateDeckFn: func(_ context.Context, p dbgen.UpdateDeckParams) (dbgen.UpdateDeckRow, error) {
			return dbgen.UpdateDeckRow{ID: p.ID, FolderID: uuid.New(), Name: p.Name, CardCount: 3}, nil
		},
	}
	resp, err := NewDeckService(q, directTx{q}).UpdateDeck(ctx, &flashcardv1.UpdateDeckRequest{Id: deckID.String(), Name: "New"})
	requireNoError(t, err)
	if resp.GetDeck().GetCardCount() != 3 {
		t.Fatalf("expected card_count 3, got %d", resp.GetDeck().GetCardCount())
	}
}

func TestDeleteDeckNotFound(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		deleteDeckFn: func(_ context.Context, _ dbgen.DeleteDeckParams) (int64, error) {
			return 0, nil
		},
	}
	_, err := NewDeckService(q, directTx{q}).DeleteDeck(ctx, &flashcardv1.DeleteDeckRequest{Id: uuid.NewString()})
	requireCode(t, err, codes.NotFound)
}
