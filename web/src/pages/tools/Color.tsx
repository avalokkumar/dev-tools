import { useEffect, useState } from "react";
import { Palette } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { colorConvert } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input } from "../../components/ui/Input";
import { CopyButton, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("color")!;

export default function ColorPage() {
  const [input, setInput] = useState("#57c4e5");
  const op = useOperation((c: string) => colorConvert(c));

  useEffect(() => {
    if (input.trim()) op.run(input.trim());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [input]);

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-[420px_1fr] gap-md">
        <Card>
          <CardHeader title="Color input" icon={<Palette className="h-5 w-5" />} />
          <div className="flex flex-col gap-md">
            <Input
              label="Color"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="#57c4e5, rgb(87,196,229), hsl(196,75%,62%)"
            />
            <div
              className="h-32 rounded-md border border-outline/20"
              style={{ backgroundColor: op.data?.hex ?? input }}
            />
            <input
              type="color"
              value={op.data?.hex ?? "#000000"}
              onChange={(e) => setInput(e.target.value)}
              className="w-full h-10 rounded cursor-pointer"
              aria-label="Color picker"
            />
          </div>
        </Card>
        <Card>
          <CardHeader title="Conversions & contrast" />
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Type a color to convert" />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-md">
              <Row label="HEX" value={op.data.hex} />
              <Row label="RGB" value={op.data.rgb} />
              <Row label="HSL" value={op.data.hsl} />
              <Row label="R" value={String(op.data.r)} />
              <Row label="G" value={String(op.data.g)} />
              <Row label="B" value={String(op.data.b)} />
            </div>
          )}
        </Card>
      </div>
    </ToolPage>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded">
      <span className="font-data-label text-data-label uppercase text-on-surface-variant w-12 shrink-0">
        {label}
      </span>
      <code className="font-code-block text-code-block text-on-surface flex-1 truncate">{value}</code>
      <CopyButton text={value} />
    </div>
  );
}

