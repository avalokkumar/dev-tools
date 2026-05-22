import { useEffect, useMemo, useState } from "react";
import { Bot, CheckCircle2, XCircle, Cpu, Sparkles } from "lucide-react";
import { listOperations, type OpInfo } from "../lib/api";
import { TOOLS, CATEGORIES } from "../lib/catalog";
import { Card, CardHeader, Badge } from "../components/ui/Card";
import { Input } from "../components/ui/Input";
import { Button } from "../components/ui/Button";
import { CodeBlock, CopyButton } from "../components/ui/Output";
import { cn } from "../lib/cn";

type ClientId = "claude-code" | "cursor" | "claude-desktop";

interface ClientCfg {
  id: ClientId;
  name: string;
  configPath: string;
  description: string;
}

const CLIENTS: ClientCfg[] = [
  {
    id: "claude-code",
    name: "Claude Code",
    configPath: ".mcp.json (project root)",
    description: "Anthropic's terminal coding assistant. Project-local config recommended.",
  },
  {
    id: "cursor",
    name: "Cursor",
    configPath: "~/.cursor/mcp.json",
    description: "Cursor IDE. Same JSON shape as Claude Code.",
  },
  {
    id: "claude-desktop",
    name: "Claude Desktop",
    configPath: "~/Library/Application Support/Claude/claude_desktop_config.json",
    description: "Claude Desktop app (macOS path shown).",
  },
];

