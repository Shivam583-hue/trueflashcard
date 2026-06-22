package service

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
)

func toTimestamp(ts pgtype.Timestamptz) *timestamppb.Timestamp {
	if !ts.Valid {
		return nil
	}
	return timestamppb.New(ts.Time)
}

func folderToProto(f dbgen.Folder) *flashcardv1.Folder {
	return &flashcardv1.Folder{
		Id:        f.ID.String(),
		OwnerId:   f.OwnerID.String(),
		Name:      f.Name,
		CreatedAt: toTimestamp(f.CreatedAt),
		UpdatedAt: toTimestamp(f.UpdatedAt),
	}
}

func buildDeck(id, folderID uuid.UUID, name, description string, created, updated pgtype.Timestamptz, cardCount int32) *flashcardv1.Deck {
	return &flashcardv1.Deck{
		Id:          id.String(),
		FolderId:    folderID.String(),
		Name:        name,
		Description: description,
		CardCount:   cardCount,
		CreatedAt:   toTimestamp(created),
		UpdatedAt:   toTimestamp(updated),
	}
}

func deckToProto(d dbgen.Deck, cardCount int32) *flashcardv1.Deck {
	return buildDeck(d.ID, d.FolderID, d.Name, d.Description, d.CreatedAt, d.UpdatedAt, cardCount)
}

func getDeckRowToProto(d dbgen.GetDeckRow) *flashcardv1.Deck {
	return buildDeck(d.ID, d.FolderID, d.Name, d.Description, d.CreatedAt, d.UpdatedAt, d.CardCount)
}

func listDeckRowToProto(d dbgen.ListDecksRow) *flashcardv1.Deck {
	return buildDeck(d.ID, d.FolderID, d.Name, d.Description, d.CreatedAt, d.UpdatedAt, d.CardCount)
}

func updateDeckRowToProto(d dbgen.UpdateDeckRow) *flashcardv1.Deck {
	return buildDeck(d.ID, d.FolderID, d.Name, d.Description, d.CreatedAt, d.UpdatedAt, d.CardCount)
}

func flashcardToProto(c dbgen.Flashcard) *flashcardv1.Flashcard {
	return &flashcardv1.Flashcard{
		Id:        c.ID.String(),
		DeckId:    c.DeckID.String(),
		Front:     c.Front,
		Back:      c.Back,
		Position:  c.Position,
		CreatedAt: toTimestamp(c.CreatedAt),
		UpdatedAt: toTimestamp(c.UpdatedAt),
	}
}
