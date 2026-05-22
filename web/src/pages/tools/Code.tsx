import { useState } from "react";
import { findTool } from "../../lib/catalog";
import { codeFmtGo, codeFmtXml, codeFmtHtml } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { FormatterShell } from "../../components/layout/FormatterShell";
import { NumberInput } from "../../components/ui/Input";
import { Tabs } from "../../components/ui/Tabs";

const tool = findTool("code")!;

export default function CodePage() {
  const [tab, setTab] = useState<"go" | "xml" | "html">("go");
  return (
    <ToolPage tool={tool}>
      <Tabs
        tabs={[
          { id: "go", label: "Go" },
          { id: "xml", label: "XML" },
          { id: "html", label: "HTML" },
        ]}
        active={tab}
        onChange={(id) => setTab(id as typeof tab)}
        className="mb-md"
      />
      {tab === "go" && <GoTab />}
      {tab === "xml" && <XmlTab />}
      {tab === "html" && <HtmlTab />}
    </ToolPage>
  );
}

function GoTab() {
  const op = useOperation((args: { input: string }) => codeFmtGo(args.input));
  return (
    <FormatterShell
      initial={`package main\n\nimport "fmt"\nfunc main(){fmt.Println("hello, devforge")}\n`}
      buildInput={(input) => ({ input })}
      op={op}
      language="go"
      downloadFilename="formatted.go"
    />
  );
}

function XmlTab() {
  const [indent, setIndent] = useState(2);
  const op = useOperation((args: { input: string; indent: number }) => codeFmtXml(args.input, args.indent));
  return (
    <FormatterShell
      initial={`<root><item id="1">a</item><item id="2">b</item></root>`}
      buildInput={(input) => ({ input, indent })}
      op={op}
      language="xml"
      downloadFilename="formatted.xml"
      options={<NumberInput label="Indent" value={indent} onChange={setIndent} min={0} max={8} />}
    />
  );
}

function HtmlTab() {
  const op = useOperation((args: { input: string }) => codeFmtHtml(args.input));
  return (
    <FormatterShell
      initial={`<div class="card"><h2>Title</h2><p>body</p></div>`}
      buildInput={(input) => ({ input })}
      op={op}
      language="html"
      downloadFilename="formatted.html"
    />
  );
}
