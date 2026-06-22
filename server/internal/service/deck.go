package service

import (
	"context"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeckService struct {
	flashcardv1.UnimplementedDeckServiceServer
	q dbgen.Querier
}

func NewDeckService(q dbgen.Querier) *DeckService {
	return &DeckService{q: q}
}

func (s *DeckService) CreateDeck(ctx context.Context, req *flashcardv1.CreateDeckRequest) (*flashcardv1.CreateDeckResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	folderID, err := parseID("folder_id", req.GetFolderId())
	if err != nil {
		return nil, err
	}
	name, err := requireText("name", req.GetName(), maxNameLen)
	if err != nil {
		return nil, err
	}
	description, err := optionalText("description", req.GetDescription(), maxDescriptionLen)
	if err != nil {
		return nil, err
	}

	deck, err := s.q.CreateDeck(ctx, dbgen.CreateDeckParams{
		FolderID:    folderID,
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.CreateDeckResponse{Deck: deckToProto(deck, 0)}, nil
}

func (s *DeckService) GetDeck(ctx context.Context, req *flashcardv1.GetDeckRequest) (*flashcardv1.GetDeckResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID("id", req.GetId())
	if err != nil {
		return nil, err
	}

	deck, err := s.q.GetDeck(ctx, dbgen.GetDeckParams{ID: id, OwnerID: ownerID})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.GetDeckResponse{Deck: getDeckRowToProto(deck)}, nil
}

func (s *DeckService) ListDecks(ctx context.Context, req *flashcardv1.ListDecksRequest) (*flashcardv1.ListDecksResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	folderID, err := parseID("folder_id", req.GetFolderId())
	if err != nil {
		return nil, err
	}

	decks, err := s.q.ListDecks(ctx, dbgen.ListDecksParams{FolderID: folderID, OwnerID: ownerID})
	if err != nil {
		return nil, translateError(err)
	}

	out := make([]*flashcardv1.Deck, 0, len(decks))
	for _, d := range decks {
		out = append(out, listDeckRowToProto(d))
	}
	return &flashcardv1.ListDecksResponse{Decks: out}, nil
}

func (s *DeckService) UpdateDeck(ctx context.Context, req *flashcardv1.UpdateDeckRequest) (*flashcardv1.UpdateDeckResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID("id", req.GetId())
	if err != nil {
		return nil, err
	}
	name, err := requireText("name", req.GetName(), maxNameLen)
	if err != nil {
		return nil, err
	}
	description, err := optionalText("description", req.GetDescription(), maxDescriptionLen)
	if err != nil {
		return nil, err
	}

	deck, err := s.q.UpdateDeck(ctx, dbgen.UpdateDeckParams{
		ID:          id,
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.UpdateDeckResponse{Deck: updateDeckRowToProto(deck)}, nil
}

func (s *DeckService) DeleteDeck(ctx context.Context, req *flashcardv1.DeleteDeckRequest) (*flashcardv1.DeleteDeckResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID("id", req.GetId())
	if err != nil {
		return nil, err
	}

	rows, err := s.q.DeleteDeck(ctx, dbgen.DeleteDeckParams{ID: id, OwnerID: ownerID})
	if err != nil {
		return nil, translateError(err)
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, "deck not found")
	}
	return &flashcardv1.DeleteDeckResponse{}, nil
}
