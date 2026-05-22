// DevForge HTTP API client.
// Generic callOp + typed wrappers for every backend operation.

export type Severity = 0 | 1 | 2;

export interface Diagnostic {
  code: string;
  message: string;
  severity: Severity;
}

const BASE = "";

export class ApiError extends Error {
  status: number;
  body: string;
  constructor(status: number, body: string) {
    super(`HTTP ${status}: ${body}`);
    this.status = status;
    this.body = body;
  }
}

export async function callOp<I, O>(
  tool: string,
  op: string,
  body: I,
  signal?: AbortSignal,
): Promise<O> {
  const r = await fetch(`${BASE}/api/v1/${tool}/${op}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
    signal,
  });
  const text = await r.text();
  if (!r.ok) throw new ApiError(r.status, text);
  return text ? (JSON.parse(text) as O) : (undefined as O);
}

export interface OpInfo {
  tool: string;
  op: string;
  name: string;
  path: string;
  description: string;
  inputSchema?: Record<string, unknown>;
}

let opCache: OpInfo[] | null = null;
export async function listOperations(force = false): Promise<OpInfo[]> {
  if (opCache && !force) return opCache;
  const r = await fetch(`${BASE}/api/v1/operations`);
  if (!r.ok) throw new ApiError(r.status, await r.text());
  opCache = (await r.json()) as OpInfo[];
  return opCache;
}

// ───────────────────────────── Generators ─────────────────────────────

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

export interface UuidHashRequest {
  input: string;
  algos?: ("md5" | "sha1" | "sha256" | "sha512")[];
  encoding?: "hex" | "base64";
}
export interface UuidHashResponse {
  digests: Record<string, string>;
  diagnostics?: Diagnostic[];
}
export const uuidHash = (r: UuidHashRequest) =>
  callOp<UuidHashRequest, UuidHashResponse>("uuid", "hash", r);

export interface FakerField {
  name: string;
  kind: string;
  locale?: string;
  params?: Record<string, unknown>;
}
export interface FakerGenerateRequest {
  spec: { fields: FakerField[] };
  count?: number;
  format?: "json" | "csv" | "sql";
  table?: string;
  seed?: number;
}
export interface FakerGenerateResponse {
  output: string;
  rows: number;
  diagnostics?: Diagnostic[];
}
export const fakerGenerate = (r: FakerGenerateRequest) =>
  callOp<FakerGenerateRequest, FakerGenerateResponse>("faker", "generate", r);
export const fakerKinds = () =>
  callOp<Record<string, never>, { kinds: { name: string; description: string; params?: string[] }[] }>(
    "faker",
    "kinds",
    {},
  );

export const idUlid = (count = 1, lowercase = false) =>
  callOp<{ count: number; lowercase: boolean }, { values: string[] }>(
    "id",
    "ulid",
    { count, lowercase },
  );
export const idSlug = (input: string, maxLen = 60, locale = "en") =>
  callOp<{ input: string; maxLen: number; locale: string }, { output: string }>(
    "id",
    "slug",
    { input, maxLen, locale },
  );

export interface TotpGenerateRequest {
  secret: string;
  encoding?: "raw" | "hex" | "base32";
  algorithm?: "sha1" | "sha256" | "sha512";
  digits?: 6 | 8;
  period?: number;
  time?: number;
}
export interface TotpGenerateResponse {
  code: string;
  counter: number;
  remainingMs: number;
  periodSec: number;
  algorithm: string;
  digits: number;
  diagnostics?: Diagnostic[];
}
export const totpGenerate = (r: TotpGenerateRequest) =>
  callOp<TotpGenerateRequest, TotpGenerateResponse>("totp", "generate", r);
export const totpVerify = (r: TotpGenerateRequest & { code: string; skew?: number }) =>
  callOp<typeof r, { valid: boolean; diagnostics?: Diagnostic[] }>("totp", "verify", r);

// Crypto
export const cryptoAesEncrypt = (input: {
  plaintext: string;
  passphrase?: string;
  key?: string;
  keySize?: 128 | 192 | 256;
  iterations?: number;
}) => callOp<typeof input, { ciphertext: string; diagnostics?: Diagnostic[] }>("crypto", "aes_encrypt", input);
export const cryptoAesDecrypt = (input: {
  ciphertext: string;
  passphrase?: string;
  key?: string;
}) => callOp<typeof input, { plaintext: string; diagnostics?: Diagnostic[] }>("crypto", "aes_decrypt", input);
export const cryptoRsaKeygen = (bits: 2048 | 3072 | 4096 = 3072) =>
  callOp<{ bits: number }, { privatePem: string; publicPem: string }>("crypto", "rsa_keygen", { bits });
export const cryptoHmac = (input: {
  input: string;
  key: string;
  keyEncoding?: "raw" | "hex" | "base64";
  algorithm?: "sha256" | "sha384" | "sha512";
}) => callOp<typeof input, { mac: string }>("crypto", "hmac", input);
export const cryptoPasswordHash = (input: {
  password: string;
  algorithm?: "bcrypt" | "argon2id";
  cost?: number;
}) => callOp<typeof input, { hash: string }>("crypto", "password_hash", input);
export const cryptoPasswordStrength = (password: string) =>
  callOp<{ password: string }, {
    score: 0 | 1 | 2 | 3 | 4;
    crackTime: string;
    feedback: { warning?: string; suggestions: string[] };
  }>("crypto", "password_strength", { password });

// ───────────────────────────── Formatters ─────────────────────────────

export interface JsonFormatResponse {
  output: string;
  diagnostics?: Diagnostic[];
}
export const jsonFormat = (input: string, indent = 2, sortKeys = false, trailingNewline = false) =>
  callOp<{ input: string; indent: number; sortKeys: boolean; trailingNewline: boolean }, JsonFormatResponse>(
    "json",
    "format",
    { input, indent, sortKeys, trailingNewline },
  );
export const jsonValidate = (input: string, schema?: string) =>
  callOp<{ input: string; schema?: string }, { valid: boolean; diagnostics?: Diagnostic[] }>(
    "json",
    "validate",
    { input, schema },
  );

export const yamlFormat = (input: string, indent = 2) =>
  callOp<{ input: string; indent: number }, JsonFormatResponse>("yaml", "format", { input, indent });
export const yamlValidate = (input: string) =>
  callOp<{ input: string }, { valid: boolean; diagnostics?: Diagnostic[] }>("yaml", "validate", { input });
export const yamlConvert = (input: string, to: "yaml" | "json", indent = 2) =>
  callOp<{ input: string; to: string; indent: number }, JsonFormatResponse>(
    "yaml",
    "convert",
    { input, to, indent },
  );

export const csvFormat = (input: string, delimiter = ",", header = true, alignColumns = false) =>
  callOp<typeof input extends string ? { input: string; delimiter: string; header: boolean; alignColumns: boolean } : never, JsonFormatResponse>(
    "csv",
    "format",
    { input, delimiter, header, alignColumns } as never,
  );
export const csvValidate = (input: string, delimiter = ",", expectedColumns?: string[], strict = false) =>
  callOp<{ input: string; delimiter: string; expectedColumns?: string[]; strict: boolean }, { valid: boolean; diagnostics?: Diagnostic[] }>(
    "csv",
    "validate",
    { input, delimiter, expectedColumns, strict },
  );

export const sqlFormat = (input: string, uppercase = true, indentWidth = 2) =>
  callOp<{ input: string; uppercase: boolean; indentWidth: number }, JsonFormatResponse>(
    "sql",
    "format",
    { input, uppercase, indentWidth },
  );
export const sqlValidate = (input: string) =>
  callOp<{ input: string }, { valid: boolean; diagnostics?: Diagnostic[] }>("sql", "validate", { input });

export const codeFmtGo = (input: string) =>
  callOp<{ input: string }, JsonFormatResponse>("code", "fmt_go", { input });
export const codeFmtXml = (input: string, indent = 2) =>
  callOp<{ input: string; indent: number }, JsonFormatResponse>("code", "fmt_xml", { input, indent });
export const codeFmtHtml = (input: string) =>
  callOp<{ input: string }, JsonFormatResponse>("code", "fmt_html", { input });

export const mdToHtml = (input: string, gfm = true, unsafe = false) =>
  callOp<{ input: string; gfm: boolean; unsafe: boolean }, { output: string }>("md", "to_html", { input, gfm, unsafe });
export const mdTableFromCsv = (
  input: string,
  delimiter = ",",
  alignment: "none" | "left" | "center" | "right" = "none",
) =>
  callOp<{ input: string; delimiter: string; alignment: string }, { output: string }>(
    "md",
    "table_from_csv",
    { input, delimiter, alignment },
  );

// ───────────────────────────── Converters ─────────────────────────────

export const encBase64Encode = (input: string, urlSafe = false, noPadding = false) =>
  callOp<{ input: string; urlSafe: boolean; noPadding: boolean }, { output: string }>(
    "enc",
    "base64_encode",
    { input, urlSafe, noPadding },
  );
export const encBase64Decode = (input: string, urlSafe = false) =>
  callOp<{ input: string; urlSafe: boolean }, { output: string; diagnostics?: Diagnostic[] }>(
    "enc",
    "base64_decode",
    { input, urlSafe },
  );
export const encUrlEncode = (input: string, path = false) =>
  callOp<{ input: string; path: boolean }, { output: string }>("enc", "url_encode", { input, path });
export const encUrlDecode = (input: string) =>
  callOp<{ input: string }, { output: string }>("enc", "url_decode", { input });
export const encHtmlEncode = (input: string) =>
  callOp<{ input: string }, { output: string }>("enc", "html_encode", { input });
export const encHtmlDecode = (input: string) =>
  callOp<{ input: string }, { output: string }>("enc", "html_decode", { input });
export const encHexEncode = (input: string, uppercase = false) =>
  callOp<{ input: string; uppercase: boolean }, { output: string }>("enc", "hex_encode", { input, uppercase });
export const encHexDecode = (input: string) =>
  callOp<{ input: string }, { output: string; diagnostics?: Diagnostic[] }>("enc", "hex_decode", { input });

export const dataJsonToCsv = (input: string, flattenSeparator = ".") =>
  callOp<{ input: string; flattenSeparator: string }, { output: string }>("data", "json_to_csv", { input, flattenSeparator });
export const dataCsvToJson = (input: string, header = true, typedValues = false) =>
  callOp<{ input: string; header: boolean; typedValues: boolean }, { output: string }>(
    "data",
    "csv_to_json",
    { input, header, typedValues },
  );
export const dataJsonToXml = (input: string, root = "root", indent = 2) =>
  callOp<{ input: string; root: string; indent: number }, { output: string }>(
    "data",
    "json_to_xml",
    { input, root, indent },
  );
export const dataXmlToJson = (input: string, attrPrefix = "@", textKey = "#text") =>
  callOp<{ input: string; attrPrefix: string; textKey: string }, { output: string }>(
    "data",
    "xml_to_json",
    { input, attrPrefix, textKey },
  );
export const dataFlatten = (input: string, separator = ".") =>
  callOp<{ input: string; separator: string }, { output: string }>("data", "flatten", { input, separator });
export const dataUnflatten = (input: string, separator = ".") =>
  callOp<{ input: string; separator: string }, { output: string }>("data", "unflatten", { input, separator });
export const dataKeyRename = (input: string, rules: { from: string; to: string; regex?: boolean }[]) =>
  callOp<{ input: string; rules: typeof rules }, { output: string }>("data", "key_rename", { input, rules });

export interface ColorConvertResponse {
  hex: string;
  rgb: string;
  hsl: string;
  r: number; g: number; b: number;
  h: number; s: number; l: number;
  diagnostics?: Diagnostic[];
}
export const colorConvert = (input: string) =>
  callOp<{ input: string }, ColorConvertResponse>("color", "convert", { input });

export interface TimeConvertResponse {
  epochS: number;
  epochMS: number;
  epochUS: number;
  epochNS: number;
  rfc3339: string;
  utc: string;
  local: string;
}
export const timeConvert = (input: string, inputFormat = "auto", tz = "UTC") =>
  callOp<{ input: string; inputFormat: string; tz: string }, TimeConvertResponse>(
    "time",
    "convert",
    { input, inputFormat, tz },
  );
export const timeRelative = (from: string, to: string) =>
  callOp<{ from: string; to: string }, { phrase: string; seconds: number }>("time", "relative", { from, to });
export const timeDuration = (input: string) =>
  callOp<{ input: string }, { hours: number; minutes: number; seconds: number; totalSeconds: number }>(
    "time",
    "duration",
    { input },
  );

export const tzConvert = (time: string, fromTZ: string, toTZ: string) =>
  callOp<{ time: string; fromTZ: string; toTZ: string }, { original: string; converted: string; dstNote?: string }>(
    "tz",
    "convert",
    { time, fromTZ, toTZ },
  );
export interface TzZone {
  name: string;
  abbrev: string;
  offsetSeconds: number;
  isDst: boolean;
}
export const tzList = (filter = "") =>
  callOp<{ filter: string }, TzZone[]>("tz", "list", { filter });

export const mathEval = (expression: string) =>
  callOp<{ expression: string }, { value: number; display: string; diagnostics?: Diagnostic[] }>("math", "eval", { expression });
export const mathUnit = (value: number, from: string, to: string) =>
  callOp<{ value: number; from: string; to: string }, { value: number; from: string; to: string }>(
    "math",
    "unit",
    { value, from, to },
  );

// ───────────────────────────── Analyzers ─────────────────────────────

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
export const diffCompare = (
  left: string,
  right: string,
  mode = "auto",
  ignoreOrder = false,
  ignoreWhitespace = false,
) =>
  callOp<{ left: string; right: string; mode: string; ignoreOrder: boolean; ignoreWhitespace: boolean }, DiffResult>(
    "diff",
    "compare",
    { left, right, mode, ignoreOrder, ignoreWhitespace },
  );

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
export const regexExplain = (pattern: string) =>
  callOp<{ pattern: string }, { tree: { token: string; description: string }[] }>(
    "regex",
    "explain",
    { pattern },
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
export const jwtVerify = (
  token: string,
  key: string,
  keyFormat: "hmac" | "pem" = "hmac",
  expectedAlgs: string[] = [],
  leewaySeconds = 0,
) =>
  callOp<{ token: string; key: string; keyFormat: string; expectedAlgs: string[]; leewaySeconds: number }, { valid: boolean; diagnostics?: Diagnostic[] }>(
    "jwt",
    "verify",
    { token, key, keyFormat, expectedAlgs, leewaySeconds },
  );

export const strCase = (input: string, mode: string) =>
  callOp<{ input: string; mode: string }, { output: string }>("str", "case", { input, mode });
export const strDiff = (left: string, right: string, ignoreWhitespace = false, ignoreCase = false) =>
  callOp<{ left: string; right: string; ignoreWhitespace: boolean; ignoreCase: boolean }, DiffResult>(
    "str",
    "diff",
    { left, right, ignoreWhitespace, ignoreCase },
  );
export const strStats = (input: string) =>
  callOp<{ input: string }, { chars: number; bytes: number; words: number; lines: number; runes: number; uniqueWords: number }>(
    "str",
    "stats",
    { input },
  );
export const strSortUnique = (input: string, reverse = false, numeric = false, uniqueOnly = false) =>
  callOp<{ input: string; reverse: boolean; numeric: boolean; uniqueOnly: boolean }, { output: string }>(
    "str",
    "sort_unique",
    { input, reverse, numeric, uniqueOnly },
  );
export const strReplace = (
  input: string,
  pattern: string,
  replace: string,
  regex = false,
  ignoreCase = false,
) =>
  callOp<{ input: string; pattern: string; replace: string; regex: boolean; ignoreCase: boolean }, { output: string; replacements: number }>(
    "str",
    "replace",
    { input, pattern, replace, regex, ignoreCase },
  );

export const urlParse = (url: string) =>
  callOp<{ input: string }, {
    scheme: string;
    hostname: string;
    port: string;
    path: string;
    query: string;
    fragment: string;
    params: { key: string; value: string }[];
  }>("url", "parse", { input: url });

export function parseHeaderString(raw: string): Record<string, string> {
  const map: Record<string, string> = {};
  for (const line of raw.split(/\r?\n/)) {
    const idx = line.indexOf(":");
    if (idx < 1) continue;
    const key = line.slice(0, idx).trim();
    const value = line.slice(idx + 1).trim();
    if (key && value) map[key] = value;
  }
  return map;
}

export const headersAnalyze = (input: string) =>
  callOp<{ headers: Record<string, string> }, {
    headers: Record<string, string>;
    findings: { header: string; present: boolean; ok: boolean; note: string }[];
  }>("headers", "analyze", { headers: parseHeaderString(input) });

export interface DnsLookupResponse {
  host: string;
  a?: string[];
  aaaa?: string[];
  cname?: string[];
  mx?: { host: string; pref: number }[];
  ns?: string[];
  txt?: string[][];
  ptr?: string[];
  anyPrivate: boolean;
}
export const dnsLookup = (host: string, type = "A", allowPrivate = false) =>
  callOp<{ host: string; type: string; allowPrivate: boolean }, DnsLookupResponse>(
    "dns",
    "lookup",
    { host, type: type.toLowerCase(), allowPrivate },
  );

export const httpRequest = (input: {
  method: string;
  url: string;
  headers?: { key: string; value: string }[];
  body?: string;
  followRedirects?: boolean;
  maxRedirects?: number;
  timeoutSeconds?: number;
  allowPrivate?: boolean;
}) => {
  const headersMap: Record<string, string> = {};
  for (const h of input.headers ?? []) {
    if (h.key) headersMap[h.key] = h.value;
  }
  return callOp<{
    method: string; url: string; headers: Record<string, string>; body?: string;
    followRedirects?: boolean; maxRedirects?: number; timeoutSec?: number; allowPrivate?: boolean;
  }, {
    status: number;
    statusText: string;
    durationMs: number;
    headers: Record<string, string>;
    body: string;
    diagnostics?: Diagnostic[];
  }>("http", "request", {
    method: input.method,
    url: input.url,
    headers: headersMap,
    body: input.body,
    followRedirects: input.followRedirects,
    maxRedirects: input.maxRedirects,
    timeoutSec: input.timeoutSeconds,
    allowPrivate: input.allowPrivate,
  });
};

export const ipCalc = (cidr: string, maxList = 0) =>
  callOp<{ cidr: string; maxList: number }, {
    network: string;
    broadcast: string;
    first: string;
    last: string;
    netmask: string;
    wildcard: string;
    prefix: number;
    usableHosts: string;
    hosts?: string[];
    diagnostics?: Diagnostic[];
  }>("ip", "calc", { cidr, maxList });

// ───────────────────────────── DevOps ─────────────────────────────

export const gitPatch = (input: {
  left: string;
  right: string;
  leftPath?: string;
  rightPath?: string;
  contextLines?: number;
}) => callOp<{ left: string; right: string; leftPath?: string; rightPath?: string; context?: number }, { output: string }>("git", "patch", { left: input.left, right: input.right, leftPath: input.leftPath, rightPath: input.rightPath, context: input.contextLines });

export const gitCommitFormat = (message: string) =>
  callOp<{ input: string }, {
    type: string;
    scope: string;
    subject: string;
    body: string;
    footer: string;
    breaking: boolean;
    valid: boolean;
    diagnostics?: Diagnostic[];
  }>("git", "commit_format", { input: message });

export const gitIgnoreGen = (templates: string[]) =>
  callOp<{ templates: string[] }, { output: string; diagnostics?: Diagnostic[] }>(
    "git",
    "ignore_gen",
    { templates },
  );

export const dockerfileLint = (input: string) =>
  callOp<{ input: string }, { diagnostics: Diagnostic[]; summary: { errors: number; warnings: number; info: number } }>(
    "dockerfile",
    "lint",
    { input },
  );

export const envParse = (input: string, allowExport = false) =>
  callOp<{ input: string; allowExport: boolean }, {
    values: Record<string, string>;
    diagnostics?: Diagnostic[];
  }>("env", "parse", { input, allowExport });

export interface EnvDiffResponse {
  added: string[] | null;
  removed: string[] | null;
  changed: { key: string; left: string; right: string }[] | null;
}
export const envDiff = (left: string, right: string) =>
  callOp<{ left: string; right: string }, EnvDiffResponse>("env", "diff", { left, right });

export const k8sValidate = (input: string) =>
  callOp<{ input: string }, {
    valid: boolean;
    kind?: string;
    apiVersion?: string;
    diagnostics?: Diagnostic[];
  }>("k8s", "validate", { input });
