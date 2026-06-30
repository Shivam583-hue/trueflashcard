"use client";

import { useCallback, useEffect, useState } from "react";

import { studyClient } from "@/lib/client";

export type DeckCounts = { due: number; new: number };

export type StudyOverview = {
  dueTotal: number;
  newTotal: number;
  reviewedToday: number;
  streakDays: number;
  byDeck: Map<string, DeckCounts>;
};

export function useStudyOverview() {
  const [overview, setOverview] = useState<StudyOverview | null>(null);

  const reload = useCallback(() => {
    studyClient
      .getStudyOverview({})
      .then((res) => {
        const byDeck = new Map<string, DeckCounts>();
        for (const d of res.decks) byDeck.set(d.deckId, { due: d.due, new: d.new });
        setOverview({
          dueTotal: res.dueTotal,
          newTotal: res.newTotal,
          reviewedToday: res.reviewedToday,
          streakDays: res.streakDays,
          byDeck,
        });
      })
      .catch(() => {});
  }, []);

  useEffect(reload, [reload]);

  return { overview, reload };
}
