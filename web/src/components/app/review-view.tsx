"use client";

import { useCallback, useEffect, useReducer, useState } from "react";
import Link from "next/link";
import { motion, useReducedMotion } from "motion/react";
import {
  ArrowLeftIcon,
  ArrowRightIcon,
  ArrowUDownLeftIcon,
  CardsThreeIcon,
} from "@phosphor-icons/react";

import { deckClient, flashcardClient } from "@/lib/client";
import { useCollection } from "@/lib/use-collection";
import { easeEmphasized } from "@/lib/motion";
import { Button, ButtonLink } from "@/components/ui/button";
import { EmptyState, ErrorState } from "@/components/app/states";

type Card = { id: string; front: string; back: string };

type SessionState = { index: number; flipped: boolean; completed: boolean };
type SessionAction =
  | { type: "flip" }
  | { type: "next"; total: number }
  | { type: "prev"; total: number }
  | { type: "restart" };

function sessionReducer(state: SessionState, action: SessionAction): SessionState {
  switch (action.type) {
    case "flip":
      return state.completed ? state : { ...state, flipped: !state.flipped };
    case "next":
      if (state.completed) return state;
      if (state.index >= action.total - 1) {
        return { ...state, flipped: false, completed: true };
      }
      return { index: state.index + 1, flipped: false, completed: false };
    case "prev":
      if (state.completed) {
        return { index: action.total - 1, flipped: false, completed: false };
      }
      if (state.index === 0) return state;
      return { index: state.index - 1, flipped: false, completed: false };
    case "restart":
      return { index: 0, flipped: false, completed: false };
  }
}

const initialSession: SessionState = {
  index: 0,
  flipped: false,
  completed: false,
};

export function ReviewView({
  folderId,
  deckId,
}: {
  folderId: string;
  deckId: string;
}) {
  const [deckName, setDeckName] = useState<string | null>(null);
  const deckHref = `/home/${folderId}/${deckId}`;

  const { status, items, needsAuth, reload } = useCollection<Card>(
    useCallback(
      () =>
        flashcardClient.listFlashcards({ deckId }).then((res) =>
          res.flashcards.map((c) => ({ id: c.id, front: c.front, back: c.back })),
        ),
      [deckId],
    ),
  );

  useEffect(() => {
    let active = true;
    deckClient
      .getDeck({ id: deckId })
      .then((res) => active && setDeckName(res.deck?.name ?? null))
      .catch(() => {});
    return () => {
      active = false;
    };
  }, [deckId]);

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
            title="Could not load cards"
            needsAuth={needsAuth}
            onRetry={reload}
          />
        </div>
      </Centered>
    );
  }

  if (items.length === 0) {
    return (
      <Centered>
        <div className="w-full max-w-xl">
          <EmptyState
            icon={CardsThreeIcon}
            title="Nothing to review"
            body="This deck has no cards yet."
            action={<ButtonLink href={deckHref}>Back to deck</ButtonLink>}
          />
        </div>
      </Centered>
    );
  }

  return (
    <Session
      cards={items}
      deckName={deckName}
      deckHref={deckHref}
    />
  );
}

