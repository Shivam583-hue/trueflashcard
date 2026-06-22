package connectapi

import (
	"context"

	"connectrpc.com/connect"

	v1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/service"
)

type folderAPI struct{ svc *service.FolderService }

func (a folderAPI) CreateFolder(ctx context.Context, r *connect.Request[v1.CreateFolderRequest]) (*connect.Response[v1.CreateFolderResponse], error) {
	return invoke(ctx, r, a.svc.CreateFolder)
}

func (a folderAPI) GetFolder(ctx context.Context, r *connect.Request[v1.GetFolderRequest]) (*connect.Response[v1.GetFolderResponse], error) {
	return invoke(ctx, r, a.svc.GetFolder)
}

func (a folderAPI) ListFolders(ctx context.Context, r *connect.Request[v1.ListFoldersRequest]) (*connect.Response[v1.ListFoldersResponse], error) {
	return invoke(ctx, r, a.svc.ListFolders)
}

func (a folderAPI) UpdateFolder(ctx context.Context, r *connect.Request[v1.UpdateFolderRequest]) (*connect.Response[v1.UpdateFolderResponse], error) {
	return invoke(ctx, r, a.svc.UpdateFolder)
}

func (a folderAPI) DeleteFolder(ctx context.Context, r *connect.Request[v1.DeleteFolderRequest]) (*connect.Response[v1.DeleteFolderResponse], error) {
	return invoke(ctx, r, a.svc.DeleteFolder)
}

type deckAPI struct{ svc *service.DeckService }

func (a deckAPI) CreateDeck(ctx context.Context, r *connect.Request[v1.CreateDeckRequest]) (*connect.Response[v1.CreateDeckResponse], error) {
	return invoke(ctx, r, a.svc.CreateDeck)
}

func (a deckAPI) GetDeck(ctx context.Context, r *connect.Request[v1.GetDeckRequest]) (*connect.Response[v1.GetDeckResponse], error) {
	return invoke(ctx, r, a.svc.GetDeck)
}

func (a deckAPI) ListDecks(ctx context.Context, r *connect.Request[v1.ListDecksRequest]) (*connect.Response[v1.ListDecksResponse], error) {
	return invoke(ctx, r, a.svc.ListDecks)
}

func (a deckAPI) UpdateDeck(ctx context.Context, r *connect.Request[v1.UpdateDeckRequest]) (*connect.Response[v1.UpdateDeckResponse], error) {
	return invoke(ctx, r, a.svc.UpdateDeck)
}

func (a deckAPI) DeleteDeck(ctx context.Context, r *connect.Request[v1.DeleteDeckRequest]) (*connect.Response[v1.DeleteDeckResponse], error) {
	return invoke(ctx, r, a.svc.DeleteDeck)
}

func (a deckAPI) ImportDeck(ctx context.Context, r *connect.Request[v1.ImportDeckRequest]) (*connect.Response[v1.ImportDeckResponse], error) {
	return invoke(ctx, r, a.svc.ImportDeck)
}

type flashcardAPI struct{ svc *service.FlashcardService }

func (a flashcardAPI) CreateFlashcard(ctx context.Context, r *connect.Request[v1.CreateFlashcardRequest]) (*connect.Response[v1.CreateFlashcardResponse], error) {
	return invoke(ctx, r, a.svc.CreateFlashcard)
}

func (a flashcardAPI) GetFlashcard(ctx context.Context, r *connect.Request[v1.GetFlashcardRequest]) (*connect.Response[v1.GetFlashcardResponse], error) {
	return invoke(ctx, r, a.svc.GetFlashcard)
}

func (a flashcardAPI) ListFlashcards(ctx context.Context, r *connect.Request[v1.ListFlashcardsRequest]) (*connect.Response[v1.ListFlashcardsResponse], error) {
	return invoke(ctx, r, a.svc.ListFlashcards)
}

func (a flashcardAPI) UpdateFlashcard(ctx context.Context, r *connect.Request[v1.UpdateFlashcardRequest]) (*connect.Response[v1.UpdateFlashcardResponse], error) {
	return invoke(ctx, r, a.svc.UpdateFlashcard)
}

func (a flashcardAPI) DeleteFlashcard(ctx context.Context, r *connect.Request[v1.DeleteFlashcardRequest]) (*connect.Response[v1.DeleteFlashcardResponse], error) {
	return invoke(ctx, r, a.svc.DeleteFlashcard)
}

func (a flashcardAPI) ImportFlashcards(ctx context.Context, r *connect.Request[v1.ImportFlashcardsRequest]) (*connect.Response[v1.ImportFlashcardsResponse], error) {
	return invoke(ctx, r, a.svc.ImportFlashcards)
}
