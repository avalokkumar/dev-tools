import { useState } from "react";
import { ClipboardList, Play, Plus, Minus, Pencil } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { envParse, envDiff } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Textarea, Toggle } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("env")!;
const SAMPLE = `DATABASE_URL=postgres://localhost:5432/app\nLOG_LEVEL=info\nFEATURE_FLAG_X=true\n`;

export default function EnvPage() {
  const [tab, setTab] = useState<"parse" | "diff">("parse");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "parse", label: "Parse" },
          { id: "diff", label: "Diff" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "parse" ? <ParseTab /> : <DiffTab />}
    </ToolPage>
  );
}

function ParseTab() {
  const [input, setInput] = useState(SAMPLE);
  const [allowExport, setAllowExport] = useState(true);
  const op = useOperation((args: { input: string; allowExport: boolean }) =>
    envParse(args.input, args.allowExport),
  );
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title=".env" icon={<ClipboardList className="h-5 w-5" />} />
        <Textarea rows={14} value={input} onChange={(e) => setInput(e.target.value)} />
        <div className="mt-md">
          <Toggle checked={allowExport} onChange={setAllowExport} label="Allow `export VAR=…` syntax" />
        </div>
        <Button onClick={() => op.run({ input, allowExport })} loading={op.loading} fullWidth className="mt-md">
          <Play className="h-4 w-4" /> Parse
        </Button>
      </Card>
      <Card>
        <CardHeader title="Variables" trailing={op.data && <Badge>{Object.keys(op.data.values).length}</Badge>} />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Parsed variables appear here" />
        ) : (
          <ul className="flex flex-col gap-1">
            {Object.entries(op.data.values).map(([k, v]) => (
              <li
                key={k}
                className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded font-code-block text-code-block"
              >
                <span className="text-secondary-container shrink-0">{k}</span>
                <span className="text-on-surface-variant">=</span>
                <span className="text-on-surface flex-1 truncate">{v}</span>
                <CopyButton text={`${k}=${v}`} />
              </li>
            ))}
          </ul>
        )}
        <Diagnostics items={op.data?.diagnostics} className="mt-md" />
      </Card>
    </div>
  );
}

function DiffTab() {
  const [left, setLeft] = useState("FOO=1\nBAR=2\n");
  const [right, setRight] = useState("FOO=2\nBAZ=3\n");
  const op = useOperation((args: { left: string; right: string }) => envDiff(args.left, args.right));
  return (
    <div className="grid grid-cols-1 gap-md">
      <Card>
        <CardHeader title="Inputs" />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-md">
          <Textarea label="Left" rows={8} value={left} onChange={(e) => setLeft(e.target.value)} />
          <Textarea label="Right" rows={8} value={right} onChange={(e) => setRight(e.target.value)} />
        </div>
        <Button onClick={() => op.run({ left, right })} loading={op.loading} className="mt-md">
          <Play className="h-4 w-4" /> Compare
        </Button>
      </Card>
      <Card padded={false}>
        {(() => {
          const added = op.data?.added ?? [];
          const removed = op.data?.removed ?? [];
          const changed = op.data?.changed ?? [];
          const total = added.length + removed.length + changed.length;
          return (
            <>
              <div className="flex items-center justify-between p-md border-b border-outline/10">
                <h3 className="font-body-md text-body-md font-medium text-on-surface">Differences</h3>
                {op.data && (
                  <div className="flex gap-2">
                    <Badge tone="success">+{added.length}</Badge>
                    <Badge tone="error">−{removed.length}</Badge>
                    <Badge tone="info">~{changed.length}</Badge>
                  </div>
                )}
              </div>
              <div className="p-md">
                {op.error ? (
                  <ErrorBanner error={op.error} />
                ) : !op.data ? (
                  <EmptyState title="Run a comparison" />
                ) : total === 0 ? (
                  <EmptyState title="No differences" />
                ) : (
                  <ul className="flex flex-col gap-1">
                    {added.map((k) => (
                      <li key={`a-${k}`} className="flex items-start gap-2 px-3 py-2 rounded bg-tertiary-fixed">
                        <Plus className="h-4 w-4 mt-0.5 text-tertiary-container" />
                        <code className="font-code-block text-code-block flex-1 break-all">{k}</code>
                      </li>
                    ))}
                    {removed.map((k) => (
                      <li key={`r-${k}`} className="flex items-start gap-2 px-3 py-2 rounded bg-error/10">
                        <Minus className="h-4 w-4 mt-0.5 text-error" />
                        <code className="font-code-block text-code-block flex-1 break-all">{k}</code>
                      </li>
                    ))}
                    {changed.map((c) => (
                      <li key={`c-${c.key}`} className="flex items-start gap-2 px-3 py-2 rounded bg-sky-aqua/10">
                        <Pencil className="h-4 w-4 mt-0.5 text-primary" />
                        <code className="font-code-block text-code-block flex-1 break-all">
                          <span className="opacity-60 mr-2">{c.key}</span>
                          <span className="text-error mr-1">{c.left}</span>
                          <span className="text-tertiary-container">→ {c.right}</span>
                        </code>
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            </>
          );
        })()}
      </Card>
    </div>
  );
}
