"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import {
  CardsThreeIcon,
  PencilSimpleIcon,
  PlusIcon,
  UploadSimpleIcon,
} from "@phosphor-icons/react";

import { deckClient, folderClient } from "@/lib/client";
import { useCollection } from "@/lib/use-collection";
import { easeEmphasized } from "@/lib/motion";
import { Button } from "@/components/ui/button";
import { Breadcrumbs } from "@/components/app/breadcrumbs";
import { ConfirmDelete } from "@/components/app/confirm-delete";
import { ImportDialog } from "@/components/app/import-dialog";
import { InlineNameForm } from "@/components/app/inline-name-form";
import {
  CardSkeletons,
  EmptyState,
  ErrorState,
} from "@/components/app/states";

type Deck = { id: string; name: string; description: string; cardCount: number };

export function DecksView({ folderId }: { folderId: string }) {
  const reduce = useReducedMotion();
  const [folderName, setFolderName] = useState<string | null>(null);

  const { status, items, setItems, needsAuth, reload } = useCollection<Deck>(
    useCallback(
      () =>
        deckClient.listDecks({ folderId }).then((res) =>
          res.decks.map((d) => ({
            id: d.id,
            name: d.name,
            description: d.description,
            cardCount: d.cardCount,
          })),
        ),
      [folderId],
    ),
  );

  useEffect(() => {
    let active = true;
    folderClient
      .getFolder({ id: folderId })
      .then((res) => active && setFolderName(res.folder?.name ?? null))
      .catch(() => {});
    return () => {
      active = false;
    };
  }, [folderId]);

  const [creating, setCreating] = useState(false);
  const [importing, setImporting] = useState(false);

  const addDeck = useCallback(
    (deck: Deck) => setItems((prev) => [...prev, deck]),
    [setItems],
  );

  const create = useCallback(
    async (name: string) => {
      const res = await deckClient.createDeck({ folderId, name, description: "" });
      if (res.deck) {
        const deck = {
          id: res.deck.id,
          name: res.deck.name,
          description: res.deck.description,
          cardCount: res.deck.cardCount,
        };
        setItems((prev) => [...prev, deck]);
      }
      setCreating(false);
    },
    [folderId, setItems],
  );

  const rename = useCallback(
    (id: string, name: string) =>
      setItems((prev) => prev.map((d) => (d.id === id ? { ...d, name } : d))),
    [setItems],
  );

  const remove = useCallback(
    (id: string) => setItems((prev) => prev.filter((d) => d.id !== id)),
    [setItems],
  );

  return (
    <div className="mx-auto max-w-5xl">
      <Breadcrumbs
        items={[
          { label: "Folders", href: "/home" },
          { label: folderName ?? "Folder" },
        ]}
      />

      <div className="mt-4 flex min-h-9 items-center justify-between gap-4">
        <h1 className="text-base font-medium text-neutral-100">Decks</h1>
        {creating ? (
          <InlineNameForm
            placeholder="Deck name"
            submitLabel="Add"
            onSubmit={create}
            onCancel={() => setCreating(false)}
          />
        ) : (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              onClick={() => setImporting(true)}
              disabled={status === "error"}
            >
              <UploadSimpleIcon size={16} />
              Import
            </Button>
            <Button onClick={() => setCreating(true)} disabled={status === "error"}>
              <PlusIcon size={16} />
              New deck
            </Button>
          </div>
        )}
      </div>

      <ImportDialog
        open={importing}
        onClose={() => setImporting(false)}
        mode="deck"
        folderId={folderId}
        onImportedDeck={addDeck}
      />

      <div className="mt-6">
        {status === "loading" && <CardSkeletons />}
        {status === "error" && (
          <ErrorState
            title="Could not load decks"
            needsAuth={needsAuth}
            onRetry={reload}
          />
        )}
        {status === "ready" && items.length === 0 && (
          <EmptyState
            icon={CardsThreeIcon}
            title="No decks yet"
            body="Add a deck to start filling it with flashcards."
          />
        )}
        {status === "ready" && items.length > 0 && (
          <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <AnimatePresence initial={!reduce}>
              {items.map((deck, i) => (
                <motion.li
                  key={deck.id}
                  layout={!reduce}
                  initial={reduce ? false : { opacity: 0, y: 12 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={reduce ? undefined : { opacity: 0, scale: 0.97 }}
                  transition={{
                    duration: 0.4,
                    delay: Math.min(i * 0.04, 0.3),
                    ease: easeEmphasized,
                  }}
                >
                  <DeckCard
                    folderId={folderId}
                    deck={deck}
                    onRenamed={rename}
                    onDeleted={remove}
                  />
                </motion.li>
              ))}
            </AnimatePresence>
          </ul>
        )}
      </div>
    </div>
  );
}

function DeckCard({
  folderId,
  deck,
  onRenamed,
  onDeleted,
}: {
  folderId: string;
  deck: Deck;
  onRenamed: (id: string, name: string) => void;
  onDeleted: (id: string) => void;
}) {
  const [editing, setEditing] = useState(false);

  async function rename(name: string) {
    const res = await deckClient.updateDeck({
      id: deck.id,
      name,
      description: deck.description,
    });
    if (res.deck) onRenamed(deck.id, res.deck.name);
    setEditing(false);
  }

  async function remove() {
    await deckClient.deleteDeck({ id: deck.id });
    onDeleted(deck.id);
  }

  if (editing) {
    return (
      <div className="rounded-xl border border-neutral-800 bg-[#0a0b0c] p-3">
        <InlineNameForm
          placeholder="Deck name"
          initial={deck.name}
          submitLabel="Save"
          onSubmit={rename}
          onCancel={() => setEditing(false)}
        />
      </div>
    );
  }

  return (
    <div className="group relative flex flex-col gap-4 rounded-xl border border-neutral-900 bg-[#0a0b0c] p-4 transition-[transform,border-color] duration-150 [transition-timing-function:var(--ease-out)] hover:border-neutral-700 [@media(hover:hover)]:hover:-translate-y-0.5">
      <Link
        href={`/home/${folderId}/${deck.id}`}
        className="absolute inset-0 rounded-xl focus:outline-none focus-visible:ring-2 focus-visible:ring-neutral-600"
        aria-label={`Open ${deck.name}`}
      />
      <div className="flex items-start justify-between">
        <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-neutral-900 text-neutral-300">
          <CardsThreeIcon size={18} weight="duotone" />
        </span>
        <div className="relative z-10 flex items-center opacity-70 transition-opacity duration-150 group-hover:opacity-100">
          <button
            aria-label={`Rename ${deck.name}`}
            onClick={() => setEditing(true)}
            className="rounded-md p-1.5 text-neutral-600 transition-colors duration-150 [transition-timing-function:var(--ease-out)] hover:text-neutral-300 active:scale-[0.95]"
          >
            <PencilSimpleIcon size={16} />
          </button>
          <ConfirmDelete label={`Delete ${deck.name}`} onConfirm={remove} />
        </div>
      </div>
      <div className="min-w-0">
        <h3 className="truncate text-sm font-medium text-neutral-100">
          {deck.name}
        </h3>
        <p className="mt-0.5 text-xs text-neutral-500">
          {deck.cardCount} {deck.cardCount === 1 ? "card" : "cards"}
        </p>
      </div>
    </div>
  );
}
