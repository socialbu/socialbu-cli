import { Command } from 'commander';
import { api, apiPaginated } from '../api';
import { outputJson, handleApiError, printTable, printSuccess } from '../output';
import chalk from 'chalk';

function parsePage(value: string): number {
  const parsed = Number.parseInt(value, 10);
  if (!Number.isInteger(parsed) || parsed < 1) {
    throw new Error('Page must be a positive integer.');
  }
  return parsed;
}

export function registerAccountCommand(program: Command): void {
  const cmd = program
    .command('account')
    .description('Manage connected social media accounts');

  cmd
    .command('list')
    .description('List all connected social accounts')
    .option('--type <type>', 'Filter by type: all, shared, user', 'all')
    .option('--page <n>', 'Page number', parsePage, 1)
    .option('--all', 'Fetch all pages')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = { type: [opts.type] };
      if (!opts.all) query.page = opts.page;

      const res = opts.all
        ? await apiPaginated('GET', '/accounts', undefined, query)
        : await api('GET', '/accounts', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const items = res.data.items || [];
      if (items.length === 0) {
        console.log(chalk.dim('No accounts found.'));
        return;
      }

      printTable(
        ['ID', 'Name', 'Type', 'Platform', 'Status'],
        items.map((a: any) => [
          String(a.id),
          a.name || '-',
          a.type || '-',
          a.provider || a.platform || '-',
          a.status || (a.is_active ? 'active' : 'inactive'),
        ])
      );
      console.log(chalk.dim(`\nPage ${res.data.currentPage ?? 1} of ${res.data.lastPage ?? 1} | Total: ${res.data.total ?? items.length}`));
    });

  cmd
    .command('get <id>')
    .description('Get details of a social account by ID')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('GET', `/accounts/${id}`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const a = res.data;
      console.log(chalk.bold('\nAccount Details\n'));
      for (const [key, val] of Object.entries(a)) {
        if (val !== null && val !== undefined && typeof val !== 'object') {
          console.log(`  ${chalk.cyan(key)}: ${val}`);
        }
      }
    });

  cmd
    .command('delete <id>')
    .description('Delete a social account')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('DELETE', `/accounts/${id}`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess(res.data.message || 'Account deleted.');
    });

  cmd
    .command('update <id>')
    .description('Update a social account name')
    .requiredOption('--name <name>', 'New name for the account')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('PUT', `/accounts/${id}`, { name: opts.name });
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess(res.data.message || 'Account updated.');
    });
}
