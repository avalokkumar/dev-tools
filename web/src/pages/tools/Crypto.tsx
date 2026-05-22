import { useState } from "react";
import { Play, ShieldCheck, KeyRound, Lock, Unlock, FileKey, Hash } from "lucide-react";
import { findTool } from "../../lib/catalog";
import {
  cryptoAesEncrypt,
  cryptoAesDecrypt,
  cryptoRsaKeygen,
  cryptoHmac,
  cryptoPasswordHash,
  cryptoPasswordStrength,
} from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader } from "../../components/ui/Card";
import { Input, Select, Textarea, NumberInput } from "../../components/ui/Input";
import { Button } from "../../components/ui/Button";
import { CodeBlock, CopyButton, EmptyState, ErrorBanner } from "../../components/ui/Output";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("crypto")!;

type TabId = "encrypt" | "decrypt" | "rsa" | "hmac" | "password" | "strength";

export default function CryptoPage() {
  const [tab, setTab] = useState<TabId>("encrypt");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "encrypt", label: "AES Encrypt" },
          { id: "decrypt", label: "AES Decrypt" },
          { id: "rsa", label: "RSA Keygen" },
          { id: "hmac", label: "HMAC" },
          { id: "password", label: "Password Hash" },
          { id: "strength", label: "Password Strength" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as TabId)}
        className="mb-md"
      />
      {tab === "encrypt" && <EncryptTab />}
      {tab === "decrypt" && <DecryptTab />}
      {tab === "rsa" && <RsaTab />}
      {tab === "hmac" && <HmacTab />}
      {tab === "password" && <PasswordTab />}
      {tab === "strength" && <StrengthTab />}
    </ToolPage>
  );
}

function EncryptTab() {
  const [plaintext, setPlaintext] = useState("");
  const [passphrase, setPassphrase] = useState("");
  const [keySize, setKeySize] = useState<128 | 192 | 256>(256);
  const op = useOperation(cryptoAesEncrypt);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Plaintext" icon={<Lock className="h-5 w-5" />} />
        <div className="flex flex-col gap-md">
          <Textarea
            label="Message"
            mono={false}
            rows={6}
            value={plaintext}
            onChange={(e) => setPlaintext(e.target.value)}
            placeholder="Enter text to encrypt…"
          />
          <Input
            label="Passphrase"
            type="password"
            value={passphrase}
            onChange={(e) => setPassphrase(e.target.value)}
            placeholder="Strong passphrase"
          />
          <Select
            label="Key size"
            value={String(keySize)}
            onChange={(e) => setKeySize(Number(e.target.value) as 128 | 192 | 256)}
            options={[
              { value: "128", label: "AES-128" },
              { value: "192", label: "AES-192" },
              { value: "256", label: "AES-256 (recommended)" },
            ]}
          />
          <Button
            onClick={() => op.run({ plaintext, passphrase, keySize })}
            loading={op.loading}
            fullWidth
            disabled={!plaintext || !passphrase}
          >
            <Play className="h-4 w-4" /> Encrypt
          </Button>
        </div>
      </Card>
      <Card padded={false}>
        <CardHeader title="Ciphertext" />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Encrypted output appears here" />
          ) : (
            <CodeBlock code={op.data.ciphertext} language="ciphertext" />
          )}
        </div>
      </Card>
    </div>
  );
}

function DecryptTab() {
  const [ciphertext, setCiphertext] = useState("");
  const [passphrase, setPassphrase] = useState("");
  const op = useOperation(cryptoAesDecrypt);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Ciphertext" icon={<Unlock className="h-5 w-5" />} />
        <div className="flex flex-col gap-md">
          <Textarea
            label="Encrypted text"
            rows={6}
            value={ciphertext}
            onChange={(e) => setCiphertext(e.target.value)}
            placeholder="Paste encrypted output…"
          />
          <Input
            label="Passphrase"
            type="password"
            value={passphrase}
            onChange={(e) => setPassphrase(e.target.value)}
          />
          <Button
            onClick={() => op.run({ ciphertext, passphrase })}
            loading={op.loading}
            fullWidth
            disabled={!ciphertext || !passphrase}
          >
            <Play className="h-4 w-4" /> Decrypt
          </Button>
        </div>
      </Card>
      <Card padded={false}>
        <CardHeader title="Plaintext" />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Decrypted output appears here" />
          ) : (
            <CodeBlock code={op.data.plaintext} language="text" />
          )}
        </div>
      </Card>
    </div>
  );
}

