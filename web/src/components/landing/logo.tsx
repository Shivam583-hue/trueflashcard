import { cn } from "@/lib/cn";

export function Logo({ className }: { className?: string }) {
  return (
    <span className={cn("inline-flex items-center gap-2", className)}>
      <span className="relative inline-block h-5 w-5" aria-hidden="true">
        <span className="absolute inset-x-0 top-0 h-3.5 w-full -rotate-6 rounded-[5px] border border-neutral-700 bg-neutral-800" />
        <span className="absolute inset-x-0 bottom-0 h-3.5 w-full rotate-3 rounded-[5px] border border-neutral-600 bg-neutral-200" />
      </span>
      <span className="text-sm font-semibold tracking-tight text-neutral-100">
        Flashcards
      </span>
    </span>
  );
}
