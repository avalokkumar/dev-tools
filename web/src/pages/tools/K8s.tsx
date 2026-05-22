import { useState } from "react";
import { Boxes, Play, CheckCircle2, XCircle } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { k8sValidate } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Textarea } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("k8s")!;
const SAMPLE = `apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
data:
  LOG_LEVEL: info
`;

export default function K8sPage() {
  const [input, setInput] = useState(SAMPLE);
  const op = useOperation((s: string) => k8sValidate(s));
  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
        <Card>
          <CardHeader title="Manifest" icon={<Boxes className="h-5 w-5" />} />
          <Textarea rows={16} value={input} onChange={(e) => setInput(e.target.value)} />
          <Button onClick={() => op.run(input)} loading={op.loading} fullWidth className="mt-md" disabled={!input.trim()}>
            <Play className="h-4 w-4" /> Validate
          </Button>
        </Card>
        <Card>
          <CardHeader
            title="Result"
            trailing={
              op.data && (
                <div className="flex gap-2">
                  {op.data.kind && <Badge>{op.data.kind}</Badge>}
                  {op.data.apiVersion && <Badge tone="info">{op.data.apiVersion}</Badge>}
                </div>
              )
            }
          />
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Validation result appears here" />
          ) : (
            <div className="flex flex-col gap-md">
              <div
                className={`flex items-center gap-3 p-md rounded-lg ${
                  op.data.valid ? "bg-tertiary-fixed text-on-tertiary-fixed" : "bg-error/10 text-error"
                }`}
              >
                {op.data.valid ? <CheckCircle2 className="h-6 w-6" /> : <XCircle className="h-6 w-6" />}
                <span className="font-body-md font-semibold">
                  {op.data.valid ? "Valid manifest" : "Invalid manifest"}
                </span>
              </div>
              <Diagnostics items={op.data.diagnostics} />
            </div>
          )}
        </Card>
      </div>
    </ToolPage>
  );
}
