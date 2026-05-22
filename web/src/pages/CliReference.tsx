import { useEffect, useMemo, useState } from "react";
import { Terminal, Search } from "lucide-react";
import { Link } from "react-router-dom";
import { listOperations, type OpInfo } from "../lib/api";
import { TOOLS } from "../lib/catalog";
import { Card, CardHeader, Badge } from "../components/ui/Card";
import { CodeBlock, CopyButton, EmptyState, ErrorBanner } from "../components/ui/Output";
import { Input } from "../components/ui/Input";
import { cn } from "../lib/cn";

export default function CliReferencePage() {
  const [ops, setOps] = useState<OpInfo[] | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [filter, setFilter] = useState("");
  const [active, setActive] = useState<OpInfo | null>(null);

  useEffect(() => {
    listOperations()
      .then((r) => {
        setOps(r);
        setActive(r[0] ?? null);
      })
      .catch((e: Error) => setError(e));
  }, []);

  const grouped = useMemo(() => {
    if (!ops) return new Map<string, OpInfo[]>();
    const q = filter.trim().toLowerCase();
    const g = new Map<string, OpInfo[]>();
    for (const o of ops) {
      if (q && !`${o.tool} ${o.op} ${o.name} ${o.description}`.toLowerCase().includes(q)) continue;
      const arr = g.get(o.tool) ?? [];
      arr.push(o);
      g.set(o.tool, arr);
    }
    return g;
  }, [ops, filter]);

  return (
    <div className="px-md md:px-lg py-md md:py-lg max-w-[1400px] mx-auto w-full">
      <header className="flex items-start gap-md mb-lg">
        <div className="h-12 w-12 rounded-md bg-primary text-on-primary flex items-center justify-center shrink-0">
          <Terminal className="h-6 w-6" />
        </div>
        <div>
          <h1 className="font-display-brand text-display-brand text-on-surface leading-none">
            CLI Reference
          </h1>
          <p className="font-body-md text-body-md text-on-surface-variant mt-1">
            Every operation usable as <code className="font-code-block bg-surface-container-low px-1 rounded">devforge run</code>{" "}
            from the command line or via the MCP server.
          </p>
        </div>
      </header>

      {error && <ErrorBanner error={error} />}

      <div className="grid grid-cols-1 lg:grid-cols-[300px_1fr] gap-md">
        <Card padded={false}>
          <div className="p-md border-b border-outline/10">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-on-surface-variant" />
              <Input
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                placeholder="Filter…"
                className="pl-9"
              />
            </div>
          </div>
          <div className="max-h-[70vh] overflow-y-auto thin-scrollbar p-2">
            {ops === null ? (
              <div className="p-md text-on-surface-variant font-body-sm text-body-sm">Loading…</div>
            ) : (
              [...grouped.entries()].map(([toolKey, list]) => {
                const meta = TOOLS.find((t) => t.tool === toolKey);
                const Icon = meta?.icon;
                return (
                  <div key={toolKey} className="mb-2">
                    <div className="flex items-center gap-2 px-2 py-1 font-data-label text-data-label uppercase text-on-surface-variant">
                      {Icon && <Icon className="h-3.5 w-3.5" />}
                      {toolKey}
                      <span className="opacity-60">({list.length})</span>
                    </div>
                    <ul className="flex flex-col">
                      {list.map((o) => (
                        <li key={o.path}>
                          <button
                            onClick={() => setActive(o)}
                            className={cn(
                              "w-full text-left px-3 py-1.5 rounded text-body-sm font-body-sm transition-colors",
                              active?.path === o.path
                                ? "bg-secondary-container text-on-secondary"
                                : "hover:bg-surface-container-low text-on-surface",
                            )}
                          >
                            <code className="font-code-block">{o.tool}.{o.op}</code>
                          </button>
                        </li>
                      ))}
                    </ul>
                  </div>
                );
              })
            )}
            {grouped.size === 0 && ops && (
              <div className="p-md text-on-surface-variant font-body-sm text-body-sm text-center">
                No operations match "{filter}".
              </div>
            )}
          </div>
        </Card>

        <Card>
          {!active ? (
            <EmptyState title="Select an operation" icon={<Terminal className="h-6 w-6" />} />
          ) : (
            <OperationDetails op={active} />
          )}
        </Card>
      </div>
    </div>
  );
}

function OperationDetails({ op }: { op: OpInfo }) {
  const meta = TOOLS.find((t) => t.tool === op.tool);
  const slug = meta?.slug;
  const cmd = `devforge run ${op.tool}.${op.op} --args '{}'`;
  return (
    <div className="flex flex-col gap-md">
      <CardHeader
        title={`${op.tool}.${op.op}`}
        subtitle={op.description}
        trailing={
          <div className="flex items-center gap-2">
            {meta && <Badge>{meta.category}</Badge>}
            {slug && (
              <Link
                to={`/tools/${slug}`}
                className="font-body-sm text-body-sm text-secondary-container hover:underline"
              >
                Open in UI →
              </Link>
            )}
          </div>
        }
      />
      <div>
        <div className="font-data-label text-data-label uppercase text-on-surface-variant mb-1">
          API path
        </div>
        <div className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded">
          <code className="font-code-block text-code-block text-on-surface flex-1 break-all">
            POST {op.path}
          </code>
          <CopyButton text={op.path} />
        </div>
      </div>
      <div>
        <div className="font-data-label text-data-label uppercase text-on-surface-variant mb-1">
          CLI invocation
        </div>
        <CodeBlock code={cmd} language="bash" />
      </div>
      {op.inputSchema && (
        <div>
          <div className="font-data-label text-data-label uppercase text-on-surface-variant mb-1">
            Input schema
          </div>
          <CodeBlock code={JSON.stringify(op.inputSchema, null, 2)} language="json" />
        </div>
      )}
    </div>
  );
}
