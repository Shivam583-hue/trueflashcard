import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import {
  DeckService,
  FlashcardService,
  FolderService,
  StudyService,
} from "@/gen/flashcard/v1/flashcard_pb";

const transport = createConnectTransport({
  baseUrl: "/rpc",
  fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
});

export const folderClient = createClient(FolderService, transport);
export const deckClient = createClient(DeckService, transport);
export const flashcardClient = createClient(FlashcardService, transport);
export const studyClient = createClient(StudyService, transport);
