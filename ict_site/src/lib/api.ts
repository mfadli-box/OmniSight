import { BE_POOL } from "./backend";

type FetchOptions = {
  method?: string;
  body?: unknown;
  token?: string;
  params?: Record<string, string>;
};

export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export async function apiFetch<T = unknown>(
  path: string,
  options: FetchOptions = {},
): Promise<T> {
  const { method = "GET", body, token, params } = options;

  let url = `${BE_POOL}/rest/pages${path}`;
  if (params) {
    const qs = new URLSearchParams(params).toString();
    if (qs) url += `?${qs}`;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(url, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
    cache: "no-store",
  });

  const data = await res.json().catch(() => null);

  if (!res.ok) {
    const message =
      data?.error || data?.message || `Request failed (${res.status})`;
    throw new ApiError(message, res.status);
  }

  return data as T;
}
