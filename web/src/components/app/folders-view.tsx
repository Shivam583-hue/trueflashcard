"use client";

import { useCallback, useState } from "react";
import Link from "next/link";
import { AnimatePresence, motion, useReducedMotion } from "motion/react";
import {
  FolderSimpleIcon,
  PencilSimpleIcon,
  PlusIcon,
} from "@phosphor-icons/react";

import { folderClient } from "@/lib/client";
import { useCollection } from "@/lib/use-collection";
import { easeEmphasized } from "@/lib/motion";
import { Button } from "@/components/ui/button";
import { ConfirmDelete } from "@/components/app/confirm-delete";
import { InlineNameForm } from "@/components/app/inline-name-form";
import { StudyBanner } from "@/components/app/study-banner";
import {
  CardSkeletons,
  EmptyState,
  ErrorState,
} from "@/components/app/states";

type Folder = { id: string; name: string };

const listFolders = () =>
  folderClient
    .listFolders({})
    .then((res) => res.folders.map((f) => ({ id: f.id, name: f.name })));

export function FoldersView() {
  const reduce = useReducedMotion();
  const { status, items, setItems, needsAuth, reload } =
    useCollection<Folder>(listFolders);
  const [creating, setCreating] = useState(false);

  const create = useCallback(
    async (name: string) => {
      const res = await folderClient.createFolder({ name });
      if (res.folder) {
        const folder = { id: res.folder.id, name: res.folder.name };
        setItems((prev) => [...prev, folder]);
      }
      setCreating(false);
    },
    [setItems],
  );

  const rename = useCallback(
    (id: string, name: string) =>
      setItems((prev) => prev.map((f) => (f.id === id ? { ...f, name } : f))),
    [setItems],
  );

  const remove = useCallback(
    (id: string) => setItems((prev) => prev.filter((f) => f.id !== id)),
    [setItems],
  );

  return (
    <div className="mx-auto max-w-5xl">
      <div className="mb-8">
        <StudyBanner />
      </div>

      <div className="flex min-h-9 items-center justify-between gap-4">
        <h1 className="text-base font-medium text-neutral-100">Your folders</h1>
        {creating ? (
          <InlineNameForm
            placeholder="Folder name"
            submitLabel="Add"
            onSubmit={create}
            onCancel={() => setCreating(false)}
          />
        ) : (
          <Button
            variant="ghost"
            onClick={() => setCreating(true)}
            disabled={status === "error"}
          >
            <PlusIcon size={16} />
            New folder
          </Button>
        )}
      </div>

      <div className="mt-6">
        {status === "loading" && <CardSkeletons />}
        {status === "error" && (
          <ErrorState
            title="Could not load folders"
            needsAuth={needsAuth}
            onRetry={reload}
          />
        )}
        {status === "ready" && items.length === 0 && (
          <EmptyState
            icon={FolderSimpleIcon}
            title="No folders yet"
            body="Create your first folder to start organizing decks."
          />
        )}
        {status === "ready" && items.length > 0 && (
          <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <AnimatePresence initial={!reduce}>
              {items.map((folder, i) => (
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
                  <FolderCard folder={folder} onRenamed={rename} onDeleted={remove} />
                </motion.li>
              ))}
            </AnimatePresence>
          </ul>
        )}
      </div>
    </div>
  );
}

function FolderCard({
  folder,
  onRenamed,
  onDeleted,
}: {
  folder: Folder;
  onRenamed: (id: string, name: string) => void;
  onDeleted: (id: string) => void;
}) {
  const [editing, setEditing] = useState(false);

  async function rename(name: string) {
    const res = await folderClient.updateFolder({ id: folder.id, name });
    if (res.folder) onRenamed(folder.id, res.folder.name);
    setEditing(false);
  }

  async function remove() {
    await folderClient.deleteFolder({ id: folder.id });
    onDeleted(folder.id);
  }

  if (editing) {
    return (
      <div className="rounded-xl border border-neutral-800 bg-[#0a0b0c] p-3">
        <InlineNameForm
          placeholder="Folder name"
          initial={folder.name}
          submitLabel="Save"
          onSubmit={rename}
          onCancel={() => setEditing(false)}
        />
      </div>
    );
  }

  return (
    <div className="group relative flex items-center gap-3 rounded-xl border border-neutral-900 bg-[#0a0b0c] p-4 transition-[transform,border-color] duration-150 [transition-timing-function:var(--ease-out)] hover:border-neutral-700 [@media(hover:hover)]:hover:-translate-y-0.5">
      <Link
        href={`/home/${folder.id}`}
        className="absolute inset-0 rounded-xl focus:outline-none focus-visible:ring-2 focus-visible:ring-neutral-600"
        aria-label={`Open ${folder.name}`}
      />
      <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-neutral-900 text-neutral-300">
        <FolderSimpleIcon size={18} weight="duotone" />
      </span>
      <span className="min-w-0 flex-1 truncate text-sm font-medium text-neutral-100">
        {folder.name}
      </span>
      <div className="relative z-10 flex items-center opacity-70 transition-opacity duration-150 group-hover:opacity-100">
        <button
          aria-label={`Rename ${folder.name}`}
          onClick={() => setEditing(true)}
          className="rounded-md p-1.5 text-neutral-600 transition-colors duration-150 [transition-timing-function:var(--ease-out)] hover:text-neutral-300 active:scale-[0.95]"
        >
          <PencilSimpleIcon size={16} />
        </button>
        <ConfirmDelete label={`Delete ${folder.name}`} onConfirm={remove} />
      </div>
    </div>
  );
}
