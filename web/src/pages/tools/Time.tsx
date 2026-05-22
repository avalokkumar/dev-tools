import { useState } from "react";
import { Play, Timer } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { timeConvert, timeRelative, timeDuration } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, Select } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("time")!;

export default function TimePage() {
  const [tab, setTab] = useState<"convert" | "relative" | "duration">("convert");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "convert", label: "Convert" },
          { id: "relative", label: "Relative" },
          { id: "duration", label: "Duration" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "convert" && <ConvertTab />}
      {tab === "relative" && <RelativeTab />}
      {tab === "duration" && <DurationTab />}
    </ToolPage>
  );
}

function ConvertTab() {
  const [input, setInput] = useState(String(Math.floor(Date.now() / 1000)));
  const [inputFormat, setInputFormat] = useState("auto");
  const [tz, setTz] = useState("UTC");
  const op = useOperation((args: { input: string; inputFormat: string; tz: string }) =>
    timeConvert(args.input, args.inputFormat, args.tz),
  );

  const formats = [
    { value: "auto", label: "Auto-detect" },
    { value: "epoch_s", label: "Unix seconds" },
    { value: "epoch_ms", label: "Unix milliseconds" },
    { value: "epoch_us", label: "Unix microseconds" },
    { value: "epoch_ns", label: "Unix nanoseconds" },
    { value: "rfc3339", label: "RFC 3339" },
    { value: "iso8601", label: "ISO 8601" },
  ];

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Input" icon={<Timer className="h-5 w-5" />} />
        <div className="flex flex-col gap-md">
          <Input label="Time" value={input} onChange={(e) => setInput(e.target.value)} />
          <Select
            label="Input format"
            value={inputFormat}
            onChange={(e) => setInputFormat(e.target.value)}
            options={formats}
          />
          <Input label="Timezone" value={tz} onChange={(e) => setTz(e.target.value)} placeholder="UTC, America/Los_Angeles" />
          <Button onClick={() => op.run({ input, inputFormat, tz })} loading={op.loading} fullWidth>
            <Play className="h-4 w-4" /> Convert
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader title="All representations" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Output appears here" />
        ) : (
          <div className="flex flex-col gap-2">
            {[
              ["RFC 3339", op.data.rfc3339],
              ["UTC", op.data.utc],
              ["Local", op.data.local],
              ["Epoch (s)", String(op.data.epochS)],
              ["Epoch (ms)", String(op.data.epochMS)],
              ["Epoch (µs)", String(op.data.epochUS)],
              ["Epoch (ns)", String(op.data.epochNS)],
            ].map(([k, v]) => (
              <div key={k} className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded">
                <span className="font-data-label text-data-label uppercase text-on-surface-variant w-24 shrink-0">
                  {k}
                </span>
                <code className="font-code-block text-code-block flex-1 truncate">{v}</code>
                <CopyButton text={String(v)} />
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}

function RelativeTab() {
  const [from, setFrom] = useState(new Date().toISOString());
  const [to, setTo] = useState("");
  const op = useOperation((args: { from: string; to: string }) => timeRelative(args.from, args.to));
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Time pair" />
        <div className="flex flex-col gap-md">
          <Input label="From" value={from} onChange={(e) => setFrom(e.target.value)} />
          <Input label="To (blank = now)" value={to} onChange={(e) => setTo(e.target.value)} />
          <Button onClick={() => op.run({ from, to })} loading={op.loading} fullWidth>
            <Play className="h-4 w-4" /> Compute
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader title="Result" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Output appears here" />
        ) : (
          <div className="flex flex-col gap-md">
            <div className="font-display-brand text-display-brand text-on-surface">
              {op.data.phrase}
            </div>
            <div className="font-data-label text-data-label uppercase text-on-surface-variant">
              {op.data.seconds.toLocaleString()} seconds
            </div>
          </div>
        )}
      </Card>
    </div>
  );
}

function DurationTab() {
  const [input, setInput] = useState("1h30m45s");
  const op = useOperation((args: { input: string }) => timeDuration(args.input));
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Duration string" />
        <div className="flex flex-col gap-md">
          <Input
            label="Input"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            hint="e.g. 1h30m, 90s, 2d4h, ISO8601 PT1H30M"
          />
          <Button onClick={() => op.run({ input })} loading={op.loading} fullWidth>
            <Play className="h-4 w-4" /> Parse
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader title="Components" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Output appears here" />
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
            <Stat label="Hours" value={op.data.hours} />
            <Stat label="Minutes" value={op.data.minutes} />
            <Stat label="Seconds" value={op.data.seconds} />
            <Stat label="Total s" value={op.data.totalSeconds} />
          </div>
        )}
      </Card>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <div className="px-3 py-3 bg-surface-container-low rounded">
      <div className="font-data-label text-data-label uppercase text-on-surface-variant">{label}</div>
      <div className="font-display-brand text-on-surface text-2xl tabular-nums">{value}</div>
    </div>
  );
}
