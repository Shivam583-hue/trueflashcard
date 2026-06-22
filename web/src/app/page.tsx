export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-neutral-950 via-neutral-900 to-black px-6">
      <div className="flex flex-col items-center gap-4 text-center">
        <span className="text-xs font-medium uppercase tracking-[0.2em] text-neutral-500">
          Flashcards
        </span>
        <h1 className="text-4xl font-semibold tracking-tight text-neutral-100 sm:text-5xl">
          Hello, world.
        </h1>
        <p className="max-w-sm text-sm leading-relaxed text-neutral-400">
          The skeleton is up. Folders, decks, and study sessions are on the way.
        </p>
      </div>
    </main>
  );
}
