package service

import (
	"context"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FolderService struct {
	flashcardv1.UnimplementedFolderServiceServer
	q dbgen.Querier
}

func NewFolderService(q dbgen.Querier) *FolderService {
	return &FolderService{q: q}
}

func (s *FolderService) CreateFolder(ctx context.Context, req *flashcardv1.CreateFolderRequest) (*flashcardv1.CreateFolderResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	name, err := requireText("name", req.GetName(), maxNameLen)
	if err != nil {
		return nil, err
	}

	folder, err := s.q.CreateFolder(ctx, dbgen.CreateFolderParams{OwnerID: ownerID, Name: name})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.CreateFolderResponse{Folder: folderToProto(folder)}, nil
}

func (s *FolderService) GetFolder(ctx context.Context, req *flashcardv1.GetFolderRequest) (*flashcardv1.GetFolderResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID("id", req.GetId())
	if err != nil {
		return nil, err
	}

	folder, err := s.q.GetFolder(ctx, dbgen.GetFolderParams{ID: id, OwnerID: ownerID})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.GetFolderResponse{Folder: folderToProto(folder)}, nil
}

func (s *FolderService) ListFolders(ctx context.Context, _ *flashcardv1.ListFoldersRequest) (*flashcardv1.ListFoldersResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}

	folders, err := s.q.ListFolders(ctx, ownerID)
	if err != nil {
		return nil, translateError(err)
	}

	out := make([]*flashcardv1.Folder, 0, len(folders))
	for _, f := range folders {
		out = append(out, folderToProto(f))
	}
	return &flashcardv1.ListFoldersResponse{Folders: out}, nil
}

func (s *FolderService) UpdateFolder(ctx context.Context, req *flashcardv1.UpdateFolderRequest) (*flashcardv1.UpdateFolderResponse, error) {
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

	folder, err := s.q.UpdateFolder(ctx, dbgen.UpdateFolderParams{ID: id, OwnerID: ownerID, Name: name})
	if err != nil {
		return nil, translateError(err)
	}
	return &flashcardv1.UpdateFolderResponse{Folder: folderToProto(folder)}, nil
}

func (s *FolderService) DeleteFolder(ctx context.Context, req *flashcardv1.DeleteFolderRequest) (*flashcardv1.DeleteFolderResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseID("id", req.GetId())
	if err != nil {
		return nil, err
	}

	rows, err := s.q.DeleteFolder(ctx, dbgen.DeleteFolderParams{ID: id, OwnerID: ownerID})
	if err != nil {
		return nil, translateError(err)
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, "folder not found")
	}
	return &flashcardv1.DeleteFolderResponse{}, nil
}
