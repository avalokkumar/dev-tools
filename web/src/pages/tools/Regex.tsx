import { useEffect, useState } from "react";
import { Regex as RegexIcon } from "lucide-react";
import { findTool } from "../../lib/catalog";
import { regexTest, regexExplain } from "../../lib/api";
import { useOperation } from "../../lib/useOperation";
import { ToolPage } from "../../components/layout/ToolPage";
import { Card, CardHeader, Badge } from "../../components/ui/Card";
import { Input, Textarea, Toggle } from "../../components/ui/Input";
import { Diagnostics, EmptyState, ErrorBanner } from "../../components/ui/Output";

const tool = findTool("regex")!;

export default function RegexPage() {
  const [pattern, setPattern] = useState(`(?P<word>\\w+)@(?P<host>[\\w.]+)`);
  const [input, setInput] = useState("alice@example.com, bob@dev.io, carol@forge.dev");
  const [flagI, setFlagI] = useState(false);
  const [flagM, setFlagM] = useState(false);
  const [flagS, setFlagS] = useState(false);
  const flags = `${flagI ? "i" : ""}${flagM ? "m" : ""}${flagS ? "s" : ""}`;

  const test = useOperation((args: { pattern: string; input: string; flags: string }) =>
    regexTest(args.pattern, args.input, args.flags),
  );
  const explain = useOperation((p: string) => regexExplain(p));

  useEffect(() => {
    if (pattern.trim()) {
      test.run({ pattern, input, flags });
      explain.run(pattern);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pattern, input, flags]);

  return (
    <ToolPage tool={tool}>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
        <Card>
          <CardHeader title="Pattern" icon={<RegexIcon className="h-5 w-5" />} />
          <div className="flex flex-col gap-md">
            <Input label="Regex (RE2)" value={pattern} onChange={(e) => setPattern(e.target.value)} />
            <div className="flex gap-md flex-wrap">
              <Toggle checked={flagI} onChange={setFlagI} label="Case-insensitive (i)" />
              <Toggle checked={flagM} onChange={setFlagM} label="Multi-line (m)" />
              <Toggle checked={flagS} onChange={setFlagS} label="Dot all (s)" />
            </div>
            <Textarea label="Input" rows={6} value={input} onChange={(e) => setInput(e.target.value)} />
          </div>
        </Card>
        <div className="flex flex-col gap-md">
          <Card>
            <CardHeader
              title="Matches"
              trailing={test.data && <Badge>{test.data.matches.length} matches</Badge>}
            />
            {test.error ? (
              <ErrorBanner error={test.error} />
            ) : !test.data ? (
              <EmptyState title="Type a pattern to highlight matches" />
            ) : test.data.matches.length === 0 ? (
              <EmptyState title="No matches" />
            ) : (
              <ul className="flex flex-col gap-1">
                {test.data.matches.map((m, i) => (
                  <li
                    key={i}
                    className="flex items-start gap-2 px-3 py-2 bg-surface-container-low rounded font-code-block text-code-block"
                  >
                    <span className="font-data-label text-data-label uppercase opacity-60 w-12 shrink-0">
                      {m.start}-{m.end}
                    </span>
                    <code className="flex-1 break-all">
                      <mark className="bg-tertiary-fixed text-on-tertiary-fixed rounded px-0.5">
                        {m.value}
                      </mark>
                      {m.groups && m.groups.length > 0 && (
                        <span className="text-on-surface-variant ml-2">
                          {m.groups.map((g, gi) => (
                            <span key={gi} className="ml-2">
                              <span className="opacity-60">{g.name || gi + 1}:</span>{" "}
                              <span className="text-sky-aqua">{g.value}</span>
                            </span>
                          ))}
                        </span>
                      )}
                    </code>
                  </li>
                ))}
              </ul>
            )}
            <Diagnostics items={test.data?.diagnostics} className="mt-md" />
          </Card>
          <Card>
            <CardHeader title="Explanation" />
            {!explain.data ? (
              <EmptyState title="Pattern explanation appears here" />
            ) : (
              <div className="flex flex-col gap-2">
                <ul className="flex flex-col gap-1">
                  {explain.data.tree.map((t, i) => (
                    <li key={i} className="flex items-start gap-2 px-3 py-2 bg-surface-container-low rounded">
                      <code className="font-code-block text-code-block text-on-surface bg-surface-container px-1 rounded shrink-0">
                        {t.token}
                      </code>
                      <span className="font-body-sm text-body-sm text-on-surface-variant">
                        {t.description}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </Card>
        </div>
      </div>
    </ToolPage>
  );
}