function Session({
  cards,
  deckName,
  deckHref,
}: {
  cards: Card[];
  deckName: string | null;
  deckHref: string;
}) {
  const reduce = useReducedMotion();
  const total = cards.length;
  const [state, dispatch] = useReducer(sessionReducer, initialSession);
  const { index, flipped, completed } = state;

  const reviewed = completed ? total : index;
  const remaining = total - reviewed;
  const fraction = reviewed / total;

  const next = useCallback(() => dispatch({ type: "next", total }), [total]);
  const prev = useCallback(() => dispatch({ type: "prev", total }), [total]);

  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === " " || e.key === "Enter") {
        e.preventDefault();
        dispatch({ type: "flip" });
      } else if (e.key === "ArrowRight") {
        next();
      } else if (e.key === "ArrowLeft") {
        prev();
      }
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [next, prev]);

  return (
    <div className="mx-auto flex min-h-[calc(100vh-8rem)] max-w-xl flex-col">
      <div className="flex items-center justify-between">
        <Link
          href={deckHref}
          className="inline-flex items-center gap-1.5 text-sm text-neutral-500 transition-colors duration-150 hover:text-neutral-200"
        >
          <ArrowLeftIcon size={15} />
          {deckName ?? "Back to deck"}
        </Link>
        <span className="text-sm tabular-nums text-neutral-500">
          {completed ? total : index + 1} / {total}
        </span>
      </div>

      <div className="flex flex-1 flex-col items-center justify-center py-8">
        {completed ? (
          <Completion total={total} onRestart={() => dispatch({ type: "restart" })} deckHref={deckHref} />
        ) : (
          <motion.div
            key={index}
            initial={reduce ? false : { opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.25, ease: easeEmphasized }}
            className="w-full"
          >
            <FlipCard
              card={cards[index]}
              flipped={flipped}
              onFlip={() => dispatch({ type: "flip" })}
            />
          </motion.div>
        )}
      </div>

      {!completed && (
        <div className="flex items-center justify-between gap-3">
          <Button variant="ghost" onClick={prev} disabled={index === 0}>
            <ArrowLeftIcon size={16} />
            Previous
          </Button>
          <span className="text-xs text-neutral-600">
            {reduce ? "Tap card to flip" : "Space to flip · ← → to move"}
          </span>
          <Button onClick={next}>
            {index === total - 1 ? "Finish" : "Next"}
            <ArrowRightIcon size={16} />
          </Button>
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
          <span className="tabular-nums">{reviewed} reviewed</span>
          <span className="tabular-nums">{remaining} remaining</span>
        </div>
      </div>
    </div>
  );
}

function FlipCard({
  card,
  flipped,
  onFlip,
}: {
  card: Card;
  flipped: boolean;
  onFlip: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onFlip}
      aria-label={flipped ? "Show question" : "Show answer"}
      className="w-full [perspective:1600px] focus:outline-none"
    >
      <div
        className="relative h-72 w-full [transform-style:preserve-3d] transition-transform duration-[450ms] [transition-timing-function:var(--ease-in-out)] motion-reduce:transition-none"
        style={{ transform: flipped ? "rotateY(180deg)" : "rotateY(0deg)" }}
      >
        <Face className="bg-gradient-to-br from-neutral-900 to-neutral-950">
          <Tag>Question</Tag>
          <p className="mt-4 text-xl leading-snug text-neutral-50">{card.front}</p>
        </Face>
        <Face className="bg-gradient-to-br from-neutral-800 to-neutral-900 [transform:rotateY(180deg)]">
          <Tag>Answer</Tag>
          <p className="mt-4 text-xl leading-snug text-neutral-50">{card.back}</p>
        </Face>
      </div>
    </button>
  );
}

function Face({
  className,
  children,
}: {
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <div
      className={`absolute inset-0 flex flex-col items-center justify-center rounded-2xl border border-neutral-800 p-8 text-center [backface-visibility:hidden] ${className ?? ""}`}
    >
      {children}
    </div>
  );
}

function Tag({ children }: { children: React.ReactNode }) {
  return (
    <span className="text-[11px] font-medium uppercase tracking-[0.18em] text-neutral-500">
      {children}
    </span>
  );
}

function Completion({
  total,
  onRestart,
  deckHref,
}: {
  total: number;
  onRestart: () => void;
  deckHref: string;
}) {
  const reduce = useReducedMotion();
  return (
    <motion.div
      initial={reduce ? false : { opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, ease: easeEmphasized }}
      className="flex flex-col items-center text-center"
    >
      <h2 className="text-2xl font-semibold tracking-tight text-neutral-50">
        Deck complete
      </h2>
      <p className="mt-2 text-sm text-neutral-400">
        You reviewed all {total} {total === 1 ? "card" : "cards"}.
      </p>
      <div className="mt-7 flex items-center gap-3">
        <Button onClick={onRestart}>
          <ArrowUDownLeftIcon size={16} />
          Review again
        </Button>
        <ButtonLink href={deckHref} variant="ghost">
          Back to deck
        </ButtonLink>
      </div>
    </motion.div>
  );
}

function Centered({ children }: { children: React.ReactNode }) {
  return (
    <div className="mx-auto flex min-h-[calc(100vh-8rem)] max-w-xl items-center justify-center">
      {children}
    </div>
  );
}
