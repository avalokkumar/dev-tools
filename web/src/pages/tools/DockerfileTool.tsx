import { useState } from "react";
import { Container, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { dockerfileLint } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Textarea } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("dockerfile")!;
const SAMPLE = `FROM ubuntu:latest
RUN apt-get update && apt-get install -y curl
COPY . /app
CMD ["/app/start.sh"]
`;

export default function DockerfilePage() {
  const [input, setInput] = useState(SAMPLE);
  const op = useOperation((s: string) => dockerfileLint(s));
  const diagnostics = op.data?.diagnostics ?? [];
  const summary = {
    errors: diagnostics.filter((d) => d.severity === 2).length,
    warnings: diagnostics.filter((d) => d.severity === 1).length,
    info: diagnostics.filter((d) => d.severity === 0).length,
  };
  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
        <Card>
          <CardHeader title="Dockerfile" icon={<Container className="h-5 w-5" />} />
          <Textarea rows={16} value={input} onChange={(e) => setInput(e.target.value)} />
          <Button onClick={() => op.run(input)} loading={op.loading} fullWidth className="mt-md" disabled={!input.trim()}>
            <Play className="h-4 w-4" /> Lint
          </Button>
        </Card>
        <Card>
          <CardHeader
            title="Diagnostics"
            trailing={
              op.data && (
                <div className="flex items-center gap-1">
                  <Badge tone="error">{summary.errors}</Badge>
                  <Badge tone="warning">{summary.warnings}</Badge>
                  <Badge tone="info">{summary.info}</Badge>
                </div>
              )
            }
          />
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Lint output appears here" />
          ) : diagnostics.length === 0 ? (
            <EmptyState title="No issues found" description="Best practices look good." />
          ) : (
            <Diagnostics items={diagnostics} />
          )}
        </Card>
      </div>
    </ToolPage>
  );
}
