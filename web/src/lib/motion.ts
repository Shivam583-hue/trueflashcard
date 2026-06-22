import type { Transition, Variants } from "motion/react";

export const easeOut = [0.23, 1, 0.32, 1] as const;
export const easeEmphasized = [0.16, 1, 0.3, 1] as const;

export const fadeUp: Variants = {
  hidden: { opacity: 0, y: 16 },
  show: { opacity: 1, y: 0 },
};

export const revealTransition = (delay = 0): Transition => ({
  duration: 0.5,
  delay,
  ease: easeEmphasized,
});

export const staggerContainer: Variants = {
  hidden: {},
  show: {
    transition: { staggerChildren: 0.06, delayChildren: 0.05 },
  },
};
