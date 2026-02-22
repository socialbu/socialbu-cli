/**
 * Shared utility helpers used across command files.
 */

/**
 * Parse a CLI page number argument.
 * Throws if the value is not a positive integer.
 */
export function parsePage(value: string): number {
  const parsed = Number.parseInt(value, 10);
  if (!Number.isInteger(parsed) || parsed < 1) {
    throw new Error('Page must be a positive integer.');
  }
  return parsed;
}
