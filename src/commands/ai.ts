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

  // AI Tools using DynamicForm backend (authenticated endpoints)
  const toolsCmd = cmd
    .command('tools')
    .description('AI content generation tools (powered by SocialBu\'s DynamicForm backend)');

  toolsCmd
    .command('list')
    .description('List all available AI tools with their field definitions')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/ai/tools');
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      const tools: any[] = res.data.tools || [];
      if (tools.length === 0) {
        console.log(chalk.dim('No tools found.'));
        return;
      }

      console.log(chalk.bold(`\n🛠  Available AI Tools (${tools.length})\n`));
      for (const tool of tools) {
        const fieldNames = tool.fields?.map((f: any) => f.id).join(', ') || '';
        console.log(`  ${chalk.cyan(tool.slug.padEnd(40))} ${chalk.dim(tool.name)}`);
        if (fieldNames) console.log(`    ${chalk.dim('Fields: ' + fieldNames)}`);
      }
      console.log();
      console.log(chalk.dim('Run: socialbu ai tools run <slug> --fields \'{"field":"value"}\''));
    });

  toolsCmd
    .command('run <slug>')
    .description('Run any AI tool by its slug (use "list" to see available tools)')
    .requiredOption('--fields <json>', 'JSON object with the tool\'s required fields')
    .option('--json', 'Output as JSON')
    .action(async (slug: string, opts) => {
      let fields: any;
      try {
        fields = JSON.parse(opts.fields);
      } catch {
        console.error(chalk.red('--fields must be valid JSON, e.g. \'{"topic":"AI","tone":"professional"}\''));
        process.exit(1);
      }

      const res = await api('POST', `/ai/tools/${slug}`, fields);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\n✨ Tool Output\n'));
      if (res.data.text) {
        console.log(res.data.text);
      } else {
        console.log(JSON.stringify(res.data, null, 2));
      }
    });

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

      const res = await api('POST', '/ai/tools/tweet_generator', body);
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
        instagram: 'instagram_caption_generator',
        linkedin: 'linkedin_post_generator',
        facebook: 'facebook_post_generator',
        tiktok: 'tiktok_caption_generator',
      };

      const formSlug = platformMap[opts.platform.toLowerCase()];
      if (!formSlug) {
        console.error(chalk.red(`Unknown platform: ${opts.platform}. Use: instagram, linkedin, facebook, tiktok`));
        process.exit(1);
      }

      const body: any = {
        post_description: opts.description,
        tone: opts.tone,
        variant_count: opts.variants,
      };

      const res = await api('POST', `/ai/tools/${formSlug}`, body);
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

      const res = await api('POST', '/ai/tools/hashtag_cluster_generator', body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      console.log(chalk.bold('\n#️⃣ Generated Hashtags\n'));
      console.log(res.data.text || JSON.stringify(res.data));
    });
}
