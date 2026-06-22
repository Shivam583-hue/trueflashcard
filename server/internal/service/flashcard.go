package service

import (
	"context"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FlashcardService struct {
	flashcardv1.UnimplementedFlashcardServiceServer
	q dbgen.Querier
}

func NewFlashcardService(q dbgen.Querier) *FlashcardService {
	return &FlashcardService{q: q}
}

func (s *FlashcardService) CreateFlashcard(ctx context.Context, req *flashcardv1.CreateFlashcardRequest) (*flashcardv1.CreateFlashcardResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	deckID, err := parseID("deck_id", req.GetDeckId())
	if err != nil {
		return nil, err
	}
	front, err := requireText("front", req.GetFront(), maxCardTextLen)
	if err != nil {
		return nil, err
	}
	back, err := requireText("back", req.GetBack(), maxCardTextLen)
	if err != nil {
		return nil, err
	}

	card, err := s.q.CreateFlashcard(ctx, dbgen.CreateFlashcardParams{
		DeckID:  deckID,
		OwnerID: ownerID,
		Front:   front,
		Back:    back,
	})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.CreateFlashcardResponse{Flashcard: flashcardToProto(card)}, nil
}

func (s *FlashcardService) GetFlashcard(ctx context.Context, req *flashcardv1.GetFlashcardRequest) (*flashcardv1.GetFlashcardResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID("id", req.GetId())
	if err != nil {
		return nil, err
	}

	card, err := s.q.GetFlashcard(ctx, dbgen.GetFlashcardParams{ID: id, OwnerID: ownerID})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.GetFlashcardResponse{Flashcard: flashcardToProto(card)}, nil
}

func (s *FlashcardService) ListFlashcards(ctx context.Context, req *flashcardv1.ListFlashcardsRequest) (*flashcardv1.ListFlashcardsResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	deckID, err := parseID("deck_id", req.GetDeckId())
	if err != nil {
		return nil, err
	}

	cards, err := s.q.ListFlashcards(ctx, dbgen.ListFlashcardsParams{DeckID: deckID, OwnerID: ownerID})
	if err != nil {
		return nil, translateError(err)
	}

	out := make([]*flashcardv1.Flashcard, 0, len(cards))
	for _, c := range cards {
		out = append(out, flashcardToProto(c))
	}
	return &flashcardv1.ListFlashcardsResponse{Flashcards: out}, nil
}

func (s *FlashcardService) UpdateFlashcard(ctx context.Context, req *flashcardv1.UpdateFlashcardRequest) (*flashcardv1.UpdateFlashcardResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID("id", req.GetId())
	if err != nil {
		return nil, err
	}
	front, err := requireText("front", req.GetFront(), maxCardTextLen)
	if err != nil {
		return nil, err
	}
	back, err := requireText("back", req.GetBack(), maxCardTextLen)
	if err != nil {
		return nil, err
	}
	if req.GetPosition() < 0 {
		return nil, status.Error(codes.InvalidArgument, "position must not be negative")
	}

	card, err := s.q.UpdateFlashcard(ctx, dbgen.UpdateFlashcardParams{
		ID:       id,
		OwnerID:  ownerID,
		Front:    front,
		Back:     back,
		Position: req.GetPosition(),
	})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.UpdateFlashcardResponse{Flashcard: flashcardToProto(card)}, nil
}

func (s *FlashcardService) DeleteFlashcard(ctx context.Context, req *flashcardv1.DeleteFlashcardRequest) (*flashcardv1.DeleteFlashcardResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID("id", req.GetId())
	if err != nil {
		return nil, err
	}

	rows, err := s.q.DeleteFlashcard(ctx, dbgen.DeleteFlashcardParams{ID: id, OwnerID: ownerID})
	if err != nil {
		return nil, translateError(err)
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, "flashcard not found")
	}
	return &flashcardv1.DeleteFlashcardResponse{}, nil
}
