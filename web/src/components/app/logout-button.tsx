"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { SignOutIcon } from "@phosphor-icons/react";

export function LogoutButton() {
  const router = useRouter();
  const [pending, setPending] = useState(false);

  async function logout() {
    setPending(true);
    try {
      await fetch("/api/auth/logout", { method: "POST" });
    } catch {}
    router.push("/");
  }

  return (
    <button
      onClick={logout}
      disabled={pending}
      className="inline-flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-neutral-400 transition-colors duration-150 [transition-timing-function:var(--ease-out)] hover:text-neutral-100 active:scale-[0.97] disabled:opacity-50"
    >
      <SignOutIcon size={18} />
      Sign out
    </button>
  );
}
