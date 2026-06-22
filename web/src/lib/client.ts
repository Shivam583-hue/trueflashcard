import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import {
  DeckService,
  FlashcardService,
  FolderService,
} from "@/gen/flashcard/v1/flashcard_pb";

const baseUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

const transport = createConnectTransport({
  baseUrl,
  fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
});

export const folderClient = createClient(FolderService, transport);
export const deckClient = createClient(DeckService, transport);
export const flashcardClient = createClient(FlashcardService, transport);
