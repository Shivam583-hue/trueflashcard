package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
	"github.com/Shivam583-hue/trueflashcard/server/internal/srs"
)

const (
	defaultNewLimit = 20
	maxNewLimit     = 100
	maxDueBatch     = 500
)

type StudyService struct {
	flashcardv1.UnimplementedStudyServiceServer
	q         dbgen.Querier
	tx        Transactor
	scheduler *srs.Scheduler
}

func NewStudyService(q dbgen.Querier, tx Transactor) *StudyService {
	return &StudyService{q: q, tx: tx, scheduler: srs.NewScheduler()}
}

func (s *StudyService) GetDueCards(ctx context.Context, req *flashcardv1.GetDueCardsRequest) (*flashcardv1.GetDueCardsResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	now := toPgTimestamp(time.Now().UTC())
	newLimit := clampNewLimit(req.GetNewLimit())

	var due, fresh []dbgen.Flashcard
	if deckIDStr := strings.TrimSpace(req.GetDeckId()); deckIDStr != "" {
		deckID, err := parseID("deck_id", deckIDStr)
		if err != nil {
			return nil, err
		}
		due, err = s.q.ListDueCardsInDeck(ctx, dbgen.ListDueCardsInDeckParams{OwnerID: ownerID, AsOf: now, DeckID: deckID, Lim: maxDueBatch})
		if err != nil {
			return nil, translateError(err)
		}
		fresh, err = s.q.ListNewCardsInDeck(ctx, dbgen.ListNewCardsInDeckParams{OwnerID: ownerID, DeckID: deckID, Lim: newLimit})
		if err != nil {
			return nil, translateError(err)
		}
	} else {
		due, err = s.q.ListDueCards(ctx, dbgen.ListDueCardsParams{OwnerID: ownerID, AsOf: now, Lim: maxDueBatch})
		if err != nil {
			return nil, translateError(err)
		}
		fresh, err = s.q.ListNewCards(ctx, dbgen.ListNewCardsParams{OwnerID: ownerID, Lim: newLimit})
		if err != nil {
			return nil, translateError(err)
		}
	}

	cards := make([]*flashcardv1.ReviewCard, 0, len(due)+len(fresh))
	for _, c := range due {
		cards = append(cards, &flashcardv1.ReviewCard{Card: flashcardToProto(c), IsNew: false})
	}
	for _, c := range fresh {
		cards = append(cards, &flashcardv1.ReviewCard{Card: flashcardToProto(c), IsNew: true})
	}
	return &flashcardv1.GetDueCardsResponse{Cards: cards}, nil
}

func (s *StudyService) SubmitReview(ctx context.Context, req *flashcardv1.SubmitReviewRequest) (*flashcardv1.SubmitReviewResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	cardID, err := parseID("card_id", req.GetCardId())
	if err != nil {
		return nil, err
	}
	rating := srs.Rating(req.GetRating())
	if !rating.Valid() {
		return nil, status.Error(codes.InvalidArgument, "rating must be again, hard, good, or easy")
	}

	now := time.Now().UTC()
	var result srs.Result
	err = s.tx.WithTx(ctx, func(q dbgen.Querier) error {
		if _, err := q.GetFlashcard(ctx, dbgen.GetFlashcardParams{ID: cardID, OwnerID: ownerID}); err != nil {
			return err
		}

		current, elapsed, err := loadCard(ctx, q, cardID, ownerID, now)
		if err != nil {
			return err
		}

		result = s.scheduler.Next(current, rating, now)
		next := result.Card

		if _, err := q.UpsertCardReviewState(ctx, dbgen.UpsertCardReviewStateParams{
			CardID:         cardID,
			OwnerID:        ownerID,
			State:          int16(next.State),
			DueAt:          toPgTimestamp(next.DueAt),
			Stability:      next.Stability,
			Difficulty:     next.Difficulty,
			Reps:           int32(next.Reps),
			Lapses:         int32(next.Lapses),
			LastReviewedAt: toPgTimestamp(next.LastReviewedAt),
		}); err != nil {
			return err
		}

		return q.InsertCardReview(ctx, dbgen.InsertCardReviewParams{
			CardID:        cardID,
			OwnerID:       ownerID,
			Rating:        int16(rating),
			StateBefore:   int16(current.State),
			ElapsedDays:   elapsed,
			Stability:     next.Stability,
			Difficulty:    next.Difficulty,
			ScheduledDays: result.IntervalDays,
		})
	})
	if err != nil {
		return nil, translateError(err)
	}

	return &flashcardv1.SubmitReviewResponse{
		DueAt:         timestamppb.New(result.Card.DueAt),
		ScheduledDays: int32(result.ScheduledDays),
	}, nil
}

