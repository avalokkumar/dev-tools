import { useState } from "react";
import { Play, RefreshCw } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { uuidGenerate, uuidHash } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, Select, Toggle, NumberInput, Field } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("uuid")!;

export default function UuidPage() {
  const [tab, setTab] = useState<"generate" | "hash">("generate");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "generate", label: "Generate" },
          { id: "hash", label: "Hash digest" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "generate" ? <GenerateTab /> : <HashTab />}
    </ToolPage>
  );
}

function GenerateTab() {
  const [version, setVersion] = useState<"4" | "7">("4");
  const [count, setCount] = useState(5);
  const [format, setFormat] = useState<"std" | "compact" | "urn">("std");
  const op = useOperation(uuidGenerate);

  function generate() {
    op.run({ version: Number(version) as 4 | 7, count, format });
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[360px_1fr] gap-md">
      <Card>
        <CardHeader title="Options" />
        <div className="flex flex-col gap-md">
          <Select
            label="Version"
            value={version}
            onChange={(e) => setVersion(e.target.value as "4" | "7")}
            options={[
              { value: "4", label: "v4 — random" },
              { value: "7", label: "v7 — time-ordered" },
            ]}
          />
          <NumberInput label="Count" value={count} onChange={setCount} min={1} max={1024} />
          <Select
            label="Format"
            value={format}
            onChange={(e) => setFormat(e.target.value as typeof format)}
            options={[
              { value: "std", label: "Standard (8-4-4-4-12)" },
              { value: "compact", label: "Compact (no hyphens)" },
              { value: "urn", label: "URN (urn:uuid:…)" },
            ]}
          />
          <Button onClick={generate} loading={op.loading} fullWidth>
            <Play className="h-4 w-4" />
            Generate
          </Button>
        </div>
      </Card>
      <Card padded={false}>
        <div className="flex items-center justify-between p-md border-b border-outline/10">
          <h3 className="font-body-md text-body-md font-medium text-on-surface">
            Generated UUIDs {op.data && <span className="text-on-surface-variant font-data-label text-data-label uppercase ml-2">{op.data.values.length} ids</span>}
          </h3>
          {op.data && <CopyButton text={op.data.values.join("\n")} label="Copy all" />}
        </div>
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState
              title="Press Generate"
              description="Configure your options on the left, then click Generate to produce UUIDs."
              icon={<RefreshCw className="h-6 w-6" />}
            />
          ) : (
            <>
              <ul className="flex flex-col gap-1">
                {op.data.values.map((v) => (
                  <li
                    key={v}
                    className="flex items-center justify-between px-3 py-2 bg-surface-container-low rounded font-code-block text-code-block group"
                  >
                    <code className="truncate">{v}</code>
                    <CopyButton text={v} label="Copy" />
                  </li>
                ))}
              </ul>
              <Diagnostics items={op.data.diagnostics} className="mt-md" />
            </>
          )}
        </div>
      </Card>
    </div>
  );
}

function HashTab() {
  const [input, setInput] = useState("");
  const [encoding, setEncoding] = useState<"hex" | "base64">("hex");
  const [algos, setAlgos] = useState<("md5" | "sha1" | "sha256" | "sha512")[]>([
    "sha256",
  ]);
  const op = useOperation(uuidHash);

  function toggleAlgo(a: "md5" | "sha1" | "sha256" | "sha512") {
    setAlgos((prev) => (prev.includes(a) ? prev.filter((x) => x !== a) : [...prev, a]));
  }

  function run() {
    op.run({ input, encoding, algos });
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[1fr_1fr] gap-md">
      <Card>
        <CardHeader title="Input" />
        <Input
          label="Plaintext"
          placeholder="Enter text to hash…"
          value={input}
          onChange={(e) => setInput(e.target.value)}
        />
        <Field label="Algorithms" className="mt-md">
          <div className="grid grid-cols-2 gap-2">
            {(["md5", "sha1", "sha256", "sha512"] as const).map((a) => (
              <Toggle
                key={a}
                checked={algos.includes(a)}
                onChange={() => toggleAlgo(a)}
                label={a.toUpperCase()}
              />
            ))}
          </div>
        </Field>
        <Select
          label="Encoding"
          className="mt-md"
          value={encoding}
          onChange={(e) => setEncoding(e.target.value as typeof encoding)}
          options={[
            { value: "hex", label: "Hex" },
            { value: "base64", label: "Base64" },
          ]}
        />
        <Button onClick={run} loading={op.loading} fullWidth className="mt-md" disabled={!input || algos.length === 0}>
          <Play className="h-4 w-4" />
          Compute Digests
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title="Digests" />
        <div className="p-md flex flex-col gap-2">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="No digest yet" description="Enter input and select algorithms." />
          ) : (
            Object.entries(op.data.digests).map(([k, v]) => (
              <div key={k} className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded">
                <span className="font-data-label text-data-label uppercase text-on-surface-variant w-16 shrink-0">
                  {k}
                </span>
                <code className="font-code-block text-code-block flex-1 truncate text-on-surface">
                  {v}
                </code>
                <CopyButton text={v} />
              </div>
            ))
          )}
        </div>
      </Card>
    </div>
  );
}
