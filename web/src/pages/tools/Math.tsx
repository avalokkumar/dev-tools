import { useState } from "react";
import { Calculator, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { mathEval, mathUnit } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, NumberInput } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("math")!;

export default function MathPage() {
  const [tab, setTab] = useState<"calc" | "unit">("calc");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[{ id: "calc", label: "Calculator" }, { id: "unit", label: "Unit converter" }]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "calc" ? <CalcTab /> : <UnitTab />}
    </ToolPage>
  );
}

function CalcTab() {
  const [expr, setExpr] = useState("(2 + 3) * sin(pi/4)");
  const op = useOperation((e: string) => mathEval(e));
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Expression" icon={<Calculator className="h-5 w-5" />} />
        <div className="flex flex-col gap-md">
          <Input label="Expression" value={expr} onChange={(e) => setExpr(e.target.value)} />
          <Button onClick={() => op.run(expr)} loading={op.loading} fullWidth disabled={!expr.trim()}>
            <Play className="h-4 w-4" /> Evaluate
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader title="Result" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Result appears here" />
        ) : (
          <div className="flex flex-col gap-md">
            <div className="flex items-center gap-2 px-3 py-3 bg-surface-container-low rounded">
              <code className="font-code-block text-code-block text-on-surface flex-1 truncate">
                {op.data.display ?? op.data.value}
              </code>
              <CopyButton text={String(op.data.display ?? op.data.value)} />
            </div>
            <Diagnostics items={op.data.diagnostics} />
          </div>
        )}
      </Card>
    </div>
  );
}

function UnitTab() {
  const [value, setValue] = useState(100);
  const [from, setFrom] = useState("km");
  const [to, setTo] = useState("mi");
  const op = useOperation((args: { value: number; from: string; to: string }) =>
    mathUnit(args.value, args.from, args.to),
  );
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Convert" />
        <div className="flex flex-col gap-md">
          <NumberInput label="Value" value={value} onChange={setValue} />
          <div className="grid grid-cols-2 gap-md">
            <Input label="From unit" value={from} onChange={(e) => setFrom(e.target.value)} />
            <Input label="To unit" value={to} onChange={(e) => setTo(e.target.value)} />
          </div>
          <Button onClick={() => op.run({ value, from, to })} loading={op.loading} fullWidth>
            <Play className="h-4 w-4" /> Convert
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader title="Result" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Result appears here" />
        ) : (
          <div className="flex flex-col gap-md">
            <div className="font-display-brand text-display-brand text-on-surface tabular-nums">
              {op.data.value}
              <span className="font-body-md text-body-md text-on-surface-variant ml-2">{op.data.to}</span>
            </div>
            <code className="font-code-block text-code-block text-on-surface-variant">
              {value} {op.data.from} = {op.data.value} {op.data.to}
            </code>
          </div>
        )}
      </Card>
    </div>
  );
}
