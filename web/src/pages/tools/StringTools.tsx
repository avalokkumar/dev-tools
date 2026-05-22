import { useState } from "react";
import { Scissors, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { strCase, strDiff, strStats, strSortUnique, strReplace } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Textarea, Input, Select, Toggle } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("string")!;

type TabId = "case" | "stats" | "sort" | "replace" | "diff";

export default function StringPage() {
  const [tab, setTab] = useState<TabId>("case");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "case", label: "Case" },
          { id: "stats", label: "Stats" },
          { id: "sort", label: "Sort/Unique" },
          { id: "replace", label: "Replace" },
          { id: "diff", label: "Line diff" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as TabId)}
        className="mb-md"
      />
      {tab === "case" && <CaseTab />}
      {tab === "stats" && <StatsTab />}
      {tab === "sort" && <SortTab />}
      {tab === "replace" && <ReplaceTab />}
      {tab === "diff" && <DiffTab />}
    </ToolPage>
  );
}

function CaseTab() {
  const [input, setInput] = useState("HelloWorld example_string");
  const [mode, setMode] = useState("camel");
  const op = useOperation((args: { input: string; mode: string }) => strCase(args.input, args.mode));
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Input" icon={<Scissors className="h-5 w-5" />} />
        <Textarea rows={6} value={input} onChange={(e) => setInput(e.target.value)} />
        <Select
          label="Case"
          className="mt-md"
          value={mode}
          onChange={(e) => setMode(e.target.value)}
          options={[
            { value: "camel", label: "camelCase" },
            { value: "pascal", label: "PascalCase" },
            { value: "snake", label: "snake_case" },
            { value: "kebab", label: "kebab-case" },
            { value: "constant", label: "CONSTANT_CASE" },
            { value: "title", label: "Title Case" },
            { value: "upper", label: "UPPERCASE" },
            { value: "lower", label: "lowercase" },
          ]}
        />
        <Button onClick={() => op.run({ input, mode })} loading={op.loading} fullWidth className="mt-md">
          <Play className="h-4 w-4" /> Convert
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title="Output" />
        <div className="p-md">
          {op.error ? <ErrorBanner error={op.error} /> : !op.data ? <EmptyState title="Output appears here" /> : <CodeBlock code={op.data.output} />}
        </div>
      </Card>
    </div>
  );
}

