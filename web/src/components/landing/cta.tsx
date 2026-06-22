"use client";

import { motion, useReducedMotion } from "motion/react";

import { ButtonLink } from "@/components/ui/button";
import { revealTransition } from "@/lib/motion";

export function CTA() {
  const reduce = useReducedMotion();

  return (
    <section className="border-t border-neutral-900 py-24">
      <motion.div
        initial={reduce ? false : { opacity: 0, y: 16 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true, amount: 0.5 }}
        transition={revealTransition()}
        className="mx-auto flex max-w-6xl flex-col items-center px-6 text-center"
      >
        <h2 className="max-w-xl text-3xl font-semibold tracking-tight text-neutral-50 sm:text-4xl">
          Start your first deck today.
        </h2>
        <p className="mt-4 max-w-sm text-sm leading-relaxed text-neutral-400">
          Sign in with Google and build a deck in under a minute.
        </p>
        <ButtonLink href="/login" className="mt-8">
          Get started
        </ButtonLink>
      </motion.div>
    </section>
  );
}
