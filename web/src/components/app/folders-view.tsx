"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import { Code, ConnectError } from "@connectrpc/connect";
import { FolderSimpleIcon, PlusIcon } from "@phosphor-icons/react";

import { folderClient } from "@/lib/client";
import { Button } from "@/components/ui/button";
import { easeEmphasized } from "@/lib/motion";

type Folder = { id: string; name: string };
type Status = "loading" | "ready" | "error";

export function FoldersView() {
  const reduce = useReducedMotion();
  const [status, setStatus] = useState<Status>("loading");
  const [folders, setFolders] = useState<Folder[]>([]);
  const [needsAuth, setNeedsAuth] = useState(false);

  const applyResult = useCallback(
    (folders: Folder[] | null, err?: unknown) => {
      if (folders) {
        setFolders(folders);
        setStatus("ready");
        return;
      }
      if (err && ConnectError.from(err).code === Code.Unauthenticated) {
        setNeedsAuth(true);
      }
      setStatus("error");
    },
    [],
  );

  const load = useCallback(async () => {
    setStatus("loading");
    setNeedsAuth(false);
    try {
      const res = await folderClient.listFolders({});
      applyResult(res.folders.map((f) => ({ id: f.id, name: f.name })));
    } catch (err) {
      applyResult(null, err);
    }
  }, [applyResult]);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const res = await folderClient.listFolders({});
        if (active) applyResult(res.folders.map((f) => ({ id: f.id, name: f.name })));
      } catch (err) {
        if (active) applyResult(null, err);
      }
    })();
    return () => {
      active = false;
    };
  }, [applyResult]);

  const onCreated = useCallback((folder: Folder) => {
    setFolders((prev) => [...prev, folder]);
  }, []);

  return (
    <div className="mx-auto max-w-5xl">
      <div className="flex items-center justify-between gap-4">
        <p className="text-sm text-neutral-500">
          {status === "ready"
            ? `${folders.length} ${folders.length === 1 ? "folder" : "folders"}`
            : " "}
        </p>
        <CreateFolder onCreated={onCreated} disabled={status === "error"} />
      </div>

      <div className="mt-6">
        {status === "loading" && <FolderSkeletons />}
        {status === "error" && (
          <ErrorState needsAuth={needsAuth} onRetry={load} />
        )}
        {status === "ready" && folders.length === 0 && <EmptyState />}
        {status === "ready" && folders.length > 0 && (
          <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <AnimatePresence initial={!reduce}>
              {folders.map((folder, i) => (
                <motion.li
                  key={folder.id}
                  layout={!reduce}
                  initial={reduce ? false : { opacity: 0, y: 12 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={reduce ? undefined : { opacity: 0, scale: 0.97 }}
                  transition={{
                    duration: 0.4,
                    delay: Math.min(i * 0.04, 0.3),
                    ease: easeEmphasized,
                  }}
                >
                  <FolderCard folder={folder} />
                </motion.li>
              ))}
            </AnimatePresence>
          </ul>
        )}
      </div>
    </div>
  );
}

function FolderCard({ folder }: { folder: Folder }) {
  return (
    <Link
      href={`/home`}
      className="group flex items-center gap-3 rounded-xl border border-neutral-900 bg-[#0a0b0c] p-4 transition-[transform,border-color] duration-150 [transition-timing-function:var(--ease-out)] hover:border-neutral-700 active:scale-[0.99] [@media(hover:hover)]:hover:-translate-y-0.5"
    >
      <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-neutral-900 text-neutral-300">
        <FolderSimpleIcon size={18} weight="duotone" />
      </span>
      <span className="truncate text-sm font-medium text-neutral-100">
        {folder.name}
      </span>
    </Link>
  );
}

function CreateFolder({
  onCreated,
  disabled,
}: {
  onCreated: (folder: Folder) => void;
  disabled: boolean;
}) {
  const reduce = useReducedMotion();
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [pending, setPending] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    const trimmed = name.trim();
    if (!trimmed || pending) return;
    setPending(true);
    try {
      const res = await folderClient.createFolder({ name: trimmed });
      if (res.folder) onCreated({ id: res.folder.id, name: res.folder.name });
      setName("");
      setOpen(false);
    } catch {
      // Surface nothing destructive; keep the input open so the user can retry.
    } finally {
      setPending(false);
    }
  }

  if (!open) {
    return (
      <Button variant="ghost" onClick={() => setOpen(true)} disabled={disabled}>
        <PlusIcon size={16} />
        New folder
      </Button>
    );
  }

  return (
    <motion.form
      onSubmit={submit}
      initial={reduce ? false : { opacity: 0, scale: 0.98 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.2, ease: easeEmphasized }}
      className="flex items-center gap-2"
      style={{ transformOrigin: "right center" }}
    >
      <input
        autoFocus
        value={name}
        onChange={(e) => setName(e.target.value)}
        onBlur={() => !name && setOpen(false)}
        placeholder="Folder name"
        maxLength={200}
        className="h-9 w-44 rounded-lg border border-neutral-800 bg-neutral-900 px-3 text-sm text-neutral-100 placeholder:text-neutral-600 focus:border-neutral-600 focus:outline-none"
      />
      <Button type="submit" disabled={pending || !name.trim()}>
        Add
      </Button>
    </motion.form>
  );
}

function FolderSkeletons() {
  return (
    <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {Array.from({ length: 6 }).map((_, i) => (
        <li
          key={i}
          className="skeleton flex h-[68px] items-center gap-3 rounded-xl border border-neutral-900 bg-[#0a0b0c] p-4"
        >
          <span className="h-9 w-9 rounded-lg bg-neutral-900" />
          <span className="h-3 w-24 rounded bg-neutral-900" />
        </li>
      ))}
    </ul>
  );
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-neutral-900 py-20 text-center">
      <span className="flex h-12 w-12 items-center justify-center rounded-xl bg-neutral-900 text-neutral-400">
        <FolderSimpleIcon size={24} weight="duotone" />
      </span>
      <h3 className="mt-4 text-sm font-medium text-neutral-200">No folders yet</h3>
      <p className="mt-1 max-w-xs text-sm text-neutral-500">
        Create your first folder to start organizing decks.
      </p>
    </div>
  );
}

function ErrorState({
  needsAuth,
  onRetry,
}: {
  needsAuth: boolean;
  onRetry: () => void;
}) {
  return (
    <div className="flex flex-col items-center justify-center rounded-2xl border border-neutral-900 py-20 text-center">
      <h3 className="text-sm font-medium text-neutral-200">
        {needsAuth ? "Please sign in" : "Could not load folders"}
      </h3>
      <p className="mt-1 max-w-xs text-sm text-neutral-500">
        {needsAuth
          ? "Your session is missing or expired."
          : "The server is unreachable right now."}
      </p>
      <div className="mt-5">
        {needsAuth ? (
          <Button onClick={() => (window.location.href = "/login")}>
            Go to sign in
          </Button>
        ) : (
          <Button variant="ghost" onClick={onRetry}>
            Try again
          </Button>
        )}
      </div>
    </div>
  );
}
