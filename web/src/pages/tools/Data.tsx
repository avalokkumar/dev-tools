import { useState } from "react";
import { findTool } from "../../lib/catalog";
import {
  dataJsonToCsv,
  dataCsvToJson,
  dataJsonToXml,
  dataXmlToJson,
  dataFlatten,
  dataUnflatten,
  dataKeyRename,
} from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { FormatterShell } from "../../components/layout/FormatterShell";
import { Input, Toggle, Select } from "../../components/ui/Input";
import { Tabs } from "../../components/ui/Tabs";
import { Button } from "../../components/ui/Button";
import { Plus, Trash2 } from "lucide-react";

const tool = findTool("data")!;

type Op =
  | "json_to_csv"
  | "csv_to_json"
  | "json_to_xml"
  | "xml_to_json"
  | "flatten"
  | "unflatten"
  | "key_rename";

export default function DataPage() {
  const [op, setOp] = useState<Op>("json_to_csv");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "json_to_csv", label: "JSON → CSV" },
          { id: "csv_to_json", label: "CSV → JSON" },
          { id: "json_to_xml", label: "JSON → XML" },
          { id: "xml_to_json", label: "XML → JSON" },
          { id: "flatten", label: "Flatten" },
          { id: "unflatten", label: "Unflatten" },
          { id: "key_rename", label: "Key rename" },
        ]}
        active={op}
        onChange={(id) => setOp(id as Op)}
        className="mb-md"
      />
      {op === "json_to_csv" && <JsonToCsv />}
      {op === "csv_to_json" && <CsvToJson />}
      {op === "json_to_xml" && <JsonToXml />}
      {op === "xml_to_json" && <XmlToJson />}
      {op === "flatten" && <Flatten />}
      {op === "unflatten" && <Unflatten />}
      {op === "key_rename" && <KeyRename />}
    </ToolPage>
  );
}

const J = `[{"id":1,"name":"a","tags":["x","y"]},{"id":2,"name":"b","tags":["z"]}]`;
const X = `<root><item id="1">a</item><item id="2">b</item></root>`;
const C = `id,name\n1,a\n2,b\n`;

function JsonToCsv() {
  const [sep, setSep] = useState(".");
  const op = useOperation((args: { input: string; flattenSeparator: string }) =>
    dataJsonToCsv(args.input, args.flattenSeparator),
  );
  return (
    <FormatterShell
      initial={J}
      buildInput={(input) => ({ input, flattenSeparator: sep })}
      op={op}
      language="csv"
      buttonLabel="Convert"
      options={<Input label="Flatten separator" value={sep} onChange={(e) => setSep(e.target.value)} />}
    />
  );
}

function CsvToJson() {
  const [header, setHeader] = useState(true);
  const [typedValues, setTypedValues] = useState(true);
  const op = useOperation((args: { input: string; header: boolean; typedValues: boolean }) =>
    dataCsvToJson(args.input, args.header, args.typedValues),
  );
  return (
    <FormatterShell
      initial={C}
      buildInput={(input) => ({ input, header, typedValues })}
      op={op}
      language="json"
      buttonLabel="Convert"
      options={
        <div className="grid grid-cols-2 gap-md">
          <Toggle checked={header} onChange={setHeader} label="Has header" />
          <Toggle checked={typedValues} onChange={setTypedValues} label="Typed values (number/bool)" />
        </div>
      }
    />
  );
}

function JsonToXml() {
  const [root, setRoot] = useState("root");
  const op = useOperation((args: { input: string; root: string; indent: number }) =>
    dataJsonToXml(args.input, args.root, args.indent),
  );
  return (
    <FormatterShell
      initial={`{"item":[{"id":1},{"id":2}]}`}
      buildInput={(input) => ({ input, root, indent: 2 })}
      op={op}
      language="xml"
      buttonLabel="Convert"
      options={<Input label="Root element" value={root} onChange={(e) => setRoot(e.target.value)} />}
    />
  );
}

function XmlToJson() {
  const op = useOperation((args: { input: string }) =>
    dataXmlToJson(args.input, "@", "#text"),
  );
  return (
    <FormatterShell
      initial={X}
      buildInput={(input) => ({ input })}
      op={op}
      language="json"
      buttonLabel="Convert"
    />
  );
}

function Flatten() {
  const [sep, setSep] = useState(".");
  const op = useOperation((args: { input: string; separator: string }) =>
    dataFlatten(args.input, args.separator),
  );
  return (
    <FormatterShell
      initial={`{"user":{"name":"alice","prefs":{"theme":"dark"}}}`}
      buildInput={(input) => ({ input, separator: sep })}
      op={op}
      language="json"
      buttonLabel="Flatten"
      options={<Input label="Separator" value={sep} onChange={(e) => setSep(e.target.value)} />}
    />
  );
}

function Unflatten() {
  const [sep, setSep] = useState(".");
  const op = useOperation((args: { input: string; separator: string }) =>
    dataUnflatten(args.input, args.separator),
  );
  return (
    <FormatterShell
      initial={`{"user.name":"alice","user.prefs.theme":"dark"}`}
      buildInput={(input) => ({ input, separator: sep })}
      op={op}
      language="json"
      buttonLabel="Unflatten"
      options={<Input label="Separator" value={sep} onChange={(e) => setSep(e.target.value)} />}
    />
  );
}

function KeyRename() {
  const [rules, setRules] = useState<{ from: string; to: string; regex?: boolean }[]>([
    { from: "name", to: "fullName" },
  ]);
  const op = useOperation((args: { input: string; rules: typeof rules }) =>
    dataKeyRename(args.input, args.rules),
  );
  return (
    <FormatterShell
      initial={`[{"id":1,"name":"alice"}]`}
      buildInput={(input) => ({ input, rules })}
      op={op}
      language="json"
      buttonLabel="Apply"
      options={
        <div className="flex flex-col gap-2">
          {rules.map((r, i) => (
            <div key={i} className="grid grid-cols-[1fr_1fr_auto_auto] gap-2 items-center">
              <Input
                value={r.from}
                onChange={(e) => setRules((rs) => rs.map((x, j) => (i === j ? { ...x, from: e.target.value } : x)))}
                placeholder="from"
              />
              <Input
                value={r.to}
                onChange={(e) => setRules((rs) => rs.map((x, j) => (i === j ? { ...x, to: e.target.value } : x)))}
                placeholder="to"
              />
              <Select
                value={r.regex ? "regex" : "exact"}
                onChange={(e) =>
                  setRules((rs) =>
                    rs.map((x, j) => (i === j ? { ...x, regex: e.target.value === "regex" } : x)),
                  )
                }
                options={[
                  { value: "exact", label: "exact" },
                  { value: "regex", label: "regex" },
                ]}
              />
              <Button
                size="sm"
                variant="ghost"
                aria-label="Remove rule"
                onClick={() => setRules((rs) => rs.filter((_, j) => j !== i))}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          ))}
          <Button
            size="sm"
            variant="outline"
            onClick={() => setRules((rs) => [...rs, { from: "", to: "" }])}
          >
            <Plus className="h-4 w-4" /> Rule
          </Button>
        </div>
      }
    />
  );
}
