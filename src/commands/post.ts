import { Command } from 'commander';
import { api, apiPaginated } from '../api';
import { outputJson, handleApiError, printError, printTable, printSuccess, printPaginationFooter } from '../output';
import { parsePage } from '../utils';
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
    .option('--comment <text>', 'First comment')
    .option('--post-as-story', 'Post as story')
    .option('--post-as-reel', 'Post as reel')
    .option('--share-reel-to-feed', 'Share reel to feed')
    .option('--video-title <text>', 'Video title')
    .option('--media-alt-text <texts...>', 'Alt text for media items')
    .option('--link <url>', 'Link for post')
    .option('--trim-link', 'Trim link from content')
    .option('--customize-link', 'Customize link preview')
    .option('--link-title <text>', 'Custom link title')
    .option('--link-description <text>', 'Custom link description')
    .option('--document-title <text>', 'Document title')
    .option('--reply-control <value>', 'Who can reply: everyone, accounts_you_follow, mentioned_only')
    .option('--threaded-replies <json>', 'JSON array of threaded replies')
    .option('--title <text>', 'Post title')
    .option('--is-spoiler', 'Mark as spoiler')
    .option('--is-nsfw', 'Mark as NSFW')
    .option('--flair-id <id>', 'Flair ID')
    .option('--mark-sensitive', 'Mark as sensitive')
    .option('--spoiler-text <text>', 'Spoiler/CW text')
    .option('--pin-title <text>', 'Pin title')
    .option('--board-name <text>', 'Board name')
    .option('--pin-link <url>', 'Pin destination URL')
    .option('--topic-type <type>', 'GBP post type: EVENT, OFFER')
    .option('--event-title <text>', 'Event title')
    .option('--event-start <datetime>', 'Event start')
    .option('--event-end <datetime>', 'Event end')
    .option('--offer-coupon <text>', 'Coupon code')
    .option('--offer-link <url>', 'Offer redeem link')
    .option('--offer-terms <text>', 'Offer terms')
    .option('--call-to-action <type>', 'CTA: BOOK, ORDER, SHOP, LEARN_MORE, SIGN_UP, CALL')
    .option('--call-to-action-url <url>', 'CTA URL')
    .option('--save-media-to-gallery', 'Save media to gallery')
    .option('--video-tags <tags>', 'Comma-separated tags')
    .option('--category-id <id>', 'YouTube category ID')
    .option('--privacy-status <status>', 'Privacy status')
    .option('--post-as-short', 'Post as YouTube short')
    .option('--made-for-kids', 'Made for kids')
    .option('--allow-stitch', 'Allow stitch')
    .option('--allow-duet', 'Allow duet')
    .option('--allow-comment', 'Allow comments')
    .option('--disclose-content', 'Disclose promotional content')
    .option('--branded-content', 'Branded/paid content')
    .option('--own-brand', 'Promoting own brand')
    .option('--auto-add-music', 'Auto add music')
    .option('--attachment-tokens <tokens...>', 'Upload tokens for existing attachments')
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
      if (opts.attachmentTokens) {
        body.existing_attachments = opts.attachmentTokens.map((uploadToken: string) => ({
          upload_token: uploadToken,
        }));
      }

      const options: any = {};
      if (opts.comment) options.comment = opts.comment;
      if (opts.postAsStory) options.post_as_story = true;
      if (opts.postAsReel) options.post_as_reel = true;
      if (opts.shareReelToFeed) options.share_reel_to_feed = true;
      if (opts.videoTitle) options.video_title = opts.videoTitle;
      if (opts.mediaAltText) options.media_alt_text = opts.mediaAltText;
      if (opts.link) options.link = opts.link;
      if (opts.trimLink) options.trim_link_from_content = true;
      if (opts.customizeLink) options.customize_link = true;
      if (opts.linkTitle) options.link_title = opts.linkTitle;
      if (opts.linkDescription) options.link_description = opts.linkDescription;
      if (opts.documentTitle) options.document_title = opts.documentTitle;
      if (opts.replyControl) options.reply_control = opts.replyControl;
      if (opts.title) options.title = opts.title;
      if (opts.isSpoiler) options.is_spoiler = true;
      if (opts.isNsfw) options.is_nsfw = true;
      if (opts.flairId) options.flair_id = Number(opts.flairId);
      if (opts.markSensitive) options.mark_sensitive = true;
      if (opts.spoilerText) options.spoiler = opts.spoilerText;
      if (opts.pinTitle) options.pin_title = opts.pinTitle;
      if (opts.boardName) options.board_name = opts.boardName;
      if (opts.pinLink) options.pin_link = opts.pinLink;
      if (opts.topicType) options.topic_type = opts.topicType;
      if (opts.eventTitle) options.event_title = opts.eventTitle;
      if (opts.eventStart) options.event_start = opts.eventStart;
      if (opts.eventEnd) options.event_end = opts.eventEnd;
      if (opts.offerCoupon) options.offer_coupon = opts.offerCoupon;
      if (opts.offerLink) options.offer_link = opts.offerLink;
      if (opts.offerTerms) options.offer_terms = opts.offerTerms;
      if (opts.callToAction) options.call_to_action = opts.callToAction;
      if (opts.callToActionUrl) options.call_to_action_url = opts.callToActionUrl;
      if (opts.saveMediaToGallery) options.save_media_to_gallery = true;
      if (opts.videoTags) options.video_tags = opts.videoTags;
      if (opts.categoryId) options.category_id = Number(opts.categoryId);
      if (opts.privacyStatus) options.privacy_status = opts.privacyStatus;
      if (opts.postAsShort) options.post_as_short = true;
      if (opts.madeForKids) options.made_for_kids = true;
      if (opts.allowStitch) options.allow_stitch = true;
      if (opts.allowDuet) options.allow_duet = true;
      if (opts.allowComment) options.allow_comment = true;
      if (opts.discloseContent) options.disclose_content = true;
      if (opts.brandedContent) options.branded_content = true;
      if (opts.ownBrand) options.own_brand = true;
      if (opts.autoAddMusic) options.auto_add_music = true;

      if (opts.threadedReplies) {
        try {
          const parsed = JSON.parse(opts.threadedReplies);
          if (!Array.isArray(parsed)) {
            printError('--threaded-replies must be a JSON array.');
            process.exit(1);
          }
          options.threaded_replies = parsed;
        } catch {
          printError('--threaded-replies must be valid JSON.');
          process.exit(1);
        }
      }

      if (Object.keys(options).length > 0) {
        body.options = options;
      }

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
    .option('--page <n>', 'Page number', parsePage, 1)
    .option('--all', 'Fetch all pages')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const query: any = { type: opts.type };
      if (!opts.all) query.page = opts.page;

      const res = opts.all
        ? await apiPaginated('GET', '/posts', undefined, query)
        : await api('GET', '/posts', undefined, query);
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
      printPaginationFooter(res.data);
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
    .option('--attachment-tokens <tokens...>', 'Upload tokens for attachments')
    .option('--video-title <text>', 'Video title (options.video_title)')
    .option('--trim-link', 'Trim link from content (options.trim_link_from_content)')
    .option('--link <url>', 'Link for post (options.link)')
    .option('--media-alt-text <texts...>', 'Alt text for media items (options.media_alt_text)')
    .option('--comment <text>', 'First comment (options.comment)')
    .option('--json', 'Output as JSON')
    .action(async (id: string, opts) => {
      const body: any = {};
      if (opts.content) body.content = opts.content;
      if (opts.publishAt) body.publish_at = opts.publishAt;
      if (opts.accounts) body.accounts = opts.accounts.map(Number);
      if (opts.team) body.team_id = Number(opts.team);
      if (opts.attachmentTokens) {
        body.existing_attachments = opts.attachmentTokens.map((uploadToken: string) => ({
          upload_token: uploadToken,
        }));
      }

      // PATCH endpoint supports limited options per spec
      const options: any = {};
      if (opts.videoTitle) options.video_title = opts.videoTitle;
      if (opts.trimLink) options.trim_link_from_content = true;
      if (opts.link) options.link = opts.link;
      if (opts.mediaAltText) options.media_alt_text = opts.mediaAltText;
      if (opts.comment) options.comment = opts.comment;

      if (Object.keys(options).length > 0) {
        body.options = options;
      }

      const res = await api('PATCH', `/posts/${id}`, body);
      handleApiError(res);

      if (opts.json) {
        outputJson(res.data);
        return;
      }

      printSuccess(res.data.message || 'Post updated.');
    });
}
