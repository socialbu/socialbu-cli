import fetch from 'node-fetch';
import { getApiKey, getBaseUrl } from './config';

export interface ApiResponse {
  ok: boolean;
  status: number;
  data: any;
}

export async function api(
  method: string,
  endpoint: string,
  body?: any,
  query?: Record<string, string | string[]>
): Promise<ApiResponse> {
  const baseUrl = getBaseUrl();
  const apiKey = getApiKey();

  let url = `${baseUrl}/api${endpoint}`;

  if (query) {
    const params = new URLSearchParams();
    for (const [key, value] of Object.entries(query)) {
      if (Array.isArray(value)) {
        value.forEach((v) => params.append(`${key}[]`, v));
      } else if (value !== undefined && value !== '') {
        params.append(key, value);
      }
    }
    const qs = params.toString();
    if (qs) url += `?${qs}`;
  }

  const headers: Record<string, string> = {
    Authorization: `Bearer ${apiKey}`,
    Accept: 'application/json',
  };

  const opts: any = { method, headers };

  if (body) {
    headers['Content-Type'] = 'application/json';
    opts.body = JSON.stringify(body);
  }

  const res = await fetch(url, opts);
  let data: any;
  try {
    data = await res.json();
  } catch {
    data = null;
  }

  return { ok: res.ok, status: res.status, data };
}
