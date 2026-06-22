import Link from "next/link";

import { Logo } from "@/components/landing/logo";

export function LandingNav() {
  return (
    <header className="sticky top-0 z-20 border-b border-neutral-900/80 bg-[#08090a]/70 backdrop-blur-md">
      <nav className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
        <Link href="/" className="transition-opacity hover:opacity-80">
          <Logo />
        </Link>
        <Link
          href="/login"
          className="text-sm font-medium text-neutral-400 transition-colors duration-150 hover:text-neutral-100"
        >
          Sign in
        </Link>
      </nav>
    </header>
  );
}