export default function McpConfigPage() {
  const [client, setClient] = useState<ClientId>("claude-code");
  const [binary, setBinary] = useState("/usr/local/bin/devforge");
  const [pluginDir, setPluginDir] = useState("");
  const [smoke, setSmoke] = useState<{ ok: boolean; count?: number; error?: string } | null>(null);
  const [loading, setLoading] = useState(false);
  const [ops, setOps] = useState<OpInfo[] | null>(null);

  useEffect(() => {
    listOperations()
      .then(setOps)
      .catch(() => setOps([]));
  }, []);

  const config = useMemo(() => {
    const env: Record<string, string> = {};
    if (pluginDir) env.DEVFORGE_PLUGIN_DIR = pluginDir;
    const server: Record<string, unknown> = { command: binary, args: ["mcp"] };
    if (Object.keys(env).length) server.env = env;
    return JSON.stringify({ mcpServers: { devforge: server } }, null, 2);
  }, [binary, pluginDir]);

  async function runSmokeTest() {
    setLoading(true);
    setSmoke(null);
    try {
      const r = await listOperations(true);
      setSmoke({ ok: true, count: r.length });
    } catch (e) {
      setSmoke({ ok: false, error: e instanceof Error ? e.message : String(e) });
    } finally {
      setLoading(false);
    }
  }

  function downloadConfig() {
    const blob = new Blob([config], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = ".mcp.json";
    a.click();
    URL.revokeObjectURL(url);
  }

  return (
    <div className="px-md md:px-lg py-md md:py-lg max-w-[1400px] mx-auto w-full flex flex-col gap-lg">
      <header className="flex items-start gap-md">
        <div className="h-12 w-12 rounded-md bg-primary text-on-primary flex items-center justify-center shrink-0">
          <Bot className="h-6 w-6" />
        </div>
        <div>
          <h1 className="font-display-brand text-display-brand text-on-surface leading-none">
            MCP Configuration
          </h1>
          <p className="font-body-md text-body-md text-on-surface-variant mt-1">
            Connect DevForge as a Model Context Protocol server to your AI coding assistant.
          </p>
        </div>
      </header>

      {/* Step 1: client */}
      <Card>
        <CardHeader title="Step 1 — Select your client" />
        <div className="grid grid-cols-1 md:grid-cols-3 gap-md">
          {CLIENTS.map((c) => {
            const active = c.id === client;
            return (
              <button
                key={c.id}
                onClick={() => setClient(c.id)}
                className={cn(
                  "text-left rounded-lg border p-md transition-colors",
                  active
                    ? "border-secondary-container bg-secondary-container/10"
                    : "border-outline/20 hover:border-outline/40 bg-surface-container-lowest",
                )}
              >
                <div className="flex items-center justify-between mb-2">
                  <Cpu className="h-5 w-5 text-on-surface-variant" />
                  {active && <Badge tone="info">Selected</Badge>}
                </div>
                <div className="font-body-md text-body-md font-medium text-on-surface">{c.name}</div>
                <div className="font-body-sm text-body-sm text-on-surface-variant mt-0.5">
                  {c.description}
                </div>
                <code className="font-data-label text-data-label uppercase text-on-surface-variant mt-2 block break-all">
                  {c.configPath}
                </code>
              </button>
            );
          })}
        </div>
      </Card>

      {/* Step 2: paths */}
      <Card>
        <CardHeader title="Step 2 — Binary path & plugins" />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-md">
          <Input
            label="DevForge binary path"
            value={binary}
            onChange={(e) => setBinary(e.target.value)}
            placeholder="/usr/local/bin/devforge"
          />
          <Input
            label="Plugin directory (optional)"
            value={pluginDir}
            onChange={(e) => setPluginDir(e.target.value)}
            placeholder="~/.devforge/plugins"
          />
        </div>
      </Card>

      {/* Step 3: config */}
      <Card padded={false}>
        <CardHeader
          title="Step 3 — Copy or download config"
          trailing={
            <div className="flex items-center gap-2 px-md">
              <Button size="sm" variant="outline" onClick={downloadConfig}>
                Download .mcp.json
              </Button>
              <CopyButton text={config} label="Copy config" />
            </div>
          }
        />
        <div className="p-md">
          <CodeBlock code={config} language="json" showCopy={false} />
        </div>
      </Card>

      {/* Step 4: smoke test */}
      <Card>
        <CardHeader
          title="Step 4 — Verify connection"
          trailing={
            <Button onClick={runSmokeTest} loading={loading} variant="primary" size="sm">
              <Sparkles className="h-4 w-4" /> Run smoke test
            </Button>
          }
        />
        {smoke ? (
          smoke.ok ? (
            <div className="flex items-center gap-3 p-md rounded-lg bg-tertiary-fixed text-on-tertiary-fixed">
              <CheckCircle2 className="h-6 w-6" />
              <div>
                <div className="font-body-md font-semibold">
                  Connected — {smoke.count} operations registered.
                </div>
                <div className="font-body-sm text-body-sm">
                  Your AI client will see all of these as MCP tools.
                </div>
              </div>
            </div>
          ) : (
            <div className="flex items-center gap-3 p-md rounded-lg bg-error/10 text-error">
              <XCircle className="h-6 w-6" />
              <div>
                <div className="font-body-md font-semibold">Server unreachable</div>
                <div className="font-body-sm text-body-sm break-all">{smoke.error}</div>
              </div>
            </div>
          )
        ) : (
          <div className="font-body-sm text-body-sm text-on-surface-variant">
            Click "Run smoke test" to verify the running server responds.
          </div>
        )}
      </Card>

      {/* Tools list */}
      <Card>
        <CardHeader
          title="Registered Tools"
          trailing={ops && <Badge>{ops.length} ops</Badge>}
        />
        {!ops ? (
          <div className="font-body-sm text-body-sm text-on-surface-variant">Loading…</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-md">
            {CATEGORIES.map((c) => {
              const tools = TOOLS.filter((t) => t.category === c.key);
              return (
                <div key={c.key} className="bg-surface-container-low rounded-md p-3">
                  <div className="font-data-label text-data-label uppercase text-on-surface-variant mb-2">
                    {c.key} ({tools.length})
                  </div>
                  <ul className="flex flex-col gap-1">
                    {tools.map((t) => (
                      <li key={t.slug} className="flex items-center gap-2 text-body-sm font-body-sm">
                        <t.icon className="h-3.5 w-3.5 text-on-surface-variant" />
                        <span className="text-on-surface flex-1 truncate">{t.name}</span>
                        <span className="font-data-label text-data-label text-on-surface-variant">
                          {t.ops.length}
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              );
            })}
          </div>
        )}
      </Card>

      {/* Sample prompts */}
      <Card>
        <CardHeader title="Sample prompts" />
        <ul className="grid grid-cols-1 md:grid-cols-2 gap-2">
          {SAMPLE_PROMPTS.map((p) => (
            <li
              key={p}
              className="flex items-start justify-between gap-2 px-3 py-2 bg-surface-container-low rounded"
            >
              <span className="font-body-sm text-body-sm text-on-surface flex-1">{p}</span>
              <CopyButton text={p} />
            </li>
          ))}
        </ul>
      </Card>
    </div>
  );
}

const SAMPLE_PROMPTS = [
  "Generate three v7 UUIDs.",
  "Decode this JWT and tell me what role the user has.",
  "Format this JSON and validate it.",
  "What's the strength of password 'hunter2'?",
  "Diff these two JSON objects semantically.",
  "Compute the usable hosts in 10.0.0.0/24.",
  "Convert this ISO timestamp to America/Los_Angeles.",
  "Generate a .gitignore for a Node + Docker project.",
];
