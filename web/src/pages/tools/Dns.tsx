import { useState } from "react";
import { Globe, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { dnsLookup } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Input, Select, Toggle } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, EmptyState, ErrorBanner } from "../../components/ui/Output";
import type { DnsLookupResponse } from "../../lib/api";

type Record = { type: string; value: string };
function flattenRecords(d: DnsLookupResponse): Record[] {
  const out: Record[] = [];
  d.a?.forEach((v) => out.push({ type: "A", value: v }));
  d.aaaa?.forEach((v) => out.push({ type: "AAAA", value: v }));
  d.cname?.forEach((v) => out.push({ type: "CNAME", value: v }));
  d.mx?.forEach((m) => out.push({ type: "MX", value: `${m.pref} ${m.host}` }));
  d.ns?.forEach((v) => out.push({ type: "NS", value: v }));
  d.txt?.forEach((v) => out.push({ type: "TXT", value: v.join(" ") }));
  d.ptr?.forEach((v) => out.push({ type: "PTR", value: v }));
  return out;
}

const tool = findTool("dns")!;

export default function DnsPage() {
  const [host, setHost] = useState("example.com");
  const [type, setType] = useState("A");
  const [allowPrivate, setAllowPrivate] = useState(false);
  const op = useOperation((args: { host: string; type: string; allowPrivate: boolean }) =>
    dnsLookup(args.host, args.type, args.allowPrivate),
  );

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-[420px_1fr] gap-md">
        <Card>
          <CardHeader title="Lookup" icon={<Globe className="h-5 w-5" />} />
          <div className="flex flex-col gap-md">
            <Input label="Host" value={host} onChange={(e) => setHost(e.target.value)} />
            <Select
              label="Record type"
              value={type}
              onChange={(e) => setType(e.target.value)}
              options={["A", "AAAA", "CNAME", "MX", "NS", "TXT", "SOA", "PTR"].map((t) => ({
                value: t,
                label: t,
              }))}
            />
            <Toggle
              checked={allowPrivate}
              onChange={setAllowPrivate}
              label="Allow private IPs in answers"
              hint="Disabled by default for SSRF safety"
            />
            <Button onClick={() => op.run({ host, type, allowPrivate })} loading={op.loading} fullWidth disabled={!host}>
              <Play className="h-4 w-4" /> Lookup
            </Button>
          </div>
        </Card>
        <Card>
          {(() => {
            const records = op.data ? flattenRecords(op.data) : [];
            return (
              <>
                <CardHeader
                  title="Records"
                  trailing={op.data && <Badge>{records.length} records</Badge>}
                />
                {op.error ? (
                  <ErrorBanner error={op.error} />
                ) : !op.data ? (
                  <EmptyState title="DNS records appear here" />
                ) : records.length === 0 ? (
                  <EmptyState title="No records" />
                ) : (
                  <ul className="flex flex-col gap-1">
                    {records.map((r, i) => (
                      <li
                        key={i}
                        className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded"
                      >
                        <Badge>{r.type}</Badge>
                        <code className="font-code-block text-code-block text-on-surface flex-1 truncate">
                          {r.value}
                        </code>
                        <CopyButton text={r.value} />
                      </li>
                    ))}
                  </ul>
                )}
              </>
            );
          })()}
        </Card>
      </div>
    </ToolPage>
  );
}
