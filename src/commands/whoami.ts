import { Command } from 'commander';
import { api } from '../api';
import { outputJson, handleApiError, printTable } from '../output';
import chalk from 'chalk';

export function registerWhoamiCommand(program: Command): void {
  program
    .command('whoami')
    .description('Show the currently authenticated user')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/user');
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const u = res.data;
      console.log(chalk.bold('\n👤 Authenticated User\n'));
      printTable(
        ['Field', 'Value'],
        [
          ['ID', String(u.id ?? '-')],
          ['Name', u.name ?? '-'],
          ['Email', u.email ?? '-'],
          ['Company', u.company ?? '-'],
          ['Verified', u.verified ? chalk.green('Yes') : chalk.yellow('No')],
        ]
      );
    });
}
