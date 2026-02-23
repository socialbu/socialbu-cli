import { Command } from 'commander';
import { api } from '../api';
import { outputJson, handleApiError, printSuccess, printTable } from '../output';
import chalk from 'chalk';
import fs from 'fs';
import path from 'path';
import fetch from 'node-fetch';

export function registerMediaCommand(program: Command): void {
  const cmd = program
    .command('media')
    .description('Upload media files for posts');

  cmd
    .command('upload <filePath>')
    .description('Upload a media file (image, video, or document)')
    .option('--json', 'Output as JSON')
    .action(async (filePath: string, opts) => {
      // Check if file exists
      if (!fs.existsSync(filePath)) {
        console.error(chalk.red(`File not found: ${filePath}`));
        process.exit(1);
      }

      const fileName = path.basename(filePath);
      const fileExt = path.extname(filePath).toLowerCase();
      
      // Determine MIME type based on extension
      const mimeTypes: Record<string, string> = {
        '.jpg': 'image/jpeg',
        '.jpeg': 'image/jpeg',
        '.png': 'image/png',
        '.gif': 'image/gif',
        '.webp': 'image/webp',
        '.mp4': 'video/mp4',
        '.mov': 'video/quicktime',
        '.pdf': 'application/pdf',
        '.doc': 'application/msword',
        '.docx': 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
        '.ppt': 'application/vnd.ms-powerpoint',
        '.pptx': 'application/vnd.openxmlformats-officedocument.presentationml.presentation',
      };
      
      const mimeType = mimeTypes[fileExt];
      if (!mimeType) {
        console.error(chalk.red(`Unsupported file type: ${fileExt}`));
        console.log(chalk.dim('Supported: jpg, png, gif, webp, mp4, mov, pdf, doc, docx, ppt, pptx'));
        process.exit(1);
      }

      // Step 1: Get signed URL
      const initRes = await api('POST', '/upload_media', {
        name: fileName,
        mime_type: mimeType,
      });
      handleApiError(initRes);

      const { signed_url, key, url } = initRes.data;

      // Step 2: Upload file to signed URL
      const fileBuffer = fs.readFileSync(filePath);
      const uploadRes = await fetch(signed_url, {
        method: 'PUT',
        headers: {
          'Content-Type': mimeType,
          'Content-Length': String(fileBuffer.length),
          'x-amz-acl': 'private',
        },
        body: fileBuffer,
      });

      if (!uploadRes.ok) {
        console.error(chalk.red(`Upload failed: ${uploadRes.status} ${uploadRes.statusText}`));
        process.exit(1);
      }

      // Step 3: Check status and get upload token
      const statusRes = await api('GET', '/upload_media/status', undefined, { key });
      handleApiError(statusRes);

      if (opts.json) {
        outputJson({
          ...statusRes.data,
          file_url: url,
        });
        return;
      }

      if (statusRes.data.upload_token) {
        printSuccess('Media uploaded successfully!');
        console.log(chalk.bold('\n📎 Media Details\n'));
        printTable(
          ['Property', 'Value'],
          [
            ['File Name', fileName],
            ['MIME Type', mimeType],
            ['Upload Token', statusRes.data.upload_token],
            ['File URL', url],
          ]
        );
        console.log();
        console.log(chalk.dim('Use the upload token with --attachment-tokens when creating a post.'));
      } else {
        console.log(chalk.yellow('Upload pending. Check status with:'));
        console.log(chalk.dim(`  socialbu media status --key ${key}`));
      }
    });

  cmd
    .command('status')
    .description('Check media upload status')
    .requiredOption('--key <key>', 'Media key from upload response')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const res = await api('GET', '/upload_media/status', undefined, { key: opts.key });
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      if (res.data.upload_token) {
        printSuccess('Media ready!');
        console.log(chalk.bold('\n📎 Upload Token:\n'));
        console.log(res.data.upload_token);
      } else {
        console.log(chalk.yellow('Media still processing...'));
        console.log(chalk.dim(`Success: ${res.data.success}`));
        console.log(chalk.dim(`Message: ${res.data.message}`));
      }
    });
}
