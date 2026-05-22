import { useEffect, useState } from "react";
import { Link as LinkIcon } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { urlParse } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input } from "../../components/ui/Input";
import { CopyButton, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("url")!;

export default function UrlPage() {
  const [u, setU] = useState("https://devforge.io/path/to/page?q=hello&lang=en#frag");
  const op = useOperation((url: string) => urlParse(url));

  useEffect(() => {
    if (u.trim()) op.run(u.trim());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [u]);

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
        <Card>
          <CardHeader title="URL" icon={<LinkIcon className="h-5 w-5" />} />
          <Input value={u} onChange={(e) => setU(e.target.value)} />
        </Card>
        <Card>
          <CardHeader title="Components" />
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Type a URL" />
          ) : (
            <div className="flex flex-col gap-2">
              {[
                ["Scheme", op.data.scheme],
                ["Hostname", op.data.hostname],
                ["Port", op.data.port],
                ["Path", op.data.path],
                ["Query", op.data.query],
                ["Fragment", op.data.fragment],
              ].map(([k, v]) => (
                <div key={k} className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded">
                  <span className="font-data-label text-data-label uppercase text-on-surface-variant w-20 shrink-0">
                    {k}
                  </span>
                  <code className="font-code-block text-code-block text-on-surface flex-1 truncate">
                    {String(v) || "—"}
                  </code>
                  {v && <CopyButton text={String(v)} />}
                </div>
              ))}
              {op.data.params.length > 0 && (
                <div className="mt-md">
                  <div className="font-data-label text-data-label uppercase text-on-surface-variant mb-1">
                    Query parameters
                  </div>
                  <ul className="flex flex-col gap-1">
                    {op.data.params.map((p, i) => (
                      <li
                        key={i}
                        className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded font-code-block text-code-block"
                      >
                        <span className="text-secondary-container">{p.key}</span>
                        <span className="text-on-surface-variant">=</span>
                        <span className="text-on-surface flex-1 truncate">{p.value}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          )}
        </Card>
      </div>
    </ToolPage>
  );
}
