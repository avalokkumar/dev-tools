// Tool catalog. Single source of truth for sidebar, home grid, search.
import {
  Key,
  Drama,
  Tag,
  Hash,
  ShieldCheck,
  Braces,
  FileText,
  Table,
  Database,
  Terminal,
  PenTool,
  Binary,
  Shuffle,
  Palette,
  Timer,
  Globe2,
  Calculator,
  GitCompareArrows,
  Regex,
  Clock,
  Ticket,
  Scissors,
  Link as LinkIcon,
  Mail,
  Globe,
  Radio,
  MapPin,
  GitBranch,
  Container,
  ClipboardList,
  Boxes,
  type LucideIcon,
} from "lucide-react";

export type Category =
  | "Generators"
  | "Formatters"
  | "Converters"
  | "Analyzers"
  | "DevOps";

export interface ToolMeta {
  /** Slug for URL: /tools/<slug> */
  slug: string;
  /** Backend tool prefix (used for /api/v1/<tool>/<op>) */
  tool: string;
  /** Display name */
  name: string;
  /** Short tagline shown on cards */
  tagline: string;
  /** Lucide icon */
  icon: LucideIcon;
  /** Category bucket */
  category: Category;
  /** Backend operation names this page exposes */
  ops: string[];
}

export const CATEGORIES: { key: Category; tagline: string }[] = [
  { key: "Generators", tagline: "Create IDs, secrets, mock data" },
  { key: "Formatters", tagline: "Pretty-print and validate" },
  { key: "Converters", tagline: "Transform between formats" },
  { key: "Analyzers", tagline: "Inspect, debug, explain" },
  { key: "DevOps", tagline: "Git, Docker, Kubernetes" },
];

