import chalk from 'chalk';
import Table from 'cli-table3';

export function outputJson(data: any): void {
  console.log(JSON.stringify(data, null, 2));
}

export function printError(msg: string): void {
  console.error(chalk.red('Error:'), msg);
}

export function printSuccess(msg: string): void {
  console.log(chalk.green('✓'), msg);
}

export function printWarning(msg: string): void {
  console.log(chalk.yellow('⚠'), msg);
}

export function printTable(headers: string[], rows: string[][]): void {
  const table = new Table({
    head: headers.map((h) => chalk.cyan(h)),
    style: { head: [], border: [] },
  });
  rows.forEach((row) => table.push(row));
  console.log(table.toString());
}

export function handleApiError(res: { ok: boolean; status: number; data: any }): void {
  if (!res.ok) {
    const msg =
      res.data?.message || res.data?.error || `API returned status ${res.status}`;
    printError(msg);
    process.exit(1);
  }
}

export function isJsonFlag(opts: any): boolean {
  // Walk up to find --json from parent command
  return opts?.json === true || opts?.parent?.json === true;
}

/**
 * Print standard pagination footer after a list command.
 */
export function printPaginationFooter(data: {
  currentPage?: number;
  lastPage?: number;
  total?: number;
  items?: any[];
}): void {
  const current = data.currentPage ?? 1;
  const last = data.lastPage ?? 1;
  const total = data.total ?? data.items?.length ?? 0;
  console.log(chalk.dim(`\nPage ${current} of ${last} | Total: ${total}`));
}
