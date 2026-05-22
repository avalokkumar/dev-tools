import { useState } from "react";
import { Play, Braces, CheckCircle2 } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { jsonFormat, jsonValidate } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Textarea, NumberInput, Toggle } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("json")!;

const SAMPLE = `{"name":"DevForge","tools":75,"local":true,"surfaces":["cli","web","mcp"]}`;

export default function JsonPage() {
  const [tab, setTab] = useState<"format" | "validate">("format");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "format", label: "Format" },
          { id: "validate", label: "Validate (Schema)" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "format" ? <FormatTab /> : <ValidateTab />}
    </ToolPage>
  );
}

function FormatTab() {
  const [input, setInput] = useState(SAMPLE);
  const [indent, setIndent] = useState(2);
  const [sortKeys, setSortKeys] = useState(false);
  const [trailingNewline, setTrailingNewline] = useState(false);
  const op = useOperation((args: { input: string; indent: number; sortKeys: boolean; trailingNewline: boolean }) =>
    jsonFormat(args.input, args.indent, args.sortKeys, args.trailingNewline),
  );

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Input JSON" icon={<Braces className="h-5 w-5" />} />
        <Textarea
          rows={14}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder='{"key":"value"}'
        />
        <div className="grid grid-cols-3 gap-md mt-md">
          <NumberInput label="Indent" value={indent} onChange={setIndent} min={0} max={8} />
          <Toggle checked={sortKeys} onChange={setSortKeys} label="Sort keys" />
          <Toggle checked={trailingNewline} onChange={setTrailingNewline} label="Trailing \\n" />
        </div>
        <Button
          onClick={() => op.run({ input, indent, sortKeys, trailingNewline })}
          loading={op.loading}
          fullWidth
          className="mt-md"
          disabled={!input.trim()}
        >
          <Play className="h-4 w-4" /> Format
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title="Output" />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Press Format" />
          ) : (
            <>
              <CodeBlock code={op.data.output} language="json" download={{ filename: "formatted.json" }} />
              <Diagnostics items={op.data.diagnostics} className="mt-md" />
            </>
          )}
        </div>
      </Card>
    </div>
  );
}

function ValidateTab() {
  const [input, setInput] = useState(SAMPLE);
  const [schema, setSchema] = useState("");
  const op = useOperation((args: { input: string; schema?: string }) => jsonValidate(args.input, args.schema));

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="JSON" />
        <Textarea rows={8} value={input} onChange={(e) => setInput(e.target.value)} />
        <div className="mt-md">
          <Textarea
            label="JSON Schema (optional)"
            rows={6}
            value={schema}
            onChange={(e) => setSchema(e.target.value)}
            placeholder='{"type":"object"}'
          />
        </div>
        <Button
          onClick={() => op.run({ input, schema: schema || undefined })}
          loading={op.loading}
          fullWidth
          className="mt-md"
          disabled={!input.trim()}
        >
          <Play className="h-4 w-4" /> Validate
        </Button>
      </Card>
      <Card>
        <CardHeader title="Result" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Validation result appears here" />
        ) : (
          <div className="flex flex-col gap-md">
            <div
              className={`flex items-center gap-3 p-md rounded-lg ${
                op.data.valid
                  ? "bg-tertiary-fixed text-on-tertiary-fixed"
                  : "bg-error/10 text-error"
              }`}
            >
              <CheckCircle2 className="h-6 w-6" />
              <span className="font-body-md font-semibold">
                {op.data.valid ? "Valid JSON" : "Invalid JSON"}
              </span>
            </div>
            <Diagnostics items={op.data.diagnostics} />
          </div>
        )}
      </Card>
    </div>
  );
}
