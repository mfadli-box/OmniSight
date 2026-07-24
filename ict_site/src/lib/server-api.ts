const BE_POOL = process.env.BE_POOL ?? "http://localhost:36665";

export class ServerApiError extends Error {
  code: string;
  status: number;
  requestId: string;

  constructor(code: string, message: string, status: number, requestId: string) {
    super(message);
    this.name = "ServerApiError";
    this.code = code;
    this.status = status;
    this.requestId = requestId;
  }
}

export async function serverApi<T>(
  path: string,
  options?: RequestInit,
): Promise<T> {
  const url = `${BE_POOL}/rest/pages${path}`;
  const res = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  const body = await res.json().catch(() => null);
  const requestId = res.headers.get("X-Request-ID") ?? "";

  if (!res.ok) {
    throw new ServerApiError(
      body?.code ?? "INTERNAL_ERROR",
      body?.error ?? "Unknown error",
      res.status,
      requestId,
    );
  }

  return body as T;
}
