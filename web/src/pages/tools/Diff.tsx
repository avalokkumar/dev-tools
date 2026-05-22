import { useState } from "react";
import { GitCompareArrows, Play, Plus, Minus, Pencil } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { diffCompare } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Textarea, Select, Toggle } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("diff")!;

const L = `{"name":"DevForge","tools":50,"local":true}`;
const R = `{"name":"DevForge","tools":75,"local":true,"version":"0.1"}`;

export default function DiffPage() {
  const [left, setLeft] = useState(L);
  const [right, setRight] = useState(R);
  const [mode, setMode] = useState("auto");
  const [ignoreOrder, setIgnoreOrder] = useState(false);
  const [ignoreWhitespace, setIgnoreWhitespace] = useState(false);
  const op = useOperation((args: { left: string; right: string; mode: string; ignoreOrder: boolean; ignoreWhitespace: boolean }) =>
    diffCompare(args.left, args.right, args.mode, args.ignoreOrder, args.ignoreWhitespace),
  );

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 gap-md">
        <Card>
          <CardHeader title="Inputs" icon={<GitCompareArrows className="h-5 w-5" />} />
          <div className="grid grid-cols-1 md:grid-cols-2 gap-md">
            <Textarea label="Left" rows={10} value={left} onChange={(e) => setLeft(e.target.value)} />
            <Textarea label="Right" rows={10} value={right} onChange={(e) => setRight(e.target.value)} />
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-md mt-md">
            <Select
              label="Mode"
              value={mode}
              onChange={(e) => setMode(e.target.value)}
              options={[
                { value: "auto", label: "Auto-detect" },
                { value: "json", label: "JSON" },
                { value: "ini", label: "INI" },
                { value: "sql", label: "SQL" },
              ]}
            />
            <Toggle checked={ignoreOrder} onChange={setIgnoreOrder} label="Ignore array order" />
            <Toggle checked={ignoreWhitespace} onChange={setIgnoreWhitespace} label="Ignore whitespace" />
          </div>
          <Button
            onClick={() => op.run({ left, right, mode, ignoreOrder, ignoreWhitespace })}
            loading={op.loading}
            className="mt-md"
            disabled={!left.trim() || !right.trim()}
          >
            <Play className="h-4 w-4" /> Compare
          </Button>
        </Card>
        <Card padded={false}>
          <div className="flex items-center justify-between p-md border-b border-outline/10">
            <h3 className="font-body-md text-body-md font-medium text-on-surface">Hunks</h3>
            {op.data && (
              <div className="flex gap-2">
                <Badge tone="success">+{op.data.summary.adds}</Badge>
                <Badge tone="error">−{op.data.summary.removes}</Badge>
                <Badge tone="info">~{op.data.summary.changes}</Badge>
              </div>
            )}
          </div>
          <div className="p-md">
            {op.error ? (
              <ErrorBanner error={op.error} />
            ) : !op.data ? (
              <EmptyState title="Run a comparison to see semantic hunks" />
            ) : op.data.hunks.length === 0 ? (
              <EmptyState title="No differences" description="Inputs are semantically equivalent." />
            ) : (
              <ul className="flex flex-col gap-1">
                {op.data.hunks.map((h, i) => {
                  const cfg =
                    h.op === "add"
                      ? { Icon: Plus, color: "text-tertiary-container", bg: "bg-tertiary-fixed" }
                      : h.op === "remove"
                        ? { Icon: Minus, color: "text-error", bg: "bg-error/10" }
                        : { Icon: Pencil, color: "text-primary", bg: "bg-sky-aqua/10" };
                  return (
                    <li
                      key={i}
                      className={`flex items-start gap-2 px-3 py-2 rounded ${cfg.bg}`}
                    >
                      <cfg.Icon className={`h-4 w-4 mt-0.5 ${cfg.color}`} />
                      <code className="font-code-block text-code-block text-on-surface flex-1 break-all">
                        <span className="font-data-label text-data-label uppercase opacity-60 mr-2">
                          {h.path || "/"}
                        </span>
                        {h.op !== "add" && (
                          <span className="text-error mr-1">{JSON.stringify(h.left)}</span>
                        )}
                        {h.op !== "remove" && (
                          <span className="text-tertiary-container">→ {JSON.stringify(h.right)}</span>
                        )}
                      </code>
                    </li>
                  );
                })}
              </ul>
            )}
          </div>
        </Card>
      </div>
    </ToolPage>
  );
}
