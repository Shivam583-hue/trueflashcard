"use client";

import { useState } from "react";
import { motion, useReducedMotion } from "motion/react";

import { Button } from "@/components/ui/button";
import { easeEmphasized } from "@/lib/motion";

export function InlineNameForm({
  placeholder,
  initial = "",
  submitLabel,
  maxLength = 200,
  onSubmit,
  onCancel,
}: {
  placeholder: string;
  initial?: string;
  submitLabel: string;
  maxLength?: number;
  onSubmit: (value: string) => Promise<void>;
  onCancel: () => void;
}) {
  const reduce = useReducedMotion();
  const [value, setValue] = useState(initial);
  const [pending, setPending] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    const trimmed = value.trim();
    if (!trimmed || pending) return;
    setPending(true);
    try {
      await onSubmit(trimmed);
    } finally {
      setPending(false);
    }
  }

  return (
    <motion.form
      onSubmit={submit}
      initial={reduce ? false : { opacity: 0, scale: 0.98 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.2, ease: easeEmphasized }}
      className="flex items-center gap-2"
      style={{ transformOrigin: "right center" }}
    >
      <input
        autoFocus
        value={value}
        onChange={(e) => setValue(e.target.value)}
        onKeyDown={(e) => e.key === "Escape" && onCancel()}
        placeholder={placeholder}
        maxLength={maxLength}
        className="h-9 w-52 rounded-lg border border-neutral-800 bg-neutral-900 px-3 text-sm text-neutral-100 placeholder:text-neutral-600 focus:border-neutral-600 focus:outline-none"
      />
      <Button type="submit" disabled={pending || !value.trim()}>
        {submitLabel}
      </Button>
      <Button type="button" variant="ghost" onClick={onCancel}>
        Cancel
      </Button>
    </motion.form>
  );
}
