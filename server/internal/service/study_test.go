package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
)

type studyStub struct {
	dbgen.Querier

	getFlashcardFn  func(context.Context, dbgen.GetFlashcardParams) (dbgen.Flashcard, error)
	getStateFn      func(context.Context, dbgen.GetCardReviewStateParams) (dbgen.CardReviewState, error)
	upsertStateFn   func(context.Context, dbgen.UpsertCardReviewStateParams) (dbgen.CardReviewState, error)
	insertReviewFn  func(context.Context, dbgen.InsertCardReviewParams) error
	listDueFn       func(context.Context, dbgen.ListDueCardsParams) ([]dbgen.Flashcard, error)
	listNewFn       func(context.Context, dbgen.ListNewCardsParams) ([]dbgen.Flashcard, error)
	countDueTotalFn func(context.Context, dbgen.CountDueTotalParams) (int32, error)
	countNewTotalFn func(context.Context, uuid.UUID) (int32, error)
	countReviewsFn  func(context.Context, dbgen.CountReviewsSinceParams) (int32, error)
	countDueDeckFn  func(context.Context, dbgen.CountDueByDeckParams) ([]dbgen.CountDueByDeckRow, error)
	countNewDeckFn  func(context.Context, uuid.UUID) ([]dbgen.CountNewByDeckRow, error)
	listReviewDayFn func(context.Context, uuid.UUID) ([]pgtype.Date, error)
	upsertedState   dbgen.UpsertCardReviewStateParams
	insertedReview  dbgen.InsertCardReviewParams
}

func (s *studyStub) GetFlashcard(ctx context.Context, p dbgen.GetFlashcardParams) (dbgen.Flashcard, error) {
	return s.getFlashcardFn(ctx, p)
}

func (s *studyStub) GetCardReviewState(ctx context.Context, p dbgen.GetCardReviewStateParams) (dbgen.CardReviewState, error) {
	return s.getStateFn(ctx, p)
}

func (s *studyStub) UpsertCardReviewState(ctx context.Context, p dbgen.UpsertCardReviewStateParams) (dbgen.CardReviewState, error) {
	s.upsertedState = p
	if s.upsertStateFn != nil {
		return s.upsertStateFn(ctx, p)
	}
	return dbgen.CardReviewState{}, nil
}

func (s *studyStub) InsertCardReview(ctx context.Context, p dbgen.InsertCardReviewParams) error {
	s.insertedReview = p
	if s.insertReviewFn != nil {
		return s.insertReviewFn(ctx, p)
	}
	return nil
}

func (s *studyStub) ListDueCards(ctx context.Context, p dbgen.ListDueCardsParams) ([]dbgen.Flashcard, error) {
	return s.listDueFn(ctx, p)
}

func (s *studyStub) ListNewCards(ctx context.Context, p dbgen.ListNewCardsParams) ([]dbgen.Flashcard, error) {
	return s.listNewFn(ctx, p)
}

func (s *studyStub) CountDueTotal(ctx context.Context, p dbgen.CountDueTotalParams) (int32, error) {
	return s.countDueTotalFn(ctx, p)
}

func (s *studyStub) CountNewTotal(ctx context.Context, id uuid.UUID) (int32, error) {
	return s.countNewTotalFn(ctx, id)
}

func (s *studyStub) CountReviewsSince(ctx context.Context, p dbgen.CountReviewsSinceParams) (int32, error) {
	return s.countReviewsFn(ctx, p)
}

func (s *studyStub) CountDueByDeck(ctx context.Context, p dbgen.CountDueByDeckParams) ([]dbgen.CountDueByDeckRow, error) {
	return s.countDueDeckFn(ctx, p)
}

func (s *studyStub) CountNewByDeck(ctx context.Context, id uuid.UUID) ([]dbgen.CountNewByDeckRow, error) {
	return s.countNewDeckFn(ctx, id)
}

func (s *studyStub) ListReviewDays(ctx context.Context, id uuid.UUID) ([]pgtype.Date, error) {
	return s.listReviewDayFn(ctx, id)
}

