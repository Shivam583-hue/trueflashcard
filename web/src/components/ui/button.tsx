import Link from "next/link";
import type { ComponentProps, ReactNode } from "react";

import { cn } from "@/lib/cn";

type Variant = "primary" | "ghost";

const base =
  "inline-flex items-center justify-center gap-2 rounded-lg px-5 py-2.5 text-sm font-medium " +
  "transition-[transform,background-color,border-color,color] duration-150 [transition-timing-function:var(--ease-out)] " +
  "active:scale-[0.97] focus:outline-none focus-visible:ring-2 focus-visible:ring-neutral-500 " +
  "disabled:pointer-events-none disabled:opacity-50";

const variants: Record<Variant, string> = {
  primary:
    "bg-neutral-100 text-neutral-950 hover:bg-white",
  ghost:
    "border border-neutral-800 bg-neutral-900/50 text-neutral-100 hover:border-neutral-700 hover:bg-neutral-800/60",
};

function classesFor(variant: Variant, className?: string) {
  return cn(base, variants[variant], className);
}

export function Button({
  variant = "primary",
  className,
  ...props
}: ComponentProps<"button"> & { variant?: Variant }) {
  return <button className={classesFor(variant, className)} {...props} />;
}

export function ButtonLink({
  href,
  variant = "primary",
  className,
  children,
}: {
  href: string;
  variant?: Variant;
  className?: string;
  children: ReactNode;
}) {
  return (
    <Link href={href} className={classesFor(variant, className)}>
      {children}
    </Link>
  );
}

export function AnchorButton({
  variant = "primary",
  className,
  ...props
}: ComponentProps<"a"> & { variant?: Variant }) {
  return <a className={classesFor(variant, className)} {...props} />;
}
