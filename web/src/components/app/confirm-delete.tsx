"use client";

import { useState } from "react";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { TrashIcon } from "@phosphor-icons/react";

import { easeEmphasized } from "@/lib/motion";

export function ConfirmDelete({
  label,
  onConfirm,
}: {
  label: string;
  onConfirm: () => Promise<void>;
}) {
  const reduce = useReducedMotion();
  const [confirming, setConfirming] = useState(false);
  const [pending, setPending] = useState(false);

  async function confirm() {
    if (pending) return;
    setPending(true);
    try {
      await onConfirm();
    } finally {
      setPending(false);
      setConfirming(false);
    }
  }

  return (
    <AnimatePresence mode="wait" initial={false}>
      {confirming ? (
        <motion.div
          key="confirm"
          initial={reduce ? false : { opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.12, ease: easeEmphasized }}
          className="flex items-center gap-1"
        >
          <button
            onClick={confirm}
            disabled={pending}
            className="rounded-md px-2 py-1 text-xs font-medium text-red-300 transition-colors duration-150 [transition-timing-function:var(--ease-out)] hover:text-red-200 active:scale-[0.97] disabled:opacity-50"
          >
            Delete
          </button>
          <button
            onClick={() => setConfirming(false)}
            className="rounded-md px-2 py-1 text-xs text-neutral-500 transition-colors duration-150 hover:text-neutral-300"
          >
            Cancel
          </button>
        </motion.div>
      ) : (
        <motion.button
          key="trash"
          initial={reduce ? false : { opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.12, ease: easeEmphasized }}
          aria-label={label}
          onClick={() => setConfirming(true)}
          className="rounded-md p-1.5 text-neutral-600 transition-colors duration-150 [transition-timing-function:var(--ease-out)] hover:text-neutral-300 active:scale-[0.95]"
        >
          <TrashIcon size={16} />
        </motion.button>
      )}
    </AnimatePresence>
  );
}
