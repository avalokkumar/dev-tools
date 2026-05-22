import { type ReactNode, useState } from "react";
import { Play } from "lucide-react";
import type { Diagnostic } from "../../lib/api";
import { Card, CardHeader } from "../ui/Card";
import { Textarea } from "../ui/Input";
import { Button } from "../ui/Button";
import { CodeBlock, Diagnostics, EmptyState, ErrorBanner } from "../ui/Output";

interface FormatterShellProps<I, O extends { output: string; diagnostics?: Diagnostic[] }> {
  inputLabel?: string;
  outputLabel?: string;
  language?: string;
  downloadFilename?: string;
  initial: string;
  rows?: number;
  /** Map textarea contents + extra options into the operation request */
  buildInput: (text: string) => I;
  /** Operation to invoke */
  op: { run: (input: I) => Promise<O>; loading: boolean; error: unknown; data: O | null };
  /** Optional inputs panel additions (e.g. extra options) */
  options?: ReactNode;
  /** Optional render override for the output content. If omitted, renders CodeBlock(output). */
  renderOutput?: (data: O) => ReactNode;
  buttonLabel?: string;
  inputIcon?: ReactNode;
}

export function FormatterShell<I, O extends { output: string; diagnostics?: Diagnostic[] }>({
  inputLabel = "Input",
  outputLabel = "Output",
  language,
  downloadFilename,
  initial,
  rows = 14,
  buildInput,
  op,
  options,
  renderOutput,
  buttonLabel = "Format",
  inputIcon,
}: FormatterShellProps<I, O>) {
  const [text, setText] = useState(initial);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <Card>
        <CardHeader title={inputLabel} icon={inputIcon} />
        <Textarea rows={rows} value={text} onChange={(e) => setText(e.target.value)} />
        {options && <div className="mt-md">{options}</div>}
        <Button
          onClick={() => op.run(buildInput(text))}
          loading={op.loading}
          fullWidth
          className="mt-md"
          disabled={!text.trim()}
        >
          <Play className="h-4 w-4" /> {buttonLabel}
        </Button>
      </Card>
      <Card padded={false}>
        <CardHeader title={outputLabel} />
        <div className="p-md">
          {op.error ? (
            <ErrorBanner error={op.error} />
          ) : !op.data ? (
            <EmptyState title="Output appears here" />
          ) : renderOutput ? (
            renderOutput(op.data)
          ) : (
            <>
              <CodeBlock
                code={op.data.output}
                language={language}
                download={downloadFilename ? { filename: downloadFilename } : undefined}
              />
              <Diagnostics items={op.data.diagnostics} className="mt-md" />
            </>
          )}
        </div>
      </Card>
    </div>
  );
}
