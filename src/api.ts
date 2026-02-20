import fetch from 'node-fetch';
import { getApiKey, getBaseUrl } from './config';

export interface ApiResponse {
  ok: boolean;
  status: number;
  data: any;
}

export type QueryValue = string | number | boolean;
export type QueryParams = Record<string, QueryValue | QueryValue[] | null | undefined>;

const MAX_RETRIES = 3;
const RETRY_BASE_DELAY_MS = 500;

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function buildUrl(endpoint: string, query?: QueryParams): string {
  const baseUrl = getBaseUrl();
  let url = `${baseUrl}/api${endpoint}`;

  if (query) {
    const params = new URLSearchParams();
    for (const [key, value] of Object.entries(query)) {
      if (value === undefined || value === null || value === '') continue;

      if (Array.isArray(value)) {
        value.forEach((v) => params.append(`${key}[]`, String(v)));
      } else {
        params.append(key, String(value));
      }
    }
    const qs = params.toString();
    if (qs) url += `?${qs}`;
  }

  return url;
}

function normalizeErrorData(status: number, data: any): any {
  if (data && typeof data === 'object') {
    const base =
      data.message || data.error || data.detail || data.title || `API request failed with status ${status}`;
    const details = data.errors || data.details;
    if (details) {
      const detailsText =
        typeof details === 'string' ? details : JSON.stringify(details);
      return {
        ...data,
        message: `${base} (${detailsText})`,
      };
    }
    return {
      ...data,
      message: base,
    };
  }

  if (typeof data === 'string' && data.trim()) {
    return { message: data };
  }

  return { message: `API request failed with status ${status}` };
}

export async function api(
  method: string,
  endpoint: string,
  body?: any,
  query?: QueryParams
): Promise<ApiResponse> {
  const apiKey = getApiKey();
  const url = buildUrl(endpoint, query);
  const headers: Record<string, string> = {
    Authorization: `Bearer ${apiKey}`,
    Accept: 'application/json',
  };

  const requestOptions: any = { method, headers };

  if (body) {
    headers['Content-Type'] = 'application/json';
    requestOptions.body = JSON.stringify(body);
  }

  for (let attempt = 0; attempt <= MAX_RETRIES; attempt++) {
    try {
      const res = await fetch(url, requestOptions);

      let data: any = null;
      const raw = await res.text();
      if (raw) {
        try {
          data = JSON.parse(raw);
        } catch {
          data = raw;
        }
      }

      if (res.status === 429 && attempt < MAX_RETRIES) {
        const retryAfterSeconds = Number(res.headers.get('retry-after'));
        const backoff = Number.isFinite(retryAfterSeconds) && retryAfterSeconds > 0
          ? retryAfterSeconds * 1000
          : RETRY_BASE_DELAY_MS * 2 ** attempt;
        await sleep(backoff);
        continue;
      }

      if (!res.ok) {
        return {
          ok: false,
          status: res.status,
          data: normalizeErrorData(res.status, data),
        };
      }

      return { ok: true, status: res.status, data };
    } catch (error: any) {
      if (attempt < MAX_RETRIES) {
        const backoff = RETRY_BASE_DELAY_MS * 2 ** attempt;
        await sleep(backoff);
        continue;
      }

      return {
        ok: false,
        status: 0,
        data: {
          message: `Network error while contacting API: ${error?.message || 'Unknown error'}`,
        },
      };
    }
  }

  return {
    ok: false,
    status: 0,
    data: { message: 'Unexpected API error.' },
  };
}

export async function apiPaginated(
  method: string,
  endpoint: string,
  body?: any,
  query?: QueryParams
): Promise<ApiResponse> {
  const baseQuery: QueryParams = { ...(query || {}) };
  delete baseQuery.page;

  const firstRes = await api(method, endpoint, body, { ...baseQuery, page: 1 });
  if (!firstRes.ok) return firstRes;

  const firstData = (firstRes.data && typeof firstRes.data === 'object') ? firstRes.data : {};
  const allItems: any[] = Array.isArray(firstData.items) ? [...firstData.items] : [];

  let currentPage = Number(firstData.currentPage ?? 1);
  let lastPage = Number(firstData.lastPage ?? currentPage);
  let nextPage = firstData.nextPage as number | null | undefined;

  while (nextPage !== null && nextPage !== undefined && currentPage < lastPage) {
    const pageRes = await api(method, endpoint, body, { ...baseQuery, page: nextPage });
    if (!pageRes.ok) return pageRes;

    const pageData = (pageRes.data && typeof pageRes.data === 'object') ? pageRes.data : {};
    if (Array.isArray(pageData.items)) {
      allItems.push(...pageData.items);
    }

    currentPage = Number(pageData.currentPage ?? nextPage);
    lastPage = Number(pageData.lastPage ?? lastPage);
    nextPage = pageData.nextPage as number | null | undefined;
  }

  return {
    ok: true,
    status: firstRes.status,
    data: {
      ...firstData,
      items: allItems,
      currentPage: 1,
      lastPage: Number.isFinite(lastPage) ? lastPage : 1,
      nextPage: null,
      total:
        typeof firstData.total === 'number' ? firstData.total : allItems.length,
    },
  };
}
