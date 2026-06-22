import Link from "next/link";
import { FolderSimpleIcon } from "@phosphor-icons/react/dist/ssr";

import { Logo } from "@/components/landing/logo";

export function Sidebar() {
  return (
    <aside className="hidden w-60 shrink-0 flex-col border-r border-neutral-900 bg-[#0a0b0c] md:flex">
      <div className="flex h-16 items-center px-6">
        <Link href="/" className="transition-opacity hover:opacity-80">
          <Logo />
        </Link>
      </div>
      <nav className="px-3 py-2">
        <Link
          href="/home"
          className="flex items-center gap-3 rounded-lg bg-neutral-900 px-3 py-2 text-sm font-medium text-neutral-100"
        >
          <FolderSimpleIcon size={18} weight="duotone" />
          Folders
        </Link>
      </nav>
    </aside>
  );
}
