"use client";

import { useEffect, useMemo, useState } from "react";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { ConnectError } from "@connectrpc/connect";
import {
  CheckCircleIcon,
  CopyIcon,
  WarningCircleIcon,
  XIcon,
} from "@phosphor-icons/react";

import { deckClient, flashcardClient } from "@/lib/client";
import {
  CARDS_PROMPT,
  DECK_PROMPT,
  type DeckParseResult,
  parseCardsImport,
  parseDeckImport,
} from "@/lib/import-parse";
import { Button } from "@/components/ui/button";

type DeckSummary = {
  id: string;
  name: string;
  description: string;
  cardCount: number;
};
type CardSummary = { id: string; front: string; back: string; position: number };

type ModeProps =
  | { mode: "deck"; folderId: string; onImportedDeck: (deck: DeckSummary) => void }
  | { mode: "cards"; deckId: string; onImportedCards: (cards: CardSummary[]) => void };

type Props = { open: boolean; onClose: () => void } & ModeProps;

const PLACEHOLDER_DECK = `{
  "name": "Photosynthesis",
  "cards": [
    { "front": "What is photosynthesis?", "back": "Conversion of light into chemical energy." }
  ]
}`;

const PLACEHOLDER_CARDS = `[
  { "front": "What is photosynthesis?", "back": "Conversion of light into chemical energy." }
]`;

export function ImportDialog({ open, onClose, ...mode }: Props) {
  const reduce = useReducedMotion();
  return (
    <AnimatePresence>
      {open && (
        <motion.div
          className="fixed inset-0 z-50 flex items-center justify-center p-4"
          initial={reduce ? false : { opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.15 }}
        >
          <button
            aria-label="Close"
            onClick={onClose}
            className="absolute inset-0 bg-black/60 backdrop-blur-sm"
          />
          <DialogBody onClose={onClose} {...(mode as ModeProps)} />
        </motion.div>
      )}
    </AnimatePresence>
  );
}

function DialogBody({ onClose, ...props }: { onClose: () => void } & ModeProps) {
  const reduce = useReducedMotion();
  const mode = props.mode;
  const [text, setText] = useState("");
  const [pending, setPending] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const result = useMemo(
    () => (mode === "deck" ? parseDeckImport(text) : parseCardsImport(text)),
    [text, mode],
  );
  const deckResult = mode === "deck" ? (result as DeckParseResult) : null;

  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") onClose();
    }
    window.addEventListener("keydown", onKey);
    document.body.style.overflow = "hidden";
    return () => {
      window.removeEventListener("keydown", onKey);
      document.body.style.overflow = "";
    };
  }, [onClose]);

  async function copyPrompt() {
    await navigator.clipboard.writeText(mode === "deck" ? DECK_PROMPT : CARDS_PROMPT);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  async function runImport() {
    if (!result.ok || pending) return;
    setPending(true);
    setServerError(null);
    try {
      if (props.mode === "deck") {
        const r = parseDeckImport(text);
        if (!r.ok) return;
        const res = await deckClient.importDeck({
          folderId: props.folderId,
          name: r.name,
          description: r.description,
          cards: r.cards,
        });
        if (res.deck) {
          props.onImportedDeck({
            id: res.deck.id,
            name: res.deck.name,
            description: res.deck.description,
            cardCount: res.deck.cardCount,
          });
        }
      } else {
        const r = parseCardsImport(text);
        if (!r.ok) return;
        const res = await flashcardClient.importFlashcards({
          deckId: props.deckId,
          cards: r.cards,
        });
        props.onImportedCards(
          res.flashcards.map((c) => ({
            id: c.id,
            front: c.front,
            back: c.back,
            position: c.position,
          })),
        );
      }
      onClose();
    } catch (e) {
      setServerError(ConnectError.from(e).message);
    } finally {
      setPending(false);
    }
  }

  const count = result.ok ? result.cards.length : 0;

  return (
    <motion.div
      role="dialog"
      aria-modal="true"
      initial={reduce ? false : { opacity: 0, scale: 0.96 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.2, ease: [0.16, 1, 0.3, 1] }}
      className="relative z-10 flex max-h-[85vh] w-full max-w-lg flex-col rounded-2xl border border-neutral-800 bg-[#0c0d0e] shadow-2xl"
    >
      <div className="flex items-start justify-between border-b border-neutral-900 p-5">
        <div>
          <h2 className="text-base font-medium text-neutral-100">
            {mode === "deck" ? "Import a deck" : "Import cards"}
          </h2>
          <p className="mt-1 text-sm text-neutral-500">
            Paste JSON{mode === "deck" ? " with a name and cards" : ""}.
          </p>
        </div>
        <button
          aria-label="Close"
          onClick={onClose}
          className="rounded-md p-1.5 text-neutral-500 transition-colors duration-150 hover:text-neutral-200 active:scale-[0.95]"
        >
          <XIcon size={18} />
        </button>
      </div>

      <div className="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto p-5">
        <div className="flex justify-end">
          <button
            onClick={copyPrompt}
            className="inline-flex items-center gap-1.5 text-xs text-neutral-500 transition-colors duration-150 hover:text-neutral-300 active:scale-[0.97]"
          >
            {copied ? <CheckCircleIcon size={14} weight="fill" /> : <CopyIcon size={14} />}
            {copied ? "Prompt copied" : "Copy LLM prompt"}
          </button>
        </div>

        <textarea
          autoFocus
          value={text}
          onChange={(e) => setText(e.target.value)}
          spellCheck={false}
          placeholder={mode === "deck" ? PLACEHOLDER_DECK : PLACEHOLDER_CARDS}
          className="h-56 w-full resize-none rounded-lg border border-neutral-800 bg-neutral-950 p-3 font-mono text-xs leading-relaxed text-neutral-100 placeholder:text-neutral-700 focus:border-neutral-600 focus:outline-none"
        />

        {text.trim() !== "" && !result.ok && (
          <ul className="space-y-1 rounded-lg border border-red-950 bg-red-950/20 p-3">
            {result.errors.slice(0, 6).map((err, i) => (
              <li key={i} className="flex items-start gap-2 text-xs text-red-300">
                <WarningCircleIcon size={14} className="mt-px shrink-0" />
                {err}
              </li>
            ))}
            {result.errors.length > 6 && (
              <li className="pl-6 text-xs text-red-400/70">
                and {result.errors.length - 6} more…
              </li>
            )}
          </ul>
        )}

        {result.ok && (
          <p className="flex items-center gap-2 rounded-lg border border-emerald-950 bg-emerald-950/20 p-3 text-xs text-emerald-300">
            <CheckCircleIcon size={14} weight="fill" className="shrink-0" />
            {deckResult?.ok
              ? `Ready: “${deckResult.name}” with ${count} ${count === 1 ? "card" : "cards"}.`
              : `Ready: ${count} ${count === 1 ? "card" : "cards"} to add.`}
          </p>
        )}

        {serverError && (
          <p className="rounded-lg border border-red-950 bg-red-950/20 p-3 text-xs text-red-300">
            {serverError}
          </p>
        )}
      </div>

      <div className="flex items-center justify-end gap-2 border-t border-neutral-900 p-5">
        <Button variant="ghost" onClick={onClose}>
          Cancel
        </Button>
        <Button onClick={runImport} disabled={!result.ok || pending}>
          {pending ? "Importing…" : "Import"}
        </Button>
      </div>
    </motion.div>
  );
}
