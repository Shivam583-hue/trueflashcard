"use client";

import { motion, useReducedMotion } from "motion/react";
import { FlameIcon, GraduationCapIcon, PlayIcon } from "@phosphor-icons/react";

import { useStudyOverview } from "@/lib/use-study-overview";
import { easeEmphasized } from "@/lib/motion";
import { ButtonLink } from "@/components/ui/button";

export function StudyBanner() {
  const reduce = useReducedMotion();
  const { overview } = useStudyOverview();

  if (!overview) {
    return <div className="skeleton h-[92px] rounded-2xl border border-neutral-900 bg-[#0a0b0c]" />;
  }

  const pending = overview.dueTotal + overview.newTotal;
  const caughtUp = pending === 0;

  return (
    <motion.div
      initial={reduce ? false : { opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, ease: easeEmphasized }}
      className="flex flex-col gap-4 rounded-2xl border border-neutral-900 bg-gradient-to-br from-neutral-900/60 to-[#0a0b0c] p-5 sm:flex-row sm:items-center sm:justify-between"
    >
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-neutral-900 text-neutral-300">
          <GraduationCapIcon size={22} weight="duotone" />
        </span>
        <div>
          <h2 className="text-sm font-medium text-neutral-100">
            {caughtUp ? "You're all caught up" : `${pending} ${pending === 1 ? "card" : "cards"} to study`}
          </h2>
          <p className="mt-0.5 text-xs text-neutral-500">
            {caughtUp
              ? "Nothing is due right now. Great work."
              : `${overview.dueTotal} due · ${overview.newTotal} new`}
          </p>
        </div>
      </div>

      <div className="flex items-center gap-2.5">
        {overview.streakDays > 0 && (
          <span className="inline-flex items-center gap-1.5 rounded-full border border-neutral-800 bg-neutral-900/60 px-3 py-1.5 text-xs font-medium text-neutral-300">
            <FlameIcon size={14} weight="fill" className="text-amber-300/80" />
            {overview.streakDays} day{overview.streakDays === 1 ? "" : "s"}
          </span>
        )}
        {overview.reviewedToday > 0 && (
          <span className="hidden rounded-full border border-neutral-800 bg-neutral-900/60 px-3 py-1.5 text-xs font-medium text-neutral-400 sm:inline-flex">
            {overview.reviewedToday} reviewed today
          </span>
        )}
        {!caughtUp && (
          <ButtonLink href="/study">
            <PlayIcon size={16} weight="fill" />
            Study
          </ButtonLink>
        )}
      </div>
    </motion.div>
  );
}
