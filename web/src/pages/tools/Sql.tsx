import { useState } from "react";
import { CheckCircle2 } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { sqlFormat, sqlValidate } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { FormatterShell } from "../../components/layout/FormatterShell";
import { NumberInput, Toggle } from "../../components/ui/Input";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("sql")!;
const SAMPLE = `select id, name from users where status='active' order by created_at desc limit 10;`;

export default function SqlPage() {
  const [tab, setTab] = useState<"format" | "validate">("format");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[{ id: "format", label: "Format" }, { id: "validate", label: "Validate" }]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "format" ? <FormatTab /> : <ValidateTab />}
    </ToolPage>
  );
}

function FormatTab() {
  const [uppercase, setUppercase] = useState(true);
  const [indent, setIndent] = useState(2);
  const op = useOperation((args: { input: string; uppercase: boolean; indentWidth: number }) =>
    sqlFormat(args.input, args.uppercase, args.indentWidth),
  );
  return (
    <FormatterShell
      initial={SAMPLE}
      buildInput={(input) => ({ input, uppercase, indentWidth: indent })}
      op={op}
      language="sql"
      downloadFilename="formatted.sql"
      options={
        <div className="grid grid-cols-2 gap-md">
          <Toggle checked={uppercase} onChange={setUppercase} label="Uppercase keywords" />
          <NumberInput label="Indent" value={indent} onChange={setIndent} min={1} max={8} />
        </div>
      }
    />
  );
}

function ValidateTab() {
  const op = useOperation((args: { input: string }) => sqlValidate(args.input));
  return (
    <FormatterShell
      initial={SAMPLE}
      buildInput={(input) => ({ input })}
      op={op as never}
      buttonLabel="Validate"
      renderOutput={(data: { valid: boolean; output: string }) => (
        <div
          className={`flex items-center gap-3 p-md rounded-lg ${
            data.valid ? "bg-tertiary-fixed text-on-tertiary-fixed" : "bg-error/10 text-error"
          }`}
        >
          <CheckCircle2 className="h-6 w-6" />
          <span className="font-body-md font-semibold">{data.valid ? "Valid SQL" : "Invalid SQL"}</span>
        </div>
      )}
    />
  );
}
