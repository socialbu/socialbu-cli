import { Command } from 'commander';
import { api } from '../api';
import { outputJson, handleApiError, printTable } from '../output';
import chalk from 'chalk';

export function registerWhoamiCommand(program: Command): void {
  program
    .command('whoami')
    .description('Show current user info and account stats')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/stats');
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\nSocialBu Account Stats\n'));
      printTable(
        ['Metric', 'Value'],
        [
          ['Unread Feeds', String(res.data.unreadFeeds ?? '-')],
          ['Active Automations', String(res.data.userAutomations ?? '-')],
          ['Pending Posts', String(res.data.userPendingPosts ?? '-')],
          ['Failed Posts', String(res.data.userFailedPosts ?? '-')],
          ['Inactive Accounts', String(res.data.inactiveAccounts ?? '-')],
        ]
      );
    });
}
