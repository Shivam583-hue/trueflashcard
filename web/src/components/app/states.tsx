import type { Icon } from "@phosphor-icons/react";

import { Button } from "@/components/ui/button";

export function EmptyState({
  icon: IconComponent,
  title,
  body,
  action,
}: {
  icon: Icon;
  title: string;
  body: string;
  action?: React.ReactNode;
}) {
  return (
    <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-neutral-900 py-20 text-center">
      <span className="flex h-12 w-12 items-center justify-center rounded-xl bg-neutral-900 text-neutral-400">
        <IconComponent size={24} weight="duotone" />
      </span>
      <h3 className="mt-4 text-sm font-medium text-neutral-200">{title}</h3>
      <p className="mt-1 max-w-xs text-sm text-neutral-500">{body}</p>
      {action && <div className="mt-5">{action}</div>}
    </div>
  );
}

export function ErrorState({
  title,
  needsAuth,
  onRetry,
}: {
  title: string;
  needsAuth: boolean;
  onRetry: () => void;
}) {
  return (
    <div className="flex flex-col items-center justify-center rounded-2xl border border-neutral-900 py-20 text-center">
      <h3 className="text-sm font-medium text-neutral-200">
        {needsAuth ? "Please sign in" : title}
      </h3>
      <p className="mt-1 max-w-xs text-sm text-neutral-500">
        {needsAuth
          ? "Your session is missing or expired."
          : "The server is unreachable right now."}
      </p>
      <div className="mt-5">
        {needsAuth ? (
          <Button onClick={() => (window.location.href = "/login")}>
            Go to sign in
          </Button>
        ) : (
          <Button variant="ghost" onClick={onRetry}>
            Try again
          </Button>
        )}
      </div>
    </div>
  );
}

export function CardSkeletons({
  count = 6,
  className = "grid gap-3 sm:grid-cols-2 lg:grid-cols-3",
}: {
  count?: number;
  className?: string;
}) {
  return (
    <ul className={className}>
      {Array.from({ length: count }).map((_, i) => (
        <li
          key={i}
          className="skeleton flex h-[68px] items-center gap-3 rounded-xl border border-neutral-900 bg-[#0a0b0c] p-4"
        >
          <span className="h-9 w-9 rounded-lg bg-neutral-900" />
          <span className="h-3 w-24 rounded bg-neutral-900" />
        </li>
      ))}
    </ul>
  );
}
