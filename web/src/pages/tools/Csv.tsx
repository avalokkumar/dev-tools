import { useState } from "react";
import { CheckCircle2 } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { csvFormat, csvValidate } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { FormatterShell } from "../../components/layout/FormatterShell";
import { Input, Toggle } from "../../components/ui/Input";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("csv")!;
const SAMPLE = `id,name,status\n1,Alice,active\n2,Bob,pending\n3,Carol,active\n`;

export default function CsvPage() {
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
  const [delimiter, setDelimiter] = useState(",");
  const [header, setHeader] = useState(true);
  const [alignColumns, setAlignColumns] = useState(true);
  const op = useOperation((args: { input: string; delimiter: string; header: boolean; alignColumns: boolean }) =>
    csvFormat(args.input, args.delimiter, args.header, args.alignColumns),
  );
  return (
    <FormatterShell
      initial={SAMPLE}
      buildInput={(input) => ({ input, delimiter, header, alignColumns })}
      op={op as never}
      language="csv"
      downloadFilename="formatted.csv"
      options={
        <div className="grid grid-cols-3 gap-md">
          <Input label="Delimiter" value={delimiter} onChange={(e) => setDelimiter(e.target.value)} maxLength={1} />
          <Toggle checked={header} onChange={setHeader} label="Has header" />
          <Toggle checked={alignColumns} onChange={setAlignColumns} label="Align columns" />
        </div>
      }
    />
  );
}

function ValidateTab() {
  const [delimiter, setDelimiter] = useState(",");
  const [strict, setStrict] = useState(false);
  const op = useOperation((args: { input: string; delimiter: string; strict: boolean }) =>
    csvValidate(args.input, args.delimiter, undefined, args.strict),
  );
  return (
    <FormatterShell
      initial={SAMPLE}
      buildInput={(input) => ({ input, delimiter, strict })}
      op={op as never}
      buttonLabel="Validate"
      options={
        <div className="grid grid-cols-2 gap-md">
          <Input label="Delimiter" value={delimiter} onChange={(e) => setDelimiter(e.target.value)} maxLength={1} />
          <Toggle checked={strict} onChange={setStrict} label="Strict shape" />
        </div>
      }
      renderOutput={(data: { valid: boolean; output: string }) => (
        <div
          className={`flex items-center gap-3 p-md rounded-lg ${
            data.valid ? "bg-tertiary-fixed text-on-tertiary-fixed" : "bg-error/10 text-error"
          }`}
        >
          <CheckCircle2 className="h-6 w-6" />
          <span className="font-body-md font-semibold">{data.valid ? "Valid CSV" : "Invalid CSV"}</span>
        </div>
      )}
    />
  );
}
