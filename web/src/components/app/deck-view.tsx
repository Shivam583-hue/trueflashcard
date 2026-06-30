"use client";

import { useCallback, useEffect, useState } from "react";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import {
  CardsThreeIcon,
  LightningIcon,
  PencilSimpleIcon,
  PlayIcon,
  PlusIcon,
  UploadSimpleIcon,
} from "@phosphor-icons/react";

import { deckClient, flashcardClient, folderClient } from "@/lib/client";
import { useCollection } from "@/lib/use-collection";
import { useStudyOverview } from "@/lib/use-study-overview";
import { easeEmphasized } from "@/lib/motion";
import { Button, ButtonLink } from "@/components/ui/button";
import { Breadcrumbs } from "@/components/app/breadcrumbs";
import { ConfirmDelete } from "@/components/app/confirm-delete";
import { FlashcardForm } from "@/components/app/flashcard-form";
import { ImportDialog } from "@/components/app/import-dialog";
import { EmptyState, ErrorState } from "@/components/app/states";

type Card = { id: string; front: string; back: string; position: number };

export function DeckView({
  folderId,
  deckId,
}: {
  folderId: string;
  deckId: string;
}) {
  const reduce = useReducedMotion();
  const [folderName, setFolderName] = useState<string | null>(null);
  const [deckName, setDeckName] = useState<string | null>(null);
  const { overview } = useStudyOverview();
  const counts = overview?.byDeck.get(deckId);

  const { status, items, setItems, needsAuth, reload } = useCollection<Card>(
    useCallback(
      () =>
        flashcardClient.listFlashcards({ deckId }).then((res) =>
          res.flashcards.map((c) => ({
            id: c.id,
            front: c.front,
            back: c.back,
            position: c.position,
          })),
        ),
      [deckId],
    ),
  );

  useEffect(() => {
    let active = true;
    folderClient
      .getFolder({ id: folderId })
      .then((res) => active && setFolderName(res.folder?.name ?? null))
      .catch(() => {});
    deckClient
      .getDeck({ id: deckId })
      .then((res) => active && setDeckName(res.deck?.name ?? null))
      .catch(() => {});
    return () => {
      active = false;
    };
  }, [folderId, deckId]);

  const [creating, setCreating] = useState(false);
  const [importing, setImporting] = useState(false);

  const addCards = useCallback(
    (cards: Card[]) => setItems((prev) => [...prev, ...cards]),
    [setItems],
  );

  const create = useCallback(
    async (front: string, back: string) => {
      const res = await flashcardClient.createFlashcard({ deckId, front, back });
      if (res.flashcard) {
        const card = {
          id: res.flashcard.id,
          front: res.flashcard.front,
          back: res.flashcard.back,
          position: res.flashcard.position,
        };
        setItems((prev) => [...prev, card]);
      }
      setCreating(false);
    },
    [deckId, setItems],
  );

  const update = useCallback(
    (card: Card) =>
      setItems((prev) => prev.map((c) => (c.id === card.id ? card : c))),
    [setItems],
  );

  const remove = useCallback(
    (id: string) => setItems((prev) => prev.filter((c) => c.id !== id)),
    [setItems],
  );

  const count = items.length;

  return (
    <div className="mx-auto max-w-3xl">
      <Breadcrumbs
        items={[
          { label: "Folders", href: "/home" },
          { label: folderName ?? "Folder", href: `/home/${folderId}` },
          { label: deckName ?? "Deck" },
        ]}
      />

      <div className="mt-4 flex min-h-9 flex-wrap items-center justify-between gap-3">
        <h1 className="flex items-center gap-2.5 text-base font-medium text-neutral-100">
          Cards
          {status === "ready" && (
            <span className="text-sm font-normal text-neutral-500">{count}</span>
          )}
          {counts && counts.due + counts.new > 0 && (
            <span className="rounded-full border border-neutral-800 bg-neutral-900/60 px-2.5 py-0.5 text-[11px] font-medium tabular-nums text-neutral-400">
              {counts.due} due · {counts.new} new
            </span>
          )}
        </h1>
        <div className="flex flex-wrap items-center gap-2">
          <Button
            variant="ghost"
            onClick={() => setImporting(true)}
            disabled={status === "error"}
          >
            <UploadSimpleIcon size={16} />
            Import
          </Button>
          <Button
            variant="ghost"
            onClick={() => setCreating(true)}
            disabled={status === "error"}
          >
            <PlusIcon size={16} />
            Add card
          </Button>
          {count > 0 && (
            <ButtonLink
              href={`/home/${folderId}/${deckId}/review`}
              variant="ghost"
            >
              <LightningIcon size={16} />
              Cram
            </ButtonLink>
          )}
          {count > 0 && (
            <ButtonLink href={`/home/${folderId}/${deckId}/study`}>
              <PlayIcon size={16} weight="fill" />
              Study
            </ButtonLink>
          )}
        </div>
      </div>

      <ImportDialog
        open={importing}
        onClose={() => setImporting(false)}
        mode="cards"
        deckId={deckId}
        onImportedCards={addCards}
      />

      <div className="mt-6 space-y-3">
        {creating && (
          <FlashcardForm
            submitLabel="Add card"
            onSubmit={create}
            onCancel={() => setCreating(false)}
          />
        )}

        {status === "loading" && <CardListSkeleton />}
        {status === "error" && (
          <ErrorState
            title="Could not load cards"
            needsAuth={needsAuth}
            onRetry={reload}
          />
        )}
        {status === "ready" && count === 0 && !creating && (
          <EmptyState
            icon={CardsThreeIcon}
            title="No cards yet"
            body="Add your first flashcard with a question and an answer."
            action={
              <Button onClick={() => setCreating(true)}>
                <PlusIcon size={16} />
                Add card
              </Button>
            }
          />
        )}
        {status === "ready" && count > 0 && (
          <ul className="space-y-3">
            <AnimatePresence initial={!reduce}>
              {items.map((card, i) => (
                <motion.li
                  key={card.id}
                  layout={!reduce}
                  initial={reduce ? false : { opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={reduce ? undefined : { opacity: 0, scale: 0.98 }}
                  transition={{
                    duration: 0.35,
                    delay: Math.min(i * 0.03, 0.24),
                    ease: easeEmphasized,
                  }}
                >
                  <CardRow card={card} index={i} onUpdated={update} onDeleted={remove} />
                </motion.li>
              ))}
            </AnimatePresence>
          </ul>
        )}
      </div>
    </div>
  );
}

function CardRow({
  card,
  index,
  onUpdated,
  onDeleted,
}: {
  card: Card;
  index: number;
  onUpdated: (card: Card) => void;
  onDeleted: (id: string) => void;
}) {
  const [editing, setEditing] = useState(false);

  async function save(front: string, back: string) {
    const res = await flashcardClient.updateFlashcard({
      id: card.id,
      front,
      back,
      position: card.position,
    });
    if (res.flashcard) {
      onUpdated({
        id: res.flashcard.id,
        front: res.flashcard.front,
        back: res.flashcard.back,
        position: res.flashcard.position,
      });
    }
    setEditing(false);
  }

  async function remove() {
    await flashcardClient.deleteFlashcard({ id: card.id });
    onDeleted(card.id);
  }

  if (editing) {
    return (
      <FlashcardForm
        initialFront={card.front}
        initialBack={card.back}
        submitLabel="Save"
        onSubmit={save}
        onCancel={() => setEditing(false)}
      />
    );
  }

  return (
    <div className="group flex items-start gap-4 rounded-xl border border-neutral-900 bg-[#0a0b0c] p-4">
      <span className="mt-0.5 w-6 shrink-0 text-xs tabular-nums text-neutral-600">
        {index + 1}
      </span>
      <div className="grid min-w-0 flex-1 gap-1.5 sm:grid-cols-2 sm:gap-4">
        <p className="min-w-0 break-words text-sm text-neutral-100">{card.front}</p>
        <p className="min-w-0 break-words text-sm text-neutral-400 sm:border-l sm:border-neutral-900 sm:pl-4">
          {card.back}
        </p>
      </div>
      <div className="flex shrink-0 items-center opacity-70 transition-opacity duration-150 group-hover:opacity-100">
        <button
          aria-label="Edit card"
          onClick={() => setEditing(true)}
          className="rounded-md p-1.5 text-neutral-600 transition-colors duration-150 [transition-timing-function:var(--ease-out)] hover:text-neutral-300 active:scale-[0.95]"
        >
          <PencilSimpleIcon size={16} />
        </button>
        <ConfirmDelete label="Delete card" onConfirm={remove} />
      </div>
    </div>
  );
}

function CardListSkeleton() {
  return (
    <ul className="space-y-3">
      {Array.from({ length: 4 }).map((_, i) => (
        <li
          key={i}
          className="skeleton h-[76px] rounded-xl border border-neutral-900 bg-[#0a0b0c]"
        />
      ))}
    </ul>
  );
}