export const TOOLS: ToolMeta[] = [
  // Generators (5)
  {
    slug: "uuid",
    tool: "uuid",
    name: "UUID Generator",
    tagline: "Generate v4/v7 UUIDs",
    icon: Key,
    category: "Generators",
    ops: ["generate", "hash"],
  },
  {
    slug: "faker",
    tool: "faker",
    name: "Data Faker",
    tagline: "Mock data in JSON, CSV, SQL",
    icon: Drama,
    category: "Generators",
    ops: ["generate", "kinds"],
  },
  {
    slug: "id",
    tool: "id",
    name: "ID Generator",
    tagline: "ULID + URL-safe slugs",
    icon: Tag,
    category: "Generators",
    ops: ["ulid", "slug"],
  },
  {
    slug: "totp",
    tool: "totp",
    name: "TOTP",
    tagline: "Generate & verify TOTP codes",
    icon: Hash,
    category: "Generators",
    ops: ["generate", "verify"],
  },
  {
    slug: "crypto",
    tool: "crypto",
    name: "Crypto Toolkit",
    tagline: "AES, RSA, HMAC, password tools",
    icon: ShieldCheck,
    category: "Generators",
    ops: [
      "aes_encrypt",
      "aes_decrypt",
      "rsa_keygen",
      "hmac",
      "password_hash",
      "password_strength",
    ],
  },

  // Formatters (6)
  {
    slug: "json",
    tool: "json",
    name: "JSON Formatter",
    tagline: "Format & validate JSON",
    icon: Braces,
    category: "Formatters",
    ops: ["format", "validate"],
  },
  {
    slug: "yaml",
    tool: "yaml",
    name: "YAML Formatter",
    tagline: "Format, validate, convert",
    icon: FileText,
    category: "Formatters",
    ops: ["format", "validate", "convert"],
  },
  {
    slug: "csv",
    tool: "csv",
    name: "CSV Formatter",
    tagline: "Format & validate CSV",
    icon: Table,
    category: "Formatters",
    ops: ["format", "validate"],
  },
  {
    slug: "sql",
    tool: "sql",
    name: "SQL Formatter",
    tagline: "Format & lint SQL",
    icon: Database,
    category: "Formatters",
    ops: ["format", "validate"],
  },
  {
    slug: "code",
    tool: "code",
    name: "Code Formatter",
    tagline: "Go, XML, HTML",
    icon: Terminal,
    category: "Formatters",
    ops: ["fmt_go", "fmt_xml", "fmt_html"],
  },
  {
    slug: "markdown",
    tool: "md",
    name: "Markdown",
    tagline: "Preview & CSV→table",
    icon: PenTool,
    category: "Formatters",
    ops: ["to_html", "table_from_csv"],
  },

  // Converters (6)
  {
    slug: "encoding",
    tool: "enc",
    name: "Encoding Converter",
    tagline: "Base64, URL, HTML, Hex",
    icon: Binary,
    category: "Converters",
    ops: [
      "base64_encode",
      "base64_decode",
      "url_encode",
      "url_decode",
      "html_encode",
      "html_decode",
      "hex_encode",
      "hex_decode",
    ],
  },
  {
    slug: "data",
    tool: "data",
    name: "Data Transformer",
    tagline: "JSON ⇄ CSV ⇄ XML, flatten",
    icon: Shuffle,
    category: "Converters",
    ops: [
      "json_to_csv",
      "csv_to_json",
      "json_to_xml",
      "xml_to_json",
      "flatten",
      "unflatten",
      "key_rename",
    ],
  },
  {
    slug: "color",
    tool: "color",
    name: "Color Converter",
    tagline: "HEX ⇄ RGB ⇄ HSL + a11y",
    icon: Palette,
    category: "Converters",
    ops: ["convert"],
  },
  {
    slug: "time",
    tool: "time",
    name: "Time Converter",
    tagline: "Epoch, ISO, durations",
    icon: Timer,
    category: "Converters",
    ops: ["convert", "relative", "duration"],
  },
  {
    slug: "timezone",
    tool: "tz",
    name: "Timezone Converter",
    tagline: "DST-aware zone math",
    icon: Globe2,
    category: "Converters",
    ops: ["convert", "list"],
  },
  {
    slug: "math",
    tool: "math",
    name: "Math & Units",
    tagline: "Calculator + unit converter",
    icon: Calculator,
    category: "Converters",
    ops: ["eval", "unit"],
  },

  // Analyzers (10)
  {
    slug: "diff",
    tool: "diff",
    name: "Smart Diff",
    tagline: "Semantic diff JSON/INI/SQL",
    icon: GitCompareArrows,
    category: "Analyzers",
    ops: ["compare"],
  },
  {
    slug: "regex",
    tool: "regex",
    name: "Regex Tester",
    tagline: "Test & explain patterns",
    icon: Regex,
    category: "Analyzers",
    ops: ["test", "explain"],
  },
  {
    slug: "cron",
    tool: "cron",
    name: "Cron Builder",
    tagline: "Parse expressions & next runs",
    icon: Clock,
    category: "Analyzers",
    ops: ["parse", "next"],
  },
  {
    slug: "jwt",
    tool: "jwt",
    name: "JWT Debugger",
    tagline: "Decode & verify tokens",
    icon: Ticket,
    category: "Analyzers",
    ops: ["decode", "verify"],
  },
  {
    slug: "string",
    tool: "str",
    name: "String Tools",
    tagline: "Case, diff, stats, replace",
    icon: Scissors,
    category: "Analyzers",
    ops: ["case", "diff", "stats", "sort_unique", "replace"],
  },
  {
    slug: "url",
    tool: "url",
    name: "URL Parser",
    tagline: "Inspect URL components",
    icon: LinkIcon,
    category: "Analyzers",
    ops: ["parse"],
  },
  {
    slug: "headers",
    tool: "headers",
    name: "HTTP Headers",
    tagline: "Security audit",
    icon: Mail,
    category: "Analyzers",
    ops: ["analyze"],
  },
  {
    slug: "dns",
    tool: "dns",
    name: "DNS Lookup",
    tagline: "A, AAAA, MX, NS, TXT…",
    icon: Globe,
    category: "Analyzers",
    ops: ["lookup"],
  },
  {
    slug: "http",
    tool: "http",
    name: "HTTP Client",
    tagline: "Send requests, inspect responses",
    icon: Radio,
    category: "Analyzers",
    ops: ["request"],
  },
  {
    slug: "ip",
    tool: "ip",
    name: "IP Calculator",
    tagline: "CIDR & subnet math",
    icon: MapPin,
    category: "Analyzers",
    ops: ["calc"],
  },

  // DevOps (4)
  {
    slug: "git",
    tool: "git",
    name: "Git Tools",
    tagline: "Patch, commit, .gitignore",
    icon: GitBranch,
    category: "DevOps",
    ops: ["patch", "commit_format", "ignore_gen"],
  },
  {
    slug: "dockerfile",
    tool: "dockerfile",
    name: "Dockerfile Lint",
    tagline: "Best-practice checks",
    icon: Container,
    category: "DevOps",
    ops: ["lint"],
  },
  {
    slug: "env",
    tool: "env",
    name: "Env Files",
    tagline: "Parse & diff .env",
    icon: ClipboardList,
    category: "DevOps",
    ops: ["parse", "diff"],
  },
  {
    slug: "k8s",
    tool: "k8s",
    name: "K8s Validator",
    tagline: "Validate manifests",
    icon: Boxes,
    category: "DevOps",
    ops: ["validate"],
  },
];

export const TOOLS_BY_CATEGORY: Record<Category, ToolMeta[]> = CATEGORIES.reduce(
  (acc, c) => {
    acc[c.key] = TOOLS.filter((t) => t.category === c.key);
    return acc;
  },
  {} as Record<Category, ToolMeta[]>,
);

export function findTool(slug: string): ToolMeta | undefined {
  return TOOLS.find((t) => t.slug === slug);
}

export const TOTAL_OPS = TOOLS.reduce((sum, t) => sum + t.ops.length, 0);
