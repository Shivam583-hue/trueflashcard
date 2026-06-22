import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-neutral-950 via-neutral-900 to-black px-6">
      <div className="flex flex-col items-center gap-4 text-center">
        <span className="text-xs font-medium uppercase tracking-[0.2em] text-neutral-500">
          Flashcards
        </span>
        <h1 className="text-4xl font-semibold tracking-tight text-neutral-100 sm:text-5xl">
          Learn anything, one card at a time.
        </h1>
        <p className="max-w-sm text-sm leading-relaxed text-neutral-400">
          Organize folders and decks, then study with focus and track your
          progress.
        </p>
        <Link
          href="/login"
          className="mt-4 rounded-lg border border-neutral-800 bg-neutral-900/60 px-5 py-2.5 text-sm font-medium text-neutral-100 transition-colors hover:border-neutral-700 hover:bg-neutral-800/60"
        >
          Sign in
        </Link>
      </div>
    </main>
  );
}
