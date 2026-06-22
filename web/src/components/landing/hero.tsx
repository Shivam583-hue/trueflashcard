"use client";

import { motion, useReducedMotion } from "motion/react";

import { FlashcardDemo } from "@/components/landing/flashcard-demo";
import { ButtonLink } from "@/components/ui/button";
import { easeEmphasized } from "@/lib/motion";

export function Hero() {
  const reduce = useReducedMotion();

  const enter = (delay: number) =>
    reduce
      ? { initial: false as const }
      : {
          initial: { opacity: 0, y: 16 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.55, delay, ease: easeEmphasized },
        };

  return (
    <section className="relative overflow-hidden">
      <div
        aria-hidden="true"
        className="aurora pointer-events-none absolute -top-40 left-1/2 h-[480px] w-[820px] -translate-x-1/2 rounded-full blur-2xl [background:radial-gradient(50%_50%_at_50%_50%,rgba(120,130,150,0.18),transparent_70%)]"
      />
      <div className="relative mx-auto grid max-w-6xl items-center gap-12 px-6 pt-20 pb-24 md:grid-cols-2 md:pt-28">
        <div className="flex flex-col items-start">
          <motion.span
            {...enter(0)}
            className="text-[11px] font-medium uppercase tracking-[0.2em] text-neutral-500"
          >
            Spaced, focused study
          </motion.span>
          <motion.h1
            {...enter(0.06)}
            className="mt-4 text-4xl font-semibold leading-[1.05] tracking-tight text-neutral-50 sm:text-5xl"
          >
            Learn anything,
            <br />
            one card at a time.
          </motion.h1>
          <motion.p
            {...enter(0.12)}
            className="mt-5 max-w-md text-base leading-relaxed text-neutral-400"
          >
            Organize folders and decks, then study with a clear sense of where
            you are — current card, reviewed, and what is left.
          </motion.p>
          <motion.div {...enter(0.18)} className="mt-8 flex items-center gap-3">
            <ButtonLink href="/login">Get started</ButtonLink>
            <ButtonLink href="#features" variant="ghost">
              See how it works
            </ButtonLink>
          </motion.div>
        </div>

        <motion.div
          {...(reduce
            ? { initial: false as const }
            : {
                initial: { opacity: 0, y: 24, scale: 0.96 },
                animate: { opacity: 1, y: 0, scale: 1 },
                transition: { duration: 0.6, delay: 0.2, ease: easeEmphasized },
              })}
          className="flex justify-center md:justify-end"
        >
          <FlashcardDemo />
        </motion.div>
      </div>
    </section>
  );
}
