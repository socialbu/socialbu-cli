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

  cmd
    .command('posts-metrics')
    .description('Get post metrics (likes, comments, etc.) over time')
    .requiredOption('--start <date>', 'Start date (YYYY-MM-DD)')
    .requiredOption('--end <date>', 'End date (YYYY-MM-DD)')
    .requiredOption('--post-type <type>', 'Post type: image, text, video')
    .requiredOption('--metrics <list>', 'Comma-separated metrics (likes,comments,views,etc)')
    .option('--accounts <ids...>', 'Filter by account IDs')
    .option('--team <id>', 'Team ID')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = {
        start: opts.start,
        end: opts.end,
        post_type: opts.postType,
        metrics: opts.metrics,
      };
      if (opts.accounts) query.accounts = opts.accounts;
      if (opts.team) query.team = opts.team;

      const res = await api('GET', '/insights/posts/metrics', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data?.data || res.data;
      console.log(chalk.bold('\n📊 Posts Metrics\n'));
      if (data?.items) {
        for (const [metric, values] of Object.entries(data.items)) {
          if (Array.isArray(values) && values.length > 0) {
            console.log(chalk.cyan(`${metric}:`));
            printTable(
              ['Date', 'Value'],
              values.map((v: any) => [v.date, String(v.value)])
            );
            console.log();
          }
        }
      }
    });

  cmd
    .command('top-posts')
    .description('Get top performing posts')
    .requiredOption('--start <date>', 'Start date (YYYY-MM-DD)')
    .requiredOption('--end <date>', 'End date (YYYY-MM-DD)')
    .requiredOption('--metrics <list>', 'Comma-separated metrics to rank by')
    .option('--accounts <ids...>', 'Filter by account IDs')
    .option('--team <id>', 'Team ID')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = {
        start: opts.start,
        end: opts.end,
        metrics: opts.metrics,
      };
      if (opts.accounts) query.accounts = opts.accounts;
      if (opts.team) query.team = opts.team;

      const res = await api('GET', '/insights/posts/top_posts', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data?.data || [];
      if (data.length === 0) {
        console.log(chalk.dim('No top posts found.'));
        return;
      }

      console.log(chalk.bold(`\n🔥 Top Posts (${data.length})\n`));
      printTable(
        ['ID', 'Content', 'Published'],
        data.map((p: any) => [
          String(p.id || '-'),
          (p.content || '').slice(0, 50) + ((p.content || '').length > 50 ? '…' : ''),
          p.published_at || '-',
        ])
      );
    });

  cmd
    .command('team-metrics')
    .description('Get team performance metrics')
    .requiredOption('--start <date>', 'Start date (YYYY-MM-DD)')
    .requiredOption('--end <date>', 'End date (YYYY-MM-DD)')
    .requiredOption('--metrics <list>', 'Comma-separated metrics to fetch')
    .option('--accounts <ids...>', 'Filter by account IDs')
    .option('--team <id>', 'Team ID')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = {
        start: opts.start,
        end: opts.end,
        metrics: opts.metrics,
      };
      if (opts.accounts) query.accounts = opts.accounts;
      if (opts.team) query.team = opts.team;

      const res = await api('GET', '/insights/teams/metrics', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data?.data || [];
      if (data.length === 0) {
        console.log(chalk.dim('No team metrics found.'));
        return;
      }

      console.log(chalk.bold(`\n👥 Team Metrics (${data.length} members)\n`));
      printTable(
        ['Member', 'Total Engagements'],
        data.map((m: any) => [
          m.member_name || `User ${m.member_id}`,
          String(m.total_engagements ?? 0),
        ])
      );
    });

  cmd
    .command('team-activity')
    .description('Get team activity logs')
    .option('--limit <n>', 'Number of logs to fetch (1-100)', '5')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/teams/activity', undefined, {
        limit: Number(opts.limit),
      });
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data?.data || [];
      if (data.length === 0) {
        console.log(chalk.dim('No team activity found.'));
        return;
      }

      console.log(chalk.bold(`\n📋 Team Activity (${data.length})\n`));
      printTable(
        ['Type', 'Title', 'When'],
        data.map((a: any) => [
          a.type || '-',
          (a.title || '-').slice(0, 50),
          a.timestamp || '-',
        ])
      );
    });

  cmd
    .command('automation-logs')
    .description('Get automation activity logs')
    .option('--limit <n>', 'Number of logs to fetch (1-100)', '5')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/insights/automations/logs', undefined, {
        limit: Number(opts.limit),
      });
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const data = res.data?.data || [];
      if (data.length === 0) {
        console.log(chalk.dim('No automation logs found.'));
        return;
      }

      console.log(chalk.bold(`\n🤖 Automation Logs (${data.length})\n`));
      printTable(
        ['Title', 'Description', 'When'],
        data.map((a: any) => [
          (a.title || '-').slice(0, 30),
          (a.description || '-').slice(0, 50),
          a.timestamp || '-',
        ])
      );
    });
}
