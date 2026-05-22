import { useEffect, useState } from "react";
import { Globe2, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { tzConvert, tzList } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Input, Select } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("timezone")!;

export default function TimezonePage() {
  const [zones, setZones] = useState<{ id: string; offset: string; abbr: string }[]>([]);
  function fmtOffset(s: number) {
    const sign = s < 0 ? "-" : "+";
    const abs = Math.abs(s);
    const h = Math.floor(abs / 3600).toString().padStart(2, "0");
    const m = Math.floor((abs % 3600) / 60).toString().padStart(2, "0");
    return `${sign}${h}:${m}`;
  }
  const [filter, setFilter] = useState("");
  const [time, setTime] = useState(new Date().toISOString().slice(0, 19));
  const [from, setFrom] = useState("UTC");
  const [to, setTo] = useState("America/Los_Angeles");
  const op = useOperation((args: { time: string; fromTZ: string; toTZ: string }) =>
    tzConvert(args.time, args.fromTZ, args.toTZ),
  );

  useEffect(() => {
    tzList(filter)
      .then((r) => setZones(r.map((z) => ({ id: z.name, abbr: z.abbrev, offset: fmtOffset(z.offsetSeconds) }))))
      .catch(() => setZones([]));
  }, [filter]);

  const opts = zones.length
    ? zones.map((z) => ({ value: z.id, label: `${z.id} (${z.abbr} ${z.offset})` }))
    : [
        { value: "UTC", label: "UTC" },
        { value: "America/Los_Angeles", label: "America/Los_Angeles" },
        { value: "Europe/Berlin", label: "Europe/Berlin" },
        { value: "Asia/Tokyo", label: "Asia/Tokyo" },
        { value: "Asia/Kolkata", label: "Asia/Kolkata" },
      ];

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-[1fr_1fr] gap-md">
        <Card>
          <CardHeader title="Convert" icon={<Globe2 className="h-5 w-5" />} />
          <div className="flex flex-col gap-md">
            <Input label="Time" value={time} onChange={(e) => setTime(e.target.value)} hint="ISO 8601 or YYYY-MM-DDTHH:MM:SS" />
            <Input label="Filter zones" value={filter} onChange={(e) => setFilter(e.target.value)} placeholder="america" />
            <div className="grid grid-cols-2 gap-md">
              <Select label="From" value={from} onChange={(e) => setFrom(e.target.value)} options={opts} />
              <Select label="To" value={to} onChange={(e) => setTo(e.target.value)} options={opts} />
            </div>
            <Button onClick={() => op.run({ time, fromTZ: from, toTZ: to })} loading={op.loading} fullWidth>
              <Play className="h-4 w-4" /> Convert
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
              <Row label={from} value={op.data.original} />
              <Row label={to} value={op.data.converted} />
              {op.data.dstNote && <Badge tone="warning">DST: {op.data.dstNote}</Badge>}
            </div>
          )}
        </Card>
      </div>
    </ToolPage>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center gap-2 px-3 py-3 bg-surface-container-low rounded">
      <span className="font-data-label text-data-label uppercase text-on-surface-variant w-44 shrink-0 truncate">
        {label}
      </span>
      <code className="font-code-block text-code-block flex-1 break-all text-on-surface">{value}</code>
      <CopyButton text={value} />
    </div>
  );
}
