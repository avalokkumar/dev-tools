import { useEffect, useState } from "react";
import { Ticket, Play, ShieldCheck, ShieldX } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { jwtDecode, jwtVerify } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Textarea, Input, Select } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("jwt")!;
const SAMPLE =
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzEiLCJyb2xlIjoiYWRtaW4iLCJpYXQiOjE3MTUwMDAwMDB9.HxSr6fJq2W2_W1cN9xq3D8kqR1g0qz2H6oZ7QfHqKqM";

export default function JwtPage() {
  const [tab, setTab] = useState<"decode" | "verify">("decode");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "decode", label: "Decode" },
          { id: "verify", label: "Verify" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "decode" ? <DecodeTab /> : <VerifyTab />}
    </ToolPage>
  );
}

function DecodeTab() {
  const [token, setToken] = useState(SAMPLE);
  const op = useOperation((t: string) => jwtDecode(t));

  useEffect(() => {
    if (token.trim()) op.run(token.trim());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token]);

  const parts = token.split(".");

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Token" icon={<Ticket className="h-5 w-5" />} />
        <Textarea rows={10} value={token} onChange={(e) => setToken(e.target.value)} />
        {parts.length === 3 && (
          <div className="mt-md font-code-block text-code-block break-all leading-relaxed bg-surface-container-low rounded p-3">
            <span className="text-secondary-container">{parts[0]}</span>
            <span className="text-on-surface-variant">.</span>
            <span className="text-sky-aqua">{parts[1]}</span>
            <span className="text-on-surface-variant">.</span>
            <span className="text-tertiary-container">{parts[2]}</span>
          </div>
        )}
      </Card>
      <div className="flex flex-col gap-md">
        {op.error ? <ErrorBanner error={op.error} /> : null}
        {!op.data && !op.error && <EmptyState title="Decoded token appears here" />}
        {op.data && (
          <>
            <Card padded={false}>
              <CardHeader title="Header" />
              <div className="p-md">
                <CodeBlock code={JSON.stringify(op.data.header, null, 2)} language="json" />
              </div>
            </Card>
            <Card padded={false}>
              <CardHeader
                title="Payload"
                trailing={<ExpiryBadge payload={op.data.payload} />}
              />
              <div className="p-md">
                <CodeBlock code={JSON.stringify(op.data.payload, null, 2)} language="json" />
              </div>
            </Card>
            <Diagnostics items={op.data.diagnostics} />
          </>
        )}
      </div>
    </div>
  );
}

function ExpiryBadge({ payload }: { payload: Record<string, unknown> }) {
  const exp = typeof payload.exp === "number" ? payload.exp : null;
  if (!exp) return null;
  const now = Math.floor(Date.now() / 1000);
  if (exp < now) return <Badge tone="error">Expired</Badge>;
  const remain = exp - now;
  const human = remain > 86400 ? `${Math.floor(remain / 86400)}d` : remain > 3600 ? `${Math.floor(remain / 3600)}h` : `${Math.floor(remain / 60)}m`;
  return <Badge tone="success">Valid · expires in {human}</Badge>;
}

function VerifyTab() {
  const [token, setToken] = useState(SAMPLE);
  const [key, setKey] = useState("");
  const [keyFormat, setKeyFormat] = useState<"hmac" | "pem">("hmac");
  const [leeway, setLeeway] = useState(0);
  const op = useOperation((args: { token: string; key: string; keyFormat: "hmac" | "pem"; leeway: number }) =>
    jwtVerify(args.token, args.key, args.keyFormat, [], args.leeway),
  );

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Token & Key" />
        <div className="flex flex-col gap-md">
          <Textarea label="Token" rows={6} value={token} onChange={(e) => setToken(e.target.value)} />
          <Textarea label="Key" rows={6} value={key} onChange={(e) => setKey(e.target.value)} placeholder="HMAC secret or PEM-encoded key" />
          <div className="grid grid-cols-2 gap-md">
            <Select
              label="Key format"
              value={keyFormat}
              onChange={(e) => setKeyFormat(e.target.value as typeof keyFormat)}
              options={[
                { value: "hmac", label: "HMAC secret" },
                { value: "pem", label: "PEM (RSA/ECDSA)" },
              ]}
            />
            <Input
              label="Leeway (s)"
              type="number"
              value={leeway}
              onChange={(e) => setLeeway(Number(e.target.value))}
            />
          </div>
          <Button
            onClick={() => op.run({ token, key, keyFormat, leeway })}
            loading={op.loading}
            fullWidth
            disabled={!token || !key}
          >
            <Play className="h-4 w-4" /> Verify
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader title="Result" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Verification result appears here" />
        ) : op.data.valid ? (
          <div className="flex items-center gap-3 p-md rounded-lg bg-tertiary-fixed text-on-tertiary-fixed">
            <ShieldCheck className="h-6 w-6" />
            <div>
              <div className="font-body-md font-semibold">Signature valid</div>
              <div className="font-body-sm text-body-sm">Token signature matches the supplied key.</div>
            </div>
          </div>
        ) : (
          <div className="flex items-center gap-3 p-md rounded-lg bg-error/10 text-error">
            <ShieldX className="h-6 w-6" />
            <div>
              <div className="font-body-md font-semibold">Signature invalid</div>
              <div className="font-body-sm text-body-sm">See diagnostics for details.</div>
            </div>
          </div>
        )}
        <Diagnostics items={op.data?.diagnostics} className="mt-md" />
      </Card>
    </div>
  );
}