func TestSubmitReviewNewCardSchedulesFuture(t *testing.T) {
	ctx, _ := authedContext()
	cardID := uuid.New()
	stub := &studyStub{
		getFlashcardFn: func(_ context.Context, _ dbgen.GetFlashcardParams) (dbgen.Flashcard, error) {
			return dbgen.Flashcard{ID: cardID}, nil
		},
		getStateFn: func(_ context.Context, _ dbgen.GetCardReviewStateParams) (dbgen.CardReviewState, error) {
			return dbgen.CardReviewState{}, pgx.ErrNoRows
		},
	}
	svc := NewStudyService(stub, directTx{stub})

	resp, err := svc.SubmitReview(ctx, &flashcardv1.SubmitReviewRequest{
		CardId: cardID.String(),
		Rating: flashcardv1.Rating_RATING_GOOD,
	})
	requireNoError(t, err)
	if resp.GetScheduledDays() <= 0 {
		t.Fatalf("expected positive interval, got %d", resp.GetScheduledDays())
	}
	if !stub.upsertedState.DueAt.Valid || !stub.upsertedState.DueAt.Time.After(stub.upsertedState.LastReviewedAt.Time) {
		t.Fatalf("expected due_at after last_reviewed_at")
	}
	if stub.insertedReview.Rating != int16(flashcardv1.Rating_RATING_GOOD) {
		t.Fatalf("review log rating mismatch: %d", stub.insertedReview.Rating)
	}
}

func TestSubmitReviewRejectsInvalidRating(t *testing.T) {
	ctx, _ := authedContext()
	svc := NewStudyService(&studyStub{}, directTx{&studyStub{}})

	_, err := svc.SubmitReview(ctx, &flashcardv1.SubmitReviewRequest{
		CardId: uuid.NewString(),
		Rating: flashcardv1.Rating_RATING_UNSPECIFIED,
	})
	requireCode(t, err, codes.InvalidArgument)
}

func TestGetDueCardsMarksNewCards(t *testing.T) {
	ctx, _ := authedContext()
	due := dbgen.Flashcard{ID: uuid.New(), Front: "due"}
	fresh := dbgen.Flashcard{ID: uuid.New(), Front: "fresh"}
	stub := &studyStub{
		listDueFn: func(_ context.Context, _ dbgen.ListDueCardsParams) ([]dbgen.Flashcard, error) {
			return []dbgen.Flashcard{due}, nil
		},
		listNewFn: func(_ context.Context, _ dbgen.ListNewCardsParams) ([]dbgen.Flashcard, error) {
			return []dbgen.Flashcard{fresh}, nil
		},
	}
	svc := NewStudyService(stub, directTx{stub})

	resp, err := svc.GetDueCards(ctx, &flashcardv1.GetDueCardsRequest{})
	requireNoError(t, err)
	if len(resp.GetCards()) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(resp.GetCards()))
	}
	if resp.GetCards()[0].GetIsNew() || !resp.GetCards()[1].GetIsNew() {
		t.Fatalf("expected first card due and second new")
	}
}

func TestGetStudyOverviewMergesDeckCounts(t *testing.T) {
	ctx, _ := authedContext()
	deckA := uuid.New()
	stub := &studyStub{
		countDueTotalFn: func(_ context.Context, _ dbgen.CountDueTotalParams) (int32, error) { return 3, nil },
		countNewTotalFn: func(_ context.Context, _ uuid.UUID) (int32, error) { return 5, nil },
		countReviewsFn:  func(_ context.Context, _ dbgen.CountReviewsSinceParams) (int32, error) { return 7, nil },
		countDueDeckFn: func(_ context.Context, _ dbgen.CountDueByDeckParams) ([]dbgen.CountDueByDeckRow, error) {
			return []dbgen.CountDueByDeckRow{{DeckID: deckA, Due: 3}}, nil
		},
		countNewDeckFn: func(_ context.Context, _ uuid.UUID) ([]dbgen.CountNewByDeckRow, error) {
			return []dbgen.CountNewByDeckRow{{DeckID: deckA, New: 5}}, nil
		},
		listReviewDayFn: func(_ context.Context, _ uuid.UUID) ([]pgtype.Date, error) {
			return nil, nil
		},
	}
	svc := NewStudyService(stub, directTx{stub})

	resp, err := svc.GetStudyOverview(ctx, &flashcardv1.GetStudyOverviewRequest{})
	requireNoError(t, err)
	if resp.GetDueTotal() != 3 || resp.GetNewTotal() != 5 || resp.GetReviewedToday() != 7 {
		t.Fatalf("unexpected totals: %+v", resp)
	}
	if len(resp.GetDecks()) != 1 || resp.GetDecks()[0].GetDue() != 3 || resp.GetDecks()[0].GetNew() != 5 {
		t.Fatalf("deck counts not merged: %+v", resp.GetDecks())
	}
}
