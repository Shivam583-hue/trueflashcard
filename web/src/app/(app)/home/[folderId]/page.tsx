import { DecksView } from "@/components/app/decks-view";

export default async function FolderPage({
  params,
}: {
  params: Promise<{ folderId: string }>;
}) {
  const { folderId } = await params;
  return <DecksView folderId={folderId} />;
}
