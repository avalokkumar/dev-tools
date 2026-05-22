import { useState } from "react";
import { findTool } from "../../lib/catalog";
import { mdToHtml, mdTableFromCsv } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { FormatterShell } from "../../components/layout/FormatterShell";
import { Card, CardHeader } from "../../components/ui/Card";
import { Toggle, Select, Textarea } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";
import { Play, Eye } from "lucide-react";

const tool = findTool("markdown")!;
const SAMPLE_MD = `# DevForge\n\nA **local-first** developer toolkit.\n\n- 75 operations\n- CLI · Web · MCP\n\n\`\`\`go\nfmt.Println("hello")\n\`\`\`\n`;
const SAMPLE_CSV = `name,role\nAlice,Engineer\nBob,Designer\n`;

export default function MarkdownPage() {
  const [tab, setTab] = useState<"preview" | "csv">("preview");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "preview", label: "Markdown → HTML" },
          { id: "csv", label: "CSV → Table" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "preview" ? <PreviewTab /> : <CsvTableTab />}
    </ToolPage>
  );
}

function PreviewTab() {
  const [input, setInput] = useState(SAMPLE_MD);
  const [gfm, setGfm] = useState(true);
  const [unsafe, setUnsafe] = useState(false);
  const op = useOperation((args: { input: string; gfm: boolean; unsafe: boolean }) =>
    mdToHtml(args.input, args.gfm, args.unsafe),
  );

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Markdown" />
        <Textarea rows={14} value={input} onChange={(e) => setInput(e.target.value)} />
        <div className="grid grid-cols-2 gap-md mt-md">
          <Toggle checked={gfm} onChange={setGfm} label="GitHub Flavored" />
          <Toggle checked={unsafe} onChange={setUnsafe} label="Allow raw HTML" hint="Unsafe" />
        </div>
        <Button
          onClick={() => op.run({ input, gfm, unsafe })}
          loading={op.loading}
          fullWidth
          className="mt-md"
          disabled={!input.trim()}
        >
          <Play className="h-4 w-4" /> Render
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title="Preview" icon={<Eye className="h-5 w-5" />} />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Render to preview" />
          ) : (
            <>
              <div
                className="prose-devforge bg-surface-container-lowest border border-outline/10 rounded p-md max-h-[60vh] overflow-auto thin-scrollbar"
                // eslint-disable-next-line react/no-danger
                dangerouslySetInnerHTML={{ __html: op.data.output }}
              />
              <details className="mt-md">
                <summary className="font-data-label text-data-label uppercase text-on-surface-variant cursor-pointer">
                  View HTML
                </summary>
                <CodeBlock code={op.data.output} language="html" download={{ filename: "preview.html" }} />
              </details>
            </>
          )}
        </div>
      </Card>
    </div>
  );
}

function CsvTableTab() {
  const [alignment, setAlignment] = useState<"none" | "left" | "center" | "right">("none");
  const [delimiter, setDelimiter] = useState(",");
  const op = useOperation((args: { input: string; delimiter: string; alignment: typeof alignment }) =>
    mdTableFromCsv(args.input, args.delimiter, args.alignment),
  );
  return (
    <FormatterShell
      initial={SAMPLE_CSV}
      buildInput={(input) => ({ input, delimiter, alignment })}
      op={op}
      language="markdown"
      downloadFilename="table.md"
      buttonLabel="Build Table"
      options={
        <div className="grid grid-cols-2 gap-md">
          <Select
            label="Alignment"
            value={alignment}
            onChange={(e) => setAlignment(e.target.value as typeof alignment)}
            options={[
              { value: "none", label: "Default" },
              { value: "left", label: "Left" },
              { value: "center", label: "Center" },
              { value: "right", label: "Right" },
            ]}
          />
          <Select
            label="Delimiter"
            value={delimiter}
            onChange={(e) => setDelimiter(e.target.value)}
            options={[
              { value: ",", label: "Comma" },
              { value: ";", label: "Semicolon" },
              { value: "\t", label: "Tab" },
            ]}
          />
        </div>
      }
    />
  );
}
