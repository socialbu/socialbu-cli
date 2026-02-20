import { Command } from 'commander';
import { api } from '../api';
import { outputJson, handleApiError, printSuccess } from '../output';
import chalk from 'chalk';

export function registerAiCommand(program: Command): void {
  const cmd = program
    .command('ai')
    .description('Generate AI-powered content');

  cmd
    .command('generate')
    .description('Generate AI content for a topic')
    .requiredOption('--topic <topic>', 'Topic to generate content about (e.g. "Mother\'s Day")')
    .option('--type <type>', 'Content type: generic, instagram_caption, linkedin_post, tweet', 'generic')
    .option('--account <id>', 'Social account ID for context')
    .option('--team <id>', 'Team ID')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const body: any = {
        topic: opts.topic,
        type: opts.type,
      };
      if (opts.account) body.account_id = Number(opts.account);
      if (opts.team) body.team_id = Number(opts.team);

      const res = await api('POST', '/generated_content', body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\n🤖 Generated Content\n'));
      console.log(res.data.content || JSON.stringify(res.data));
      console.log();
      if (res.data.id) console.log(chalk.dim(`Content ID: ${res.data.id}`));
    });

  cmd
    .command('autocomplete')
    .description('Autocomplete existing post content')
    .requiredOption('--account <id>', 'Social account ID')
    .option('--content <text>', 'Partial content to autocomplete')
    .option('--team <id>', 'Team ID')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const body: any = {
        account_id: Number(opts.account),
      };
      if (opts.content) body.content = opts.content;
      if (opts.team) body.team_id = Number(opts.team);

      const res = await api('POST', '/generated_content/autocomplete_post', body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\n✨ Autocompleted Content\n'));
      console.log(res.data.content || JSON.stringify(res.data));
    });

  cmd
    .command('from-post <postId>')
    .description('Generate similar content based on an existing post')
    .option('--json', 'Output as JSON')
    .action(async (postId: string, opts) => {
      const res = await api('POST', '/generated_content/generate_by_post', {
        post_id: Number(postId),
      });
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\n🤖 Generated Content (from post)\n'));
      console.log(res.data.content || JSON.stringify(res.data));
      console.log();
      if (res.data.id) console.log(chalk.dim(`Content ID: ${res.data.id}`));
    });
}
