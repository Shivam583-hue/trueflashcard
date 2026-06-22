export const MAX_IMPORT_CARDS = 1000;
const MAX_TEXT = 10000;

export type CardData = { front: string; back: string };

function stripFences(raw: string): string {
  const t = raw.trim();
  const fenced = t.match(/^```(?:json)?\s*([\s\S]*?)\s*```$/i);
  return fenced ? fenced[1].trim() : t;
}

function isRecord(v: unknown): v is Record<string, unknown> {
  return typeof v === "object" && v !== null && !Array.isArray(v);
}

function parseJson(text: string): { value: unknown } | { error: string } {
  const cleaned = stripFences(text);
  if (cleaned === "") return { error: "Paste the JSON for your deck." };
  try {
    return { value: JSON.parse(cleaned) };
  } catch (e) {
    return { error: `That isn't valid JSON — ${(e as Error).message}` };
  }
}

function validateCards(value: unknown): { cards: CardData[] } | { errors: string[] } {
  if (!Array.isArray(value)) {
    return { errors: ['"cards" must be an array of { front, back }.'] };
  }
  if (value.length === 0) return { errors: ["Add at least one card."] };
  if (value.length > MAX_IMPORT_CARDS) {
    return { errors: [`Too many cards: ${value.length} (max ${MAX_IMPORT_CARDS}).`] };
  }

  const errors: string[] = [];
  const cards: CardData[] = [];
  value.forEach((raw, i) => {
    const n = i + 1;
    if (!isRecord(raw)) {
      errors.push(`Card ${n}: must be an object with "front" and "back".`);
      return;
    }
    const { front, back } = raw;
    if (typeof front !== "string" || front.trim() === "") {
      errors.push(`Card ${n}: "front" is required.`);
      return;
    }
    if (typeof back !== "string" || back.trim() === "") {
      errors.push(`Card ${n}: "back" is required.`);
      return;
    }
    if (front.length > MAX_TEXT || back.length > MAX_TEXT) {
      errors.push(`Card ${n}: front and back must be under ${MAX_TEXT} characters.`);
      return;
    }
    cards.push({ front: front.trim(), back: back.trim() });
  });

  if (errors.length > 0) return { errors };
  return { cards };
}

export type DeckParseResult =
  | { ok: true; name: string; description: string; cards: CardData[] }
  | { ok: false; errors: string[] };

export function parseDeckImport(text: string): DeckParseResult {
  const parsed = parseJson(text);
  if ("error" in parsed) return { ok: false, errors: [parsed.error] };

  const root = parsed.value;
  if (!isRecord(root)) {
    return {
      ok: false,
      errors: ['Expected an object like { "name": "…", "cards": [ … ] }.'],
    };
  }

  const errors: string[] = [];
  const name = root.name;
  if (typeof name !== "string" || name.trim() === "") {
    errors.push('"name" is required.');
  }
  const description =
    typeof root.description === "string" ? root.description.trim() : "";

  const cardsResult = validateCards(root.cards);
  if ("errors" in cardsResult) errors.push(...cardsResult.errors);

  if (errors.length > 0) return { ok: false, errors };
  return {
    ok: true,
    name: (name as string).trim(),
    description,
    cards: (cardsResult as { cards: CardData[] }).cards,
  };
}

export type CardsParseResult =
  | { ok: true; cards: CardData[] }
  | { ok: false; errors: string[] };

export function parseCardsImport(text: string): CardsParseResult {
  const parsed = parseJson(text);
  if ("error" in parsed) return { ok: false, errors: [parsed.error] };

  const root = parsed.value;
  const cardsValue = Array.isArray(root)
    ? root
    : isRecord(root)
      ? root.cards
      : undefined;

  const cardsResult = validateCards(cardsValue);
  if ("errors" in cardsResult) return { ok: false, errors: cardsResult.errors };
  return { ok: true, cards: cardsResult.cards };
}

export const DECK_PROMPT = `Create flashcards about <TOPIC>.
Respond with ONLY a JSON object (no markdown, no commentary) in exactly this shape:
{
  "name": "Deck title",
  "description": "one-line description (optional)",
  "cards": [
    { "front": "question or term", "back": "answer or definition" }
  ]
}
Use 10-20 cards. Keep each front concise and each back short. Output JSON only.`;

export const CARDS_PROMPT = `Create flashcards about <TOPIC>.
Respond with ONLY a JSON array (no markdown, no commentary) in exactly this shape:
[
  { "front": "question or term", "back": "answer or definition" }
]
Use 10-20 cards. Keep each front concise and each back short. Output JSON only.`;
