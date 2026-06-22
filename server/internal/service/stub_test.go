package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
)

type stubQuerier struct {
	dbgen.Querier
	calls []string

	createFolderFn func(context.Context, dbgen.CreateFolderParams) (dbgen.Folder, error)
	getFolderFn    func(context.Context, dbgen.GetFolderParams) (dbgen.Folder, error)
	listFoldersFn  func(context.Context, uuid.UUID) ([]dbgen.Folder, error)
	updateFolderFn func(context.Context, dbgen.UpdateFolderParams) (dbgen.Folder, error)
	deleteFolderFn func(context.Context, dbgen.DeleteFolderParams) (int64, error)

	createDeckFn func(context.Context, dbgen.CreateDeckParams) (dbgen.Deck, error)
	getDeckFn    func(context.Context, dbgen.GetDeckParams) (dbgen.GetDeckRow, error)
	listDecksFn  func(context.Context, dbgen.ListDecksParams) ([]dbgen.ListDecksRow, error)
	updateDeckFn func(context.Context, dbgen.UpdateDeckParams) (dbgen.UpdateDeckRow, error)
	deleteDeckFn func(context.Context, dbgen.DeleteDeckParams) (int64, error)

	lockDeckFn        func(context.Context, dbgen.LockDeckForUpdateParams) (uuid.UUID, error)
	createFlashcardFn func(context.Context, dbgen.CreateFlashcardParams) (dbgen.Flashcard, error)
	getFlashcardFn    func(context.Context, dbgen.GetFlashcardParams) (dbgen.Flashcard, error)
	listFlashcardsFn  func(context.Context, dbgen.ListFlashcardsParams) ([]dbgen.Flashcard, error)
	updateFlashcardFn func(context.Context, dbgen.UpdateFlashcardParams) (dbgen.Flashcard, error)
	deleteFlashcardFn func(context.Context, dbgen.DeleteFlashcardParams) (int64, error)
}

func (s *stubQuerier) record(name string) { s.calls = append(s.calls, name) }

// directTx runs the transaction body immediately against a querier, simulating
// a transaction without a real database.
type directTx struct{ q dbgen.Querier }

func (d directTx) WithTx(_ context.Context, fn func(dbgen.Querier) error) error {
	return fn(d.q)
}

func (s *stubQuerier) CreateFolder(ctx context.Context, p dbgen.CreateFolderParams) (dbgen.Folder, error) {
	s.record("CreateFolder")
	return s.createFolderFn(ctx, p)
}

func (s *stubQuerier) GetFolder(ctx context.Context, p dbgen.GetFolderParams) (dbgen.Folder, error) {
	s.record("GetFolder")
	return s.getFolderFn(ctx, p)
}

func (s *stubQuerier) ListFolders(ctx context.Context, ownerID uuid.UUID) ([]dbgen.Folder, error) {
	s.record("ListFolders")
	return s.listFoldersFn(ctx, ownerID)
}

func (s *stubQuerier) UpdateFolder(ctx context.Context, p dbgen.UpdateFolderParams) (dbgen.Folder, error) {
	s.record("UpdateFolder")
	return s.updateFolderFn(ctx, p)
}

func (s *stubQuerier) DeleteFolder(ctx context.Context, p dbgen.DeleteFolderParams) (int64, error) {
	s.record("DeleteFolder")
	return s.deleteFolderFn(ctx, p)
}

func (s *stubQuerier) CreateDeck(ctx context.Context, p dbgen.CreateDeckParams) (dbgen.Deck, error) {
	s.record("CreateDeck")
	return s.createDeckFn(ctx, p)
}

func (s *stubQuerier) GetDeck(ctx context.Context, p dbgen.GetDeckParams) (dbgen.GetDeckRow, error) {
	s.record("GetDeck")
	return s.getDeckFn(ctx, p)
}

func (s *stubQuerier) ListDecks(ctx context.Context, p dbgen.ListDecksParams) ([]dbgen.ListDecksRow, error) {
	s.record("ListDecks")
	return s.listDecksFn(ctx, p)
}

func (s *stubQuerier) UpdateDeck(ctx context.Context, p dbgen.UpdateDeckParams) (dbgen.UpdateDeckRow, error) {
	s.record("UpdateDeck")
	return s.updateDeckFn(ctx, p)
}

func (s *stubQuerier) DeleteDeck(ctx context.Context, p dbgen.DeleteDeckParams) (int64, error) {
	s.record("DeleteDeck")
	return s.deleteDeckFn(ctx, p)
}

func (s *stubQuerier) LockDeckForUpdate(ctx context.Context, p dbgen.LockDeckForUpdateParams) (uuid.UUID, error) {
	s.record("LockDeckForUpdate")
	return s.lockDeckFn(ctx, p)
}

func (s *stubQuerier) CreateFlashcard(ctx context.Context, p dbgen.CreateFlashcardParams) (dbgen.Flashcard, error) {
	s.record("CreateFlashcard")
	return s.createFlashcardFn(ctx, p)
}

func (s *stubQuerier) GetFlashcard(ctx context.Context, p dbgen.GetFlashcardParams) (dbgen.Flashcard, error) {
	s.record("GetFlashcard")
	return s.getFlashcardFn(ctx, p)
}

func (s *stubQuerier) ListFlashcards(ctx context.Context, p dbgen.ListFlashcardsParams) ([]dbgen.Flashcard, error) {
	s.record("ListFlashcards")
	return s.listFlashcardsFn(ctx, p)
}

func (s *stubQuerier) UpdateFlashcard(ctx context.Context, p dbgen.UpdateFlashcardParams) (dbgen.Flashcard, error) {
	s.record("UpdateFlashcard")
	return s.updateFlashcardFn(ctx, p)
}

func (s *stubQuerier) DeleteFlashcard(ctx context.Context, p dbgen.DeleteFlashcardParams) (int64, error) {
	s.record("DeleteFlashcard")
	return s.deleteFlashcardFn(ctx, p)
}
