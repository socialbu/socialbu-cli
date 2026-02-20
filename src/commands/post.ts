import { Command } from 'commander';
import { api } from '../api';
import { outputJson, handleApiError, printTable, printSuccess } from '../output';
import chalk from 'chalk';

export function registerPostCommand(program: Command): void {
  const cmd = program
    .command('post')
    .description('Create, list, and manage posts');

  cmd
    .command('create')
    .description('Create a new post')
    .requiredOption('--accounts <ids...>', 'Account ID(s) to post to (space-separated)')
    .requiredOption('--publish-at <datetime>', 'Publish time in UTC (YYYY-MM-DD HH:mm:ss)')
    .option('--content <text>', 'Post content/text')
    .option('--draft', 'Save as draft instead of scheduling')
    .option('--queue <ids...>', 'Queue ID(s) to add this post to')
    .option('--team <id>', 'Team ID')
    .option('--postback-url <url>', 'URL to receive post status callbacks')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const body: any = {
        accounts: opts.accounts.map(Number),
        publish_at: opts.publishAt,
      };
      if (opts.content) body.content = opts.content;
      if (opts.draft) body.draft = true;
      if (opts.queue) body.queue_ids = opts.queue.map(Number);
      if (opts.team) body.team_id = Number(opts.team);
      if (opts.postbackUrl) body.postback_url = opts.postbackUrl;

      const res = await api('POST', '/posts', body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess('Post created!');
      if (res.data.posts && res.data.posts.length > 0) {
        printTable(
          ['Post ID', 'Account ID', 'Status'],
          res.data.posts.map((p: any) => [
            String(p.id || '-'),
            String(p.account_id || '-'),
            p.status || '-',
          ])
        );
      }
    });

  cmd
    .command('list')
    .description('List scheduled or published posts')
    .option('--type <types...>', 'Post types: scheduled, published, draft, awaiting_approval', ['scheduled'])
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/posts', undefined, { type: opts.type });
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const items = res.data.items || [];
      if (items.length === 0) {
        console.log(chalk.dim('No posts found.'));
        return;
      }

      printTable(
        ['ID', 'Content', 'Publish At', 'Status'],
        items.map((p: any) => [
          String(p.id),
          (p.content || '').slice(0, 60) + ((p.content || '').length > 60 ? '…' : ''),
          p.publish_at || '-',
          p.status || '-',
        ])
      );
      console.log(chalk.dim(`\nPage ${res.data.currentPage ?? 1} of ${res.data.lastPage ?? 1} | Total: ${res.data.total ?? items.length}`));
    });

  cmd
    .command('get <id>')
    .description('Get a post by ID')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('GET', `/posts/${id}`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\nPost Details\n'));
      for (const [key, val] of Object.entries(res.data)) {
        if (val !== null && val !== undefined && typeof val !== 'object') {
          console.log(`  ${chalk.cyan(key)}: ${val}`);
        }
      }
    });

  cmd
    .command('delete <id>')
    .description('Delete a post')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('DELETE', `/posts/${id}`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess(res.data.message || 'Post deleted.');
    });

  cmd
    .command('update <id>')
    .description('Update a post')
    .option('--content <text>', 'Updated content')
    .option('--publish-at <datetime>', 'Updated publish time (YYYY-MM-DD HH:mm:ss UTC)')
    .option('--accounts <ids...>', 'Updated account IDs')
    .option('--team <id>', 'Team ID')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const body: any = {};
      if (opts.content) body.content = opts.content;
      if (opts.publishAt) body.publish_at = opts.publishAt;
      if (opts.accounts) body.accounts = opts.accounts.map(Number);
      if (opts.team) body.team_id = Number(opts.team);

      const res = await api('PUT', `/posts/${id}`, body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess(res.data.message || 'Post updated.');
    });
}
