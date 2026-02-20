import { Command } from 'commander';
import { loadConfig, saveConfig } from '../config';
import { printSuccess, printWarning, outputJson } from '../output';
import chalk from 'chalk';

export function registerConfigCommand(program: Command): void {
  const cmd = program
    .command('config')
    .description('Manage CLI configuration');

  cmd
    .command('set-key <key>')
    .description('Store your SocialBu API key')
    .action((key: string) => {
      const config = loadConfig();
      config.apiKey = key;
      saveConfig(config);
      printSuccess('API key saved to ~/.socialbu/config.json');
    });

  cmd
    .command('get-key')
    .description('Show stored API key')
    .action(() => {
      const config = loadConfig();
      if (config.apiKey) {
        const masked = config.apiKey.slice(0, 8) + '...' + config.apiKey.slice(-4);
        console.log(chalk.bold('API Key:'), masked);
      } else {
        printWarning('No API key configured. Run: socialbu config set-key <key>');
      }
    });

  cmd
    .command('set-url <url>')
    .description('Set custom API base URL (advanced)')
    .action((url: string) => {
      const config = loadConfig();
      config.baseUrl = url;
      saveConfig(config);
      printSuccess(`Base URL set to ${url}`);
    });

  cmd
    .command('reset')
    .description('Remove all stored configuration')
    .action(() => {
      saveConfig({});
      printSuccess('Configuration reset.');
    });

  cmd
    .command('show')
    .description('Show current configuration')
    .option('--json', 'Output as JSON')
    .action((opts) => {
      const config = loadConfig();
      if (opts.json) {
        outputJson(config);
      } else {
        console.log(chalk.bold('Configuration:'));
        console.log(`  API Key: ${config.apiKey ? config.apiKey.slice(0, 8) + '...' : chalk.dim('not set')}`);
        console.log(`  Base URL: ${config.baseUrl || chalk.dim('https://socialbu.com (default)')}`);
      }
    });
}
