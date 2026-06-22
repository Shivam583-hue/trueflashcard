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

func TestCreateFolderSuccess(t *testing.T) {
	ctx, userID := authedContext()
	var captured dbgen.CreateFolderParams
	folderID := uuid.New()
	q := &stubQuerier{
		createFolderFn: func(_ context.Context, p dbgen.CreateFolderParams) (dbgen.Folder, error) {
			captured = p
			return dbgen.Folder{
				ID:        folderID,
				OwnerID:   userID,
				Name:      p.Name,
				CreatedAt: validTimestamp(),
				UpdatedAt: validTimestamp(),
			}, nil
		},
	}

	resp, err := NewFolderService(q).CreateFolder(ctx, &flashcardv1.CreateFolderRequest{Name: "  Biology  "})
	requireNoError(t, err)

	if captured.OwnerID != userID {
		t.Fatalf("owner id not scoped to caller: got %v want %v", captured.OwnerID, userID)
	}
	if captured.Name != "Biology" {
		t.Fatalf("name not trimmed: got %q", captured.Name)
	}
	if resp.GetFolder().GetId() != folderID.String() {
		t.Fatalf("unexpected folder id: %q", resp.GetFolder().GetId())
	}
	if resp.GetFolder().GetCreatedAt() == nil {
		t.Fatal("expected created_at to be set")
	}
}

func TestCreateFolderEmptyName(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	_, err := NewFolderService(q).CreateFolder(ctx, &flashcardv1.CreateFolderRequest{Name: "   "})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestCreateFolderUnauthenticated(t *testing.T) {
	q := &stubQuerier{}
	_, err := NewFolderService(q).CreateFolder(context.Background(), &flashcardv1.CreateFolderRequest{Name: "Biology"})
	requireCode(t, err, codes.Unauthenticated)
	requireNoCalls(t, q)
}

func TestGetFolderInvalidID(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	_, err := NewFolderService(q).GetFolder(ctx, &flashcardv1.GetFolderRequest{Id: "not-a-uuid"})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestGetFolderNotFound(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		getFolderFn: func(_ context.Context, _ dbgen.GetFolderParams) (dbgen.Folder, error) {
			return dbgen.Folder{}, pgx.ErrNoRows
		},
	}
	_, err := NewFolderService(q).GetFolder(ctx, &flashcardv1.GetFolderRequest{Id: uuid.NewString()})
	requireCode(t, err, codes.NotFound)
}

func TestListFoldersSuccess(t *testing.T) {
	ctx, userID := authedContext()
	var capturedOwner uuid.UUID
	q := &stubQuerier{
		listFoldersFn: func(_ context.Context, ownerID uuid.UUID) ([]dbgen.Folder, error) {
			capturedOwner = ownerID
			return []dbgen.Folder{
				{ID: uuid.New(), OwnerID: ownerID, Name: "A"},
				{ID: uuid.New(), OwnerID: ownerID, Name: "B"},
			}, nil
		},
	}
	resp, err := NewFolderService(q).ListFolders(ctx, &flashcardv1.ListFoldersRequest{})
	requireNoError(t, err)
	if capturedOwner != userID {
		t.Fatalf("list not scoped to caller")
	}
	if len(resp.GetFolders()) != 2 {
		t.Fatalf("expected 2 folders, got %d", len(resp.GetFolders()))
	}
}

func TestUpdateFolderEmptyName(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{}
	_, err := NewFolderService(q).UpdateFolder(ctx, &flashcardv1.UpdateFolderRequest{Id: uuid.NewString(), Name: ""})
	requireCode(t, err, codes.InvalidArgument)
	requireNoCalls(t, q)
}

func TestDeleteFolderNotFound(t *testing.T) {
	ctx, _ := authedContext()
	q := &stubQuerier{
		deleteFolderFn: func(_ context.Context, _ dbgen.DeleteFolderParams) (int64, error) {
			return 0, nil
		},
	}
	_, err := NewFolderService(q).DeleteFolder(ctx, &flashcardv1.DeleteFolderRequest{Id: uuid.NewString()})
	requireCode(t, err, codes.NotFound)
}

func TestDeleteFolderSuccess(t *testing.T) {
	ctx, userID := authedContext()
	var captured dbgen.DeleteFolderParams
	q := &stubQuerier{
		deleteFolderFn: func(_ context.Context, p dbgen.DeleteFolderParams) (int64, error) {
			captured = p
			return 1, nil
		},
	}
	_, err := NewFolderService(q).DeleteFolder(ctx, &flashcardv1.DeleteFolderRequest{Id: uuid.NewString()})
	requireNoError(t, err)
	if captured.OwnerID != userID {
		t.Fatalf("delete not scoped to caller")
	}
}
