import { Command } from 'commander';
import { api } from '../api';
import { outputJson, handleApiError, printTable } from '../output';
import chalk from 'chalk';

export function registerAnalyticsCommand(program: Command): void {
  const cmd = program
    .command('analytics')
    .description('View analytics and insights');

  cmd
    .command('stats')
    .description('Get user stats overview')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/stats');
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\n📊 Account Stats\n'));
      printTable(
        ['Metric', 'Value'],
        [
          ['Unread Feeds', String(res.data.unreadFeeds ?? 0)],
          ['Automations', String(res.data.userAutomations ?? 0)],
          ['Pending Posts', String(res.data.userPendingPosts ?? 0)],
          ['Failed Posts', String(res.data.userFailedPosts ?? 0)],
          ['Inactive Accounts', String(res.data.inactiveAccounts ?? 0)],
        ]
      );
    });

  cmd
    .command('followers')
    .description('Get followers count for all accounts')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/accounts/followers');
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data.data || res.data;
      console.log(chalk.bold(`\n👥 Total Followers: ${data.total_followers ?? '-'}\n`));

      if (data.followers_by_account?.length) {
        printTable(
          ['Account ID', 'Followers'],
          data.followers_by_account.map((a: any) => [
            String(a.account_id),
            String(a.followers),
          ])
        );
      }
    });

  cmd
    .command('engagement')
    .description('Get total engagement rate')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/accounts/engagement/rate');
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const rate = res.data?.data?.total_engagement_rate ?? res.data?.total_engagement_rate;
      console.log(chalk.bold('\n📈 Engagement Rate:'), `${rate ?? '-'}%`);
    });

  cmd
    .command('engagement-trend')
    .description('Get engagement trend over time')
    .requiredOption('--start <date>', 'Start date (YYYY-MM-DD)')
    .requiredOption('--end <date>', 'End date (YYYY-MM-DD)')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/accounts/engagement/trend', undefined, {
        start: opts.start,
        end: opts.end,
      });
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data.data || [];
      if (data.length === 0) {
        console.log(chalk.dim('No engagement data.'));
        return;
      }

      printTable(
        ['Date', 'Engagements'],
        data.map((d: any) => [d.date, String(d.engagements)])
      );
    });

  cmd
    .command('followers-growth')
    .description('Get followers growth over time')
    .requiredOption('--start <date>', 'Start date (YYYY-MM-DD)')
    .requiredOption('--end <date>', 'End date (YYYY-MM-DD)')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/accounts/followers/growth', undefined, {
        start: opts.start,
        end: opts.end,
      });
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data.data || [];
      if (data.length === 0) {
        console.log(chalk.dim('No growth data.'));
        return;
      }

      printTable(
        ['Date', 'Total Followers'],
        data.map((d: any) => [d.date, String(d.total_followers)])
      );
    });

  cmd
    .command('posts-count')
    .description('Get posts count over time')
    .requiredOption('--start <date>', 'Start date (YYYY-MM-DD)')
    .requiredOption('--end <date>', 'End date (YYYY-MM-DD)')
    .option('--accounts <ids...>', 'Filter by account IDs')
    .option('--post-type <type>', 'Filter: image, text, video')
    .option('--team <id>', 'Team ID')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = { start: opts.start, end: opts.end };
      if (opts.accounts) query.accounts = opts.accounts;
      if (opts.postType) query.post_type = opts.postType;
      if (opts.team) query.team = opts.team;

      const res = await api('GET', '/insights/posts/counts', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data.data || [];
      printTable(
        ['Date', 'Count'],
        data.map((d: any) => [d.date, String(d.count)])
      );
    });

  cmd
    .command('inbox')
    .description('Get open conversations count')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/inbox/unread-count');
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const count = res.data?.data?.open_msgs_count ?? '-';
      console.log(chalk.bold('\n💬 Open Conversations:'), count);
    });
}
