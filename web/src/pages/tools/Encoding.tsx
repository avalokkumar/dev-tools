import { useState } from "react";
import { ArrowRightLeft, Play } from "lucide-react";
import { findTool } from "../../lib/catalog";
import {
  encBase64Encode,
  encBase64Decode,
  encUrlEncode,
  encUrlDecode,
  encHtmlEncode,
  encHtmlDecode,
  encHexEncode,
  encHexDecode,
} from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Textarea, Toggle } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("encoding")!;

type Codec = "base64" | "url" | "html" | "hex";

export default function EncodingPage() {
  const [codec, setCodec] = useState<Codec>("base64");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "base64", label: "Base64" },
          { id: "url", label: "URL" },
          { id: "html", label: "HTML Entities" },
          { id: "hex", label: "Hex" },
        ]}
        active={codec}
        onChange={(id) => setCodec(id as Codec)}
        className="mb-md"
      />
      <CodecPanel codec={codec} />
    </ToolPage>
  );
}

function CodecPanel({ codec }: { codec: Codec }) {
  const [direction, setDirection] = useState<"encode" | "decode">("encode");
  const [input, setInput] = useState("");
  const [urlSafe, setUrlSafe] = useState(false);
  const [noPadding, setNoPadding] = useState(false);
  const [pathMode, setPathMode] = useState(false);
  const [uppercase, setUppercase] = useState(false);
  const op = useOperation(async (args: { input: string }): Promise<{ output: string }> => {
    const text = args.input;
    if (codec === "base64")
      return direction === "encode"
        ? encBase64Encode(text, urlSafe, noPadding)
        : encBase64Decode(text, urlSafe);
    if (codec === "url")
      return direction === "encode" ? encUrlEncode(text, pathMode) : encUrlDecode(text);
    if (codec === "html")
      return direction === "encode" ? encHtmlEncode(text) : encHtmlDecode(text);
    return direction === "encode" ? encHexEncode(text, uppercase) : encHexDecode(text);
  });

  function swap() {
    setDirection((d) => (d === "encode" ? "decode" : "encode"));
    if (op.data) setInput(op.data.output);
    op.reset();
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[1fr_auto_1fr] gap-md items-start">
      <Card>
        <CardHeader
          title={direction === "encode" ? "Plain text" : `${codec.toUpperCase()} input`}
        />
        <Textarea
          rows={10}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={direction === "encode" ? "Enter text…" : "Enter encoded data…"}
        />
        {codec === "base64" && (
          <div className="grid grid-cols-2 gap-md mt-md">
            <Toggle checked={urlSafe} onChange={setUrlSafe} label="URL-safe alphabet" />
            <Toggle checked={noPadding} onChange={setNoPadding} label="No padding" disabled={direction === "decode"} />
          </div>
        )}
        {codec === "url" && (
          <div className="mt-md">
            <Toggle checked={pathMode} onChange={setPathMode} label="Path mode (preserve /)" disabled={direction === "decode"} />
          </div>
        )}
        {codec === "hex" && (
          <div className="mt-md">
            <Toggle checked={uppercase} onChange={setUppercase} label="Uppercase output" disabled={direction === "decode"} />
          </div>
        )}
        <div className="grid grid-cols-2 gap-2 mt-md">
          <Button
            variant={direction === "encode" ? "primary" : "outline"}
            onClick={() => setDirection("encode")}
          >
            Encode
          </Button>
          <Button
            variant={direction === "decode" ? "primary" : "outline"}
            onClick={() => setDirection("decode")}
          >
            Decode
          </Button>
        </div>
        <Button
          onClick={() => op.run({ input })}
          loading={op.loading}
          fullWidth
          className="mt-md"
          disabled={!input.trim()}
        >
          <Play className="h-4 w-4" /> Run
        </Button>
      </Card>
      <div className="hidden lg:flex flex-col items-center justify-center pt-12">
        <Button variant="ghost" size="sm" onClick={swap} aria-label="Swap" title="Swap direction">
          <ArrowRightLeft className="h-5 w-5" />
        </Button>
      </div>
      <Card padded={false}>
        <CardHeader title={direction === "encode" ? `${codec.toUpperCase()} output` : "Plain text"} />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Output appears here" />
          ) : (
            <CodeBlock code={op.data.output} language={direction === "encode" ? codec : "text"} />
          )}
        </div>
      </Card>
    </div>
  );
}
