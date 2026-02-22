import { Command } from 'commander';
import { api, apiPaginated } from '../api';
import { outputJson, handleApiError, printTable, printSuccess, printPaginationFooter } from '../output';
import { parsePage } from '../utils';
import chalk from 'chalk';


export function registerNotificationsCommand(program: Command): void {
  const cmd = program
    .command('notifications')
    .alias('notif')
    .description('Manage notifications');

  cmd
    .command('list')
    .description('List all notifications')
    .option('--page <n>', 'Page number', parsePage, 1)
    .option('--all', 'Fetch all pages')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = {};
      if (!opts.all) query.page = opts.page;

      const res = opts.all
        ? await apiPaginated('GET', '/notifications', undefined, query)
        : await api('GET', '/notifications', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const items = res.data.items || [];
      if (items.length === 0) {
        console.log(chalk.dim('No notifications.'));
        return;
      }

      printTable(
        ['ID', 'Title', 'Level', 'Unread', 'Created'],
        items.map((n: any) => [
          String(n.id),
          (n.title || n.body || '').slice(0, 50),
          n.level || '-',
          n.is_unread ? chalk.yellow('●') : chalk.dim('○'),
          n.created_at || '-',
        ])
      );
      printPaginationFooter(res.data);
    });

  cmd
    .command('unread')
    .description('List unread notifications')
    .option('--page <n>', 'Page number', parsePage, 1)
    .option('--all', 'Fetch all pages')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = {};
      if (!opts.all) query.page = opts.page;

      const res = opts.all
        ? await apiPaginated('GET', '/notifications/unread', undefined, query)
        : await api('GET', '/notifications/unread', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const items = res.data.items || [];
      if (items.length === 0) {
        console.log(chalk.dim('No unread notifications. 🎉'));
        return;
      }

      printTable(
        ['ID', 'Title', 'Level', 'Created'],
        items.map((n: any) => [
          String(n.id),
          (n.title || n.body || '').slice(0, 50),
          n.level || '-',
          n.created_at || '-',
        ])
      );
      printPaginationFooter(res.data);
    });

  cmd
    .command('get <id>')
    .description('Get a notification by ID')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('GET', `/notifications/${id}`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const n = res.data;
      console.log(chalk.bold('\nNotification\n'));
      console.log(`  ${chalk.cyan('Title:')} ${n.title || '-'}`);
      console.log(`  ${chalk.cyan('Body:')} ${n.body || '-'}`);
      console.log(`  ${chalk.cyan('Level:')} ${n.level || '-'}`);
      console.log(`  ${chalk.cyan('Unread:')} ${n.is_unread ? 'Yes' : 'No'}`);
      console.log(`  ${chalk.cyan('Created:')} ${n.created_at || '-'}`);
      if (n.url) console.log(`  ${chalk.cyan('URL:')} ${n.url}`);
    });

  cmd
    .command('read <id>')
    .description('Mark a notification as read')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('POST', `/notifications/${id}/mark_read`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess('Notification marked as read.');
    });

  cmd
    .command('mark-unread <id>')
    .description('Mark a notification as unread')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('POST', `/notifications/${id}/mark_unread`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess('Notification marked as unread.');
    });

  cmd
    .command('mark-all-read')
    .description('Mark all notifications as read')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('POST', '/notifications/mark_all_read');
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess('All notifications marked as read.');
    });
}
