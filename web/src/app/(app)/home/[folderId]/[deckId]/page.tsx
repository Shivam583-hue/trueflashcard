import { DeckView } from "@/components/app/deck-view";

export default async function DeckPage({
  params,
}: {
  params: Promise<{ folderId: string; deckId: string }>;
}) {
  const { folderId, deckId } = await params;
  return <DeckView folderId={folderId} deckId={deckId} />;
}
