import { useEffect, useState } from "react";
import { Play, Hash, ShieldCheck, ShieldX } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { totpGenerate, totpVerify } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, Select, NumberInput } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CopyButton, Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("totp")!;

export default function TotpPage() {
  const [tab, setTab] = useState<"generate" | "verify">("generate");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "generate", label: "Generate" },
          { id: "verify", label: "Verify" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "generate" ? <Generate /> : <Verify />}
    </ToolPage>
  );
}

function commonFields(state: ReturnType<typeof useFormState>) {
  return (
    <>
      <Input
        label="Secret"
        value={state.secret}
        onChange={(e) => state.setSecret(e.target.value)}
        placeholder="JBSWY3DPEHPK3PXP"
      />
      <div className="grid grid-cols-2 gap-md">
        <Select
          label="Encoding"
          value={state.encoding}
          onChange={(e) => state.setEncoding(e.target.value as typeof state.encoding)}
          options={[
            { value: "base32", label: "Base32" },
            { value: "hex", label: "Hex" },
            { value: "raw", label: "Raw" },
          ]}
        />
        <Select
          label="Algorithm"
          value={state.algorithm}
          onChange={(e) => state.setAlgorithm(e.target.value as typeof state.algorithm)}
          options={[
            { value: "sha1", label: "SHA-1" },
            { value: "sha256", label: "SHA-256" },
            { value: "sha512", label: "SHA-512" },
          ]}
        />
        <NumberInput label="Digits" value={state.digits} onChange={(n) => state.setDigits(n as 6 | 8)} min={6} max={8} step={2} />
        <NumberInput label="Period (s)" value={state.period} onChange={state.setPeriod} min={15} max={120} />
      </div>
    </>
  );
}

function useFormState() {
  const [secret, setSecret] = useState("JBSWY3DPEHPK3PXP");
  const [encoding, setEncoding] = useState<"raw" | "hex" | "base32">("base32");
  const [algorithm, setAlgorithm] = useState<"sha1" | "sha256" | "sha512">("sha1");
  const [digits, setDigits] = useState<6 | 8>(6);
  const [period, setPeriod] = useState(30);
  return { secret, setSecret, encoding, setEncoding, algorithm, setAlgorithm, digits, setDigits, period, setPeriod };
}

function Generate() {
  const fs = useFormState();
  const op = useOperation(totpGenerate);
  const [now, setNow] = useState(Date.now());

  useEffect(() => {
    const id = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(id);
  }, []);

  function run() {
    op.run({ secret: fs.secret, encoding: fs.encoding, algorithm: fs.algorithm, digits: fs.digits, period: fs.period });
  }

  const expiresIn = op.data
    ? Math.max(0, op.data.periodSec - Math.floor((now / 1000) % op.data.periodSec))
    : 0;

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[420px_1fr] gap-md">
      <Card>
        <CardHeader title="Configuration" />
        <div className="flex flex-col gap-md">
          {commonFields(fs)}
          <Button onClick={run} loading={op.loading} fullWidth disabled={!fs.secret}>
            <Play className="h-4 w-4" /> Generate Code
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader title="Current Code" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Press Generate" icon={<Hash className="h-6 w-6" />} />
        ) : (
          <div className="flex flex-col items-center gap-md py-md">
            <div className="font-code-block text-[64px] font-semibold tracking-widest text-on-surface tabular-nums">
              {op.data.code}
            </div>
            <div className="flex items-center gap-2 font-data-label text-data-label uppercase text-on-surface-variant">
              <span>Expires in</span>
              <span className="text-on-surface tabular-nums">{expiresIn}s</span>
            </div>
            <div className="w-full h-2 bg-surface-container rounded-full overflow-hidden">
              <div
                className="h-full bg-secondary-container transition-all"
                style={{ width: `${(expiresIn / op.data.periodSec) * 100}%` }}
              />
            </div>
            <CopyButton text={op.data.code} label="Copy code" size="md" />
          </div>
        )}
      </Card>
    </div>
  );
}

function Verify() {
  const fs = useFormState();
  const [code, setCode] = useState("");
  const [skew, setSkew] = useState(1);
  const op = useOperation(totpVerify);

  function run() {
    op.run({
      secret: fs.secret,
      encoding: fs.encoding,
      algorithm: fs.algorithm,
      digits: fs.digits,
      period: fs.period,
      code,
      skew,
    });
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[420px_1fr] gap-md">
      <Card>
        <CardHeader title="Configuration" />
        <div className="flex flex-col gap-md">
          {commonFields(fs)}
          <div className="grid grid-cols-2 gap-md">
            <Input label="Code" value={code} onChange={(e) => setCode(e.target.value)} placeholder="123456" />
            <NumberInput label="Skew (steps)" value={skew} onChange={setSkew} min={0} max={5} />
          </div>
          <Button onClick={run} loading={op.loading} fullWidth disabled={!fs.secret || !code}>
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
              <div className="font-body-md font-semibold">Code is valid</div>
              <div className="font-body-sm text-body-sm">Token matched current time window.</div>
            </div>
          </div>
        ) : (
          <div className="flex items-center gap-3 p-md rounded-lg bg-error/10 text-error">
            <ShieldX className="h-6 w-6" />
            <div>
              <div className="font-body-md font-semibold">Code is invalid</div>
              <div className="font-body-sm text-body-sm">Token does not match — check secret, code, or clock skew.</div>
            </div>
          </div>
        )}
        <Diagnostics items={op.data?.diagnostics} className="mt-md" />
      </Card>
    </div>
  );
}
