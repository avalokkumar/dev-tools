import { useState } from "react";
import { Mail, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { headersAnalyze } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Textarea } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("headers")!;
const SAMPLE = `HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Strict-Transport-Security: max-age=31536000
X-Frame-Options: DENY
Content-Security-Policy: default-src 'self'`;

export default function HeadersPage() {
  const [input, setInput] = useState(SAMPLE);
  const op = useOperation((s: string) => headersAnalyze(s));
  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
        <Card>
          <CardHeader title="Headers" icon={<Mail className="h-5 w-5" />} />
          <Textarea rows={14} value={input} onChange={(e) => setInput(e.target.value)} placeholder="Paste HTTP response headers" />
          <Button onClick={() => op.run(input)} loading={op.loading} fullWidth className="mt-md" disabled={!input.trim()}>
            <Play className="h-4 w-4" /> Analyze
          </Button>
        </Card>
        <Card>
          <CardHeader title="Security audit" />
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Audit results appear here" />
          ) : (
            <div className="flex flex-col gap-md">
              <ul className="flex flex-col gap-2">
                {op.data.findings.map((f) => (
                  <li
                    key={f.header}
                    className={`flex items-start gap-3 px-3 py-2 rounded font-body-sm text-body-sm ${
                      f.ok
                        ? "bg-tertiary-container/20 text-on-surface"
                        : "bg-error-container/20 text-on-surface"
                    }`}
                  >
                    <span className={`shrink-0 font-semibold ${
                      f.ok ? "text-tertiary-container" : "text-error"
                    }`}>
                      {f.ok ? "✓" : "✗"}
                    </span>
                    <div>
                      <div className="font-code-block text-code-block">{f.header}</div>
                      <div className="text-on-surface-variant text-xs mt-0.5">{f.note}</div>
                    </div>
                  </li>
                ))}
              </ul>
              {Object.keys(op.data.headers).length > 0 && (
                <div>
                  <div className="font-data-label text-data-label uppercase text-on-surface-variant mb-1">
                    Parsed headers
                  </div>
                  <ul className="flex flex-col gap-1">
                    {Object.entries(op.data.headers).map(([k, v]) => (
                      <li
                        key={k}
                        className="flex items-start gap-2 px-3 py-2 bg-surface-container-low rounded font-code-block text-code-block"
                      >
                        <span className="text-secondary-container shrink-0">{k}:</span>
                        <span className="text-on-surface break-all">{v}</span>
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
