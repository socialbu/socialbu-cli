import { Command } from 'commander';
import { api, apiPaginated } from '../api';
import { outputJson, handleApiError, printTable, printSuccess, printPaginationFooter } from '../output';
import { parsePage } from '../utils';
import chalk from 'chalk';


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
      printPaginationFooter(res.data);
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
    .command('update <id>')
    .description('Update a team')
    .requiredOption('--name <name>', 'Team name')
    .requiredOption('--accounts <ids...>', 'Account IDs to include')
    .option('--require-approval', 'Require content approval')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const body: any = {
        name: opts.name,
        accounts: opts.accounts.map((id: string) => ({ id: Number(id) })),
      };
      if (opts.requireApproval) body.requires_content_approval = true;

      const res = await api('PUT', `/teams/${id}`, body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess(res.data.message || 'Team updated.');
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
