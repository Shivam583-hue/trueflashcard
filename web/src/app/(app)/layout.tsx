import { Sidebar } from "@/components/app/sidebar";
import { LogoutButton } from "@/components/app/logout-button";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen bg-[#08090a]">
      <Sidebar />
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-16 items-center justify-between border-b border-neutral-900 px-6">
          <h1 className="text-sm font-medium text-neutral-300">Folders</h1>
          <LogoutButton />
        </header>
        <main className="flex-1 px-6 py-8">{children}</main>
      </div>
    </div>
  );
}
