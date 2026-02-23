#!/usr/bin/env node

import { program } from 'commander';
import { registerConfigCommand } from './commands/config';
import { registerWhoamiCommand } from './commands/whoami';
import { registerAccountCommand } from './commands/account';
import { registerPostCommand } from './commands/post';
import { registerAiCommand } from './commands/ai';
import { registerTeamCommand } from './commands/team';
import { registerAnalyticsCommand } from './commands/analytics';
import { registerNotificationsCommand } from './commands/notifications';
import { registerCurateCommand } from './commands/curate';
import { registerMediaCommand } from './commands/media';

program
  .name('socialbu')
  .description('CLI for managing SocialBu — posts, accounts, analytics, and more')
  .version('0.1.0');

registerConfigCommand(program);
registerWhoamiCommand(program);
registerAccountCommand(program);
registerPostCommand(program);
registerAiCommand(program);
registerTeamCommand(program);
registerAnalyticsCommand(program);
registerNotificationsCommand(program);
registerCurateCommand(program);
registerMediaCommand(program);

program.parse();