func (s *StudyService) GetStudyOverview(ctx context.Context, _ *flashcardv1.GetStudyOverviewRequest) (*flashcardv1.GetStudyOverviewResponse, error) {
	ownerID, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	asOf := toPgTimestamp(now)

	dueTotal, err := s.q.CountDueTotal(ctx, dbgen.CountDueTotalParams{OwnerID: ownerID, AsOf: asOf})
	if err != nil {
		return nil, translateError(err)
	}
	newTotal, err := s.q.CountNewTotal(ctx, ownerID)
	if err != nil {
		return nil, translateError(err)
	}
	reviewedToday, err := s.q.CountReviewsSince(ctx, dbgen.CountReviewsSinceParams{OwnerID: ownerID, Since: toPgTimestamp(startOfUTCDay(now))})
	if err != nil {
		return nil, translateError(err)
	}
	dueByDeck, err := s.q.CountDueByDeck(ctx, dbgen.CountDueByDeckParams{OwnerID: ownerID, AsOf: asOf})
	if err != nil {
		return nil, translateError(err)
	}
	newByDeck, err := s.q.CountNewByDeck(ctx, ownerID)
	if err != nil {
		return nil, translateError(err)
	}
	days, err := s.q.ListReviewDays(ctx, ownerID)
	if err != nil {
		return nil, translateError(err)
	}

	decks := make(map[uuid.UUID]*flashcardv1.DeckDue)
	order := make([]uuid.UUID, 0, len(dueByDeck)+len(newByDeck))
	get := func(id uuid.UUID) *flashcardv1.DeckDue {
		if d, ok := decks[id]; ok {
			return d
		}
		d := &flashcardv1.DeckDue{DeckId: id.String()}
		decks[id] = d
		order = append(order, id)
		return d
	}
	for _, row := range dueByDeck {
		get(row.DeckID).Due = row.Due
	}
	for _, row := range newByDeck {
		get(row.DeckID).New = row.New
	}
	deckDues := make([]*flashcardv1.DeckDue, 0, len(order))
	for _, id := range order {
		deckDues = append(deckDues, decks[id])
	}

	return &flashcardv1.GetStudyOverviewResponse{
		DueTotal:      dueTotal,
		NewTotal:      newTotal,
		ReviewedToday: reviewedToday,
		StreakDays:    computeStreak(days, now),
		Decks:         deckDues,
	}, nil
}

func loadCard(ctx context.Context, q dbgen.Querier, cardID, ownerID uuid.UUID, now time.Time) (srs.Card, float64, error) {
	state, err := q.GetCardReviewState(ctx, dbgen.GetCardReviewStateParams{CardID: cardID, OwnerID: ownerID})
	if errors.Is(err, pgx.ErrNoRows) {
		return srs.Card{}, 0, nil
	}
	if err != nil {
		return srs.Card{}, 0, err
	}
	card := srs.Card{
		State:          srs.State(state.State),
		Stability:      state.Stability,
		Difficulty:     state.Difficulty,
		Reps:           int(state.Reps),
		Lapses:         int(state.Lapses),
		LastReviewedAt: fromPgTimestamp(state.LastReviewedAt),
		DueAt:          fromPgTimestamp(state.DueAt),
	}
	elapsed := 0.0
	if !card.LastReviewedAt.IsZero() {
		elapsed = math.Max(0, now.Sub(card.LastReviewedAt).Hours()/24)
	}
	return card, elapsed, nil
}

func clampNewLimit(v int32) int32 {
	if v <= 0 {
		return defaultNewLimit
	}
	if v > maxNewLimit {
		return maxNewLimit
	}
	return v
}

func startOfUTCDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func computeStreak(days []pgtype.Date, now time.Time) int32 {
	present := make(map[string]bool, len(days))
	for _, d := range days {
		if d.Valid {
			present[d.Time.UTC().Format("2006-01-02")] = true
		}
	}
	if len(present) == 0 {
		return 0
	}

	cursor := startOfUTCDay(now)
	if !present[cursor.Format("2006-01-02")] {
		cursor = cursor.AddDate(0, 0, -1)
		if !present[cursor.Format("2006-01-02")] {
			return 0
		}
	}

	var streak int32
	for present[cursor.Format("2006-01-02")] {
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}
	return streak
}
