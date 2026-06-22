"use client";

import { useEffect, useRef, useState } from "react";
import { useReducedMotion } from "motion/react";

const DECK = [
  { front: "What does TCP guarantee?", back: "Ordered, reliable, error-checked delivery." },
  { front: "Big-O of binary search?", back: "O(log n) — halves the search space each step." },
  { front: "What is idempotency?", back: "Same request, same result, no extra effect." },
  { front: "ACID — the 'I'?", back: "Isolation: concurrent txns don't interfere." },
];

export function FlashcardDemo() {
  const reduce = useReducedMotion();
  const [index, setIndex] = useState(0);
  const [flipped, setFlipped] = useState(false);
  const flippedRef = useRef(false);

  useEffect(() => {
    if (reduce) return;
    const tick = () => {
      if (!flippedRef.current) {
        flippedRef.current = true;
        setFlipped(true);
      } else {
        flippedRef.current = false;
        setFlipped(false);
        setIndex((i) => (i + 1) % DECK.length);
      }
    };
    const id = setInterval(tick, 2200);
    return () => clearInterval(id);
  }, [reduce]);

  const card = DECK[index];
  const reviewed = index;
  const remaining = DECK.length - index - 1;
  const progress = (index / DECK.length) * 100;

  return (
    <div className="w-full max-w-sm">
      <div className="[perspective:1600px]">
        <div
          className="relative h-60 w-full [transform-style:preserve-3d] transition-transform duration-500 [transition-timing-function:var(--ease-in-out)] motion-reduce:transition-none"
          style={{ transform: flipped ? "rotateY(180deg)" : "rotateY(0deg)" }}
        >
          <CardFace className="bg-gradient-to-br from-neutral-900 to-neutral-950">
            <span className="text-[11px] font-medium uppercase tracking-[0.18em] text-neutral-500">
              Question
            </span>
            <p className="mt-3 text-lg leading-snug text-neutral-100">{card.front}</p>
          </CardFace>
          <CardFace
            className="bg-gradient-to-br from-neutral-800 to-neutral-900 [transform:rotateY(180deg)]"
          >
            <span className="text-[11px] font-medium uppercase tracking-[0.18em] text-neutral-400">
              Answer
            </span>
            <p className="mt-3 text-lg leading-snug text-neutral-50">{card.back}</p>
          </CardFace>
        </div>
      </div>

      <div className="mt-5">
        <div className="h-1 w-full overflow-hidden rounded-full bg-neutral-800">
          <div
            className="h-full rounded-full bg-neutral-300 transition-[width] duration-500 [transition-timing-function:var(--ease-out)]"
            style={{ width: `${progress}%` }}
          />
        </div>
        <div className="mt-2.5 flex items-center justify-between text-xs text-neutral-500">
          <span className="tabular-nums">
            Card {index + 1} of {DECK.length}
          </span>
          <span className="tabular-nums">
            {reviewed} reviewed · {remaining} left
          </span>
        </div>
      </div>
    </div>
  );
}

function CardFace({
  className,
  children,
}: {
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <div
      className={`absolute inset-0 flex flex-col justify-center rounded-2xl border border-neutral-800 p-6 [backface-visibility:hidden] ${className ?? ""}`}
    >
      {children}
    </div>
  );
}
