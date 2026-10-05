export function parseEnum<T extends string>(e: Record<string, T>, value: unknown): T | null {
  return Object.values(e).find((v) => v === value) ?? null;
}
