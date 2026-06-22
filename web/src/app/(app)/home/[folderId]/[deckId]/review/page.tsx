import { ReviewView } from "@/components/app/review-view";

export default async function ReviewPage({
  params,
}: {
  params: Promise<{ folderId: string; deckId: string }>;
}) {
  const { folderId, deckId } = await params;
  return <ReviewView folderId={folderId} deckId={deckId} />;
}
