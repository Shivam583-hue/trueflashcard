import Link from "next/link";

import { Sidebar } from "@/components/app/sidebar";
import { LogoutButton } from "@/components/app/logout-button";
import { Logo } from "@/components/landing/logo";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen bg-[#08090a]">
      <Sidebar />
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-16 items-center justify-between border-b border-neutral-900 px-6">
          <Link href="/home" className="md:hidden">
            <Logo />
          </Link>
          <div className="hidden md:block" />
          <LogoutButton />
        </header>
        <main className="flex-1 px-6 py-8">{children}</main>
      </div>
    </div>
  );
}