function RsaTab() {
  const [bits, setBits] = useState<2048 | 3072 | 4096>(3072);
  const op = useOperation((b: 2048 | 3072 | 4096) => cryptoRsaKeygen(b));

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[360px_1fr] gap-md">
      <Card>
        <CardHeader title="Key Size" icon={<KeyRound className="h-5 w-5" />} />
        <div className="flex flex-col gap-md">
          <Select
            label="Bits"
            value={String(bits)}
            onChange={(e) => setBits(Number(e.target.value) as 2048 | 3072 | 4096)}
            options={[
              { value: "2048", label: "2048-bit" },
              { value: "3072", label: "3072-bit (recommended)" },
              { value: "4096", label: "4096-bit" },
            ]}
          />
          <Button onClick={() => op.run(bits)} loading={op.loading} fullWidth>
            <Play className="h-4 w-4" /> Generate Key Pair
          </Button>
        </div>
      </Card>
      <div className="flex flex-col gap-md">
        {op.error ? <ErrorBanner error={op.error} /> : null}
        {!op.data && !op.error && <EmptyState title="Generate a key pair to view PEM output" />}
        {op.data && (
          <>
            <CodeBlock
              code={op.data.privatePem}
              language="private-key.pem"
              download={{ filename: "private.pem" }}
            />
            <CodeBlock
              code={op.data.publicPem}
              language="public-key.pem"
              download={{ filename: "public.pem" }}
            />
          </>
        )}
      </div>
    </div>
  );
}

function HmacTab() {
  const [input, setInput] = useState("");
  const [key, setKey] = useState("");
  const [keyEncoding, setKeyEncoding] = useState<"raw" | "hex" | "base64">("raw");
  const [algorithm, setAlgorithm] = useState<"sha256" | "sha384" | "sha512">("sha256");
  const op = useOperation(cryptoHmac);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Input" icon={<Hash className="h-5 w-5" />} />
        <div className="flex flex-col gap-md">
          <Textarea label="Message" rows={5} value={input} onChange={(e) => setInput(e.target.value)} />
          <Input label="Key" value={key} onChange={(e) => setKey(e.target.value)} />
          <div className="grid grid-cols-2 gap-md">
            <Select
              label="Key encoding"
              value={keyEncoding}
              onChange={(e) => setKeyEncoding(e.target.value as typeof keyEncoding)}
              options={[
                { value: "raw", label: "Raw" },
                { value: "hex", label: "Hex" },
                { value: "base64", label: "Base64" },
              ]}
            />
            <Select
              label="Algorithm"
              value={algorithm}
              onChange={(e) => setAlgorithm(e.target.value as typeof algorithm)}
              options={[
                { value: "sha256", label: "HMAC-SHA-256" },
                { value: "sha384", label: "HMAC-SHA-384" },
                { value: "sha512", label: "HMAC-SHA-512" },
              ]}
            />
          </div>
          <Button
            onClick={() => op.run({ input, key, keyEncoding, algorithm })}
            loading={op.loading}
            fullWidth
            disabled={!input || !key}
          >
            <Play className="h-4 w-4" /> Compute MAC
          </Button>
        </div>
      </Card>
      <Card padded={false}>
        <CardHeader title="MAC" />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="MAC appears here" />
          ) : (
            <div className="flex items-center justify-between gap-2 px-3 py-2 bg-surface-container-low rounded">
              <code className="font-code-block text-code-block break-all">{op.data.mac}</code>
              <CopyButton text={op.data.mac} />
            </div>
          )}
        </div>
      </Card>
    </div>
  );
}

