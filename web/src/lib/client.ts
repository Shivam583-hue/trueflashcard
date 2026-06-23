import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import {
  DeckService,
  FlashcardService,
  FolderService,
} from "@/gen/flashcard/v1/flashcard_pb";

const transport = createConnectTransport({
  baseUrl: "/rpc",
});

export const folderClient = createClient(FolderService, transport);
export const deckClient = createClient(DeckService, transport);
export const flashcardClient = createClient(FlashcardService, transport);
