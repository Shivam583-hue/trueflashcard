package srs

import (
	"math"
	"time"
)

type Rating int

const (
	RatingAgain Rating = 1
	RatingHard  Rating = 2
	RatingGood  Rating = 3
	RatingEasy  Rating = 4
)

func (r Rating) Valid() bool { return r >= RatingAgain && r <= RatingEasy }

type State int

const (
	StateNew        State = 0
	StateLearning   State = 1
	StateReview     State = 2
	StateRelearning State = 3
)

const (
	decay  = -0.5
	factor = 19.0 / 81.0

	minDifficulty = 1.0
	maxDifficulty = 10.0
	minStability  = 0.1
)

type Params struct {
	W                [19]float64
	RequestRetention float64
	MaximumInterval  float64
	MinInterval      time.Duration
}

func DefaultParams() Params {
	return Params{
		W: [19]float64{
			0.40255, 1.18385, 3.173, 15.69105, 7.1949, 0.5345, 1.4604,
			0.0046, 1.54575, 0.1192, 1.01925, 1.9395, 0.11, 0.29605,
			2.2698, 0.2315, 2.9898, 0.51655, 0.6621,
		},
		RequestRetention: 0.9,
		MaximumInterval:  365,
		MinInterval:      time.Minute,
	}
}

type Scheduler struct {
	p Params
}

func NewScheduler() *Scheduler { return &Scheduler{p: DefaultParams()} }

func NewSchedulerWith(p Params) *Scheduler { return &Scheduler{p: p} }

type Card struct {
	State          State
	Stability      float64
	Difficulty     float64
	Reps           int
	Lapses         int
	LastReviewedAt time.Time
	DueAt          time.Time
}

type Result struct {
	Card          Card
	IntervalDays  float64
	ScheduledDays int
}

func (s *Scheduler) Next(card Card, rating Rating, now time.Time) Result {
	w := s.p.W
	next := card
	next.Reps = card.Reps + 1

	if card.Reps == 0 || card.LastReviewedAt.IsZero() {
		next.Difficulty = clampDifficulty(w[4] - math.Exp(w[5]*float64(rating-1)) + 1)
		next.Stability = math.Max(w[rating-1], minStability)
		next.State = firstState(rating)
		if rating == RatingAgain {
			next.Lapses = card.Lapses + 1
		}
	} else {
		elapsed := math.Max(0, now.Sub(card.LastReviewedAt).Hours()/24)
		retrievability := math.Pow(1+factor*elapsed/card.Stability, decay)
		next.Difficulty = s.nextDifficulty(card.Difficulty, rating)
		if rating == RatingAgain {
			next.Stability = s.forgetStability(next.Difficulty, card.Stability, retrievability)
			next.Lapses = card.Lapses + 1
		} else {
			next.Stability = s.recallStability(next.Difficulty, card.Stability, retrievability, rating)
		}
		next.State = nextState(card.State, rating)
	}

	next.Stability = math.Max(next.Stability, minStability)
	interval := s.interval(next.Stability)

	duration := max(time.Duration(interval*float64(24*time.Hour)), s.p.MinInterval)
	next.LastReviewedAt = now
	next.DueAt = now.Add(duration)

	return Result{
		Card:          next,
		IntervalDays:  interval,
		ScheduledDays: int(math.Round(interval)),
	}
}

func (s *Scheduler) interval(stability float64) float64 {
	iv := stability / factor * (math.Pow(s.p.RequestRetention, 1/decay) - 1)
	if iv < 0 {
		return 0
	}
	if iv > s.p.MaximumInterval {
		return s.p.MaximumInterval
	}
	return iv
}

func (s *Scheduler) nextDifficulty(difficulty float64, rating Rating) float64 {
	w := s.p.W
	deltaD := -w[6] * float64(rating-RatingGood)
	damped := difficulty + deltaD*(maxDifficulty-difficulty)/9
	target := w[4] - math.Exp(w[5]*float64(RatingEasy-1)) + 1
	reverted := w[7]*target + (1-w[7])*damped
	return clampDifficulty(reverted)
}

func (s *Scheduler) recallStability(difficulty, stability, retrievability float64, rating Rating) float64 {
	w := s.p.W
	hardPenalty := 1.0
	if rating == RatingHard {
		hardPenalty = w[15]
	}
	easyBonus := 1.0
	if rating == RatingEasy {
		easyBonus = w[16]
	}
	growth := math.Exp(w[8]) *
		(11 - difficulty) *
		math.Pow(stability, -w[9]) *
		(math.Exp((1-retrievability)*w[10]) - 1) *
		hardPenalty *
		easyBonus
	return stability * (1 + growth)
}

func (s *Scheduler) forgetStability(difficulty, stability, retrievability float64) float64 {
	w := s.p.W
	longTerm := w[11] *
		math.Pow(difficulty, -w[12]) *
		(math.Pow(stability+1, w[13]) - 1) *
		math.Exp((1-retrievability)*w[14])
	shortTerm := stability / math.Exp(w[17]*w[18])
	return math.Min(longTerm, shortTerm)
}

func clampDifficulty(d float64) float64 {
	if d < minDifficulty {
		return minDifficulty
	}
	if d > maxDifficulty {
		return maxDifficulty
	}
	return d
}

func firstState(rating Rating) State {
	if rating == RatingAgain || rating == RatingHard {
		return StateLearning
	}
	return StateReview
}

func nextState(prev State, rating Rating) State {
	if rating == RatingAgain {
		if prev == StateReview || prev == StateRelearning {
			return StateRelearning
		}
		return StateLearning
	}
	if rating == RatingHard && (prev == StateLearning || prev == StateRelearning) {
		return prev
	}
	return StateReview
}
