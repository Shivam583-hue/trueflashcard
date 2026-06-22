"use client";

import { motion, useReducedMotion } from "motion/react";
import { CardsThree, FolderSimple, GaugeIcon } from "@phosphor-icons/react";

import { fadeUp, staggerContainer, revealTransition } from "@/lib/motion";

const FEATURES = [
  {
    icon: FolderSimple,
    title: "Folders & decks",
    body: "Group decks into folders so a whole subject stays in one place, not scattered across a list.",
  },
  {
    icon: CardsThree,
    title: "Front and back",
    body: "Every card is a simple question and answer. Flip to check yourself, mark it, move on.",
  },
  {
    icon: GaugeIcon,
    title: "Always know your place",
    body: "A live count of the current card, what you have reviewed, and how much is left in the deck.",
  },
];

export function Features() {
  const reduce = useReducedMotion();

  return (
    <section id="features" className="border-t border-neutral-900 py-24">
      <div className="mx-auto max-w-6xl px-6">
        <motion.h2
          initial={reduce ? false : { opacity: 0, y: 16 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, amount: 0.6 }}
          transition={revealTransition()}
          className="max-w-md text-2xl font-semibold tracking-tight text-neutral-100 sm:text-3xl"
        >
          Built around how you actually study.
        </motion.h2>

        <motion.ul
          variants={staggerContainer}
          initial={reduce ? false : "hidden"}
          whileInView="show"
          viewport={{ once: true, amount: 0.2 }}
          className="mt-12 grid gap-px overflow-hidden rounded-2xl border border-neutral-900 bg-neutral-900 sm:grid-cols-3"
        >
          {FEATURES.map(({ icon: Icon, title, body }) => (
            <motion.li
              key={title}
              variants={fadeUp}
              transition={revealTransition()}
              className="bg-[#0a0b0c] p-7"
            >
              <Icon size={22} weight="duotone" className="text-neutral-300" />
              <h3 className="mt-5 text-base font-medium text-neutral-100">
                {title}
              </h3>
              <p className="mt-2 text-sm leading-relaxed text-neutral-400">
                {body}
              </p>
            </motion.li>
          ))}
        </motion.ul>
      </div>
    </section>
  );
}
