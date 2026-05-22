import { useState } from "react";
import { Radio, Play, Plus, Trash2 } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { httpRequest } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Input, Select, Textarea, Toggle, NumberInput } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("http")!;

export default function HttpPage() {
  const [method, setMethod] = useState("GET");
  const [url, setUrl] = useState("https://httpbin.org/get");
  const [body, setBody] = useState("");
  const [headers, setHeaders] = useState<{ key: string; value: string }[]>([]);
  const [followRedirects, setFollowRedirects] = useState(true);
  const [timeoutSeconds, setTimeoutSeconds] = useState(15);
  const [allowPrivate, setAllowPrivate] = useState(false);
  const op = useOperation(httpRequest);

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 gap-md">
        <Card>
          <CardHeader title="Request" icon={<Radio className="h-5 w-5" />} />
          <div className="flex flex-col gap-md">
            <div className="grid grid-cols-[120px_1fr_auto] gap-2">
              <Select
                value={method}
                onChange={(e) => setMethod(e.target.value)}
                options={["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"].map((m) => ({
                  value: m,
                  label: m,
                }))}
              />
              <Input value={url} onChange={(e) => setUrl(e.target.value)} placeholder="https://…" />
              <Button
                onClick={() =>
                  op.run({
                    method,
                    url,
                    headers,
                    body: body || undefined,
                    followRedirects,
                    timeoutSeconds,
                    allowPrivate,
                  })
                }
                loading={op.loading}
                disabled={!url}
              >
                <Play className="h-4 w-4" /> Send
              </Button>
            </div>
            <div>
              <div className="flex items-center justify-between mb-1">
                <span className="font-data-label text-data-label uppercase text-on-surface-variant">
                  Headers
                </span>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setHeaders((h) => [...h, { key: "", value: "" }])}
                >
                  <Plus className="h-4 w-4" /> Header
                </Button>
              </div>
              <div className="flex flex-col gap-2">
                {headers.map((h, i) => (
                  <div key={i} className="grid grid-cols-[1fr_2fr_auto] gap-2">
                    <Input
                      placeholder="Header"
                      value={h.key}
                      onChange={(e) =>
                        setHeaders((hs) => hs.map((x, j) => (i === j ? { ...x, key: e.target.value } : x)))
                      }
                    />
                    <Input
                      placeholder="value"
                      value={h.value}
                      onChange={(e) =>
                        setHeaders((hs) => hs.map((x, j) => (i === j ? { ...x, value: e.target.value } : x)))
                      }
                    />
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setHeaders((hs) => hs.filter((_, j) => j !== i))}
                      aria-label="Remove"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </div>
            {!["GET", "HEAD"].includes(method) && (
              <Textarea label="Body" rows={6} value={body} onChange={(e) => setBody(e.target.value)} />
            )}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-md">
              <Toggle checked={followRedirects} onChange={setFollowRedirects} label="Follow redirects" />
              <NumberInput label="Timeout (s)" value={timeoutSeconds} onChange={setTimeoutSeconds} min={1} max={60} />
              <Toggle
                checked={allowPrivate}
                onChange={setAllowPrivate}
                label="Allow private IPs"
                hint="SSRF guard"
              />
            </div>
          </div>
        </Card>
        <Card padded={false}>
          <div className="flex items-center justify-between p-md border-b border-outline/10">
            <h3 className="font-body-md text-body-md font-medium text-on-surface">Response</h3>
            {op.data && (
              <div className="flex items-center gap-2">
                <Badge tone={op.data.status < 400 ? "success" : "error"}>
                  {op.data.status} {op.data.statusText}
                </Badge>
                <Badge>{op.data.durationMs}ms</Badge>
              </div>
            )}
          </div>
          <div className="p-md">
            {op.error ? (
              <ErrorBanner error={op.error} />
            ) : !op.data ? (
              <EmptyState title="Send a request" />
            ) : (
              <div className="flex flex-col gap-md">
                <details open>
                  <summary className="font-data-label text-data-label uppercase text-on-surface-variant cursor-pointer mb-2">
                    Headers ({Object.keys(op.data.headers).length})
                  </summary>
                  <ul className="flex flex-col gap-1">
                    {Object.entries(op.data.headers).map(([k, v]) => (
                      <li
                        key={k}
                        className="flex items-start gap-2 px-3 py-1.5 bg-surface-container-low rounded font-code-block text-code-block"
                      >
                        <span className="text-secondary-container shrink-0">{k}:</span>
                        <span className="text-on-surface break-all">{v}</span>
                      </li>
                    ))}
                  </ul>
                </details>
                <CodeBlock code={op.data.body} language="body" />
                <Diagnostics items={op.data.diagnostics} />
              </div>
            )}
          </div>
        </Card>
      </div>
    </ToolPage>
  );
}
