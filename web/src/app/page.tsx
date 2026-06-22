import { CTA } from "@/components/landing/cta";
import { Features } from "@/components/landing/features";
import { Hero } from "@/components/landing/hero";
import { LandingNav } from "@/components/landing/nav";
import { Logo } from "@/components/landing/logo";

export default function Home() {
  return (
    <div className="min-h-screen bg-[#08090a]">
      <LandingNav />
      <main>
        <Hero />
        <Features />
        <CTA />
      </main>
      <footer className="border-t border-neutral-900 py-10">
        <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 sm:flex-row">
          <Logo />
          <p className="text-xs text-neutral-600">
            © {new Date().getFullYear()} Flashcards
          </p>
        </div>
      </footer>
    </div>
  );
}
