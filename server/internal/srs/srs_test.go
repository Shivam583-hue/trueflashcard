package srs

import (
	"math"
	"testing"
	"time"
)

var epoch = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func firstReview(rating Rating) Result {
	return NewScheduler().Next(Card{}, rating, epoch)
}

func TestFirstReviewIntervalsAreOrdered(t *testing.T) {
	again := firstReview(RatingAgain).IntervalDays
	hard := firstReview(RatingHard).IntervalDays
	good := firstReview(RatingGood).IntervalDays
	easy := firstReview(RatingEasy).IntervalDays

	if !(again < hard && hard < good && good < easy) {
		t.Fatalf("intervals not strictly ordered: again=%.3f hard=%.3f good=%.3f easy=%.3f", again, hard, good, easy)
	}
}

func TestFirstGoodIntervalMatchesInitialStability(t *testing.T) {
	res := firstReview(RatingGood)
	want := DefaultParams().W[RatingGood-1]
	if math.Abs(res.IntervalDays-want) > 1e-6 {
		t.Fatalf("interval %.6f, want %.6f (interval equals stability at retention 0.9)", res.IntervalDays, want)
	}
}

func TestRepeatedGoodGrowsStability(t *testing.T) {
	s := NewScheduler()
	card := s.Next(Card{}, RatingGood, epoch).Card

	prev := card.Stability
	now := epoch
	for i := range 5 {
		now = card.DueAt
		card = s.Next(card, RatingGood, now).Card
		if card.Stability <= prev {
			t.Fatalf("stability did not grow on rep %d: %.4f <= %.4f", i, card.Stability, prev)
		}
		prev = card.Stability
	}
}

func TestAgainShrinksStabilityAndCountsLapse(t *testing.T) {
	s := NewScheduler()
	good := s.Next(Card{}, RatingGood, epoch).Card

	again := s.Next(good, RatingAgain, good.DueAt)
	if again.Card.Stability >= good.Stability {
		t.Fatalf("again should shrink stability: %.4f >= %.4f", again.Card.Stability, good.Stability)
	}
	if again.Card.Lapses != 1 {
		t.Fatalf("expected 1 lapse, got %d", again.Card.Lapses)
	}
	if again.Card.DueAt.Sub(good.DueAt) <= 0 {
		t.Fatalf("again should still schedule a future due date")
	}
}

func TestDifficultyStaysInRange(t *testing.T) {
	s := NewScheduler()
	ratings := []Rating{RatingAgain, RatingAgain, RatingAgain, RatingEasy, RatingEasy, RatingHard, RatingGood}
	card := Card{}
	now := epoch
	for i, r := range ratings {
		res := s.Next(card, r, now)
		card = res.Card
		now = card.DueAt
		if card.Difficulty < minDifficulty-1e-9 || card.Difficulty > maxDifficulty+1e-9 {
			t.Fatalf("difficulty out of range after step %d: %.4f", i, card.Difficulty)
		}
	}
}

func TestEasyBeatsGoodBeatsHardOnReview(t *testing.T) {
	s := NewScheduler()
	base := s.Next(Card{}, RatingGood, epoch).Card
	at := base.DueAt

	hard := s.Next(base, RatingHard, at).Card.Stability
	good := s.Next(base, RatingGood, at).Card.Stability
	easy := s.Next(base, RatingEasy, at).Card.Stability

	if !(hard < good && good < easy) {
		t.Fatalf("recall stability not ordered: hard=%.4f good=%.4f easy=%.4f", hard, good, easy)
	}
}

func TestMinIntervalFloor(t *testing.T) {
	s := NewScheduler()
	res := s.Next(Card{}, RatingAgain, epoch)
	if res.Card.DueAt.Sub(epoch) < DefaultParams().MinInterval {
		t.Fatalf("due date below minimum interval floor")
	}
}
