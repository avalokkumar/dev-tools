import { useState } from "react";
import { findTool } from "../../lib/catalog";
import { yamlFormat, yamlValidate, yamlConvert } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { FormatterShell } from "../../components/layout/FormatterShell";
import { NumberInput, Select } from "../../components/ui/Input";
import { Tabs } from "../../components/ui/Tabs";
import { CheckCircle2 } from "lucide-react";

const tool = findTool("yaml")!;

const SAMPLE = `name: DevForge\nversion: 0.1\ntools: 75\nfeatures:\n  - cli\n  - web\n  - mcp\n`;

export default function YamlPage() {
  const [tab, setTab] = useState<"format" | "validate" | "convert">("format");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "format", label: "Format" },
          { id: "validate", label: "Validate" },
          { id: "convert", label: "Convert" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "format" && <FormatTab />}
      {tab === "validate" && <ValidateTab />}
      {tab === "convert" && <ConvertTab />}
    </ToolPage>
  );
}

function FormatTab() {
  const [indent, setIndent] = useState(2);
  const op = useOperation((args: { input: string; indent: number }) => yamlFormat(args.input, args.indent));
  return (
    <FormatterShell
      initial={SAMPLE}
      buildInput={(input) => ({ input, indent })}
      op={op}
      language="yaml"
      downloadFilename="formatted.yaml"
      options={<NumberInput label="Indent" value={indent} onChange={setIndent} min={1} max={8} />}
    />
  );
}

function ValidateTab() {
  const op = useOperation((args: { input: string }) => yamlValidate(args.input));
  return (
    <FormatterShell
      initial={SAMPLE}
      buildInput={(input) => ({ input })}
      op={op as never}
      buttonLabel="Validate"
      renderOutput={(data: { valid: boolean; diagnostics?: { code: string; message: string; severity: 0 | 1 | 2 }[]; output: string }) => (
        <div className="flex flex-col gap-md">
          <div
            className={`flex items-center gap-3 p-md rounded-lg ${
              data.valid ? "bg-tertiary-fixed text-on-tertiary-fixed" : "bg-error/10 text-error"
            }`}
          >
            <CheckCircle2 className="h-6 w-6" />
            <span className="font-body-md font-semibold">
              {data.valid ? "Valid YAML" : "Invalid YAML"}
            </span>
          </div>
        </div>
      )}
    />
  );
}

function ConvertTab() {
  const [to, setTo] = useState<"json" | "yaml">("json");
  const [indent, setIndent] = useState(2);
  const op = useOperation((args: { input: string; to: "json" | "yaml"; indent: number }) =>
    yamlConvert(args.input, args.to, args.indent),
  );
  return (
    <FormatterShell
      initial={SAMPLE}
      buildInput={(input) => ({ input, to, indent })}
      op={op}
      language={to}
      downloadFilename={`output.${to}`}
      buttonLabel="Convert"
      options={
        <div className="grid grid-cols-2 gap-md">
          <Select
            label="Target"
            value={to}
            onChange={(e) => setTo(e.target.value as typeof to)}
            options={[
              { value: "json", label: "→ JSON" },
              { value: "yaml", label: "→ YAML (re-emit)" },
            ]}
          />
          <NumberInput label="Indent" value={indent} onChange={setIndent} min={0} max={8} />
        </div>
      }
    />
  );
}
