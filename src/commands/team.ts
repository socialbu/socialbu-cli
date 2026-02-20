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

export function registerTeamCommand(program: Command): void {
  const cmd = program
    .command('team')
    .description('Manage teams');

  cmd
    .command('list')
    .description('List your teams')
    .option('--type <type>', 'Filter: created, joined')
    .option('--page <n>', 'Page number', parsePage, 1)
    .option('--all', 'Fetch all pages')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = {};
      if (opts.type) query.type = [opts.type];
      if (!opts.all) query.page = opts.page;

      const res = opts.all
        ? await apiPaginated('GET', '/teams', undefined, query)
        : await api('GET', '/teams', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const items = res.data.items || [];
      if (items.length === 0) {
        console.log(chalk.dim('No teams found.'));
        return;
      }

      printTable(
        ['ID', 'Name', 'Members', 'Accounts'],
        items.map((t: any) => [
          String(t.id),
          t.name || '-',
          String(t.members?.length ?? '-'),
          String(t.accounts?.length ?? '-'),
        ])
      );
      console.log(chalk.dim(`\nPage ${res.data.currentPage ?? 1} of ${res.data.lastPage ?? 1} | Total: ${res.data.total ?? items.length}`));
    });

  cmd
    .command('create')
    .description('Create a new team')
    .requiredOption('--name <name>', 'Team name')
    .requiredOption('--accounts <ids...>', 'Account IDs to include')
    .option('--require-approval', 'Require content approval')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const body: any = {
        name: opts.name,
        accounts: opts.accounts.map((id: string) => ({ id: Number(id) })),
      };
      if (opts.requireApproval) body.requires_content_approval = true;

      const res = await api('POST', '/teams', body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess(res.data.message || 'Team created.');
    });

  cmd
    .command('delete <id>')
    .description('Delete a team')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('DELETE', `/teams/${id}`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess(res.data.message || 'Team deleted.');
    });
}
