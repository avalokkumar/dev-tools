import { useState } from "react";
import { MapPin, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { ipCalc } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, NumberInput } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("ip")!;

export default function IpPage() {
  const [cidr, setCidr] = useState("10.0.0.0/24");
  const [maxHosts, setMaxHosts] = useState(16);
  const op = useOperation((args: { cidr: string; maxList: number }) =>
    ipCalc(args.cidr, args.maxList),
  );

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-[400px_1fr] gap-md">
        <Card>
          <CardHeader title="CIDR" icon={<MapPin className="h-5 w-5" />} />
          <div className="flex flex-col gap-md">
            <Input label="CIDR block" value={cidr} onChange={(e) => setCidr(e.target.value)} placeholder="10.0.0.0/24" />
            <NumberInput label="Max host list" value={maxHosts} onChange={setMaxHosts} min={0} max={4096} />
            <Button onClick={() => op.run({ cidr, maxList: maxHosts })} loading={op.loading} fullWidth disabled={!cidr}>
              <Play className="h-4 w-4" /> Calculate
            </Button>
          </div>
        </Card>
        <Card>
          <CardHeader title="Subnet info" />
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Subnet info appears here" />
          ) : (
            <div className="flex flex-col gap-md">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                {[
                  ["Network", op.data.network],
                  ["Broadcast", op.data.broadcast],
                  ["Netmask", op.data.netmask],
                  ["Wildcard", op.data.wildcard],
                  ["First host", op.data.first],
                  ["Last host", op.data.last],
                  ["Prefix", `/${op.data.prefix}`],
                  ["Usable hosts", op.data.usableHosts.toLocaleString()],
                ].map(([k, v]) => (
                  <div key={k} className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded">
                    <span className="font-data-label text-data-label uppercase text-on-surface-variant w-28 shrink-0">
                      {k}
                    </span>
                    <code className="font-code-block text-code-block text-on-surface flex-1 truncate">
                      {v}
                    </code>
                    <CopyButton text={String(v)} />
                  </div>
                ))}
              </div>
              {op.data.hosts && op.data.hosts.length > 0 && (
                <details>
                  <summary className="font-data-label text-data-label uppercase text-on-surface-variant cursor-pointer">
                    Host list ({op.data.hosts.length})
                  </summary>
                  <ul className="mt-2 grid grid-cols-2 md:grid-cols-3 gap-1 font-code-block text-code-block">
                    {op.data.hosts.map((h: string) => (
                      <li key={h} className="px-2 py-1 bg-surface-container-low rounded">
                        {h}
                      </li>
                    ))}
                  </ul>
                </details>
              )}
              <Diagnostics items={op.data.diagnostics} />
            </div>
          )}
        </Card>
      </div>
    </ToolPage>
  );
}
