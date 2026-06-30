"use client";

import type { ReactNode } from "react";

export function FlipCard({
  front,
  back,
  flipped,
  onFlip,
}: {
  front: string;
  back: string;
  flipped: boolean;
  onFlip: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onFlip}
      aria-label={flipped ? "Show question" : "Show answer"}
      className="w-full [perspective:1600px] focus:outline-none"
    >
      <div
        className="relative h-72 w-full [transform-style:preserve-3d] transition-transform duration-[450ms] [transition-timing-function:var(--ease-in-out)] motion-reduce:transition-none"
        style={{ transform: flipped ? "rotateY(180deg)" : "rotateY(0deg)" }}
      >
        <Face className="bg-gradient-to-br from-neutral-900 to-neutral-950">
          <Tag>Question</Tag>
          <p className="mt-4 text-xl leading-snug text-neutral-50">{front}</p>
        </Face>
        <Face className="bg-gradient-to-br from-neutral-800 to-neutral-900 [transform:rotateY(180deg)]">
          <Tag>Answer</Tag>
          <p className="mt-4 text-xl leading-snug text-neutral-50">{back}</p>
        </Face>
      </div>
    </button>
  );
}

function Face({
  className,
  children,
}: {
  className?: string;
  children: ReactNode;
}) {
  return (
    <div
      className={`absolute inset-0 flex flex-col items-center justify-center rounded-2xl border border-neutral-800 p-8 text-center [backface-visibility:hidden] ${className ?? ""}`}
    >
      {children}
    </div>
  );
}

function Tag({ children }: { children: ReactNode }) {
  return (
    <span className="text-[11px] font-medium uppercase tracking-[0.18em] text-neutral-500">
      {children}
    </span>
  );
}

export function Centered({ children }: { children: ReactNode }) {
  return (
    <div className="mx-auto flex min-h-[calc(100vh-8rem)] max-w-xl items-center justify-center">
      {children}
    </div>
  );
}
