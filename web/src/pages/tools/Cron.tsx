import { useEffect, useState } from "react";
import { Clock, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { cronParse, cronNext } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, NumberInput, Select } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("cron")!;

export default function CronPage() {
  const [expr, setExpr] = useState("*/5 * * * *");
  const [flavor, setFlavor] = useState("unix");
  const [n, setN] = useState(5);
  const [tz, setTz] = useState("UTC");
  const parse = useOperation((args: { expression: string; flavor: string }) =>
    cronParse(args.expression, args.flavor),
  );
  const next = useOperation((args: { expression: string; n: number; tz: string }) =>
    cronNext(args.expression, args.n, args.tz),
  );

  useEffect(() => {
    if (expr.trim()) parse.run({ expression: expr, flavor });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [expr, flavor]);

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-[400px_1fr] gap-md">
        <Card>
          <CardHeader title="Expression" icon={<Clock className="h-5 w-5" />} />
          <div className="flex flex-col gap-md">
            <Input label="Cron expression" value={expr} onChange={(e) => setExpr(e.target.value)} />
            <div className="grid grid-cols-2 gap-md">
              <Select
                label="Flavor"
                value={flavor}
                onChange={(e) => setFlavor(e.target.value)}
                options={[
                  { value: "unix", label: "Unix (5 fields)" },
                  { value: "quartz", label: "Quartz (6-7 fields)" },
                  { value: "aws", label: "AWS EventBridge" },
                ]}
              />
              <NumberInput label="Next runs" value={n} onChange={setN} min={1} max={50} />
            </div>
            <Input label="Timezone" value={tz} onChange={(e) => setTz(e.target.value)} />
            <Button
              onClick={() => next.run({ expression: expr, n, tz })}
              loading={next.loading}
              fullWidth
            >
              <Play className="h-4 w-4" /> Compute next runs
            </Button>
          </div>
        </Card>
        <div className="flex flex-col gap-md">
          <Card>
            <CardHeader title="Description" />
            {parse.error ? (
              <ErrorBanner error={parse.error} />
            ) : !parse.data ? (
              <EmptyState title="Type an expression to parse" />
            ) : (
              <div className="flex flex-col gap-md">
                <p className="font-body-md text-body-md text-on-surface">{parse.data.description}</p>
                <ul className="grid grid-cols-2 md:grid-cols-3 gap-2">
                  {parse.data.fields.map((f) => (
                    <li key={f.name} className="px-3 py-2 bg-surface-container-low rounded">
                      <div className="font-data-label text-data-label uppercase text-on-surface-variant">
                        {f.name}
                      </div>
                      <code className="font-code-block text-code-block text-on-surface">{f.value}</code>
                    </li>
                  ))}
                </ul>
                <Diagnostics items={parse.data.diagnostics} />
              </div>
            )}
          </Card>
          <Card>
            <CardHeader title="Next runs" />
            {next.error ? (
              <ErrorBanner error={next.error} />
            ) : !next.data ? (
              <EmptyState title="Press 'Compute next runs'" />
            ) : (
              <ol className="flex flex-col gap-1">
                {next.data.runs.map((r, i) => (
                  <li
                    key={i}
                    className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded"
                  >
                    <span className="font-data-label text-data-label uppercase text-on-surface-variant w-8">
                      #{i + 1}
                    </span>
                    <code className="font-code-block text-code-block text-on-surface flex-1">{r}</code>
                  </li>
                ))}
              </ol>
            )}
          </Card>
        </div>
      </div>
    </ToolPage>
  );
}
