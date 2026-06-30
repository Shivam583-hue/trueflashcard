"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { motion, useReducedMotion } from "motion/react";
import {
  ArrowLeftIcon,
  CardsThreeIcon,
  CheckCircleIcon,
  SparkleIcon,
} from "@phosphor-icons/react";

import { studyClient } from "@/lib/client";
import { Rating } from "@/gen/flashcard/v1/flashcard_pb";
import { isUnauthenticated } from "@/lib/connect-error";
import { easeEmphasized } from "@/lib/motion";
import { formatDueIn } from "@/lib/interval";
import { Button, ButtonLink } from "@/components/ui/button";
import { FlipCard, Centered } from "@/components/app/flip-card";
import { EmptyState, ErrorState } from "@/components/app/states";

type QueueCard = { id: string; front: string; back: string; isNew: boolean };
type Status = "loading" | "ready" | "error";

const ratings: { rating: Rating; label: string; key: string; tint: string }[] = [
  { rating: Rating.AGAIN, label: "Again", key: "1", tint: "text-rose-300/90" },
  { rating: Rating.HARD, label: "Hard", key: "2", tint: "text-amber-300/90" },
  { rating: Rating.GOOD, label: "Good", key: "3", tint: "text-emerald-300/90" },
  { rating: Rating.EASY, label: "Easy", key: "4", tint: "text-sky-300/90" },
];

export function StudyView({
  deckId,
  backHref,
  backLabel,
}: {
  deckId?: string;
  backHref: string;
  backLabel: string;
}) {
  const [status, setStatus] = useState<Status>("loading");
  const [needsAuth, setNeedsAuth] = useState(false);
  const [queue, setQueue] = useState<QueueCard[]>([]);
  const [total, setTotal] = useState(0);

  const fetchQueue = useCallback(
    () =>
      studyClient.getDueCards({ deckId: deckId ?? "" }).then((res) =>
        res.cards
          .filter((c) => c.card)
          .map((c) => ({
            id: c.card!.id,
            front: c.card!.front,
            back: c.card!.back,
            isNew: c.isNew,
          })),
      ),
    [deckId],
  );

  const apply = useCallback((cards: QueueCard[]) => {
    setQueue(cards);
    setTotal(cards.length);
    setStatus("ready");
  }, []);

  const fail = useCallback((err: unknown) => {
    if (isUnauthenticated(err)) setNeedsAuth(true);
    setStatus("error");
  }, []);

  const load = useCallback(() => {
    setStatus("loading");
    setNeedsAuth(false);
    fetchQueue().then(apply).catch(fail);
  }, [fetchQueue, apply, fail]);

  useEffect(() => {
    let active = true;
    fetchQueue()
      .then((cards) => active && apply(cards))
      .catch((err) => active && fail(err));
    return () => {
      active = false;
    };
  }, [fetchQueue, apply, fail]);

  if (status === "loading") {
    return (
      <Centered>
        <div className="skeleton h-64 w-full max-w-xl rounded-2xl border border-neutral-900 bg-[#0a0b0c]" />
      </Centered>
    );
  }

  if (status === "error") {
    return (
      <Centered>
        <div className="w-full max-w-xl">
          <ErrorState
            title="Could not start studying"
            needsAuth={needsAuth}
            onRetry={load}
          />
        </div>
      </Centered>
    );
  }

  if (total === 0) {
    return (
      <Centered>
        <div className="w-full max-w-xl">
          <EmptyState
            icon={CheckCircleIcon}
            title="All caught up"
            body="No cards are due right now. Come back later or add more cards."
            action={<ButtonLink href={backHref}>{backLabel}</ButtonLink>}
          />
        </div>
      </Centered>
    );
  }

  return (
    <Session
      initial={queue}
      total={total}
      backHref={backHref}
      backLabel={backLabel}
    />
  );
}

