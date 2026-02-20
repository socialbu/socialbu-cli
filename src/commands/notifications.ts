import { Command } from 'commander';
import { api } from '../api';
import { outputJson, handleApiError, printTable, printSuccess } from '../output';
import chalk from 'chalk';

export function registerNotificationsCommand(program: Command): void {
  const cmd = program
    .command('notifications')
    .alias('notif')
    .description('Manage notifications');

  cmd
    .command('list')
    .description('List all notifications')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/notifications');
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
      console.log(chalk.dim(`\nTotal: ${res.data.total ?? items.length}`));
    });

  cmd
    .command('unread')
    .description('List unread notifications')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/notifications/unread');
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
      const res = await api('PUT', `/notifications/${id}/read`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess('Notification marked as read.');
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
