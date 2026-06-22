"use client";

import { useState } from "react";
import { motion, useReducedMotion } from "motion/react";

import { Button } from "@/components/ui/button";
import { easeEmphasized } from "@/lib/motion";

export function FlashcardForm({
  initialFront = "",
  initialBack = "",
  submitLabel,
  onSubmit,
  onCancel,
}: {
  initialFront?: string;
  initialBack?: string;
  submitLabel: string;
  onSubmit: (front: string, back: string) => Promise<void>;
  onCancel: () => void;
}) {
  const reduce = useReducedMotion();
  const [front, setFront] = useState(initialFront);
  const [back, setBack] = useState(initialBack);
  const [pending, setPending] = useState(false);

  const valid = front.trim() !== "" && back.trim() !== "";

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!valid || pending) return;
    setPending(true);
    try {
      await onSubmit(front.trim(), back.trim());
    } finally {
      setPending(false);
    }
  }

  return (
    <motion.form
      onSubmit={submit}
      initial={reduce ? false : { opacity: 0, y: -8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2, ease: easeEmphasized }}
      className="rounded-xl border border-neutral-800 bg-[#0a0b0c] p-4"
    >
      <div className="grid gap-4 sm:grid-cols-2">
        <Field label="Front" value={front} onChange={setFront} autoFocus />
        <Field label="Back" value={back} onChange={setBack} />
      </div>
      <div className="mt-4 flex items-center gap-2">
        <Button type="submit" disabled={!valid || pending}>
          {submitLabel}
        </Button>
        <Button type="button" variant="ghost" onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </motion.form>
  );
}

function Field({
  label,
  value,
  onChange,
  autoFocus,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  autoFocus?: boolean;
}) {
  return (
    <label className="flex flex-col gap-2">
      <span className="text-xs font-medium uppercase tracking-wide text-neutral-500">
        {label}
      </span>
      <textarea
        autoFocus={autoFocus}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        rows={3}
        maxLength={10000}
        className="resize-none rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-sm text-neutral-100 placeholder:text-neutral-600 focus:border-neutral-600 focus:outline-none"
        placeholder={`${label} of the card`}
      />
    </label>
  );
}
