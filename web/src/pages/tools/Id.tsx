import { useState } from "react";
import { Play, Tag } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { idUlid, idSlug } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, NumberInput, Toggle } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("id")!;

export default function IdPage() {
  const [tab, setTab] = useState<"ulid" | "slug">("ulid");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "ulid", label: "ULID" },
          { id: "slug", label: "Slug" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "ulid" ? <UlidTab /> : <SlugTab />}
    </ToolPage>
  );
}

function UlidTab() {
  const [count, setCount] = useState(5);
  const [lowercase, setLowercase] = useState(false);
  const op = useOperation((args: { count: number; lowercase: boolean }) =>
    idUlid(args.count, args.lowercase),
  );

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[360px_1fr] gap-md">
      <Card>
        <CardHeader title="Options" />
        <div className="flex flex-col gap-md">
          <NumberInput label="Count" value={count} onChange={setCount} min={1} max={1024} />
          <Toggle checked={lowercase} onChange={setLowercase} label="Lowercase output" />
          <Button
            onClick={() => op.run({ count, lowercase })}
            loading={op.loading}
            fullWidth
          >
            <Play className="h-4 w-4" /> Generate
          </Button>
        </div>
      </Card>
      <Card padded={false}>
        <CardHeader title="ULIDs" />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Press Generate" icon={<Tag className="h-6 w-6" />} />
          ) : (
            <ul className="flex flex-col gap-1">
              {op.data.values.map((v) => (
                <li
                  key={v}
                  className="flex items-center justify-between px-3 py-2 bg-surface-container-low rounded font-code-block text-code-block"
                >
                  <code className="truncate">{v}</code>
                  <CopyButton text={v} />
                </li>
              ))}
            </ul>
          )}
        </div>
      </Card>
    </div>
  );
}

function SlugTab() {
  const [input, setInput] = useState("Hello, World! — A New Beginning");
  const [maxLen, setMaxLen] = useState(60);
  const [locale, setLocale] = useState("en");
  const op = useOperation((args: { input: string; maxLen: number; locale: string }) =>
    idSlug(args.input, args.maxLen, args.locale),
  );

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[1fr_1fr] gap-md">
      <Card>
        <CardHeader title="Input" />
        <div className="flex flex-col gap-md">
          <Input
            label="Source text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
          />
          <div className="grid grid-cols-2 gap-md">
            <NumberInput label="Max length" value={maxLen} onChange={setMaxLen} min={4} max={255} />
            <Input label="Locale" value={locale} onChange={(e) => setLocale(e.target.value)} />
          </div>
          <Button
            onClick={() => op.run({ input, maxLen, locale })}
            loading={op.loading}
            fullWidth
            disabled={!input}
          >
            <Play className="h-4 w-4" /> Slugify
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader title="Slug" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Output appears here" />
        ) : (
          <div className="flex items-center justify-between px-3 py-3 bg-surface-container-low rounded font-code-block text-code-block">
            <code className="truncate">{op.data.output}</code>
            <CopyButton text={op.data.output} />
          </div>
        )}
      </Card>
    </div>
  );
}
