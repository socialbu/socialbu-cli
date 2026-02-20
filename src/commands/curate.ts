import { Command } from 'commander';
import { api, apiPaginated } from '../api';
import { outputJson, handleApiError, printTable } from '../output';
import chalk from 'chalk';

function parsePage(value: string): number {
  const parsed = Number.parseInt(value, 10);
  if (!Number.isInteger(parsed) || parsed < 1) {
    throw new Error('Page must be a positive integer.');
  }
  return parsed;
}

export function registerCurateCommand(program: Command): void {
  const cmd = program
    .command('curate')
    .description('Browse curated content topics and items');

  cmd
    .command('topics')
    .description('List curation topics')
    .option('--query <q>', 'Search topic name')
    .option('--page <n>', 'Page number', parsePage, 1)
    .option('--all', 'Fetch all pages')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = {};
      if (opts.query) query.q = opts.query;
      if (!opts.all) query.page = opts.page;

      const res = opts.all
        ? await apiPaginated('GET', '/curation/topics', undefined, query)
        : await api('GET', '/curation/topics', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const items = res.data.items || [];
      if (items.length === 0) {
        console.log(chalk.dim('No topics found.'));
        return;
      }

      printTable(
        ['ID', 'Title', 'Score'],
        items.map((t: any) => [
          String(t.id),
          t.title || '-',
          String(t.score ?? '-'),
        ])
      );
      console.log(chalk.dim(`\nPage ${res.data.currentPage ?? 1} of ${res.data.lastPage ?? 1} | Total: ${res.data.total ?? items.length}`));
    });

  cmd
    .command('items')
    .description('List curated content items')
    .option('--feed <id>', 'Filter by feed ID')
    .option('--search <q>', 'Search string')
    .option('--from <date>', 'Start date (YYYY-MM-DD)')
    .option('--to <date>', 'End date (YYYY-MM-DD)')
    .option('--page <n>', 'Page number', parsePage, 1)
    .option('--all', 'Fetch all pages')
    .option('--per-page <n>', 'Items per page')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = {};
      if (opts.feed) query.feed_id = opts.feed;
      if (opts.search) query.search = opts.search;
      if (opts.from) query.from = opts.from;
      if (opts.to) query.to = opts.to;
      if (!opts.all) query.page = opts.page;
      if (opts.perPage) query.per_page = opts.perPage;

      const res = opts.all
        ? await apiPaginated('GET', '/curation/items', undefined, query)
        : await api('GET', '/curation/items', undefined, query);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const items = res.data.items || [];
      if (items.length === 0) {
        console.log(chalk.dim('No items found.'));
        return;
      }

      printTable(
        ['ID', 'Title', 'Authors', 'Published'],
        items.map((i: any) => [
          String(i.id),
          (i.title || '').slice(0, 50),
          i.authors || '-',
          i.published_at || '-',
        ])
      );
      console.log(chalk.dim(`\nPage ${res.data.currentPage ?? 1} of ${res.data.lastPage ?? 1} | Total: ${res.data.total ?? items.length}`));
    });

  cmd
    .command('get <id>')
    .description('Get a curated item by ID')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const res = await api('GET', `/curation/items/${id}`);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const item = res.data;
      console.log(chalk.bold('\n📰 Curated Item\n'));
      console.log(`  ${chalk.cyan('Title:')} ${item.title || '-'}`);
      console.log(`  ${chalk.cyan('Authors:')} ${item.authors || '-'}`);
      console.log(`  ${chalk.cyan('Published:')} ${item.published_at || '-'}`);
      console.log(`  ${chalk.cyan('Link:')} ${item.link || '-'}`);
      if (item.description) {
        console.log(`  ${chalk.cyan('Description:')}`);
        console.log(`    ${item.description.slice(0, 200)}`);
      }
    });
}