function PasswordTab() {
  const [password, setPassword] = useState("");
  const [algorithm, setAlgorithm] = useState<"bcrypt" | "argon2id">("bcrypt");
  const [cost, setCost] = useState(12);
  const op = useOperation(cryptoPasswordHash);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Password" icon={<FileKey className="h-5 w-5" />} />
        <div className="flex flex-col gap-md">
          <Input
            label="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <div className="grid grid-cols-2 gap-md">
            <Select
              label="Algorithm"
              value={algorithm}
              onChange={(e) => setAlgorithm(e.target.value as typeof algorithm)}
              options={[
                { value: "bcrypt", label: "bcrypt" },
                { value: "argon2id", label: "argon2id" },
              ]}
            />
            <NumberInput label="Cost" value={cost} onChange={setCost} min={4} max={20} />
          </div>
          <Button
            onClick={() => op.run({ password, algorithm, cost })}
            loading={op.loading}
            fullWidth
            disabled={!password}
          >
            <Play className="h-4 w-4" /> Hash
          </Button>
        </div>
      </Card>
      <Card padded={false}>
        <CardHeader title="Hash" />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Hash appears here" />
          ) : (
            <CodeBlock code={op.data.hash} language="hash" />
          )}
        </div>
      </Card>
    </div>
  );
}

function StrengthTab() {
  const [password, setPassword] = useState("");
  const op = useOperation((p: string) => cryptoPasswordStrength(p));

  function check(p: string) {
    setPassword(p);
    if (p) op.run(p);
    else op.reset();
  }

  const score = op.data?.score ?? 0;
  const labels = ["Very weak", "Weak", "Fair", "Strong", "Excellent"];
  const colors = ["bg-error", "bg-secondary-container", "bg-orange-400", "bg-tertiary-fixed-dim", "bg-tertiary-fixed-dim"];

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title="Password" icon={<ShieldCheck className="h-5 w-5" />} />
        <Input
          label="Test password"
          type="password"
          value={password}
          onChange={(e) => check(e.target.value)}
          placeholder="Type to evaluate…"
        />
      </Card>
      <Card>
        <CardHeader title="Strength" />
        {op.error ? (
          <ErrorBanner error={op.error} />
        ) : !op.data ? (
          <EmptyState title="Type a password to evaluate" />
        ) : (
          <div className="flex flex-col gap-md">
            <div className="flex items-center gap-2">
              <span className="font-data-label text-data-label uppercase text-on-surface-variant w-20">
                Score
              </span>
              <div className="flex-1 h-2 bg-surface-container rounded-full overflow-hidden flex">
                {[0, 1, 2, 3, 4].map((i) => (
                  <div
                    key={i}
                    className={`flex-1 mr-0.5 last:mr-0 ${i <= score ? colors[score] : "bg-surface-container"}`}
                  />
                ))}
              </div>
              <span className="font-body-sm text-body-sm text-on-surface w-20 text-right">
                {labels[score]}
              </span>
            </div>
            <div>
              <span className="font-data-label text-data-label uppercase text-on-surface-variant">
                Crack time
              </span>
              <div className="font-body-md text-body-md text-on-surface mt-0.5">
                {op.data.crackTime}
              </div>
            </div>
            {op.data.feedback.warning && (
              <div className="text-error font-body-sm text-body-sm">
                ⚠ {op.data.feedback.warning}
              </div>
            )}
            {op.data.feedback.suggestions.length > 0 && (
              <ul className="list-disc pl-5 font-body-sm text-body-sm text-on-surface-variant flex flex-col gap-1">
                {op.data.feedback.suggestions.map((s, i) => (
                  <li key={i}>{s}</li>
                ))}
              </ul>
            )}
          </div>
        )}
      </Card>
    </div>
  );
}
