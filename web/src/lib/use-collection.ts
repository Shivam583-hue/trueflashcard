"use client";

import { useCallback, useEffect, useState } from "react";

import { isUnauthenticated } from "@/lib/connect-error";

export type CollectionStatus = "loading" | "ready" | "error";

export function useCollection<T>(loader: () => Promise<T[]>) {
  const [status, setStatus] = useState<CollectionStatus>("loading");
  const [items, setItems] = useState<T[]>([]);
  const [needsAuth, setNeedsAuth] = useState(false);

  const apply = useCallback((data: T[] | null, err?: unknown) => {
    if (data) {
      setItems(data);
      setStatus("ready");
      return;
    }
    if (isUnauthenticated(err)) setNeedsAuth(true);
    setStatus("error");
  }, []);

  const reload = useCallback(async () => {
    setStatus("loading");
    setNeedsAuth(false);
    try {
      apply(await loader());
    } catch (err) {
      apply(null, err);
    }
  }, [apply, loader]);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await loader();
        if (active) apply(data);
      } catch (err) {
        if (active) apply(null, err);
      }
    })();
    return () => {
      active = false;
    };
  }, [apply, loader]);

  return { status, items, setItems, needsAuth, reload };
}
