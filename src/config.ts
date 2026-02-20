import * as fs from 'fs';
import * as path from 'path';
import * as os from 'os';

const CONFIG_DIR = path.join(os.homedir(), '.socialbu');
const CONFIG_FILE = path.join(CONFIG_DIR, 'config.json');

export interface SocialBuConfig {
  apiKey?: string;
  baseUrl?: string;
}

function ensureConfigDir(): void {
  if (!fs.existsSync(CONFIG_DIR)) {
    fs.mkdirSync(CONFIG_DIR, { recursive: true });
  }
}

export function loadConfig(): SocialBuConfig {
  try {
    if (fs.existsSync(CONFIG_FILE)) {
      return JSON.parse(fs.readFileSync(CONFIG_FILE, 'utf-8'));
    }
  } catch {
    // ignore parse errors
  }
  return {};
}

export function saveConfig(config: SocialBuConfig): void {
  ensureConfigDir();
  fs.writeFileSync(CONFIG_FILE, JSON.stringify(config, null, 2), 'utf-8');
}

export function getApiKey(): string {
  const config = loadConfig();
  if (!config.apiKey) {
    console.error('No API key configured. Run: socialbu config set-key <key>');
    process.exit(1);
  }
  return config.apiKey;
}

export function getBaseUrl(): string {
  const config = loadConfig();
  return config.baseUrl || 'https://socialbu.com';
}
