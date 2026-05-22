// Minimal fetch wrapper around the DevForge HTTP API.
// In Phase D this file will be regenerated from the OpenAPI doc.

export interface Diagnostic {
  code: string;
  message: string;
  severity: 0 | 1 | 2;
}

const BASE = "";

export async function callOp<I, O>(
  tool: string,
  op: string,
  body: I,
): Promise<O> {
  const r = await fetch(`${BASE}/api/v1/${tool}/${op}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!r.ok) {
    throw new Error(`HTTP ${r.status}: ${await r.text()}`);
  }
  return (await r.json()) as O;
}

export interface OpInfo {
  tool: string;
  op: string;
  name: string;
  path: string;
  description: string;
}

export async function listOperations(): Promise<OpInfo[]> {
  const r = await fetch(`${BASE}/api/v1/operations`);
  if (!r.ok) throw new Error(`HTTP ${r.status}`);
  return (await r.json()) as OpInfo[];
}

// Tool-specific wrappers.

export interface UuidGenerateRequest {
  version?: 4 | 7;
  count?: number;
  format?: "std" | "compact" | "urn";
}
export interface UuidGenerateResponse {
  values: string[];
  diagnostics?: Diagnostic[];
}
export const uuidGenerate = (r: UuidGenerateRequest) =>
  callOp<UuidGenerateRequest, UuidGenerateResponse>("uuid", "generate", r);

export interface JsonFormatResponse {
  output: string;
  diagnostics?: Diagnostic[];
}
export const jsonFormat = (input: string, indent = 2) =>
  callOp<{ input: string; indent: number }, JsonFormatResponse>("json", "format", {
    input,
    indent,
  });

export interface DiffHunk {
  path: string;
  op: "add" | "remove" | "change";
  left?: unknown;
  right?: unknown;
}
export interface DiffResult {
  mode: string;
  hunks: DiffHunk[];
  summary: { adds: number; removes: number; changes: number };
}
export const diffCompare = (left: string, right: string, mode = "auto") =>
  callOp<{ left: string; right: string; mode: string }, DiffResult>("diff", "compare", {
    left,
    right,
    mode,
  });

export interface RegexMatch {
  start: number;
  end: number;
  value: string;
  groups?: { name?: string; start: number; end: number; value: string }[];
}
export const regexTest = (pattern: string, input: string, flags = "") =>
  callOp<{ pattern: string; input: string; flags: string }, { matches: RegexMatch[]; diagnostics?: Diagnostic[] }>(
    "regex",
    "test",
    { pattern, input, flags },
  );

export const cronParse = (expression: string, flavor = "unix") =>
  callOp<{ expression: string; flavor: string }, { description: string; fields: { name: string; value: string }[]; diagnostics?: Diagnostic[] }>(
    "cron",
    "parse",
    { expression, flavor },
  );

export const cronNext = (expression: string, n = 5, tz = "UTC") =>
  callOp<{ expression: string; n: number; tz: string }, { runs: string[]; diagnostics?: Diagnostic[] }>(
    "cron",
    "next",
    { expression, n, tz },
  );

export interface JwtDecodeResponse {
  header: Record<string, unknown>;
  payload: Record<string, unknown>;
  raw: { header: string; payload: string; signature: string };
  diagnostics?: Diagnostic[];
}
export const jwtDecode = (token: string) =>
  callOp<{ token: string }, JwtDecodeResponse>("jwt", "decode", { token });

export const tzConvert = (time: string, fromTZ: string, toTZ: string) =>
  callOp<{ time: string; fromTZ: string; toTZ: string }, { original: string; converted: string; dstNote?: string }>(
    "tz",
    "convert",
    { time, fromTZ, toTZ },
  );