function Session({
  initial,
  total,
  backHref,
  backLabel,
}: {
  initial: QueueCard[];
  total: number;
  backHref: string;
  backLabel: string;
}) {
  const reduce = useReducedMotion();
  const [queue, setQueue] = useState<QueueCard[]>(initial);
  const [flipped, setFlipped] = useState(false);
  const [graduated, setGraduated] = useState(0);
  const [lastInterval, setLastInterval] = useState<string | null>(null);

  const current = queue[0];
  const completed = queue.length === 0;

  const rate = useCallback(
    (rating: Rating) => {
      const card = queue[0];
      if (!card) return;
      setFlipped(false);
      studyClient
        .submitReview({ cardId: card.id, rating })
        .then((res) => setLastInterval(formatDueIn(res.dueAt)))
        .catch(() => {});
      if (rating === Rating.AGAIN) {
        setQueue((q) => [...q.slice(1), card]);
      } else {
        setGraduated((g) => g + 1);
        setQueue((q) => q.slice(1));
      }
    },
    [queue],
  );

  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (completed) return;
      if (e.key === " " || e.key === "Enter") {
        e.preventDefault();
        setFlipped((f) => !f);
        return;
      }
      if (flipped) {
        const match = ratings.find((r) => r.key === e.key);
        if (match) {
          e.preventDefault();
          rate(match.rating);
        }
      }
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [flipped, completed, rate]);

  const fraction = total === 0 ? 0 : graduated / total;

  return (
    <div className="mx-auto flex min-h-[calc(100vh-8rem)] max-w-xl flex-col">
      <div className="flex items-center justify-between">
        <Link
          href={backHref}
          className="inline-flex items-center gap-1.5 text-sm text-neutral-500 transition-colors duration-150 hover:text-neutral-200"
        >
          <ArrowLeftIcon size={15} />
          {backLabel}
        </Link>
        <span className="text-sm tabular-nums text-neutral-500">
          {Math.min(graduated + (completed ? 0 : 1), total)} / {total}
        </span>
      </div>

      <div className="flex flex-1 flex-col items-center justify-center py-8">
        {completed ? (
          <Completion graduated={graduated} backHref={backHref} backLabel={backLabel} />
        ) : (
          <motion.div
            key={current.id}
            initial={reduce ? false : { opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.25, ease: easeEmphasized }}
            className="w-full"
          >
            {current.isNew && (
              <div className="mb-3 flex justify-center">
                <span className="inline-flex items-center gap-1.5 rounded-full border border-neutral-800 bg-neutral-900/60 px-2.5 py-1 text-[11px] font-medium text-neutral-400">
                  <SparkleIcon size={12} weight="fill" />
                  New card
                </span>
              </div>
            )}
            <FlipCard
              front={current.front}
              back={current.back}
              flipped={flipped}
              onFlip={() => setFlipped((f) => !f)}
            />
          </motion.div>
        )}
      </div>

      {!completed && (
        <div className="min-h-[88px]">
          {flipped ? (
            <div className="grid grid-cols-4 gap-2">
              {ratings.map((r) => (
                <button
                  key={r.rating}
                  onClick={() => rate(r.rating)}
                  className="flex flex-col items-center gap-1 rounded-lg border border-neutral-800 bg-neutral-900/50 px-2 py-3 text-sm font-medium text-neutral-100 transition-[transform,background-color,border-color] duration-150 [transition-timing-function:var(--ease-out)] hover:border-neutral-700 hover:bg-neutral-800/60 active:scale-[0.97] focus:outline-none focus-visible:ring-2 focus-visible:ring-neutral-500"
                >
                  <span className={r.tint}>{r.label}</span>
                  <span className="text-[11px] tabular-nums text-neutral-600">{r.key}</span>
                </button>
              ))}
            </div>
          ) : (
            <div className="flex flex-col items-center gap-3">
              <Button onClick={() => setFlipped(true)} className="w-full max-w-xs">
                Show answer
              </Button>
              <span className="text-xs text-neutral-600">
                {reduce ? "Tap card to flip" : "Space to flip · 1-4 to rate"}
              </span>
            </div>
          )}
        </div>
      )}

      <div className="mt-6">
        <div className="h-1 w-full overflow-hidden rounded-full bg-neutral-900">
          <div
            className="h-full origin-left rounded-full bg-neutral-300 transition-transform duration-500 [transition-timing-function:var(--ease-out)] motion-reduce:transition-none"
            style={{ transform: `scaleX(${fraction})` }}
          />
        </div>
        <div className="mt-2.5 flex items-center justify-between text-xs text-neutral-500">
          <span className="tabular-nums">{graduated} reviewed</span>
          <span className="tabular-nums">
            {lastInterval ? `Next in ${lastInterval}` : `${queue.length} in queue`}
          </span>
        </div>
      </div>
    </div>
  );
}

function Completion({
  graduated,
  backHref,
  backLabel,
}: {
  graduated: number;
  backHref: string;
  backLabel: string;
}) {
  const reduce = useReducedMotion();
  return (
    <motion.div
      initial={reduce ? false : { opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, ease: easeEmphasized }}
      className="flex flex-col items-center text-center"
    >
      <span className="flex h-12 w-12 items-center justify-center rounded-xl bg-neutral-900 text-emerald-300/90">
        <CheckCircleIcon size={26} weight="duotone" />
      </span>
      <h2 className="mt-4 text-2xl font-semibold tracking-tight text-neutral-50">
        Session complete
      </h2>
      <p className="mt-2 text-sm text-neutral-400">
        You reviewed {graduated} {graduated === 1 ? "card" : "cards"}. Spacing is handled for you.
      </p>
      <div className="mt-7 flex items-center gap-3">
        <ButtonLink href={backHref}>
          <CardsThreeIcon size={16} />
          {backLabel}
        </ButtonLink>
      </div>
    </motion.div>
  );
}
