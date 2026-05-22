import { useState } from "react";
import { GitBranch, Play, CheckCircle2 } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { gitPatch, gitCommitFormat, gitIgnoreGen } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Input, Textarea, NumberInput, Toggle } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("git")!;

const TEMPLATES = [
  "node",
  "python",
  "go",
  "java",
  "rust",
  "ruby",
  "php",
  "macos",
  "windows",
  "linux",
  "vscode",
  "intellij",
  "vim",
  "terraform",
  "docker",
];

export default function GitPage() {
  const [tab, setTab] = useState<"patch" | "commit" | "gitignore">("patch");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "patch", label: "Patch" },
          { id: "commit", label: "Commit format" },
          { id: "gitignore", label: ".gitignore generator" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "patch" && <PatchTab />}
      {tab === "commit" && <CommitTab />}
      {tab === "gitignore" && <IgnoreTab />}
    </ToolPage>
  );
}

function PatchTab() {
  const [left, setLeft] = useState("hello\nworld\n");
  const [right, setRight] = useState("hello\nDevForge\nworld\n");
  const [leftPath, setLeftPath] = useState("a.txt");
  const [rightPath, setRightPath] = useState("a.txt");
  const [contextLines, setContextLines] = useState(3);
  const op = useOperation(gitPatch);
  return (
    <div className="grid grid-cols-1 gap-md">
      <Card>
        <CardHeader title="Inputs" icon={<GitBranch className="h-5 w-5" />} />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-md">
          <Textarea label="Left" rows={8} value={left} onChange={(e) => setLeft(e.target.value)} />
          <Textarea label="Right" rows={8} value={right} onChange={(e) => setRight(e.target.value)} />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-md mt-md">
          <Input label="Left path" value={leftPath} onChange={(e) => setLeftPath(e.target.value)} />
          <Input label="Right path" value={rightPath} onChange={(e) => setRightPath(e.target.value)} />
          <NumberInput label="Context lines" value={contextLines} onChange={setContextLines} min={0} max={20} />
        </div>
        <Button
          onClick={() => op.run({ left, right, leftPath, rightPath, contextLines })}
          loading={op.loading}
          className="mt-md"
        >
          <Play className="h-4 w-4" /> Generate patch
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title="Patch" />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Unified diff appears here" />
          ) : (
            <CodeBlock code={op.data.output} language="diff" download={{ filename: "changes.patch" }} />
          )}
        </div>
      </Card>
    </div>
  );
}

function CommitTab() {
  const [msg, setMsg] = useState("feat(api): add new endpoint for users\n\nAdds GET /users with pagination support.");
  const op = useOperation((s: string) => gitCommitFormat(s));
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Commit message" />
        <Textarea rows={10} value={msg} onChange={(e) => setMsg(e.target.value)} />
        <Button onClick={() => op.run(msg)} loading={op.loading} fullWidth className="mt-md" disabled={!msg.trim()}>
          <Play className="h-4 w-4" /> Validate
        </Button>
      </Card>
      <Card>
        <CardHeader title="Parsed (Conventional Commits)" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Output appears here" />
        ) : (
          <div className="flex flex-col gap-md">
            <div className="flex items-center gap-2 flex-wrap">
              <Badge tone={op.data.valid ? "success" : "error"}>
                {op.data.valid ? "Valid" : "Invalid"}
              </Badge>
              {op.data.breaking && <Badge tone="error">BREAKING</Badge>}
              {op.data.type && <Badge>{op.data.type}</Badge>}
              {op.data.scope && <Badge>{op.data.scope}</Badge>}
            </div>
            <div>
              <span className="font-data-label text-data-label uppercase text-on-surface-variant">
                Subject
              </span>
              <div className="font-body-md text-body-md text-on-surface mt-1">{op.data.subject}</div>
            </div>
            {op.data.body && (
              <div>
                <span className="font-data-label text-data-label uppercase text-on-surface-variant">Body</span>
                <pre className="font-code-block text-code-block bg-surface-container-low rounded p-3 mt-1 whitespace-pre-wrap">
                  {op.data.body}
                </pre>
              </div>
            )}
            {op.data.footer && (
              <div>
                <span className="font-data-label text-data-label uppercase text-on-surface-variant">Footer</span>
                <pre className="font-code-block text-code-block bg-surface-container-low rounded p-3 mt-1 whitespace-pre-wrap">
                  {op.data.footer}
                </pre>
              </div>
            )}
            <Diagnostics items={op.data.diagnostics} />
          </div>
        )}
      </Card>
    </div>
  );
}

function IgnoreTab() {
  const [picked, setPicked] = useState<Record<string, boolean>>({ node: true });
  const op = useOperation((tpls: string[]) => gitIgnoreGen(tpls));
  function toggle(t: string) {
    setPicked((p) => ({ ...p, [t]: !p[t] }));
  }
  const selected = TEMPLATES.filter((t) => picked[t]);
  return (
    <div className="grid grid-cols-1 lg:grid-cols-[420px_1fr] gap-md">
      <Card>
        <CardHeader title="Templates" trailing={<Badge>{selected.length}</Badge>} />
        <div className="grid grid-cols-2 gap-2">
          {TEMPLATES.map((t) => (
            <Toggle key={t} checked={!!picked[t]} onChange={() => toggle(t)} label={t} />
          ))}
        </div>
        <Button
          onClick={() => op.run(selected)}
          loading={op.loading}
          className="mt-md"
          fullWidth
          disabled={selected.length === 0}
        >
          <Play className="h-4 w-4" /> Generate
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader
          title=".gitignore"
          trailing={op.data && <CheckCircle2 className="h-4 w-4 text-tertiary-container" />}
        />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Pick templates and generate" />
          ) : (
            <CodeBlock code={op.data.output} language="gitignore" download={{ filename: ".gitignore" }} />
          )}
        </div>
      </Card>
    </div>
  );
}
