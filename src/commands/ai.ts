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

  // AI Tools using DynamicForm backend
  const toolsCmd = cmd
    .command('tools')
    .description('AI content generation tools');

  toolsCmd
    .command('tweet')
    .description('Generate tweet/post content using AI')
    .requiredOption('--description <text>', 'Brief description of the post content')
    .option('--tone <tone>', 'Tone: funny, aesthetic, cool, professional', 'professional')
    .option('--variants <n>', 'Number of variants (1-5)', '3')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const body: any = {
        post_description: opts.description,
        tone: opts.tone,
        variant_count: opts.variants,
      };

      const res = await api('POST', '/generate/forms/TweetGenerator', body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\n🐦 Generated Tweets\n'));
      const posts = res.data.text?.split('---END_POST---').filter((p: string) => p.trim());
      if (posts && posts.length > 0) {
        posts.forEach((post: string, i: number) => {
          console.log(chalk.cyan(`${i + 1}.`) + ' ' + post.trim());
          console.log();
        });
      } else {
        console.log(res.data.text || JSON.stringify(res.data));
      }
    });

  toolsCmd
    .command('caption')
    .description('Generate social media captions using AI')
    .requiredOption('--platform <platform>', 'Platform: instagram, linkedin, facebook, tiktok')
    .requiredOption('--description <text>', 'Brief description of the content')
    .option('--tone <tone>', 'Tone: funny, aesthetic, cool, professional', 'professional')
    .option('--variants <n>', 'Number of variants (1-5)', '2')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const platformMap: Record<string, string> = {
        instagram: 'InstagramCaptionGenerator',
        linkedin: 'LinkedinPostGenerator',
        facebook: 'FacebookPostGenerator',
        tiktok: 'TiktokCaptionGenerator',
      };

      const formClass = platformMap[opts.platform.toLowerCase()];
      if (!formClass) {
        console.error(chalk.red(`Unknown platform: ${opts.platform}. Use: instagram, linkedin, facebook, tiktok`));
        process.exit(1);
      }

      const body: any = {
        post_description: opts.description,
        tone: opts.tone,
        variant_count: opts.variants,
      };

      const res = await api('POST', `/generate/forms/${formClass}`, body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold(`\n✨ Generated ${opts.platform.charAt(0).toUpperCase() + opts.platform.slice(1)} Caption\n`));
      console.log(res.data.text || JSON.stringify(res.data));
    });

  toolsCmd
    .command('hashtags')
    .description('Generate hashtag clusters using AI')
    .requiredOption('--topic <text>', 'Topic or keyword to generate hashtags for')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const body: any = {
        topic: opts.topic,
      };

      const res = await api('POST', '/generate/forms/HashtagClusterGenerator', body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\n#️⃣ Generated Hashtags\n'));
      console.log(res.data.text || JSON.stringify(res.data));
    });
}
