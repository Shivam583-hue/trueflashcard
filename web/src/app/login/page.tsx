import Link from "next/link";

import { GoogleIcon } from "@/components/google-icon";
import { googleLoginUrl } from "@/lib/config";

export default function LoginPage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-neutral-950 via-neutral-900 to-black px-6">
      <div className="w-full max-w-sm">
        <div className="flex flex-col items-center gap-2 text-center">
          <span className="text-xs font-medium uppercase tracking-[0.2em] text-neutral-500">
            Flashcards
          </span>
          <h1 className="text-2xl font-semibold tracking-tight text-neutral-100">
            Sign in to continue
          </h1>
          <p className="text-sm text-neutral-400">
            Use your Google account to access your folders and decks.
          </p>
        </div>

        <a
          href={googleLoginUrl}
          className="mt-8 flex w-full items-center justify-center gap-3 rounded-lg border border-neutral-800 bg-neutral-900/60 px-4 py-3 text-sm font-medium text-neutral-100 transition-colors hover:border-neutral-700 hover:bg-neutral-800/60 focus:outline-none focus-visible:ring-2 focus-visible:ring-neutral-600"
        >
          <GoogleIcon className="h-5 w-5" />
          Continue with Google
        </a>

        <p className="mt-6 text-center text-xs text-neutral-600">
          <Link href="/" className="transition-colors hover:text-neutral-400">
            Back to home
          </Link>
        </p>
      </div>
    </main>
  );
}