function StatsTab() {
  const [input, setInput] = useState("DevForge\nis\na local-first developer toolkit.");
  const op = useOperation((s: string) => strStats(s));
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Input" />
        <Textarea rows={10} value={input} onChange={(e) => setInput(e.target.value)} />
        <Button onClick={() => op.run(input)} loading={op.loading} fullWidth className="mt-md">
          <Play className="h-4 w-4" /> Analyze
        </Button>
      </Card>
      <Card>
        <CardHeader title="Stats" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Stats appear here" />
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
            {Object.entries(op.data).map(([k, v]) => (
              <div key={k} className="px-3 py-2 bg-surface-container-low rounded">
                <div className="font-data-label text-data-label uppercase text-on-surface-variant">{k}</div>
                <div className="font-display-brand text-on-surface text-xl tabular-nums">{v as number}</div>
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}

function SortTab() {
  const [input, setInput] = useState("banana\napple\ncherry\napple");
  const [reverse, setReverse] = useState(false);
  const [numeric, setNumeric] = useState(false);
  const [uniqueOnly, setUniqueOnly] = useState(true);
  const op = useOperation((args: { input: string; reverse: boolean; numeric: boolean; uniqueOnly: boolean }) =>
    strSortUnique(args.input, args.reverse, args.numeric, args.uniqueOnly),
  );
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Lines" />
        <Textarea rows={10} value={input} onChange={(e) => setInput(e.target.value)} />
        <div className="grid grid-cols-3 gap-md mt-md">
          <Toggle checked={reverse} onChange={setReverse} label="Reverse" />
          <Toggle checked={numeric} onChange={setNumeric} label="Numeric" />
          <Toggle checked={uniqueOnly} onChange={setUniqueOnly} label="Unique" />
        </div>
        <Button onClick={() => op.run({ input, reverse, numeric, uniqueOnly })} loading={op.loading} fullWidth className="mt-md">
          <Play className="h-4 w-4" /> Sort
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title="Output" />
        <div className="p-md">
          {op.error ? <ErrorBanner error={op.error} /> : !op.data ? <EmptyState title="Output appears here" /> : <CodeBlock code={op.data.output} />}
        </div>
      </Card>
    </div>
  );
}

function ReplaceTab() {
  const [input, setInput] = useState("foo bar foo baz");
  const [pattern, setPattern] = useState("foo");
  const [replace, setReplace] = useState("FOO");
  const [regex, setRegex] = useState(false);
  const [ignoreCase, setIgnoreCase] = useState(false);
  const op = useOperation((args: { input: string; pattern: string; replace: string; regex: boolean; ignoreCase: boolean }) =>
    strReplace(args.input, args.pattern, args.replace, args.regex, args.ignoreCase),
  );
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Input" />
        <Textarea rows={6} value={input} onChange={(e) => setInput(e.target.value)} />
        <div className="grid grid-cols-2 gap-md mt-md">
          <Input label="Pattern" value={pattern} onChange={(e) => setPattern(e.target.value)} />
          <Input label="Replace" value={replace} onChange={(e) => setReplace(e.target.value)} />
        </div>
        <div className="grid grid-cols-2 gap-md mt-md">
          <Toggle checked={regex} onChange={setRegex} label="Treat pattern as regex" />
          <Toggle checked={ignoreCase} onChange={setIgnoreCase} label="Case-insensitive" />
        </div>
        <Button
          onClick={() => op.run({ input, pattern, replace, regex, ignoreCase })}
          loading={op.loading}
          fullWidth
          className="mt-md"
        >
          <Play className="h-4 w-4" /> Replace
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title="Output" />
        <div className="p-md flex flex-col gap-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Output appears here" />
          ) : (
            <>
              <div className="font-data-label text-data-label uppercase text-on-surface-variant">
                {op.data.replacements} replacement{op.data.replacements === 1 ? "" : "s"}
              </div>
              <CodeBlock code={op.data.output} />
            </>
          )}
        </div>
      </Card>
    </div>
  );
}

function DiffTab() {
  const [left, setLeft] = useState("alpha\nbeta\ngamma");
  const [right, setRight] = useState("alpha\nbeta-2\ngamma\ndelta");
  const [ignoreWhitespace, setIgnoreWhitespace] = useState(false);
  const [ignoreCase, setIgnoreCase] = useState(false);
  const op = useOperation((args: { left: string; right: string; ignoreWhitespace: boolean; ignoreCase: boolean }) =>
    strDiff(args.left, args.right, args.ignoreWhitespace, args.ignoreCase),
  );
  return (
    <div className="grid grid-cols-1 gap-md">
      <Card>
        <CardHeader title="Lines" />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-md">
          <Textarea label="Left" rows={8} value={left} onChange={(e) => setLeft(e.target.value)} />
          <Textarea label="Right" rows={8} value={right} onChange={(e) => setRight(e.target.value)} />
        </div>
        <div className="grid grid-cols-2 gap-md mt-md">
          <Toggle checked={ignoreWhitespace} onChange={setIgnoreWhitespace} label="Ignore whitespace" />
          <Toggle checked={ignoreCase} onChange={setIgnoreCase} label="Ignore case" />
        </div>
        <Button onClick={() => op.run({ left, right, ignoreWhitespace, ignoreCase })} loading={op.loading} className="mt-md">
          <Play className="h-4 w-4" /> Diff
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title="Hunks" />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Output appears here" />
          ) : (
            <CodeBlock code={JSON.stringify(op.data, null, 2)} language="json" />
          )}
        </div>
      </Card>
    </div>
  );
}
