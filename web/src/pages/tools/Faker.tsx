import { useEffect, useState } from "react";
import { Play, Plus, Trash2, Drama } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { fakerGenerate, fakerKinds, type FakerField } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, Select, NumberInput } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("faker")!;

const DEFAULT_FIELDS: FakerField[] = [
  { name: "id", kind: "uuid" },
  { name: "name", kind: "name" },
  { name: "email", kind: "email" },
  { name: "created_at", kind: "date" },
];

export default function FakerPage() {
  const [fields, setFields] = useState<FakerField[]>(DEFAULT_FIELDS);
  const [count, setCount] = useState(10);
  const [format, setFormat] = useState<"json" | "csv" | "sql">("json");
  const [table, setTable] = useState("data");
  const [seed, setSeed] = useState(0);
  const [kinds, setKinds] = useState<string[]>([]);
  const op = useOperation(fakerGenerate);

  useEffect(() => {
    fakerKinds()
      .then((r) => setKinds(r.kinds.map((k) => k.name)))
      .catch(() => setKinds([]));
  }, []);

  function update(i: number, patch: Partial<FakerField>) {
    setFields((f) => f.map((x, idx) => (idx === i ? { ...x, ...patch } : x)));
  }
  function remove(i: number) {
    setFields((f) => f.filter((_, idx) => idx !== i));
  }
  function add() {
    setFields((f) => [...f, { name: `field_${f.length + 1}`, kind: "uuid" }]);
  }
  function run() {
    op.run({ spec: { fields }, count, format, table, seed: seed || 0 });
  }

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-[480px_1fr] gap-md">
        <Card>
          <CardHeader title="Schema" trailing={
            <Button size="sm" variant="outline" onClick={add}>
              <Plus className="h-4 w-4" /> Field
            </Button>
          } />
          <div className="flex flex-col gap-2 mb-md">
            {fields.map((f, i) => (
              <div key={i} className="grid grid-cols-[1fr_1fr_auto] gap-2 items-end">
                <Input
                  label={i === 0 ? "Name" : undefined}
                  value={f.name}
                  onChange={(e) => update(i, { name: e.target.value })}
                  placeholder="field"
                />
                <Select
                  label={i === 0 ? "Kind" : undefined}
                  value={f.kind}
                  onChange={(e) => update(i, { kind: e.target.value })}
                  options={(kinds.length ? kinds : ["uuid","name","email","int","float","bool","date"]).map((k) => ({
                    value: k,
                    label: k,
                  }))}
                />
                <Button variant="ghost" size="sm" onClick={() => remove(i)} aria-label="Remove">
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
          <div className="grid grid-cols-2 gap-md mb-md">
            <NumberInput label="Count" value={count} onChange={setCount} min={1} max={10000} />
            <Select
              label="Format"
              value={format}
              onChange={(e) => setFormat(e.target.value as typeof format)}
              options={[
                { value: "json", label: "JSON" },
                { value: "csv", label: "CSV" },
                { value: "sql", label: "SQL INSERT" },
              ]}
            />
            <Input label="Table" value={table} onChange={(e) => setTable(e.target.value)} />
            <NumberInput label="Seed (0 = random)" value={seed} onChange={setSeed} min={0} />
          </div>
          <Button onClick={run} loading={op.loading} fullWidth>
            <Play className="h-4 w-4" /> Generate
          </Button>
        </Card>
        <Card padded={false}>
          <CardHeader title="Output" />
          <div className="p-md">
            {op.error ? (
              <ErrorBanner error={op.error} />
            ) : !op.data ? (
              <EmptyState
                title="Configure schema and generate"
                icon={<Drama className="h-6 w-6" />}
              />
            ) : (
              <>
                <CodeBlock
                  code={op.data.output}
                  language={format}
                  download={{
                    filename: `${table}.${format === "sql" ? "sql" : format}`,
                    mime: format === "json" ? "application/json" : "text/plain",
                  }}
                />
                <Diagnostics items={op.data.diagnostics} className="mt-md" />
              </>
            )}
          </div>
        </Card>
      </div>
    </ToolPage>
  );
}
