import { StudyView } from "@/components/app/study-view";

export default async function DeckStudyPage({
  params,
}: {
  params: Promise<{ folderId: string; deckId: string }>;
}) {
  const { folderId, deckId } = await params;
  return (
    <StudyView
      deckId={deckId}
      backHref={`/home/${folderId}/${deckId}`}
      backLabel="Back to deck"
    />
  );
}
